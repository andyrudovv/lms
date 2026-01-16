package db

import (
	"context"
	"errors"
	"net"
	"net/url"
	"time"


	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	ErrFailedConnection = errors.New("failed to create db conn pool")
	ErrPingDatabaseConn = errors.New("failed ping database connection")
)

// PostgresConnector database connector struct
type PostgresConnector struct {
	pool *pgxpool.Pool
}

type Config struct {
	ConnRetry                int
	ConnTimeout              time.Duration
	MaxOpenConnTTL           time.Duration
	MaxIdleConnTTL           time.Duration
	MaxConnLifetimeJitterTTL time.Duration
	User                     string
	Password                 string
	Host                     string
	Port                     string
	Name                     string
}

func (cfg Config) GetDSN() string {
	query := make(url.Values)
	if cfg.MaxOpenConnTTL > 0 {
		query.Set("pool_max_conn_lifetime", cfg.MaxOpenConnTTL.String())
	}
	if cfg.MaxIdleConnTTL > 0 {
		query.Set("pool_max_conn_idle_time", cfg.MaxIdleConnTTL.String())
	}
	if cfg.MaxIdleConnTTL > 0 {
		query.Set("pool_max_conn_idle_time", cfg.MaxIdleConnTTL.String())
	}
	if cfg.MaxConnLifetimeJitterTTL > 0 {
		query.Set("pool_max_conn_lifetime_jitter", cfg.MaxConnLifetimeJitterTTL.String())
	}
	dsn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     net.JoinHostPort(cfg.Host, cfg.Port),
		Path:     cfg.Name,
		RawQuery: query.Encode(),
	}
	return dsn.String()
}

// NewPostgresConnector returns new postgres connector
func NewPostgresConnector(ctx context.Context, cfg Config) (*PostgresConnector, error) {
	pool, err := pgxpool.Connect(ctx, cfg.GetDSN())
	if err != nil {
		return nil, err
	}

	if err = PingConnection(ctx, &cfg, func(pingCtx context.Context) error {
		return pool.Ping(pingCtx)
	}); err != nil {
		return nil, ErrFailedConnection
	}

	return &PostgresConnector{
		pool: pool,
	}, nil
}

func PingConnection(ctx context.Context, cfg *Config, pinger func(ctx context.Context) error) error {
	ticker := time.NewTicker(cfg.ConnTimeout)
	defer ticker.Stop()

	var err error
	for i := 0; i < cfg.ConnRetry; i++ {
		switch err = pinger(ctx); err {
		case nil:
			return nil
		case context.Canceled, context.DeadlineExceeded:
			return ErrPingDatabaseConn
		default:
			// logger.ErrorKV("failed to ping database connection", logger.Err, err)
		}
		select {
		case <-ctx.Done():
			// logger.ErrorKV("database.PingConnection.Done", logger.Err, ctx.Err())
			return ErrPingDatabaseConn
		case <-ticker.C:
		}
	}
	return ErrPingDatabaseConn
}

func (c *PostgresConnector) Close() {
	c.pool.Close()
}

func (c *PostgresConnector) Exec(ctx context.Context, query string, args ...interface{}) (int64, error) {
	rows, err := c.pool.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return rows.RowsAffected(), nil
}

func (c *PostgresConnector) QueryRowContext(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return c.pool.QueryRow(ctx, query, args...)
}

func (c *PostgresConnector) QueryContext(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return c.pool.Query(ctx, query, args...)
}

func (c *PostgresConnector) Ping(ctx context.Context) error {
	return c.pool.Ping(ctx)
}

func GetErrCodeAndConstraint(err error) (string, string) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code, pgErr.ConstraintName
	}
	return "", ""
}
