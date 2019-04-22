package args

import (
    "os"
    "github.com/joho/godotenv"
    "strconv"
    u "go-contacts/utils"
)

var (
    DB_name string
    Port string
    Token_password string
    Site_root string
    Graphql_complexity int
    Test_GraphQL int
)

var SecureKey struct {
    keyvalue *string
}

func os_getenv(name string, def_value string) string {
    ret := os.Getenv(name)
    if ret == "" {
        return def_value
    }
    return ret
}

func init() {
    u.Has_error(godotenv.Load(), "Failed to load .env file")
    DB_name = os_getenv( "db_name", "gocontacts.db" )
    Port = os_getenv( "port", "8000" )
    Site_root = os_getenv( "site_root", "/go-demo" )
    Token_password = os_getenv( "token_password", "qpwoeiruty1029384756" )

    t, err := strconv.Atoi( os_getenv( "graph_complexity", "200" ) )
    if err != nil {
        t = 200
    }
    Graphql_complexity = t

    t, err = strconv.Atoi( os_getenv("test_ql", "0"))
    if err != nil {
        t = 0
    }
    Test_GraphQL = t

    SecureKey.keyvalue = &Token_password
}
