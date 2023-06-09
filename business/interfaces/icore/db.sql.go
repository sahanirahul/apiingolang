package icore

import (
	"context"
	"database/sql"
)

type IDB interface {
	Conn(context.Context) (*sql.Conn, error)
}
