package models

import (
    _ "github.com/mattn/go-sqlite3"
    u "go-contacts/utils"
    "go-contacts/args"
    "database/sql"
)

func GetDB() (*sql.DB, error) {
    conn, err := sql.Open("sqlite3", args.DB_name)
    if u.Has_error(err, "Connection to Database failed") {
        return nil, err
    }
    return conn, nil
}
