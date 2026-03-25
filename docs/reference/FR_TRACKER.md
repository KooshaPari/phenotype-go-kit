# FR Tracker - phenotype-go-kit

| FR ID | Description | Status | Test Location |
|-------|-------------|--------|---------------|
| FR-LOG-001 | Logger injection | Implemented | `logctx/logctx_test.go` |
| FR-LOG-002 | Missing logger panic | Implemented | `logctx/logctx_test.go` |
| FR-BUF-001 | Fixed capacity | Implemented | `ringbuffer/ringbuffer_test.go` |
| FR-BUF-002 | Overwrite on full | Implemented | `ringbuffer/ringbuffer_test.go` |
| FR-BUF-003 | Oldest-first retrieval | Implemented | `ringbuffer/ringbuffer_test.go` |
| FR-WAIT-001 | Exponential backoff | Implemented | `waitfor/waitfor_test.go` |
| FR-WAIT-002 | Timeout | Implemented | `waitfor/waitfor_test.go` |
| FR-WAIT-003 | Testable clock | Implemented | `waitfor/waitfor_test.go` |
| FR-REG-001 | Owner-scoped registration | Implemented | `registry/registry_test.go` |
| FR-REG-002 | Ref-counted removal | Implemented | `registry/registry_test.go` |
| FR-REG-003 | Change hooks | Implemented | `registry/registry_test.go` |
| FR-REG-004 | Thread safety | Implemented | `registry/registry_test.go` |
