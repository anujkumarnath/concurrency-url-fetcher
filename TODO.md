- [ ] Fix blocking result printing before worker can pick another job
- [ ] Fix project structue
```
    -concurrent-url-fetcher/
    │
    ├── cmd/
    │   └── fetcher/
    │       └── main.go
    │
    ├── internal/
    │   ├── app/
    │   │   └── app.go
    │   │
    │   ├── fetcher/
    │   │   ├── worker.go
    │   │   ├── pool.go
    │   │   └── result.go
    │   │
    │   ├── signal/
    │   │   └── shutdown.go
    │   │
    │   └── stats/
    │       └── stats.go
    │
    ├── go.mod
    ├── go.sum
    └── README.md
```
