"""Exercise the running Compose stack through web -> api_gateway -> PostgreSQL.

Only records created by this script are removed. Integration tests in
Docker Compose separately cover schema constraints and real FFmpeg decoding.
"""

import json
import os
import urllib.error
import urllib.request
import uuid

BASE = os.environ.get("CI_WEB_URL", "http://127.0.0.1:18088")
API = os.environ.get("CI_API_URL", "http://127.0.0.1:18080")


def request(method, path, body=None, expected=200, base=BASE):
    data = json.dumps(body).encode("utf-8") if body is not None else None
    req = urllib.request.Request(base + path, data=data, method=method,
                                 headers={"Content-Type": "application/json"})
    try:
        response = urllib.request.urlopen(req, timeout=20)
    except urllib.error.HTTPError as error:
        response = error
    with response:
        raw = response.read()
        assert response.status == expected, f"{method} {path}: {response.status}, expected {expected}"
        if not raw:
            return None
        if "application/json" in response.headers.get("Content-Type", ""):
            return json.loads(raw)
        return raw


def check():
    camera_id = object_id = event_id = None
    try:
        assert request("GET", "/health")["status"] == "ok"
        assert request("GET", "/health", base=API)["status"] == "ok"
        assert b'camera-form' in request("GET", "/")
        for path in ["/app.js", "/camera-video.js", "/camera-source.mjs", "/style.css"]:
            assert request("GET", path), f"empty asset: {path}"

        camera_id = request("POST", "/api/camera", {
            "comment": "CI camera", "ip": "192.0.2.10", "port": 80,
            "stream_url": "rtsp://192.0.2.10/live", "is_enabled": True,
        }, expected=201)
        camera_path = f"/api/camera/{camera_id}"
        original = request("GET", camera_path)
        request("PUT", camera_path, {"comment": "CI updated", "is_enabled": False})
        camera = request("GET", camera_path)
        assert camera["comment"] == "CI updated" and camera["is_enabled"] is False
        assert camera["ip"] == original["ip"] and camera["port"] == original["port"]
        assert camera["created_at"] == original["created_at"]
        assert request("GET", f"/api/camera/count?camera_id={camera_id}&is_enabled=false") == 1
        assert len(request("GET", f"/api/camera/list?camera_id={camera_id}&limit=1&offset=0")) == 1
        request("GET", f"/live/{camera_id}", expected=409)
        request("PUT", camera_path, {"is_enabled": True})

        object_id = request("POST", "/api/object", {"type": f"ci-{uuid.uuid4()}"}, expected=201)
        embedding = [0.0] * 512
        embedding[0], embedding[-1] = 0.125, -0.25
        event_id = request("POST", "/api/event", {
            "camera_id": camera_id, "object_id": object_id, "type": "ci-detection",
            "confidence": 0.9, "embedding": embedding, "photo_uuid": str(uuid.uuid4()),
        }, expected=201)
        event = request("GET", f"/api/event/{event_id}")
        assert event["embedding"] == embedding
        filters = f"camera_id={camera_id}&confidence_min=0.8"
        assert request("GET", f"/api/event/count?{filters}") == 1
        assert len(request("GET", f"/api/event/list?{filters}&limit=15&offset=0")) == 1
        request("DELETE", camera_path, expected=409)
        assert request("GET", camera_path)["id"] == camera_id
        print("PASS: live Compose, web proxy, camera CRUD, filters and pgvector round trip")
    finally:
        if event_id is not None:
            request("DELETE", f"/api/event/{event_id}")
        if object_id is not None:
            request("DELETE", f"/api/object/{object_id}")
        if camera_id is not None:
            path = f"/api/camera/{camera_id}"
            request("DELETE", path)
            request("GET", path, expected=404)
            request("DELETE", path, expected=404)


if __name__ == "__main__":
    check()
