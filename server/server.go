package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/MicahParks/templater"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	hhpostgres "github.com/MicahParks/httphandle/postgres"

	"github.com/MicahParks/aws-ses-web-ui/config"
	"github.com/MicahParks/aws-ses-web-ui/postgres"
	"github.com/MicahParks/aws-ses-web-ui/ses"
)

type Server struct {
	Conf     config.Config
	Postgres postgres.Postgres
	SES      ses.SES
	logger   *slog.Logger
	pool     *pgxpool.Pool
	tmplr    templater.Templater
}

func NewServer(conf config.Config, logger *slog.Logger, tmplr templater.Templater) (Server, error) {
	var p postgres.Postgres
	var pool *pgxpool.Pool
	var err error
	if conf.ASWU.UsePostgres {
		pool, err = hhpostgres.Pool(context.Background(), conf.Postgres)
		if err != nil {
			return Server{}, fmt.Errorf("failed to create PostgreSQL pool: %w", err)
		}
		p = postgres.New(pool)
	}
	logger.Debug("Database connection established.")

	ss, err := ses.New(conf.SES)
	if err != nil {
		return Server{}, fmt.Errorf("failed to create SES session: %w", err)
	}

	s := Server{
		Conf:     conf,
		Postgres: p,
		SES:      ss,
		logger:   logger,
		pool:     pool,
		tmplr:    tmplr,
	}

	return s, nil
}

func (s Server) Logger() *slog.Logger {
	return s.logger
}
func (s Server) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return s.pool.BeginTx(ctx, pgx.TxOptions{})
}
func (s Server) Shutdown(_ context.Context) error {
	if s.Conf.ASWU.UsePostgres {
		s.pool.Close()
	}
	return nil
}
