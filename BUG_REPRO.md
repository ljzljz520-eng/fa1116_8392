# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	gestureparticles/cmd/gesture-server	[no test files]
ok  	gestureparticles/internal/analytics	0.001s
ok  	gestureparticles/internal/api	0.010s
?   	gestureparticles/internal/catalog	[no test files]
--- FAIL: Test1116BusinessRegression (0.00s)
    processor_test.go:44: second observation reused previous content: "raise left hand"
FAIL
FAIL	gestureparticles/internal/flow021	0.001s
ok  	gestureparticles/internal/importer	0.010s
ok  	gestureparticles/internal/integration	0.015s
?   	gestureparticles/internal/lesson	[no test files]
ok  	gestureparticles/internal/model	0.002s
?   	gestureparticles/internal/particle	[no test files]
?   	gestureparticles/internal/registry	[no test files]
ok  	gestureparticles/internal/review	0.008s
ok  	gestureparticles/internal/store	0.006s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/gesture-server): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/gesture-server): exit `0`
