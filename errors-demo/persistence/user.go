package persistence

import (
	"context"
	"database/sql"
	"fmt"

	xerrors "github.com/pkg/errors"
)

type User struct {
	Id   int64  `db:"id""`
	Name string `db:"name"`
}

func (d *dbManager) GetUserByName(ctx context.Context, name string) (*User, error) {
	user := new(User)
	err := d.db.GetContext(ctx, user, "select * from user where name = ?", name)
	if err != nil {
		errMsg := fmt.Sprintf("exec: select * from user where name = %s", name)
		if xerrors.Is(sql.ErrNoRows, err) {
			return nil, xerrors.WithMessage(ErrUserNotFound, errMsg)
		}
		return nil, xerrors.Wrapf(err, errMsg)
	}
	return user, nil
}
