package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

type Database interface {
	Init()

	QueryRow(query string, args ...any) func(dest ...any) error
}

type DatabaseService struct {
	ctx             context.Context
	pool            *pgxpool.Pool
	DbConnectionUrl string
}

func (ds *DatabaseService) Init() {
	if ds.pool != nil {
		log.Println("database has already been created")
		return
	}

	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, ds.DbConnectionUrl)

	if err != nil {
		log.Println("Unable to create database connections' pool: ", err)
		os.Exit(1)
	}

	dbpool.Config().AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	ds.pool = dbpool
	ds.ctx = ctx
}

func (ds *DatabaseService) QueryRow(query string, args ...any) func(dest ...any) error {
	return func(dest ...any) error {
		conn, err := ds.pool.Acquire(ds.ctx)

		if err != nil {
			log.Println(err)
			return err
		}

		defer func() {
			conn.Conn().Close(ds.ctx)
		}()

		row := conn.QueryRow(context.Background(), query, args...)
		err = row.Scan(dest...)

		if err != nil {
			log.Println("err: ", err)
			return err
		}

		return nil
	}
}
