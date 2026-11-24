-- +goose Up
CREATE TABLE users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username TEXT NOT NULL UNIQUE CHECK (length(btrim(username)) BETWEEN 1 AND 100),
    name TEXT NOT NULL CHECK (length(btrim(name)) BETWEEN 1 AND 200),
    password_hash TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE locations (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL CHECK (length(btrim(name)) BETWEEN 1 AND 200),
    type TEXT NOT NULL CHECK (length(btrim(type)) BETWEEN 1 AND 100),
    description TEXT NOT NULL DEFAULT '',
    coordinates JSONB CHECK (coordinates IS NULL OR jsonb_typeof(coordinates) = 'object')
);

CREATE TABLE cameras (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    location_id BIGINT NOT NULL REFERENCES locations(id) ON DELETE RESTRICT,
    name TEXT NOT NULL CHECK (length(btrim(name)) BETWEEN 1 AND 200),
    stream_url TEXT,
    resolution TEXT,
    fps NUMERIC(7, 3) CHECK (fps > 0),
    status TEXT NOT NULL DEFAULT 'offline' CHECK (status IN ('offline', 'online', 'disabled')),
    UNIQUE (id, location_id)
);

CREATE TABLE persons (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT UNIQUE REFERENCES users(id) ON DELETE RESTRICT,
    visual_features JSONB CHECK (visual_features IS NULL OR jsonb_typeof(visual_features) = 'object'),
    first_seen TIMESTAMPTZ,
    last_seen TIMESTAMPTZ,
    CHECK (first_seen IS NULL OR last_seen IS NULL OR first_seen <= last_seen)
);

CREATE TABLE object_types (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL UNIQUE CHECK (length(btrim(name)) BETWEEN 1 AND 100),
    description TEXT NOT NULL DEFAULT ''
);

CREATE TABLE objects (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    object_type_id BIGINT NOT NULL REFERENCES object_types(id) ON DELETE RESTRICT,
    name TEXT NOT NULL CHECK (length(btrim(name)) BETWEEN 1 AND 200),
    status TEXT NOT NULL DEFAULT 'unknown' CHECK (status IN ('unknown', 'visible', 'lost', 'found')),
    image_path TEXT,
    first_seen TIMESTAMPTZ,
    last_seen TIMESTAMPTZ,
    CHECK (first_seen IS NULL OR last_seen IS NULL OR first_seen <= last_seen)
);

CREATE TABLE event_types (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL UNIQUE CHECK (length(btrim(name)) BETWEEN 1 AND 100),
    description TEXT NOT NULL DEFAULT ''
);

CREATE TABLE events (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    event_type_id BIGINT NOT NULL REFERENCES event_types(id) ON DELETE RESTRICT,
    camera_id BIGINT NOT NULL,
    location_id BIGINT NOT NULL REFERENCES locations(id) ON DELETE RESTRICT,
    person_id BIGINT REFERENCES persons(id) ON DELETE RESTRICT,
    object_id BIGINT REFERENCES objects(id) ON DELETE RESTRICT,
    occurred_at TIMESTAMPTZ NOT NULL,
    confidence DOUBLE PRECISION CHECK (confidence >= 0 AND confidence <= 1),
    description TEXT NOT NULL DEFAULT '',
    FOREIGN KEY (camera_id, location_id) REFERENCES cameras(id, location_id) ON DELETE RESTRICT,
    UNIQUE (id, camera_id)
);

CREATE TABLE video_segments (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    camera_id BIGINT NOT NULL REFERENCES cameras(id) ON DELETE RESTRICT,
    event_id BIGINT,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    storage_reference TEXT NOT NULL CHECK (length(btrim(storage_reference)) > 0),
    preview TEXT,
    duration_seconds DOUBLE PRECISION GENERATED ALWAYS AS (extract(epoch FROM end_time - start_time)::DOUBLE PRECISION) STORED,
    CHECK (end_time > start_time),
    FOREIGN KEY (event_id, camera_id) REFERENCES events(id, camera_id) ON DELETE RESTRICT
);

CREATE TABLE trajectories (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    person_id BIGINT NOT NULL REFERENCES persons(id) ON DELETE RESTRICT,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    start_location_id BIGINT REFERENCES locations(id) ON DELETE RESTRICT,
    end_location_id BIGINT REFERENCES locations(id) ON DELETE RESTRICT,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'completed', 'lost')),
    direction TEXT,
    CHECK (end_time IS NULL OR end_time >= start_time)
);

CREATE TABLE trajectory_points (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    trajectory_id BIGINT NOT NULL REFERENCES trajectories(id) ON DELETE RESTRICT,
    video_segment_id BIGINT REFERENCES video_segments(id) ON DELETE RESTRICT,
    observed_at TIMESTAMPTZ NOT NULL,
    x DOUBLE PRECISION NOT NULL CHECK (x > '-Infinity'::FLOAT8 AND x < 'Infinity'::FLOAT8),
    y DOUBLE PRECISION NOT NULL CHECK (y > '-Infinity'::FLOAT8 AND y < 'Infinity'::FLOAT8),
    z DOUBLE PRECISION CHECK (z > '-Infinity'::FLOAT8 AND z < 'Infinity'::FLOAT8)
);

CREATE TABLE query_history (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    text TEXT NOT NULL CHECK (length(btrim(text)) BETWEEN 1 AND 4000),
    response TEXT,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at TIMESTAMPTZ,
    CHECK (completed_at IS NULL OR completed_at >= created_at)
);

CREATE TABLE query_events (
    query_id BIGINT NOT NULL REFERENCES query_history(id) ON DELETE RESTRICT,
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE RESTRICT,
    PRIMARY KEY (query_id, event_id)
);

-- Index referencing columns and the common time-based lookup paths.
CREATE INDEX cameras_location_idx ON cameras(location_id);
CREATE INDEX objects_type_idx ON objects(object_type_id);
CREATE INDEX events_type_idx ON events(event_type_id);
CREATE INDEX events_camera_time_idx ON events(camera_id, occurred_at DESC);
CREATE INDEX events_location_idx ON events(location_id);
CREATE INDEX events_object_time_idx ON events(object_id, occurred_at DESC);
CREATE INDEX events_person_time_idx ON events(person_id, occurred_at DESC);
CREATE INDEX video_segments_camera_time_idx ON video_segments(camera_id, start_time);
CREATE INDEX video_segments_event_idx ON video_segments(event_id, camera_id);
CREATE INDEX trajectories_person_time_idx ON trajectories(person_id, start_time);
CREATE INDEX trajectories_start_location_idx ON trajectories(start_location_id);
CREATE INDEX trajectories_end_location_idx ON trajectories(end_location_id);
CREATE INDEX trajectory_points_trajectory_time_idx ON trajectory_points(trajectory_id, observed_at);
CREATE INDEX trajectory_points_video_idx ON trajectory_points(video_segment_id);
CREATE INDEX query_history_user_time_idx ON query_history(user_id, created_at DESC);
CREATE INDEX query_events_event_idx ON query_events(event_id);

-- +goose Down
DROP TABLE query_events;
DROP TABLE query_history;
DROP TABLE trajectory_points;
DROP TABLE trajectories;
DROP TABLE video_segments;
DROP TABLE events;
DROP TABLE event_types;
DROP TABLE objects;
DROP TABLE object_types;
DROP TABLE persons;
DROP TABLE cameras;
DROP TABLE locations;
DROP TABLE users;
