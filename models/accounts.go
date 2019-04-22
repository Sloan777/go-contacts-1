package models

import (
    "github.com/dgrijalva/jwt-go"
    //"strings"
    "golang.org/x/crypto/bcrypt"
    "github.com/google/uuid"
    u "go-contacts/utils"
    "go-contacts/args"
)

type Token struct {
    UserId int64
    jwt.StandardClaims
}

type Account struct {
    ID  int64
    Username string
    Email string `json:"email"`
    Password string `json:"password"`
    Token string `json:"token";sql:"-"`
}


func (account *Account) createToken() {
    tk := &Token{UserId: account.ID}
    tk.Id = uuid.New().String()
    token := jwt.NewWithClaims( jwt.GetSigningMethod("HS256"), tk )
    tokenString, _ := token.SignedString([]byte(args.Token_password))
    account.Token = tokenString
}

// username, email, password
func (account *Account) Register() (map[string] interface{}) {
    if account.Username == "" {
        return u.Message(false, "Username is missing")
    }
    if account.Email == "" {
        return u.Message(false, "Email is missing")
    }
    if account.Password == "" {
        return u.Message(false, "Password is missing")
    }

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
    account.Password = string(hashedPassword)

    conn, err := GetDB()
    defer conn.Close()
    if err != nil {
        return u.Message(false, err.Error())
    }

    stmt, err := conn.Prepare("insert into accounts(username, email, password) values(?, ?, ?)")
    if err != nil {
        return u.Message(false, err.Error())
    }

    res, err := stmt.Exec(account.Username, account.Email, account.Password)
    if err != nil {
        return u.Message(false, err.Error())
    }

    uid, err := res.LastInsertId()
    if err != nil {
        return u.Message(false, err.Error())
    }

    account.ID = uid
    account.createToken()

    account.Password = "" // delete password
    response := u.Message(true, "Account has been created")
    response["account"] = account
    return response
}

// email, password
func (account *Account) Login() (map[string] interface{}) {
    if account.Email == "" {
        return u.Message(false, "Email address is missing")
    }
    if account.Password == "" {
        return u.Message(false, "Password is missing")
    }

    conn, err := GetDB()
    defer conn.Close()
    if err != nil {
        return u.Message(false, err.Error())
    }

    stmt, err := conn.Prepare("select uid, username, password from accounts where email = ? limit 1")
    if err != nil {
        return u.Message(false, err.Error())
    }
    rows, err := stmt.Query(account.Email)
    defer rows.Close()
    if err != nil {
        return u.Message(false, err.Error())
    }

    if !rows.Next() {
        return u.Message(false, "Email does not exist")
    }

    var db_password string
    rows.Scan(&account.ID, &account.Username, &db_password)
    err = bcrypt.CompareHashAndPassword([]byte(db_password), []byte(account.Password))
    if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
        return u.Message(false, "Incorrect email or password")
    }
    account.Password = ""
    account.createToken()

    resp := u.Message(true, "Logged In")
    resp["account"] = account
    return resp
}

/*
func GetUser(u uint) *Account {
    acc := &Account{}
    stmt, err := GetDB().Prepare("select * from accounts where id = ? limit 1")
    rows, err := stmt.Query(u)
    if !rows.Next() {
        rows.Close()
        return nil
    }
    err = rows.Scan(&acc.ID, &acc.Email, &acc.Password)
    rows.Close()
    acc.Password = ""
    return acc
}
*/
