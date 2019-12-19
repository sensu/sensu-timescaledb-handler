package main

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugins-go-library/sensu"
	"github.com/stretchr/testify/assert"
)

func FixtureEventWithMetrics(entity, check string) *corev2.Event {
	event := corev2.FixtureEvent(entity, check)
	event.Metrics = corev2.FixtureMetrics()
	return event
}

func TestTimescaleDBHandler_Run(t *testing.T) {
	type fields struct {
		PluginConfig sensu.PluginConfig
		Config       TimescaleDBHandlerConfig
		DB           *sql.DB
	}
	type args struct {
		event *corev2.Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TimescaleDBHandler{
				PluginConfig: tt.fields.PluginConfig,
				Config:       tt.fields.Config,
				DB:           tt.fields.DB,
			}
			if err := tr.Run(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("TimescaleDBHandler.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTimescaleDBHandler_Validate(t *testing.T) {
	t.Parallel()

	type fields struct {
		PluginConfig sensu.PluginConfig
		Config       TimescaleDBHandlerConfig
		DB           *sql.DB
	}
	type args struct {
		event *corev2.Event
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  bool
		errMatch string
	}{
		{
			name:     "fail when no dsn specified",
			wantErr:  true,
			errMatch: "missing DSN",
		},
		{
			name: "fail when no table specified",
			fields: fields{
				Config: TimescaleDBHandlerConfig{
					DSN: "postgresql://foohost/bardb",
				},
			},
			wantErr:  true,
			errMatch: "missing Table",
		},
		{
			name: "fail when event has no metrics",
			fields: fields{
				Config: TimescaleDBHandlerConfig{
					DSN:   "postgresql://foohost/bardb",
					Table: "metrics",
				},
			},
			args: args{
				event: corev2.FixtureEvent("entity1", "check1"),
			},
			wantErr:  true,
			errMatch: "event does not contain metrics",
		},
		{
			name: "pass when required config and args are met",
			fields: fields{
				Config: TimescaleDBHandlerConfig{
					DSN:   "postgresql://foohost/bardb",
					Table: "metrics",
				},
			},
			args: args{
				event: FixtureEventWithMetrics("entity1", "check1"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TimescaleDBHandler{
				PluginConfig: tt.fields.PluginConfig,
				Config:       tt.fields.Config,
				DB:           tt.fields.DB,
			}
			err := tr.Validate(tt.args.event)
			assert.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				assert.Contains(t, tt.errMatch, err.Error())
			} else {
				assert.Contains(t, tt.errMatch, "")
			}
		})
	}
}

func TestTimescaleDBHandler_ProcessEvent(t *testing.T) {

	//mock.ExpectPrepare("INSERT INTO (time, name, value, source, tags) VALUES($1, $2, $3, $4, $5)")
	// mock.ExpectExec("INSERT INTO metrics(time, name, value, source, tags)").
	// 	WithArgs("timestamphere", "cpu.user", 84, "i-424242", "{\"foo\":\"bar\"}").
	// 	WillReturnResult(sqlmock.NewResult(1, 1))

	type fields struct {
		PluginConfig sensu.PluginConfig
		Config       TimescaleDBHandlerConfig
		DB           *sql.DB
	}
	type args struct {
		event *corev2.Event
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		dbMocks  func(sqlmock.Sqlmock)
		wantErr  bool
		errMatch string
	}{
		{
			name: "fails when timestamp is invalid",
			args: args{},
			dbMocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("INSERT INTO (time, name, value, source, tags) VALUES($1, $2, $3, $4, $5)")
			},
			wantErr:  true,
			errMatch: "abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			tr := &TimescaleDBHandler{
				PluginConfig: tt.fields.PluginConfig,
				Config:       tt.fields.Config,
				DB:           db,
			}
			if tt.dbMocks != nil {
				tt.dbMocks(mock)
			}
			err = tr.ProcessEvent(tt.args.event)
			assert.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				assert.Contains(t, tt.errMatch, err.Error())
			} else {
				assert.Contains(t, tt.errMatch, "")
			}
			if err = mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestTimescaleDBHandler_Setup(t *testing.T) {
	type fields struct {
		PluginConfig sensu.PluginConfig
		Config       TimescaleDBHandlerConfig
		DB           *sql.DB
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TimescaleDBHandler{
				PluginConfig: tt.fields.PluginConfig,
				Config:       tt.fields.Config,
				DB:           tt.fields.DB,
			}
			if err := tr.Setup(); (err != nil) != tt.wantErr {
				t.Errorf("TimescaleDBHandler.Setup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
