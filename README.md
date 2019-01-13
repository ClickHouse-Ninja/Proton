# [WIP] Proton - it's a Pinba storage server.

# @todo
- [ ] Grafana dashboards
- [X] reports (materialized views and queries) [basic reports](examples/reports/basic.md)
- [ ] [timers](https://github.com/tony2001/pinba_engine/wiki/PHP-extension#pinba_timer_start)

# Installation

### Install the ClickHouse server

```sh
sudo apt-key adv --keyserver keyserver.ubuntu.com --recv E0C56BD4    # optional

echo "deb http://repo.yandex.ru/clickhouse/deb/stable/ main/" | sudo tee /etc/apt/sources.list.d/clickhouse.list
sudo apt-get update

sudo apt-get install -y clickhouse-server clickhouse-client

sudo service clickhouse-server start
clickhouse-client
```

### Create Proton schema and the raw request table

```sh
clickhouse-client -n < schema/schema.sql
```

and then create the base report table and materialize view

```sh
clickhouse-client -n < schema/reports/base.sql
```

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