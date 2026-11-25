package main

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type liveCamera struct {
	ID        int    `json:"id"`
	StreamURL string `json:"stream_url"`
	IsEnabled bool   `json:"is_enabled"`
}

type liveHandler struct {
	gateway *url.URL
	client  *http.Client
	slots   chan struct{}
}

func newLiveHandler(gateway *url.URL) *liveHandler {
	return &liveHandler{gateway: gateway, client: &http.Client{Timeout: 5 * time.Second}, slots: make(chan struct{}, 4)}
}

func cameraStreamURL(raw string) (*url.URL, error) {
	source, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return nil, err
	}
	switch source.Scheme {
	case "rtsp", "rtsps", "http", "https":
		if source.Hostname() != "" && source.Fragment == "" {
			return source, nil
		}
	}
	return nil, &url.Error{Op: "parse", URL: "camera source", Err: errInvalidCameraSource{}}
}

type errInvalidCameraSource struct{}

func (errInvalidCameraSource) Error() string { return "use an RTSP or HTTP(S) camera URL" }

func hostCameraSource(source *url.URL, gateway string) *url.URL {
	host := source.Hostname()
	ip := net.ParseIP(host)
	if gateway == "" || (!strings.EqualFold(host, "localhost") && (ip == nil || !ip.IsLoopback())) {
		return source
	}
	copy := *source
	if port := source.Port(); port != "" {
		copy.Host = net.JoinHostPort(gateway, port)
	} else {
		copy.Host = gateway
	}
	return &copy
}

func ffmpegArguments(source *url.URL) []string {
	args := []string{"-nostdin", "-hide_banner", "-loglevel", "error", "-threads", "1", "-filter_threads", "1",
		"-rw_timeout", "12000000", "-protocol_whitelist", "http,https,tcp,tls,udp,rtp,crypto"}
	if source.Scheme == "rtsp" || source.Scheme == "rtsps" {
		args = append(args, "-rtsp_transport", "tcp", "-timeout", "12000000")
	}
	return append(args, "-i", source.String(), "-map", "0:v:0", "-an", "-vf", "fps=10,scale=960:540:force_original_aspect_ratio=decrease",
		"-threads", "1", "-c:v", "mjpeg", "-q:v", "5", "-f", "mpjpeg", "-boundary_tag", "frame", "pipe:1")
}

func (h *liveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.Error(w, "Invalid camera ID", http.StatusBadRequest)
		return
	}
	target := *h.gateway
	target.Path = "/api/camera/" + strconv.Itoa(id)
	request, err := http.NewRequestWithContext(r.Context(), http.MethodGet, target.String(), nil)
	if err != nil {
		http.Error(w, "Camera unavailable", http.StatusBadGateway)
		return
	}
	for _, header := range []string{"Authorization", "Cookie"} {
		if value := r.Header.Get(header); value != "" {
			request.Header.Set(header, value)
		}
	}
	response, err := h.client.Do(request)
	if err != nil {
		http.Error(w, "Gateway unavailable", http.StatusBadGateway)
		return
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusNotFound {
		http.Error(w, "Camera not found", http.StatusNotFound)
		return
	}
	if response.StatusCode != http.StatusOK {
		http.Error(w, "Camera unavailable", http.StatusBadGateway)
		return
	}
	var camera liveCamera
	if err := json.NewDecoder(io.LimitReader(response.Body, 65536)).Decode(&camera); err != nil {
		http.Error(w, "Invalid camera data", http.StatusBadGateway)
		return
	}
	if !camera.IsEnabled {
		http.Error(w, "Camera disabled", http.StatusConflict)
		return
	}
	source, err := cameraStreamURL(camera.StreamURL)
	if err != nil {
		http.Error(w, "No supported network stream", http.StatusBadRequest)
		return
	}
	source = hostCameraSource(source, os.Getenv("CAMERA_LOCAL_HOST"))
	_ = response.Body.Close()
	select {
	case h.slots <- struct{}{}:
		defer func() { <-h.slots }()
	default:
		http.Error(w, "Too many live viewers", http.StatusTooManyRequests)
		return
	}
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	command := exec.CommandContext(ctx, "ffmpeg", ffmpegArguments(source)...)
	output, err := command.StdoutPipe()
	if err != nil {
		http.Error(w, "Preview unavailable", http.StatusServiceUnavailable)
		return
	}
	// Do not expose the input URL (which may contain a camera password) in logs or responses.
	command.Stderr = io.Discard
	if err := command.Start(); err != nil {
		http.Error(w, "Preview unavailable", http.StatusServiceUnavailable)
		return
	}
	defer func() { cancel(); _ = output.Close(); _ = command.Wait() }()
	firstFrameTimeout := time.AfterFunc(15*time.Second, cancel)
	defer firstFrameTimeout.Stop()
	first := make([]byte, 32768)
	n, err := output.Read(first)
	firstFrameTimeout.Stop()
	if n == 0 {
		http.Error(w, "Cannot open camera stream", http.StatusBadGateway)
		return
	}
	// A preview is long-lived; it must not inherit the server's normal response deadline.
	_ = http.NewResponseController(w).SetWriteDeadline(time.Time{})
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	w.Header().Set("Cache-Control", "no-store")
	writer := liveWriter{ResponseWriter: w}
	if _, err := writer.Write(first[:n]); err != nil {
		return
	}
	_, _ = io.Copy(writer, output)
}

type liveWriter struct{ http.ResponseWriter }

func (w liveWriter) Write(data []byte) (int, error) {
	n, err := w.ResponseWriter.Write(data)
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
	return n, err
}
