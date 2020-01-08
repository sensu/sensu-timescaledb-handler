# Sensu TimescaleDB Handler
TravisCI: [![TravisCI Build Status](https://travis-ci.org/sensu/sensu-timescaledb-handler.svg?branch=master)](https://travis-ci.org/sensu/sensu-timescaledb-handler)

The Sensu TimescaleDB Handler is a [Sensu Event Handler][3] that sends metrics to
the time series database [TimescaleDB][2]. [Sensu][1] can collect metrics using
check output metric extraction or the StatsD listener. Those collected metrics
pass through the event pipeline, allowing Sensu to deliver the metrics to the
configured metric event handlers. This TimescaleDB handler will allow you to
store, instrument, and visualize the metric data from Sensu.

## Installation

Download the latest version of the sensu-timescaledb-handler from [releases][4],
or create an executable script from this source.

### Compiling

From the local path of the sensu-timescaledb-handler repository:
```
go build -o /usr/local/bin/sensu-timescaldb-handler main.go
```

## Configuration

Example Sensu Go handler definition:

```json
{
    "api_version": "core/v2",
    "type": "Handler",
    "metadata": {
        "namespace": "default",
        "name": "timescaledb"
    },
    "spec": {
        "type": "pipe",
        "command": "sensu-timescaledb-handler -d postgresql://127.0.0.1:5432/sensu",
        "timeout": 10,
        "filters": [
            "has_metrics"
        ]
    }
}
```

## Usage Examples

Help:
```
Usage:
  sensu-timescaledb-handler [flags]

Flags:
  -d, --dsn string     the DSN of the TimescaleDB database, should be in DSN (Data Source Name) format (e.g. postgresql://localhost:5432/sensu) (default "postgres://localhost:5432/sensu")
  -h, --help           help for sensu-timescaledb-handler
  -t, --table string   the PostgreSQL table to store metrics in (default "metrics")
```

## Example Database

``` sql
CREATE database sensu;
CREATE TABLE metrics (
    time    TIMESTAMPTZ         NOT NULL,
    name    TEXT                NOT NULL,
    value   DOUBLE PRECISION    NULL,
    source  TEXT                NOT NULL,
    tags    JSONB
);
```

## Contributing

See https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md

[1]: https://github.com/sensu/sensu-go
[2]: https://github.com/timescale/timescaledb
[3]: https://docs.sensu.io/sensu-go/latest/reference/handlers/#how-do-sensu-handlers-work
[4]: https://github.com/sensu/sensu-timescaledb-handler/releases
