package main

import (
	"encoding/json"
	"image/jpeg"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestLiveRejectsUnavailableSources(t *testing.T) {
	for _, tc := range []struct {
		name     string
		upstream int
		body     string
		expected int
	}{
		{"missing", 404, "", 404},
		{"gateway error", 500, "", 502},
		{"bad JSON", 200, "not-json", 502},
		{"disabled", 200, `{"id":1,"is_enabled":false,"stream_url":"rtsp://camera/live"}`, 409},
		{"webcam is browser local", 200, `{"id":1,"is_enabled":true,"stream_url":"webcam://local"}`, 400},
		{"local file forbidden", 200, `{"id":1,"is_enabled":true,"stream_url":"file:///etc/passwd"}`, 400},
		{"missing URL", 200, `{"id":1,"is_enabled":true}`, 400},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gateway := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/camera/1" {
					t.Errorf("metadata path = %s", r.URL.Path)
				}
				w.WriteHeader(tc.upstream)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer gateway.Close()
			handler, err := newHandler(gateway.URL)
			if err != nil {
				t.Fatal(err)
			}
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/live/1?src=file:///ignored", nil))
			if response.Code != tc.expected {
				t.Fatalf("status = %d; want %d", response.Code, tc.expected)
			}
		})
	}
}

func TestLiveLimitsViewers(t *testing.T) {
	gateway := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(liveCamera{ID: 1, IsEnabled: true, StreamURL: "rtsp://camera/live"})
	}))
	defer gateway.Close()
	target, _ := url.Parse(gateway.URL)
	handler := newLiveHandler(target)
	for i := 0; i < cap(handler.slots); i++ {
		handler.slots <- struct{}{}
	}
	request := httptest.NewRequest(http.MethodGet, "/live/1", nil)
	request.SetPathValue("id", "1")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d", response.Code)
	}
}

func TestLiveSourceValidation(t *testing.T) {
	for _, raw := range []string{"file:///tmp/frame.jpg", "exec:whoami", "rtsp://", "http://camera/live#fragment", "webcam://local"} {
		if _, err := cameraStreamURL(raw); err == nil {
			t.Errorf("accepted unsupported source %q", raw)
		}
	}
	for _, raw := range []string{"rtsp://user:password@192.0.2.1/live", "rtsps://camera/live", "http://camera/video", "https://camera/video"} {
		if _, err := cameraStreamURL(raw); err != nil {
			t.Errorf("valid source rejected: %v", err)
		}
	}
}

// Docker's test image provides real FFmpeg. Local Go-only runs may skip this case.
func TestLiveFFmpeg(t *testing.T) {
	binary, err := exec.LookPath("ffmpeg")
	if err != nil {
		if os.Getenv("FFMPEG_TESTS") == "1" {
			t.Fatal("FFmpeg is required for this test run")
		}
		t.Skip("FFmpeg is not installed; run the Docker test image")
	}
	folder := t.TempDir()
	video := filepath.Join(folder, "source.mp4")
	command := exec.Command(binary, "-nostdin", "-hide_banner", "-loglevel", "error", "-f", "lavfi", "-i", "testsrc=size=160x120:rate=3",
		"-frames:v", "3", "-filter_threads", "1", "-threads", "1", "-c:v", "mpeg4", "-pix_fmt", "yuv420p", "-movflags", "+faststart", video)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("create test video: %v: %s", err, output)
	}
	source := httptest.NewServer(http.FileServer(http.Dir(folder)))
	defer source.Close()
	gateway := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/camera/7" {
			t.Errorf("metadata path = %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(liveCamera{ID: 7, IsEnabled: true, StreamURL: source.URL + "/source.mp4"})
	}))
	defer gateway.Close()
	handler, err := newHandler(gateway.URL)
	if err != nil {
		t.Fatal(err)
	}
	web := httptest.NewServer(handler)
	defer web.Close()
	client := &http.Client{Timeout: 25 * time.Second}
	response, err := client.Get(web.URL + "/live/7?src=file:///ignored")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("preview status = %d", response.StatusCode)
	}
	media, params, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil || media != "multipart/x-mixed-replace" {
		t.Fatal("expected multipart video")
	}
	frames := multipart.NewReader(response.Body, params["boundary"])
	for i := 0; i < 2; i++ {
		part, err := frames.NextPart()
		if err != nil {
			t.Fatalf("frame %d: %v", i, err)
		}
		picture, err := jpeg.Decode(part)
		_ = part.Close()
		if err != nil {
			t.Fatalf("decode frame %d: %v", i, err)
		}
		if picture.Bounds().Dx() < 1 || picture.Bounds().Dx() > 960 || picture.Bounds().Dy() > 540 {
			t.Fatal("invalid frame dimensions")
		}
	}
}

func TestHostCameraSource(t *testing.T) {
	source, _ := url.Parse("rtsp://user:password@localhost:8554/live")
	mapped := hostCameraSource(source, "host.docker.internal")
	if mapped.Host != "host.docker.internal:8554" || mapped.User.String() != source.User.String() || mapped.Path != "/live" {
		t.Fatal("host source was not mapped correctly")
	}
	if source.Host != "localhost:8554" {
		t.Fatal("original URL was modified")
	}
	if hostCameraSource(source, "") != source {
		t.Fatal("local test source must remain unchanged")
	}
	remote, _ := url.Parse("rtsp://192.0.2.10/live")
	if hostCameraSource(remote, "host.docker.internal") != remote {
		t.Fatal("remote camera was remapped")
	}
}
