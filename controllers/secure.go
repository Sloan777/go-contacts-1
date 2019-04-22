package controllers

import (
    "context"
    "net/http"
    jwt "github.com/dgrijalva/jwt-go"
    "go-contacts/models"
    "go-contacts/args"
    u "go-contacts/utils"
    "strings"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func SecureMiddleware(next HandlerFunc) HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        tokenHeader := r.Header.Get("Authorization") // Get token from the header

        if tokenHeader == "" {
            u.Respond(w, u.Message(false, "Missing Auto Token"))
            return
        }

        splitted := strings.Split(tokenHeader, " ") // check if the token match for the format `Bearer {token-body}`
        if len(splitted) != 2 {
            u.Respond(w, u.Message(false, "Invalid/Malformed auth token"))
            return
        }

        tokenPart := splitted[1] // the token part
        tk := &models.Token{}

        token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
                        return []byte(args.Token_password), nil
                    })
        if err != nil {
            u.Respond(w, u.Message(false, "Malformed authentication token"))
            return
        }

        if !token.Valid {
            u.Respond(w, u.Message(false, "Token is not valid"))
            return
        }

        next(w, r.WithContext( context.WithValue(r.Context(), args.SecureKey, tk.UserId) ))
    }
}

func WeakSecureMiddleware(next http.Handler) HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        tokenHeader := r.Header.Get("Authorization") // Get token from the header

        if tokenHeader == "" {
            next.ServeHTTP(w, r)
            return
        }

        splitted := strings.Split(tokenHeader, " ") // check if the token match for the format `Bearer {token-body}`
        if len(splitted) != 2 {
            next.ServeHTTP(w, r)
            return
        }

        tokenPart := splitted[1] // the token part
        tk := &models.Token{}

        token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
                        return []byte(args.Token_password), nil
                    })
        if err != nil {
            next.ServeHTTP(w, r)
            return
        }

        if !token.Valid {
            next.ServeHTTP(w, r)
            return
        }

        u.Logger.Println("User ID:", tk.UserId)
        next.ServeHTTP(w, r.WithContext( context.WithValue(r.Context(), args.SecureKey, tk.UserId) ))
    }
}

func UserFromContext(ctx context.Context) (int64) {
    uid, ok := ctx.Value(args.SecureKey).(int64)
    if !ok {
        return -1
    }
    return uid
}
