package initializers

import (
	"context"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type Database struct {
	Client *sqlx.DB
}

func ConnectToDB() (*Database, error) {
	log.Info("Setting up a new database connection")

	connectionString := os.Getenv("DB_URL")

	db, err := sqlx.Connect("mysql", connectionString)
	if err != nil {
		return &Database{}, fmt.Errorf("failed to connect to database: %w", err)
	}
	log.Info("Database connection made sucessfully")
	return &Database{
		Client: db,
	}, nil
}

func (d *Database) Ping(ctx context.Context) error {
	return d.Client.DB.PingContext(ctx)
}
