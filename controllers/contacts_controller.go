package controllers

import (
    "net/http"
    "encoding/json"
    "go-contacts/models"
    u "go-contacts/utils"
)

var AddContactHandler = SecureMiddleware( func(w http.ResponseWriter, r *http.Request) {
    contact := &models.Contact{}
    err := json.NewDecoder(r.Body).Decode(contact)
    if err != nil {
        u.Respond(w, u.Message(false, err.Error()))
        return
    }

    uid := UserFromContext( r.Context() )
    u.Respond(w, models.InsertContact(uid, contact))
})

var DeleteContactHandler = SecureMiddleware( func(w http.ResponseWriter, r *http.Request) {
    var contactId int64
    err := json.NewDecoder(r.Body).Decode(&contactId)
    if err != nil {
        u.Respond(w, u.Message(false, err.Error()))
        return
    }

    uid := UserFromContext( r.Context() )
    u.Respond(w, models.DeleteContact(uid, contactId))
})

var UpdateContactHandler = SecureMiddleware( func(w http.ResponseWriter, r *http.Request) {
    contact := &models.Contact{}
    err := json.NewDecoder(r.Body).Decode(contact)
    if err != nil {
        u.Respond(w, u.Message(false, err.Error()))
        return
    }

    uid := UserFromContext( r.Context() )
    u.Respond(w, models.UpdateContact(uid, contact))
})

var AllContactsHandler = SecureMiddleware( func(w http.ResponseWriter, r *http.Request) {
    uid := UserFromContext( r.Context() )
    u.Respond(w, models.GetAllContacts(uid));
})
