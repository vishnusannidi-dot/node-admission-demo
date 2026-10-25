package main

import (
	"testing"
	"time"
)

func TestAdmissionBoundaries(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name    string
		edit    func(*Evidence)
		allowed bool
	}{
		{"valid", func(e *Evidence) {}, true},
		{"failed", func(e *Evidence) { e.Result = "FAIL" }, false},
		{"unknown", func(e *Evidence) { e.Result = "UNKNOWN" }, false},
		{"wrong node", func(e *Evidence) { e.Node = "other-node" }, false},
		{"old generation", func(e *Evidence) { e.Generation = 1 }, false},
		{"missing observation", func(e *Evidence) { e.Observed = time.Time{} }, false},
		{"future observation", func(e *Evidence) { e.Observed = now.Add(time.Minute) }, false},
		{"expired at boundary", func(e *Evidence) { e.Expires = now }, false},
		{"invalid interval", func(e *Evidence) { e.Expires = e.Observed }, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Evidence{"hardware", "demo-node", 2, "PASS", now.Add(-time.Minute), now.Add(time.Hour)}
			tt.edit(&e)
			got := Evaluate("demo-node", 2, []string{"hardware"}, []Evidence{e}, now)
			if got.Allowed != tt.allowed || (!got.Allowed && len(got.Reasons) == 0) {
				t.Fatalf("unexpected decision: %+v", got)
			}
		})
	}
}

func TestMissingAndConflictingEvidence(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	pass := Evidence{"hardware", "demo-node", 2, "PASS", now.Add(-time.Minute), now.Add(time.Hour)}
	fail := pass
	fail.Result = "FAIL"
	for _, records := range [][]Evidence{nil, {pass, fail}, {fail, pass}, {pass, pass}} {
		if got := Evaluate("demo-node", 2, []string{"hardware"}, records, now); got.Allowed {
			t.Fatalf("missing or duplicate evidence allowed: %+v", records)
		}
	}
	if got := Evaluate("demo-node", 2, []string{"hardware", "performance"}, []Evidence{pass}, now); got.Allowed {
		t.Fatal("missing required performance check allowed")
	}
}

func TestInvalidRequests(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	for _, required := range [][]string{nil, {""}, {"hardware", "hardware"}} {
		if Evaluate("demo-node", 2, required, nil, now).Allowed {
			t.Fatal("invalid policy allowed")
		}
	}
	for _, request := range []struct {
		node       string
		generation int
		now        time.Time
	}{
		{"", 2, now}, {"demo-node", 0, now}, {"demo-node", 2, time.Time{}},
	} {
		if Evaluate(request.node, request.generation, []string{"hardware"}, nil, request.now).Allowed {
			t.Fatal("invalid request allowed")
		}
	}
}
