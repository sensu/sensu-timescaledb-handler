package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugins-go-library/sensu"
)

type TimescaleDBHandler struct {
	sensu.PluginConfig

	Config TimescaleDBHandlerConfig
	DB     *sql.DB
}

type TimescaleDBHandlerConfig struct {
	// DSN is a data source name. It is either a URL or a postgres connection
	// string. See https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
	// for more information.
	DSN   string
	Table string
}

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

func (t *TimescaleDBHandler) Setup() error {
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
