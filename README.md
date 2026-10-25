# Node admission demo

A small Go project exploring a reliability question: when is evidence sufficient to let a compute node proceed?

The evaluator requires passing hardware, certification, and performance records. Each record must belong to the requested node and configuration generation, and remain valid at the evaluation time. Missing, duplicate, failed, mismatched, future-dated, or expired required evidence blocks admission.

## Run

Requires Go 1.22 or newer; no external dependencies, cloud account, or cluster.

```sh
go run .
go vet ./...
go test ./...
go test -race ./...
```

The race command requires a supported Go platform and C toolchain. CI runs it on Linux.

The executable prints four deterministic scenarios: fresh evidence permits admission; expired evidence, a changed configuration generation, and a missing performance check each deny it. All node identities, timestamps, and results are synthetic.

## Design choices

- A caller-provided clock makes expiration boundaries reproducible in tests.
- Exactly one record per required check prevents conflicting results from hiding behind a pass.
- Configuration generation binding prevents reusing evidence after configuration changes.
- Explicit denial reasons explain which requirement blocks admission.

The evaluator is a pure function. It trusts its input and does not collect measurements, verify signatures, authenticate producers, persist state, or implement Kubernetes admission. Unrequired records are ignored. A real service would need those controls plus concurrency handling and an authoritative policy source.

## Scope

This is an independently written educational simulation, not a production certification system. It contains no employer source code, internal design documents, real operational data, or imported project history.
