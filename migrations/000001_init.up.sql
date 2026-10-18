CREATE TYPE job_status AS ENUM (
    'pending', 'running', 'succeeded', 'failed', 'dead', 'cancelled'
);

CREATE TABLE jobs (
    id              BIGSERIAL PRIMARY KEY,
    tenant_id       TEXT        NOT NULL,
    queue           TEXT        NOT NULL DEFAULT 'default',
    kind            TEXT        NOT NULL,
    payload         JSONB       NOT NULL DEFAULT '{}'::jsonb,
    status          job_status  NOT NULL DEFAULT 'pending',
    priority        SMALLINT    NOT NULL DEFAULT 0,
    attempt         INT         NOT NULL DEFAULT 0,
    max_attempts    INT         NOT NULL DEFAULT 5,
    run_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    idempotency_key TEXT UNIQUE,
    last_error      TEXT
);

CREATE INDEX idx_jobs_tenant_created
    ON jobs (tenant_id, created_at DESC, id DESC);