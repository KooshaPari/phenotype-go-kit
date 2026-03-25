# Code Entity Map - phenotype-go-kit

## Forward Map (Code -> Requirements)

| Entity | File | FR |
|--------|------|----|
| `logctx.WithLogger` | `logctx/logctx.go` | FR-LOG-001 |
| `logctx.From` | `logctx/logctx.go` | FR-LOG-001, FR-LOG-002 |
| `ringbuffer.New` | `ringbuffer/ringbuffer.go` | FR-BUF-001 |
| `ringbuffer.Push` | `ringbuffer/ringbuffer.go` | FR-BUF-002 |
| `ringbuffer.GetAll` | `ringbuffer/ringbuffer.go` | FR-BUF-003 |
| `waitfor.WaitFor` | `waitfor/waitfor.go` | FR-WAIT-001, FR-WAIT-002, FR-WAIT-003 |
| `waitfor.After` | `waitfor/waitfor.go` | FR-WAIT-003 |
| `registry.New` | `registry/registry.go` | FR-REG-001 |
| `registry.Register` | `registry/registry.go` | FR-REG-001 |
| `registry.Unregister` | `registry/registry.go` | FR-REG-002 |
| `registry.SetHook` | `registry/registry.go` | FR-REG-003 |

## Reverse Map (Requirements -> Code)

| FR | Primary Entities |
|----|-----------------|
| FR-LOG-001 | `logctx.WithLogger`, `logctx.From` |
| FR-LOG-002 | `logctx.From` (panic path) |
| FR-BUF-001 | `ringbuffer.New` |
| FR-BUF-002 | `ringbuffer.Push` |
| FR-BUF-003 | `ringbuffer.GetAll` |
| FR-WAIT-001 | `waitfor.WaitFor` |
| FR-WAIT-002 | `waitfor.WaitFor` (timeout path) |
| FR-WAIT-003 | `waitfor.WaitFor`, `waitfor.After` |
| FR-REG-001 | `registry.Register` |
| FR-REG-002 | `registry.Unregister` |
| FR-REG-003 | `registry.SetHook` |
| FR-REG-004 | All registry methods (RWMutex) |
