CREATE TABLE places (
    id UUID PRIMARY KEY,

    title VARCHAR(120) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    address VARCHAR(255) NOT NULL DEFAULT '',
    city VARCHAR(100) NOT NULL,

    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,

    created_by_user_id UUID NOT NULL,
    status VARCHAR(32) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT places_title_not_empty CHECK (length(trim(title)) > 0),
    CONSTRAINT places_description_length CHECK (length(description) <= 2000),
    CONSTRAINT places_city_not_empty CHECK (length(trim(city)) > 0),
    CONSTRAINT places_latitude_range CHECK (latitude >= -90 AND latitude <= 90),
    CONSTRAINT places_longitude_range CHECK (longitude >= -180 AND longitude <= 180),
    CONSTRAINT places_status_valid CHECK (
        status IN ('PENDING_REVIEW', 'APPROVED', 'REJECTED', 'ARCHIVED')
    )
);