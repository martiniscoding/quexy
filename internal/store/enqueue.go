package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

const jobCols = `id, tenant_id, queue, kind, payload, status, priority,
    attempt, max_attempts, run_at, created_at, updated_at,
    idempotency_key, last_error`

func scanJob(row pgx.Row) (Job, error) {
	var j Job
	err := row.Scan(&j.ID, &j.TenantID, &j.Queue, &j.Kind, &j.Payload,
		&j.Status, &j.Priority, &j.Attempt, &j.MaxAttempts,
		&j.RunAt, &j.CreatedAt, &j.UpdatedAt,
		&j.IdempotencyKey, &j.LastError)
	return j, err
}
func (s *Store) GetByID(ctx context.Context, id int64) (Job, error) {
	q := `SELECT ` + jobCols + ` FROM jobs WHERE id = $1`

	j, err := scanJob(s.pool.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return Job{}, ErrNotFound
	}
	if err != nil {
		return Job{}, fmt.Errorf("get job: %w", err)
	}
	return j, nil
}

func (s *Store) GetByIdempotencyKey(ctx context.Context, key string) (Job, error) {
	q := `SELECT ` + jobCols + ` FROM jobs WHERE idempotency_key = $1`
	j, err := scanJob(s.pool.QueryRow(ctx, q, key))
	if errors.Is(err, pgx.ErrNoRows) {
		return Job{}, ErrNotFound
	}
	if err != nil {
		return Job{}, fmt.Errorf("get job: %w", err)
	}
	return j, nil
}

func (s *Store) EnqueuJob(ctx context.Context, tenant_id string, queue string, kind string,
	payload json.RawMessage, priority int, max_attempts int, run_at time.Time , idempotency_key string) (Job, error) {
	q := `INSERT INTO jobs (
        tenant_id, queue, kind, payload, priority,
        max_attempts, run_at, idempotency_key
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    ON CONFLICT (idempotency_key) DO NOTHING
    RETURNING ` + jobCols

	var runAt = time.Time{}

	j, err := scanJob(s.pool.QueryRow(ctx, q, tenant_id, queue, kind, payload, priority, max_attempts,
		runAt, idempotency_key))
	if errors.Is(err, pgx.ErrNoRows) {
		return Job{}, ErrNotFound
	}
	if err != nil {
		return Job{}, fmt.Errorf("could not enqueue %w", err)
	}
	return j, nil
}
