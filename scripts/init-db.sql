-- PostgreSQL runs this only while initializing an empty data volume.
-- Versioned migrations also enable the extension for existing databases.
CREATE EXTENSION IF NOT EXISTS vector WITH SCHEMA public;
