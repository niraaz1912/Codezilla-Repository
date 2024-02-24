package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
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

type LoginResonse struct {
	Sessionid *uuid.UUID `json:"sessionid"`
}

type PostUserLocation struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

type GetLocationRequest struct {
	Filter string `json:"filter"`
}

type GetLocationResponse struct {
	Locations map[string][]LocationInstance `json:"locations"`
}

type LocationInstance struct {
	Time      uint64  `json:"time"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func createAccount(c *gin.Context) {
	var req LoginRequest
	var resp Empty

	if err := c.BindJSON(&req); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, resp)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 8)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, resp)
		return
	}

	row := db.QueryRow("select username from users where username=?", req.Username)

	var username string
	err = row.Scan(&username)
	if err != sql.ErrNoRows {
		log.Println("User already exists")
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
	var resp LoginResonse

	if err := c.BindJSON(&req); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, resp)
		return
	}

	var sessionIdString sql.NullString
	var passhash string

	log.Println(req.Username)
	row := db.QueryRow("SELECT sessionID, passhash FROM users WHERE username=?", req.Username)
	err := row.Scan(&sessionIdString, &passhash)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusUnauthorized, resp)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passhash), []byte(req.Password)); err != nil {
		log.Println(passhash)
		log.Println(err)
		c.IndentedJSON(http.StatusUnauthorized, resp)
		return
	}

	sessionid := uuid.New()
	if !sessionIdString.Valid {
		cookie, err := c.Cookie("sessionid")

		if err != nil {
			cookie = sessionid.String()
			c.SetCookie("sessionid", cookie, 100000, "/", "localhost", false, true)
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
		sessionid, err = uuid.Parse(sessionIdString.String)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusBadRequest, resp)
			return
		}
	}

	resp.Sessionid = &sessionid

	c.IndentedJSON(http.StatusOK, resp)
}

func logout(c *gin.Context) {
	var resp Empty

	cookie, err := c.Cookie("sessionid")
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, resp)
		return
	}

	var cookiestr sql.NullString
	row := db.QueryRow(`SELECT sessionid FROM users WHERE sessionid=$1`, cookie)
	err = row.Scan(&cookiestr)

	if err == sql.ErrNoRows {
		log.Printf("session id '%s' does not exist\n", cookie)
		c.IndentedJSON(http.StatusUnauthorized, resp)
		return
	}

	tx, _ := db.Begin()
	_, err = tx.Exec("UPDATE users SET sessionid=null WHERE sessionid=?", cookie)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, resp)
		return
	}
	_, err = tx.Exec("DELETE FROM location WHERE sessionid=?", cookie)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, resp)
		return
	}
	tx.Commit()

	c.SetCookie("sessionid", "", -1, "/", "localhost", false, true)

	c.IndentedJSON(http.StatusOK, resp)
}

func postLocation(c *gin.Context) {
	var req PostUserLocation
	var resp Empty

	if err := c.BindJSON(&req); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, resp)
		return
	}
	if req.Latitude == nil || req.Longitude == nil {
		log.Println("invalid latitude or longitude")
		c.IndentedJSON(http.StatusBadRequest, resp)
		return
	}

	cookie, err := c.Cookie("sessionid")
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusForbidden, resp)
		return
	}

	var username sql.NullString
	row := db.QueryRow(`SELECT username FROM users WHERE sessionid=$1`, cookie)
	err = row.Scan(&username)

	if err == sql.ErrNoRows {
		log.Printf("session id '%s' does not exist\n", cookie)
		c.IndentedJSON(http.StatusUnauthorized, resp)
		return
	}

	tx, _ := db.Begin()
	_, err = tx.Exec("INSERT INTO locationhistory (username, longitude, latitude, time) VALUES (?, ?, ?, ?)", username, req.Longitude, req.Latitude, time.Now().Unix())
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, resp)
		return
	}
	tx.Commit()

	c.IndentedJSON(http.StatusCreated, resp)
}

func getLocation(c *gin.Context) {
	var req GetLocationRequest
	var resp GetLocationResponse

	if err := c.BindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, resp)
		return
	}

	_, err := c.Cookie("sessionid")
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusForbidden, resp)
		return
	}

	rows, err := db.Query("select users.username, locationhistory.longitude, locationhistory.latitude, locationhistory.time from locationhistory left join users on users.username=locationhistory.username")
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, resp)
		return
	}

	locationHistoriesMap := map[string][]LocationInstance{}
	for rows.Next() {
		var key string
		location := LocationInstance{}
		if err := rows.Scan(&key, &location.Longitude, &location.Latitude, &location.Time); err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, resp)
			return
		}

		locationHistoriesMap[key] = append(locationHistoriesMap[key], location)
	}

	resp.Locations = locationHistoriesMap

	c.IndentedJSON(http.StatusOK, resp)
}

func main() {
	var err error

	db, err = sql.Open("sqlite3", fileName)

	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/login/new", createAccount)
	router.POST("/login", login)
	router.POST("/logout", logout)
	router.POST("/location", postLocation)
	router.GET("/location", getLocation)

	router.Run(":8081")
}
