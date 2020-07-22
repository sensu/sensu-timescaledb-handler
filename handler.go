package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
	"net/url"
	
	_ "github.com/lib/pq"
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
)

// TimescaleDBHandler is a timescaledb handler
type TimescaleDBHandler struct {
	sensu.PluginConfig

	Config TimescaleDBHandlerConfig
	DB     *sql.DB
}

// TimescaleDBHandlerConfig is a timescaledb handler config
type TimescaleDBHandlerConfig struct {
	// DSN is a data source name. It is either a URL or a postgres connection
	// string. See https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
	// for more information.
	DSN     string
	Table   string
	SslMode string
}

// Run runs the timescaledb handler
func (t *TimescaleDBHandler) Run(event *corev2.Event) error {
	if err := t.Setup(); err != nil {
		return err
	}

	defer t.DB.Close()

	if err := t.ProcessEvent(event); err != nil {
		return err
	}

	return nil
}

// ProcessEvent processes the timescaledb handler event
func (t *TimescaleDBHandler) ProcessEvent(event *corev2.Event) error {
	query := fmt.Sprintf("INSERT INTO %s(time, name, value, source, tags) VALUES($1, $2, $3, $4, $5)", t.Config.Table)
	stmt, err := t.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, point := range event.Metrics.Points {
		timestamp, err := convertInt64ToTime(point.Timestamp)
		if err != nil {
			return err
		}

		jsonTags, err := json.Marshal(point.Tags)
		if err != nil {
			return err
		}

		_, err = stmt.Exec(timestamp, point.Name, point.Value, event.Entity.Name, jsonTags)
		if err != nil {
			return err
		}
	}

	return nil
}

// Setup sets up the timescaledb handler
func (t *TimescaleDBHandler) Setup() error {
	// Set connection parameter defaults
	// Reference: https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
	params := url.Values{}
	params.Set("sslmode", t.Config.SslMode)

	// Parse the Data Source Name (DSN) & add missing default parameters
	dsn,err := url.Parse(t.Config.DSN)
	if err != nil {
		return err
	} 
	q := dsn.Query()
	for k := range q {
		params.Set(k,q[k][0])
	}
	dsn.RawQuery = params.Encode()
	t.Config.DSN = dsn.String()
	
	// Connect to TimescaleDB (Postgres database)
	db, err := sql.Open("postgres", t.Config.DSN)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	t.DB = db

	return nil
}

// Validate validates the timescaledb config
func (t *TimescaleDBHandler) Validate(event *corev2.Event) error {
	if len(t.Config.DSN) == 0 {
		return errors.New("missing DSN")
	}
	if len(t.Config.Table) == 0 {
		return errors.New("missing Table")
	}
	if !event.HasMetrics() {
		return errors.New("event does not contain metrics")
	}
	var sslmodes = []string{"disable","require","verify-ca","verify-full"}
	if indexOf(t.Config.SslMode,sslmodes) < 0 {
		return errors.New(fmt.Sprintf("unsupported sslmode \"%s\"", t.Config.SslMode))
	}
	return nil
}

func convertInt64ToTime(t int64) (time.Time, error) {
	stringTimestamp := strconv.FormatInt(t, 10)
	if len(stringTimestamp) > 10 {
		stringTimestamp = stringTimestamp[:10]
	}
	t, err := strconv.ParseInt(stringTimestamp, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(t, 0), nil
}

func indexOf(k string, s []string) (int) {
	for i,v := range s {
		if k == v {
			return i
		}
	}
	return -1
}