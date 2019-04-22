# Go-Contacts

A Contact Manager App built in Go Programming Language. This app implemented the [REST API](REST.md) and the [GraphQL](gql/schema.graphql). One can register and login, and then create, update, delete contacts, and also get a list of contacts. 

# Dependencies

This project is implemented and tested using `go version go1.12 linux/amd64`

To install the dependencies list below, type `go get <package name>`

* `golang.org/x/crypto/bcrypt`
* `github.com/gorilla/mux`
* `github.com/99designs/gqlgen`
* `github.com/mattn/go-sqlite3`
* `github.com/dgrijalva/jwt-go`
* `github.com/google/uuid`
* `github.com/joho/godotenv`

# Install

Use the following command to clone this project

```
git clone git@github.com:zijuexiansheng/go-contacts.git
```

It is recommended to clone this project into your `$GOPATH/src` directory. Some other untested approaches including adding the directory to your `$GOPATH`, or use `go mod`.

# Configuration

The server can be configured by changing the `/path/to/repo/.env` file. The following parameters can be changed.

* `db_name`: the SQLite database name (default: `gocontacts.db`)
    * You can copy the `default.db` in this repository to this name and run tests on it
* `port`: (default: `8000`)
* `site_root`: the root path to the service (default: `/go-dmo`)
    * This is how we defined `$SITEROOT` that will be used in the REST API and GraphQL below.
    * Suppose that the servers listens on `http://localhost:8000`
    * `site_root` is set to be `/go-demo`
    * `$SITEROOT` will then be `http://localhost:8000/go-demo`
    * `$SITEROOT/api` will be `http://localhost:8000/go-demo/api`
* `token_password`: A secret token password. Please set one by yourself. This password should be kept secret
* `graph_complexity`: The complexity of GraphQL queries (default: `200`). This parameter is used to prevent heavy-loaded request that could overload the server, particularly the database system
* `test_ql`: If set to `1`, one can open the `$SITEROOT` and test the GraphQL, as in the examples in the [screenshots](screenshots/). If set to `0` (default), then the `$SITEROOT` will only return a greeting
    * N.B., for requests requiring authorizations, please use `curl` or `Postman` to login and acquire the token first, and then add the token to the Header. See the "Authorization" in the "HTTP HEADERS" in [gql-7](screenshots/gql-7.png) for an example (ignore the "Authorization-1" and "Authorization-3")

# Run

* Go to the repository directory (make sure the directory is listed in the `$GOPATH`)
* Copy the `default.db` to any database file name (say `cp default.db gocontacts.db`)
* Change the configuration file `/path/to/repo/.env` as you wish
    * and remember to change the `db_name`
* `go run main.go`
* Use `curl` or `Postman` to test the REST API (also GraphQL)
* If `test_ql = 1`, then you can go to `$SITEROOT` and use the Playground to try with the GraphQL.
    * See [screenshots](screenshots/) for examples

# REST API

* All the requests should go to `$SITEROOT/api/...`
* To read the REST API, click [here](REST.md)

# GraphQL

* All the requests should go to `$SITEROOT/graphql`
* To read the GraphQL schema, click [here](gql/schema.graphql)
* Visit [here](screenshots/) for some screenshots

# License

MIT License
