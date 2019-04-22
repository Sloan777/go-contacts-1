package gql

import (
    "fmt"
    "github.com/99designs/gqlgen/graphql"
    "io"
    "strconv"
)

type User struct {
    Uid         int64
    Username    string
    Email       string
}

type Contact struct {
    Cid         int64
    Name        string
    Phone       string
    Address     string
    Email       string
}

func MarshalID(id int64) graphql.Marshaler {
    return graphql.WriterFunc(func(w io.Writer) {
        io.WriteString(w, strconv.FormatInt(id, 10))
    })
}

func UnmarshalID(v interface{}) (int64, error) {
    id, ok := v.(string)
    if !ok {
        return 0, fmt.Errorf("ids must be strings")
    }
    i, e := strconv.ParseInt(id, 10, 64)
    return int64(i), e
}
