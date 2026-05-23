CREATE TABLE place_reports (
    id UUID PRIMARY KEY,
    place_id UUID NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    reported_by_user_id UUID NOT NULL,
    reason VARCHAR(50) NOT NULL,
    comment VARCHAR(500) NOT NULL DEFAULT '',
    status VARCHAR(32) NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMPTZ NOT NULL,
    resolved_at TIMESTAMPTZ,

    CONSTRAINT place_reports_unique_user_per_place UNIQUE (place_id, reported_by_user_id),
    CONSTRAINT place_reports_reason_valid CHECK (
        reason IN (
            'SPAM',
            'OFFENSIVE_CONTENT',
            'WRONG_INFO',
            'DUPLICATE',
            'CLOSED_PLACE',
            'OTHER'
        )
    ),
    CONSTRAINT place_reports_status_valid CHECK (
        status IN ('OPEN', 'RESOLVED', 'REJECTED')
    )
);

CREATE INDEX idx_place_reports_place_id ON place_reports(place_id);
CREATE INDEX idx_place_reports_status ON place_reports(status);
CREATE INDEX idx_place_reports_created_at ON place_reports(created_at);
