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
?   	hospitalportal/internal/report	[no test files]
--- FAIL: TestDepartmentLookupKeepsNotFoundContext (0.01s)
    integration_test.go:50: missing department unexpectedly succeeded: api.Result{OK:true, Kind:"department", Message:"no department found", Rows:[]string(nil), Error:""}
FAIL
FAIL	hospitalportal	0.027s
ok  	hospitalportal/cmd/portal	0.002s
ok  	hospitalportal/internal/api	0.011s
ok  	hospitalportal/internal/audit	0.011s
ok  	hospitalportal/internal/domain	0.002s
ok  	hospitalportal/internal/service	0.015s
ok  	hospitalportal/internal/store	0.014s
ok  	hospitalportal/internal/validate	0.001s
ok  	hospitalportal/internal/workflow	0.009s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/portal): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/portal): exit `0`
