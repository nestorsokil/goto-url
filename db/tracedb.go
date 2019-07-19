package db

import (
	"context"
	"github.com/opentracing/opentracing-go"
)

// traceDb is a wrapper that traces db methods
type traceDb struct {
	actual DataStorage
}

func (t *traceDb) Find(ctx context.Context, key string) (*Record, error) {
	span, newCtx := opentracing.StartSpanFromContext(ctx, "db.Find")
	defer span.Finish()
	return t.actual.Find(newCtx, key)
}

func (t *traceDb) SaveWithExpiration(ctx context.Context, record *Record, expireIn int64) error {
	span, newCtx := opentracing.StartSpanFromContext(ctx, "db.SaveWithExpiration")
	defer span.Finish()
	return t.actual.SaveWithExpiration(newCtx, record, expireIn)
}

func (t *traceDb) Exists(ctx context.Context, key string) (bool, error) {
	span, newCtx := opentracing.StartSpanFromContext(ctx, "db.Exists")
	defer span.Finish()
	return t.actual.Exists(newCtx, key)
}

func (t *traceDb) Shutdown(ctx context.Context) {
	t.actual.Shutdown(ctx)
}
