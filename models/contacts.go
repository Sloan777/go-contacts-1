package models

import (
    u "go-contacts/utils"
)

type Contact struct {
    ID int64
    Name string
    Phone string
    Address string
    Email string
}

func GetAllContacts(userId int64) (map[string] interface{}) {
    conn, err := GetDB()
    defer conn.Close()
    if err != nil {
        return u.Message(false, err.Error())
    }

    stmt, err := conn.Prepare("select cid, name, phone, address, email from contacts where uid = ?")
    if err != nil {
        return u.Message(false, err.Error())
    }

    rows, err := stmt.Query(userId)
    defer rows.Close()
    if err != nil {
        return u.Message(false, err.Error())
    }

    contacts := []*Contact{}
    for rows.Next() {
        t := &Contact{}
        rows.Scan(&t.ID, &t.Name, &t.Phone, &t.Address, &t.Email)
        contacts = append(contacts, t)
    }

    resp := u.Message(true, "Fetched all contacts")
    resp["all-contacts"] = contacts

    return resp
}

func DeleteContact(userId int64, contactId int64) (map[string] interface{}) {
    conn, err := GetDB()
    defer conn.Close()
    if err != nil {
        return u.Message(false, err.Error())
    }

    stmt, err := conn.Prepare("delete from contacts where cid = ? and uid = ?")
    if err != nil {
        return u.Message(false, err.Error())
    }
    stmt.Exec(contactId, userId)

    return u.Message(true, "Deleted")
}

func InsertContact(userId int64, contact *Contact) (map[string] interface{}) {
    conn, err := GetDB()
    defer conn.Close()
    if err != nil {
        return u.Message(false, err.Error())
    }

    stmt, err := conn.Prepare("insert into contacts(name, phone, address, email, uid) values(?, ?, ?, ?, ?)")
    if err != nil {
        return u.Message(false, err.Error())
    }
    res, err := stmt.Exec(contact.Name, contact.Phone, contact.Address, contact.Email, userId)
    if err != nil {
        return u.Message(false, err.Error())
    }
    cid, err := res.LastInsertId()
    if err != nil {
        return u.Message(false, err.Error())
    }
    contact.ID = cid

    resp := u.Message(true, "Inserted")
    resp["contact"] = contact
    return resp
}

func UpdateContact(userId int64, contact *Contact) (map[string] interface{}) {
    conn, err := GetDB()
    defer conn.Close()
    if err != nil {
        return u.Message(false, err.Error())
    }

    stmt, err := conn.Prepare("update contacts set name = ?, phone = ?, address = ?, email = ? where uid = ? and cid = ?")
    if err != nil {
        return u.Message(false, err.Error())
    }
    _, err = stmt.Exec(contact.Name, contact.Phone, contact.Address, contact.Email, userId, contact.ID)
    if err != nil {
        return u.Message(false, err.Error())
    }

    return u.Message(true, "Updated")
}
