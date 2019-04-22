package gql

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
    u "go-contacts/utils"
    "go-contacts/models"
)

func DBGetUser(user *User) error {
    conn, err := models.GetDB()
    defer conn.Close()
    if err != nil {
        return err
    }
    stmt, err := conn.Prepare("select username, email from accounts where uid = ?")
    if err != nil {
        return err
    }
    rows, err := stmt.Query(user.Uid)
    defer rows.Close()
    if err != nil {
        return err
    }

    if !rows.Next() {
        return fmt.Errorf("User ID does not exist! But is this possible?")
    }

    rows.Scan(&user.Username, &user.Email)
    return nil
}

func DBGetContacts(userId int64, limit int, offset int) ([]Contact, error) {
    contacts := []Contact{}
    conn, err := models.GetDB()
    defer conn.Close()
    if err != nil {
        return contacts, err
    }
    stmt, err := conn.Prepare("select cid, name, phone, address, email from contacts where uid = ? limit ? offset ?")
    if err != nil {
        return contacts, err
    }
    rows, err := stmt.Query(userId, limit, offset)
    defer rows.Close()
    if err != nil {
        return contacts, err
    }

    contact := Contact{}
    for rows.Next() {
        rows.Scan(&contact.Cid, &contact.Name, &contact.Phone, &contact.Address, &contact.Email)
        contacts = append(contacts, contact)
    }
    return contacts, nil
}

func DBCreateUser(input *NewUser) (int64, error) {
    if input.Username == "" {
        return -1, fmt.Errorf("Username is missing")
    }
    if input.Password == "" {
        return -1, fmt.Errorf("Password is missing")
    }
    if input.Email == "" {
        return -1, fmt.Errorf("Email is missing")
    }

    conn, err := models.GetDB()
    defer conn.Close()
    if err != nil {
        return -1, err
    }

    stmt, err := conn.Prepare("insert into accounts(username, email, password) values(?, ?, ?)")
    if err != nil {
        return -1, err
    }

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    input.Password = string(hashedPassword)
    res, err := stmt.Exec(input.Username, input.Email, input.Password)
    if err != nil {
        return -1, err
    }

    uid, err := res.LastInsertId()
    if err != nil {
        return -1, err
    }

    return uid, nil
}

func DBAddContact(input *NewContact, userId int64) (int64, error) {
    conn, err := models.GetDB()
    defer conn.Close()
    if err != nil {
        return -1, err
    }

    stmt, err := conn.Prepare("insert into contacts(name, phone, address, email, uid) values(?, ?, ?, ?, ?)")
    if err != nil {
        return -1, err
    }
    res, err := stmt.Exec(input.Name, input.Phone, input.Address, input.Email, userId)
    if err != nil {
        return -1, err
    }
    cid, err := res.LastInsertId()
    if err != nil {
        return -1, err
    }
    return cid, nil
}

func DBDeleteContact(userId, cid int64) error {
    conn, err := models.GetDB()
    defer conn.Close()
    if err != nil {
        return err
    }

    stmt, err := conn.Prepare("delete from contacts where uid = ? and cid = ?")
    if err != nil {
        return err
    }
    res, err := stmt.Exec(userId, cid)
    if err != nil {
        u.Logger.Println("[gpl.db_query] [ERROR]: delete contact Exec error")
        return err
    }
    row_cnt, err := res.RowsAffected()
    if err != nil {
        u.Logger.Println("[gpl.db_query] [ERROR]: delete contact RowsAffected error")
        return err
    }
    if row_cnt == 1 {
        return nil
    } else {
        return fmt.Errorf("No such contact ID")
    }
}

func DBUpdateContact(input *ContactUpdate, userId int64) (int64, error) {
    conn, err := models.GetDB()
    defer conn.Close()
    if err != nil {
        return input.Cid, err
    }

    stmt, err := conn.Prepare("update contacts set name = ?, phone = ?, address = ?, email = ? where uid = ? and cid = ?")
    if err != nil {
        return input.Cid, err
    }
    res, err := stmt.Exec(input.Name, input.Phone, input.Address, input.Email, userId, input.Cid)
    if err != nil {
        u.Logger.Println("[gpl.db_query] [ERROR]: Update contact Exec error")
        return input.Cid, err
    }
    row_cnt, err := res.RowsAffected()
    if err != nil {
        u.Logger.Println("[gpl.db_query] [ERROR]: Update contact RowsAffected error")
        return input.Cid, err
    }
    if row_cnt == 1 {
        return input.Cid, nil
    } else {
        return input.Cid, fmt.Errorf("Cannot update this contact")
    }
}
