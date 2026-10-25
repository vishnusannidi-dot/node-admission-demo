package main

import (
	"encoding/json"
	"os"
	"time"
)

func main() {
	// Fixed clock and invented identities make the example reproducible.
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	required := []string{"hardware", "certification", "performance"}
	records := make([]Evidence, 0, len(required))
	for _, check := range required {
		records = append(records, Evidence{Check: check, Node: "demo-node", Generation: 2,
			Result: "PASS", Observed: now.Add(-time.Minute), Expires: now.Add(time.Hour)})
	}
	report := map[string]Decision{
		"fresh":               Evaluate("demo-node", 2, required, records, now),
		"expired":             Evaluate("demo-node", 2, required, records, now.Add(time.Hour)),
		"changed_generation":  Evaluate("demo-node", 3, required, records, now),
		"missing_performance": Evaluate("demo-node", 2, required, records[:2], now),
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil {
		os.Exit(1)
	}
}
