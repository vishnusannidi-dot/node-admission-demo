package main

import "time"

// Evidence is synthetic input; this demo does not verify signatures or run checks.
type Evidence struct {
	Check      string
	Node       string
	Generation int
	Result     string
	Observed   time.Time
	Expires    time.Time
}

type Decision struct {
	Allowed bool     `json:"allowed"`
	Reasons []string `json:"reasons"`
}

// Evaluate requires one passing, current record for each required check.
// It has no side effects and uses the caller's clock for repeatable decisions.
func Evaluate(node string, generation int, required []string, records []Evidence, now time.Time) Decision {
	reasons := make([]string, 0)
	if node == "" || generation < 1 || now.IsZero() || len(required) == 0 {
		return Decision{Reasons: []string{"invalid evaluation request"}}
	}
	seen := make(map[string]bool)
	for _, check := range required {
		if check == "" || seen[check] {
			return Decision{Reasons: []string{"invalid required checks"}}
		}
		seen[check] = true
		var matches []Evidence
		for _, record := range records {
			if record.Check == check {
				matches = append(matches, record)
			}
		}
		if len(matches) != 1 {
			reasons = append(reasons, check+": expected exactly one record")
			continue
		}
		record := matches[0]
		switch {
		case record.Node != node || record.Generation != generation:
			reasons = append(reasons, check+": identity or generation mismatch")
		case record.Observed.IsZero() || record.Observed.After(now) || !record.Expires.After(record.Observed) || !now.Before(record.Expires):
			reasons = append(reasons, check+": evidence outside validity interval")
		case record.Result != "PASS":
			reasons = append(reasons, check+": result is not PASS")
		}
	}
	return Decision{Allowed: len(reasons) == 0, Reasons: reasons}
}
