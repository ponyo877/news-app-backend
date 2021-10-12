package handler

// For Checking GitOps .................
import (
    "fmt"
    "os"
    "strconv"
    "database/sql"
    _ "github.com/lib/pq"
    "io/ioutil"
    "encoding/json"
)

type DBConfig struct {
    Host    string  `json:"host"`
    Port    int     `json:"port"`
    User    string  `json:"user"`
    Dbname  string  `json:"dbname"`
    Pass    string  `json:"pass"`
}

func checkError(err error) {
	if err != nil {
	        fmt.Fprintf(os.Stderr, "fatal: error: %s", err.Error())
		os.Exit(1)
	}
}

func openDB() *sql.DB{
    jsonString, err := ioutil.ReadFile("config/config.json")
    checkError(err)
    
    var c DBConfig
    err = json.Unmarshal(jsonString, &c)
    checkError(err)

    ConnStr := "host=" + c.Host +
                " port=" + strconv.Itoa(c.Port) +
                " user=" + c.User +
                " dbname=" + c.Dbname +
                " password=" + c.Pass +
                " sslmode=disable"
    
    db, err := sql.Open("postgres", ConnStr)
    checkError(err)
    return db
}