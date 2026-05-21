CREATE TABLE categories (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE place_categories (
    place_id UUID NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    PRIMARY KEY (place_id, category_id)
);

CREATE INDEX idx_place_categories_place_id ON place_categories(place_id);
CREATE INDEX idx_place_categories_category_id ON place_categories(category_id);

INSERT INTO categories (id, name, slug, created_at)
VALUES
    ('11111111-1111-1111-1111-000000000001', 'Beautiful places', 'beautiful-places',  now()),
    ('11111111-1111-1111-1111-000000000002', 'Food', 'food',  now()),
    ('11111111-1111-1111-1111-000000000003', 'Coffee', 'coffee',  now()),
    ('11111111-1111-1111-1111-000000000004', 'Nature', 'nature', now()),
    ('11111111-1111-1111-1111-000000000005', 'Museums', 'museums',  now()),
    ('11111111-1111-1111-1111-000000000006', 'Walking', 'walking',  now()),
    ('11111111-1111-1111-1111-000000000007', 'Date places', 'date-places',  now()),
    ('11111111-1111-1111-1111-000000000008', 'Work/study places', 'work-study-places', now()),
    ('11111111-1111-1111-1111-000000000009', 'Sport', 'sport', now()),
    ('11111111-1111-1111-1111-000000000010', 'Hidden gems', 'hidden-gems',  now());