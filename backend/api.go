package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

const fileName = "db.db"

var db *sql.DB

type Empty struct{}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PostUserLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type GetLocationRequest struct {
	Filter string `json:"filter"`
}

type GetLocationResponse struct {
	Locations []LocationInstance `json:"locations"`
}

type LocationInstance struct {
	Username  string  `json:"username"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func createAccount(c *gin.Context) {
	var req LoginRequest
	var resp Empty

	if err := c.BindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, resp)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 8)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, resp)
		return
	}

	row := db.QueryRow("select username from users where username=?", req.Username)

	var username string
	err = row.Scan(&username)
	if err != sql.ErrNoRows {
		c.IndentedJSON(http.StatusUnauthorized, resp)
		return
	}

	tx, _ := db.Begin()
	_, err = tx.Exec("INSERT INTO users (username, role, passhash) VALUES (?, 'user', ?)", req.Username, string(hashed))
	tx.Commit()
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, resp)
		return
	}

	c.IndentedJSON(http.StatusCreated, resp)
}

func login(c *gin.Context) {
	var req LoginRequest
	var resp Empty

	if err := c.BindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, resp)
		return
	}

	var sessionIdString string
	var passhash string

	log.Println(req.Username)
	row := db.QueryRow("SELECT sessionID, passhash FROM users WHERE username=?", req.Username)
	err := row.Scan(&sessionIdString, &passhash)
	if err != nil {
		log.Println(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passhash), []byte(req.Password)); err != nil {
		log.Println(passhash)
		log.Println(err)
		c.IndentedJSON(http.StatusUnauthorized, resp)
		return
	}

	sessionid := uuid.New()
	if sessionIdString == "" {

		cookie, err := c.Cookie("session")

		if err != nil {
			cookie = sessionid.String()
			c.SetCookie("session", cookie, 3500, "/", "localhost", false, true)
		}

		tx, _ := db.Begin()
		_, err = tx.Exec("UPDATE users SET sessionID=? WHERE username=?", sessionid, req.Username)
		tx.Commit()
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, resp)
			return
		}
	} else {
		var err error
		sessionid, err = uuid.Parse(sessionIdString)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusBadRequest, resp)
			return
		}
	}

	c.IndentedJSON(http.StatusOK, resp)
}

func logout(c *gin.Context) {
	var resp Empty

	cookie, err := c.Cookie("session")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, resp)
		return
	}

	row := db.QueryRow("select sessionID from users where sessionID=?", cookie)

	var cookiestr string
	err = row.Scan(&cookiestr)
	if err != sql.ErrNoRows {
		c.IndentedJSON(http.StatusUnauthorized, resp)
		return
	}

	tx, _ := db.Begin()
	_, err = tx.Exec("UPDATE users SET sessionID=null WHERE sessionid=?", cookie)
	tx.Commit()
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, resp)
		return
	}

	c.SetCookie("session", "", -1, "/", "localhost", false, true)

	c.IndentedJSON(http.StatusOK, resp)
}

func postLocation(c *gin.Context) {
	var location PostUserLocation
	var resp Empty

	if err := c.BindJSON(&location); err != nil {
		c.IndentedJSON(http.StatusBadRequest, resp)
		return
	}

	_, err := c.Cookie("session")
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, resp)
		return
	}

	c.IndentedJSON(http.StatusCreated, resp)
}

func getLocation(c *gin.Context) {
	var req GetLocationRequest
	var resp GetLocationResponse

	if err := c.BindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, resp)
		return
	}

	_, err := c.Cookie("session")
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, resp)
		return
	}

	locations := []LocationInstance{
		{Username: "user1", Latitude: 123, Longitude: 456},
		{Username: "user2", Latitude: 456, Longitude: 789},
	}
	resp.Locations = locations

	c.IndentedJSON(http.StatusOK, resp)
}

func main() {
	var err error

	db, err = sql.Open("sqlite3", fileName)

	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.POST("/login/new", createAccount)
	router.POST("/login", login)
	router.POST("/location", postLocation)
	router.GET("/location", getLocation)

	router.Run("localhost:8081")
}
