package main

import (
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
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
			Default:   "postgres://localhost:5432/sensu",
			Usage:     "the DSN of the TimescaleDB database, should be in DSN (Data Source Name) format (e.g. postgresql://localhost:5432/sensu)",
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

	goHandler := sensu.NewEnterpriseGoHandler(&handler.PluginConfig, configOptions, handler.Validate, handler.Run)
	goHandler.Execute()
}
