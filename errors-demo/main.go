package main

import (
	"context"
	"flag"
	"fmt"

	xerrors "github.com/pkg/errors"

	"github/chrisyxy/geek-go/errors-demo/persistence"
)

var (
	dbUser = flag.String("db_user", "root", "mysql user")
	dbPwd  = flag.String("db_pwd", "root", "mysql password")
	dbAddr = flag.String("db_addr", "localhost:3306", "mysql address")
	dbName = flag.String("db_name", "test", "mysql database")
)

func main() {
	flag.Parse()

	dbSourceURL := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		*dbUser,
		*dbPwd,
		*dbAddr,
		*dbName)

	dbManager, err := persistence.NewDbManager(dbSourceURL)
	if err != nil {
		panic(err)
	}
	defer dbManager.Close()

	var (
		ctx      = context.Background()
		userName string
	)

	for {
		fmt.Printf("please input search username: ")
		fmt.Scanln(&userName)
		if len(userName) == 0 {
			fmt.Println("input username is empty.")
			continue
		}

		user, err := dbManager.GetUserByName(ctx, userName)
		if err != nil {
			fmt.Printf("%+v\n", err)
			if xerrors.Is(persistence.ErrUserNotFound, xerrors.Cause(err)) {
				continue
			}
			return
		}
		fmt.Printf("search user: %#v\n", user)
	}
}
