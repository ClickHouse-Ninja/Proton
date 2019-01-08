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
  0.1 rev[5c29788] master (2019-01-08.09:50:14 UTC).
USAGE:
  -addr string
        listen address (default "127.0.0.1:30002")
  -backlog int
        backlog size (default 100000)
  -concurrency int
        number of the background processes (default 8)
  -dsn string
        ClickHouse DSN (default "native://127.0.0.1:9000")
```