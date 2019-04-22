package utils

import (
    "encoding/json"
    "net/http"
    "log"
    "os"
)

var Logger *log.Logger

func init() {
    Logger = log.New(os.Stdout, "", 0)
}

func Has_error(err error, err_msg string) bool {
    if err != nil {
        Logger.Println("[Error]:", err_msg + ".", err)
        return true
    }
    return false
}

func Message(status bool, message string) (map[string]interface{}) {
    return map[string]interface{} {"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
    w.Header().Add("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}
