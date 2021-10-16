package postgres

import "time"

// Option - extended options.
type Option func(*Postgres)

// MaxPoolSize - set max pool connection to db. Default: 1.
func MaxPoolSize(size int) Option {
	return func(c *Postgres) {
		c.maxPoolSize = size
	}
}

// ConnAttempts - set max idle connection to db. Default: 10.
func ConnAttempts(attempts int) Option {
	return func(c *Postgres) {
		c.connAttempts = attempts
	}
}

// ConnTimeout - set max lifetime connection to db. Default: 1s.
func ConnTimeout(timeout time.Duration) Option {
	return func(c *Postgres) {
		c.connTimeout = timeout
	}
}

// Debug - set debug log in db activity. Default: false.
func Debug(str string) Option {
	return func(c *Postgres) {
		c.debug = str == "true"
	}
}
