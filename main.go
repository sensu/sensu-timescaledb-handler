package main

import (
	"github.com/sensu/sensu-plugins-go-library/sensu"
)

func main() {
	handler := TimescaleDBHandler{
		PluginConfig: sensu.PluginConfig{
			Name:  "sensu-timescaledb-handler",
			Short: "a timescaledb handler built for use with sensu",
		},
	}

	configOptions := []*sensu.PluginConfigOption{
		{
			Path:      "dsn",
			Env:       "TIMESCALEDB_DSN",
			Argument:  "dsn",
			Shorthand: "d",
			Default:   "postgres://localhost:5432/metrics",
			Usage:     "the DSN of the TimescaleDB database, should be in DSN (Data Source Name) format (e.g. postgresql://localhost:5432/metrics)",
			Value:     &handler.Config.DSN,
		},
		{
			Path:      "table",
			Env:       "TIMESCALEDB_TABLE",
			Argument:  "table",
			Shorthand: "t",
			Default:   "metrics",
			Usage:     "the PostgreSQL table to store metrics in",
			Value:     &handler.Config.Table,
		},
	}

	goHandler := sensu.NewGoHandler(&handler.PluginConfig, configOptions, handler.Validate, handler.Run)
	goHandler.Execute()
}
