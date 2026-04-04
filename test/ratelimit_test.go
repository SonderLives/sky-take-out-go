package test

import (
	"context"
	"goflow/internal/pkg/ratelimit"
	"testing"
	"time"
)

func TestBuildKeyIncludesMethod(t *testing.T) {
	getKey := ratelimit.BuildKey("127.0.0.1", "GET", "/admin/v1/products")
	postKey := ratelimit.BuildKey("127.0.0.1", "POST", "/admin/v1/products")

	if getKey == postKey {
		t.Fatalf("expected different keys for different methods, got same key: %s", getKey)
	}
}

func TestRetryAfterSecondsRoundsUp(t *testing.T) {
	got := ratelimit.RetryAfterSeconds(1500 * time.Millisecond)
	if got != "2" {
		t.Fatalf("expected retry-after 2, got %s", got)
	}
}

func TestAllowInvalidLimit(t *testing.T) {
	limiter := ratelimit.NewRateLimiter(nil)
	_, _, err := limiter.Allow(context.Background(), "k", 0, time.Second)
	if err == nil {
		t.Fatalf("expected error when limit is invalid")
	}
}

func TestAllowInvalidWindow(t *testing.T) {
	limiter := ratelimit.NewRateLimiter(nil)
	_, _, err := limiter.Allow(context.Background(), "k", 1, 0)
	if err == nil {
		t.Fatalf("expected error when window is invalid")
	}
}
