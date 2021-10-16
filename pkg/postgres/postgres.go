// Package postgres implements postgres connection.
package postgres

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
	_defaultDebug        = false
)

// Postgres -.
type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration
	debug        bool
	Query        *gorm.DB
}

// New -.
func New(url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
		debug:        _defaultDebug,
	}

	// set custom options
	for _, opt := range opts {
		opt(pg)
	}

	loggerLevel := logger.Silent
	if pg.debug {
		loggerLevel = logger.Info
	}

	connect, err := gorm.Open(postgres.Open(url), &gorm.Config{Logger: logger.Default.LogMode(loggerLevel)})
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	sqlDB, err := connect.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect DB: %w", err)
	}

	sqlDB.SetConnMaxLifetime(pg.connTimeout)
	sqlDB.SetMaxIdleConns(pg.connAttempts)
	sqlDB.SetMaxOpenConns(pg.maxPoolSize)
	pg.Query = connect

	return pg, nil
}
