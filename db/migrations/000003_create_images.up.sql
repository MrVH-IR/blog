CREATE TABLE images(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    imageable_type VARCHAR(50) NOT NULL,
    imageable_id BIGINT NOT NULL,
    path TEXT NOT NULL,
    alt VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE index idx_images_imageable ON images(imageable_id, imageable_type);