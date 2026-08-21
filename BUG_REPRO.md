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
?   	parkvisitor/cmd/visitor-sync	[no test files]
?   	parkvisitor/internal/analytics	[no test files]
ok  	parkvisitor/internal/clock	0.001s
ok  	parkvisitor/internal/domain	0.001s
?   	parkvisitor/internal/importer	[no test files]
?   	parkvisitor/internal/policy	[no test files]
?   	parkvisitor/internal/report	[no test files]
--- FAIL: TestBusiness37Regression (0.01s)
    business37_test.go:20: inconsistent repeated result={Batch:{ID:RB684-37 Reference:RB684-37 Source:api State:confirmed Total:1 Valid:1 Invalid:0 BusinessDate:2026-08-21 ConfirmedAt:2026-08-21T08:00:00Z} Records:[{ID:rb684-37-visitor BatchID:RB684-37 Name:Boundary Guest Company:Acme Host:Gate VisitDate:2026-08-20 Status:validated Notes: Tags:[] CreatedAt:2026-08-21T08:00:00Z UpdatedAt:2026-08-21T08:00:00Z}] Issues:map[]}
FAIL
FAIL	parkvisitor/internal/service	0.024s
ok  	parkvisitor/internal/storage	0.012s
ok  	parkvisitor/internal/transport	0.005s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/visitor-sync): exit `0`
- Frontend build (web): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/visitor-sync): exit `0`
- Frontend build (web): exit `0`
