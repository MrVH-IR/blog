CREATE TABLE tags(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name JSONB NOT NULL,
    taggable_type VARCHAR(50) NOT NULL,
    taggable_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tags_taggable ON tags (taggable_type, taggable_id);