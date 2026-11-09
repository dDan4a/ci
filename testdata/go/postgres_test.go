//go:build integration

package greeting

import (
	"net"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestPostgresAcceptsConnections(t *testing.T) {
	dsn, err := url.Parse(os.Getenv("FIXTURE_PG_DSN"))
	if err != nil || dsn.Host == "" {
		t.Fatalf("FIXTURE_PG_DSN is not a Postgres URL: %q", os.Getenv("FIXTURE_PG_DSN"))
	}

	dialer := net.Dialer{Timeout: 5 * time.Second}

	conn, err := dialer.DialContext(t.Context(), "tcp", dsn.Host)
	if err != nil {
		t.Fatalf("dial %s: %v", dsn.Host, err)
	}

	if err := conn.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
}
