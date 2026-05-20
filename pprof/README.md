# pprof

Go runtime profiling endpoints for debugging performance issues.

## Security Model

pprof runs on a **separate server** (default `127.0.0.1:6060`) that is **localhost-only** by default.

```
Public Internet → Ingress → :8080 (API with auth)
                            :8081 (health - internal)

Localhost only (not routable):
  └── 127.0.0.1:6060 (pprof)
      └── /debug/pprof/*
```

**Why a separate port?**
- Port `6060` in GitOps configs is an obvious red flag during PR reviews
- Easy to spot if someone accidentally enables in non-dev environments
- Localhost-only by default - not network accessible

**Access via kubectl port-forward:**
```bash
kubectl port-forward pod/myapp-xxx 6060:6060
```

## Benefits

- **CPU Profiling**: Identify hot spots and slow functions
- **Memory Profiling**: Find memory leaks and allocation patterns
- **Goroutine Analysis**: Detect goroutine leaks and deadlocks
- **Block Profiling**: Find contention points
- **Mutex Profiling**: Identify lock contention

## Usage

### Enable pprof

```bash
# Environment variables
export PPROF_ENABLED=true
export PPROF_ADDRESS=127.0.0.1:6060  # optional, this is the default

# Or CLI flags
./myapp --pprof-enabled --pprof-address=127.0.0.1:6060
```

### Code

```go
import (
    "github.com/kilnfi/go-utils/app"
    "github.com/kilnfi/go-utils/pprof"
)

cfg := &app.Config{
    PProf: &pprof.Config{
        Enabled: true,
        Address: "127.0.0.1:6060", // localhost-only (default)
    },
}
```

### Endpoints

| Endpoint | Description |
|----------|-------------|
| `/debug/pprof/` | Index page with links to all profiles |
| `/debug/pprof/profile?seconds=30` | CPU profile (30 second sample) |
| `/debug/pprof/heap` | Heap memory profile |
| `/debug/pprof/goroutine` | All goroutine stack traces |
| `/debug/pprof/block` | Stack traces of blocked goroutines |
| `/debug/pprof/mutex` | Stack traces of mutex contention |
| `/debug/pprof/trace?seconds=5` | Execution trace (5 second sample) |
| `/debug/pprof/allocs` | Memory allocation profile |

### Accessing in Kubernetes

```bash
# Port-forward to the pprof server
kubectl port-forward pod/myapp-xxx 6060:6060

# Then access locally
go tool pprof http://localhost:6060/debug/pprof/heap
```

### Collecting Profiles

```bash
# CPU profile (30 seconds)
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Heap profile
go tool pprof http://localhost:6060/debug/pprof/heap

# Goroutine dump (human-readable)
curl http://localhost:6060/debug/pprof/goroutine?debug=2

# Execution trace
curl -o trace.out http://localhost:6060/debug/pprof/trace?seconds=5
go tool trace trace.out
```

### Interactive Analysis

```bash
# Download and open interactive web UI
go tool pprof -http=:9090 http://localhost:6060/debug/pprof/heap

# Compare two heap profiles to find leaks
go tool pprof -http=:9090 -diff_base=before.pb.gz after.pb.gz
```

## Best Practices

1. **Never enable in production** - Port 6060 in GitOps = red flag in PR review
2. **Keep localhost-only** - Never bind to `0.0.0.0`
3. **Use debug=2 for goroutine dumps** - Provides full stack traces
4. **Profile under load** - Idle profiles are not representative
5. **Compare before/after** - Use `-diff_base` to find regressions
