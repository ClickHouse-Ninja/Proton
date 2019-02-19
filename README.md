# [WIP] Proton - it's a Pinba storage server.

# @todo
- [X] Grafana dashboards [basic reports](examples/grafana/report.json), [Proton Server](examples/grafana/proton-server.json)
- [X] reports (materialized views and queries) [basic reports](examples/reports/basic.md), [tags reports](examples/reports/tags.md)
- [X] [timers](https://github.com/tony2001/pinba_engine/wiki/PHP-extension#pinba_timer_start)

# Installation

### Install ClickHouse server

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

### Add Proton dictionary to ClickHouse server

Example:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<dictionaries>
    <dictionary>
      <name>Proton</name>
        <source>
            <clickhouse>
                <host>127.0.0.1</host>
                <port>9000</port>
                <user>default</user>
                <password></password>
                <db>proton</db>
                <table>dictionary</table>
            </clickhouse>
        </source>
      <lifetime>600</lifetime>
      <layout><hashed /></layout>
      <structure>
         <id><name>ID</name></id>
         <attribute>
               <name>Value</name>
                <type>String</type>
                <null_value></null_value>
         </attribute>
      </structure>
   </dictionary>
</dictionaries>
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

![Grafana basic reports](/assets/grafana-basic-reports.png)
