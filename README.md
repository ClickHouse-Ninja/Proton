# [WIP] Proton - it's a Pinba storage server.

# @todo
- [ ] Grafana dashboards
- [ ] reports (materialized views and queries)
- [ ] [timers](https://github.com/tony2001/pinba_engine/wiki/PHP-extension#pinba_timer_start)

# Usage:

```
NAME:
  Proton - high performance Pinba storage server.
VERSION:
  0.1 rev[9ef1982] master (2019-01-08.18:24:30 UTC).
USAGE:
  -addr string
        listen address (default ":30002")
  -backlog int
        backlog size (default 100000)
  -concurrency int
        number of the background processes (default 2)
  -dsn string
        ClickHouse DSN (default "native://127.0.0.1:9000")
  -metrics_addr string
        Address on which to expose metrics (default ":2112")
```