package client

import (
	"os"
	"regexp"
	"testing"
)

// Every o11y route is /v1/o11y/…; the host already says api, so a path never
// does. A call under /api/ answers 404.
func TestEveryCallIsUnderV1O11y(t *testing.T) {
	src, err := os.ReadFile("client.go")
	if err != nil {
		t.Fatal(err)
	}
	if m := regexp.MustCompile(`"%s/api/[^"]*"|"/api/[^"]*"`).FindAll(src, -1); len(m) > 0 {
		t.Fatalf("client.go calls routes under /api/: %q", m)
	}
	calls := regexp.MustCompile(`"%s(/[^"?]*)`).FindAllSubmatch(src, -1)
	if len(calls) == 0 {
		t.Fatal("client.go names no routes")
	}
	under := regexp.MustCompile(`^/v1/o11y/`)
	for _, c := range calls {
		if !under.Match(c[1]) {
			t.Errorf("client.go calls %s, outside /v1/o11y/", c[1])
		}
	}
}
