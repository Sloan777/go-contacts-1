package controllers

import (
    "net/http"
    u "go-contacts/utils"
    "go-contacts/models"
    "encoding/json"
)

// username, email, password
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    account := &models.Account{}
    err := json.NewDecoder(r.Body).Decode(account)
    if err != nil {
        u.Respond(w, u.Message(false, "Invalid Request"))
        return
    }
    u.Respond(w, account.Register())
}

// email, password
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    account := &models.Account{}
    err := json.NewDecoder(r.Body).Decode(account)
    if err != nil {
        u.Respond(w, u.Message(false, "Invalid Request"))
        return
    }
    u.Respond(w, account.Login())
}
