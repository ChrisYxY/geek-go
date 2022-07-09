package persistence

import (
	"context"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	xerrors "github.com/pkg/errors"
)

const (
	MySQLDriver string = "mysql"
)

type DBManager interface {
	Close() error

	GetUserByName(ctx context.Context, name string) (*User, error)
}

type dbManager struct {
	db *sqlx.DB
}

func NewDbManager(sourceUrl string) (DBManager, error) {
	db, err := sqlx.Connect(MySQLDriver, sourceUrl)
	if err != nil {
		return nil, xerrors.Wrap(err, "open db failed")
	}

	return &dbManager{db: db}, nil
}

func (d *dbManager) Close() error {
	return d.db.Close()
}
