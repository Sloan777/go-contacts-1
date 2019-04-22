package main

import (
    "net/http"
    "time"
    "github.com/gorilla/mux"
    gql_handler "github.com/99designs/gqlgen/handler"
    //"go-contacts/app"
    u "go-contacts/utils"
    "go-contacts/args"
    "go-contacts/controllers"
    "go-contacts/gql"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
    u.Logger.Println("[rootHandler]")
    msg := u.Message(true, "Hello, World")
    u.Respond(w, msg)
}

func redirectHandler(url string, status_code_optional ...int) func(http.ResponseWriter, *http.Request) {

    /*
    status  HTTP        Temp/Perma  Cacheable       Method
    301     HTTP/1.0    Permanent   Yes             GET / POST may change
    302     HTTP/1.0    Temporary   not by default  GET / POST may change
    303     HTTP/1.1    Temporary   never           always GET
    307     HTTP/1.1    Temporary   not by default  may not change
    308     HTTP/1.1    Permanent   by default      may not change
    */
    status_code := 302
    if len(status_code_optional) > 0 {
        status_code = status_code_optional[0]
    }
    return func(w http.ResponseWriter, r *http.Request) {
        u.Logger.Println("[redirect] Redirect to", url, "status_code =", status_code)
        http.Redirect(w, r, url, status_code)
    }
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Do stuff here
        u.Logger.Printf("[I] [%s] %s %s %s %s\n", time.Now().Format(time.UnixDate), r.Method, r.Proto, r.URL, r.RemoteAddr)
        next.ServeHTTP(w, r)
    })
}

func main() {
    router := mux.NewRouter()
    router.Use(loggingMiddleware)

    if args.Test_GraphQL == 1 {
        router.Handle("/", gql_handler.Playground("GraphQL Playground", args.Site_root + "/graphql"))
    } else {
        router.HandleFunc("/", redirectHandler(args.Site_root))
    }
    router.HandleFunc(args.Site_root, rootHandler)
    router.HandleFunc(args.Site_root + "/api/login", controllers.LoginHandler).Methods("POST")
    router.HandleFunc(args.Site_root + "/api/register", controllers.RegisterHandler).Methods("POST")
    router.HandleFunc(args.Site_root + "/api/contact/add", controllers.AddContactHandler).Methods("POST")
    router.HandleFunc(args.Site_root + "/api/contact/delete", controllers.DeleteContactHandler).Methods("DELETE")
    router.HandleFunc(args.Site_root + "/api/contact/update", controllers.UpdateContactHandler).Methods("PUT")
    router.HandleFunc(args.Site_root + "/api/all-contacts", controllers.AllContactsHandler).Methods("GET")

    my_graphql_handler := gql_handler.GraphQL(
                gql.NewExecutableSchema(gql.NewRootResolvers()),
                gql_handler.ComplexityLimit(args.Graphql_complexity),
            )
    router.HandleFunc(args.Site_root + "/graphql", controllers.WeakSecureMiddleware(my_graphql_handler))

    u.Logger.Println("Listening on: 127.0.0.1:" + args.Port)
    u.Has_error(http.ListenAndServe("127.0.0.1:" + args.Port, router), "Failed to Start HTTP Server")
}
