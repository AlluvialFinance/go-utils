# 🔧 Framework One-Pager: Lifecycle-Oriented Application Framework

## 🧩 Overview

This framework provides a robust structure for building modular Go applications with built-in support for **lifecycle management** and **HTTP services**. It simplifies spinning up production-ready apps by managing configuration, logging, health checks, metrics, and the orderly start/stop of services.

## 💡 Core Concepts

- **App**: Central application orchestrator that wires together services and manages their lifecycle.
- **Service**: Any component your app depends on. These can expose HTTP routes, background routines, health checks, metrics, etc.

## 🔁 Lifecycle Interfaces

Services can optionally implement the following interfaces from `app/service.go`:

| Interface       | Purpose                                                                 |
|-----------------|-------------------------------------------------------------------------|
| `Loggable`      | Receives a logger instance                                              |
| `Initializable` | Performs setup work before the app starts (DB connections, validations) |
| `Runnable`      | Starts long-running processes or goroutines                             |
| `Closable`      | Cleans up resources (e.g., closes DB connections)                       |
| `API`           | Registers HTTP routes                                                   |
| `Middleware`    | Adds HTTP middlewares                                                   |
| `Checkable`     | Exposes health checks                                                   |
| `Measurable`    | Registers Prometheus metrics                                             |

The app detects and invokes these interfaces at appropriate stages in the lifecycle.

## 🚀 App Lifecycle

```
Construct → RegisterService → Init → Start → Stop → Close
```

The framework ensures:

1. Health and metrics servers start before your services.
2. Proper signal handling (e.g., SIGINT/SIGTERM).
3. Graceful shutdown and cleanup in reverse order.

---

## 🧪 Sample Usage

Here’s a minimal example using Cobra and Viper to wire up a small service:

### `main.go`

```go
func main() {
    _ = godotenv.Load(".env")
    rootCmd := NewMyAppCommand()
    if err := rootCmd.Execute(); err != nil {
        logrus.WithError(err).Fatal("execution failed")
    }
}
```

### `cmd.go`

```go
func NewMyAppCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use: "myapp",
    }

    ctx := kilncmdutils.WithViper(context.Background(), viper.New())
    cmd.AddCommand(NewRunCmd(ctx))

    return cmd
}
```

### `run.go`

```go
func NewRunCmd(ctx context.Context) *cobra.Command {
    return &cobra.Command{
        Use:   "run",
        Short: "Run the service",
        RunE: func(cmd *cobra.Command, _ []string) error {
            app, err := NewKilnApp(ctx)
            if err != nil {
                return err
            }

            svc := NewMyService() // <- notice we don't explicitly call SetLogger, Init, Start, Stop, Close
            app.RegisterService(svc)

            return app.Run() // <- this will iterate over all registered services and call their lifecycle methods
        },
    }
}
```

### `myservice.go`

```go
// Implement any combination of lifecycle interfaces as needed for your service
// The Go compiler will indicate which methods are required based on the interfaces you choose
// Select interfaces according to the responsibilities of your service
// Possible interfaces include: Loggable, Initializable, Runnable, Closable, API, Middleware, Checkable, Measurable
type Service interface {
    app.Runnable
    app.Initializable
    app.Loggable
}
var _ Service = (*MyService)(nil)

type MyService struct{
    log *logrus.Logger
}

func NewMyService() *MyService {
    return &MyService{}
}

func (s *MyService) SetLogger(logger *logrus.Logger) {
    s.log = logger
}

func (s *MyService) Init(ctx context.Context) error {
    s.log.Info("Initializing service...")
    // Perform any initialization logic here
    return nil
}

func (s *MyService) Start(ctx context.Context) error {
    go func() {
        <-ctx.Done()
        fmt.Println("Stopped.")
    }()
    s.log.Info("Service started")
    return nil
}

func (s *MyService) Stop(ctx context.Context) error {
    s.log.Info("Stopping service...")
    return nil
}

func (s *MyService) Close() error {
   s.log.Info("Closing service...")
    return nil
}
```


## 🔄 Automation Highlights

Once a service is registered via `app.RegisterService(...)`, the framework **automatically** handles:

- Calling `SetLogger`, `Init`, `Start`, `Stop`, and `Close` on services implementing the corresponding interfaces.
- Context propagation across all services.
- OS signal handling (e.g., SIGINT, SIGTERM) to trigger graceful shutdown.

> ✅ Your service implementations remain clean and focused — no need to explicitly invoke lifecycle methods.


---

## ✅ Benefits

- **Fast bootstrapping**: Get up and running with minimal boilerplate.
- **Reduced bugs**: Standardized lifecycle reduces missed initialization or cleanup.
- **Extendable**: Plug in services without tight coupling.

> ⚠️ While not perfect, this framework offers a clean, maintainable path to building services with clarity and consistency.