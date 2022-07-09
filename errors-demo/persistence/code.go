package persistence

import (
	"fmt"
	"net/http"
)

var (
	ErrUserNotFound = fmt.Errorf("[%d]: user not found", http.StatusNotFound)
)
