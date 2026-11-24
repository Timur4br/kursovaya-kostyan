CREATE EXTENSION IF NOT EXISTS vector;

-- ============================================
-- camera
-- ============================================
CREATE TABLE camera (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    comment VARCHAR,
    stream_url VARCHAR,
    ip VARCHAR NOT NULL,
    http_port INTEGER NOT NULL DEFAULT 80,
    is_enabled BOOLEAN DEFAULT TRUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- ============================================
-- users
-- ============================================
CREATE TABLE users (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    login VARCHAR NOT NULL UNIQUE,
    password VARCHAR NOT NULL,
    username VARCHAR NOT NULL,
    last_login TIMESTAMP WITH TIME ZONE,
    role_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    FOREIGN KEY (role_id) REFERENCES role(id)
);

-- ============================================
-- role
-- ============================================
CREATE TABLE role (
    id INTEGER GENERATED ALWAYS AS IDENTITY NOT NULL PRIMARY KEY,
    name TEXT NOT NULL
);

-- ============================================
-- authority
-- ============================================
CREATE TABLE authority (
    id INTEGER GENERATED ALWAYS AS IDENTITY NOT NULL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- ============================================
-- user_authority
-- ============================================
CREATE TABLE user_authority (
    id INTEGER GENERATED ALWAYS AS IDENTITY NOT NULL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    authority_id INTEGER NOT NULL,
    CONSTRAINT fk_user_authority_user
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE,

    CONSTRAINT fk_user_authority_authority
        FOREIGN KEY (authority_id)
        REFERENCES authority (id)
        ON DELETE CASCADE
);

-- ============================================
-- objects
-- ============================================
CREATE TABLE objects (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    type VARCHAR NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- ============================================
-- events
-- ============================================
CREATE TABLE events (
    "uuid" UUID PRIMARY KEY gen_random_uuid(),
    camera_id INTEGER NOT NULL,
    object_id INTEGER NOT NULL,
    type VARCHAR NOT NULL,
    confidence FLOAT NOT NULL,
    embedding vector(512) NOT NULL,
    photo_uuid UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    FOREIGN KEY (camera_id) REFERENCES camera(id),
    FOREIGN KEY (object_id) REFERENCES objects(id)
)
