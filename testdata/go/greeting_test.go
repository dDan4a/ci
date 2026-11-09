package greeting

import "testing"

func TestGreeting(t *testing.T) {
	if got := Greeting("ci"); got != "hello, ci" {
		t.Fatalf("Greeting() = %q, want %q", got, "hello, ci")
	}
}
