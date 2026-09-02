package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func Connect(databaseURL string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		conn.Close(context.Background())
		return nil, err
	}

	return conn, nil
}
