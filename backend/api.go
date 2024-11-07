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

type LogoutRequest struct {
	Sessionid *uuid.UUID `json:"sessionid"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResonse struct {
	Sessionid *uuid.UUID `json:"sessionid"`
}

type PostUserLocation struct {
	UserID    *uuid.UUID `json:"sessionid"`
	Latitude  *float64   `json:"latitude"`
	Longitude *float64   `json:"longitude"`
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
    log.Println("Received signup request")

    var req LoginRequest
    if err := c.BindJSON(&req); err != nil {
        log.Println("Error binding JSON:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    log.Printf("Checking if username %s exists", req.Username)

    // Check if the user already exists
    var existingUsername string
    row := db.QueryRow("SELECT username FROM users WHERE username = ?", req.Username)
    err := row.Scan(&existingUsername)

    if err == nil {
        log.Println("User already exists:", req.Username)
        c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
        return
    } else if err != sql.ErrNoRows {
        log.Println("Error checking for existing user:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
        return
    }

    log.Println("Username available, creating new account")

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        log.Println("Error hashing password:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
        return
    }

    // Insert user with a default role of "user"
    _, err = db.Exec("INSERT INTO users (username, passhash, role) VALUES (?, ?, ?)", req.Username, hashedPassword, "user")
    if err != nil {
        log.Println("Error inserting new user:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
        return
    }

    log.Println("Account created successfully for", req.Username)
    c.JSON(http.StatusOK, gin.H{"message": "Account created successfully"})
}




func login(c *gin.Context) {
    log.Println("Received login request")

    var req LoginRequest
    if err := c.BindJSON(&req); err != nil {
        log.Println("Error binding JSON:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    log.Printf("Checking if username %s exists", req.Username)

    // Retrieve the user's hashed password from the database
    var storedHashedPassword string
    row := db.QueryRow("SELECT passhash FROM users WHERE username = ?", req.Username)
    err := row.Scan(&storedHashedPassword)
    
    if err == sql.ErrNoRows {
        log.Println("User not found:", req.Username)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    } else if err != nil {
        log.Println("Error retrieving user data:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
        return
    }

    log.Println("User found, comparing password")

    // Compare the hashed password with the provided password
    err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(req.Password))
    if err != nil {
        log.Println("Password mismatch for user:", req.Username)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    // Login successful
    log.Println("Login successful for user:", req.Username)
    c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}


func logout(c *gin.Context) {
	var req LogoutRequest
	var resp Empty

	if err := c.BindJSON(&req); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, resp)
		return
	}

	var cookiestr sql.NullString
	row := db.QueryRow(`SELECT sessionid FROM users WHERE sessionid=$1`, req.Sessionid)
	err := row.Scan(&cookiestr)

	if err == sql.ErrNoRows {
		log.Printf("session id '%s' does not exist\n", req.Sessionid)
		c.IndentedJSON(http.StatusUnauthorized, resp)
		return
	}

	tx, _ := db.Begin()
	_, err = tx.Exec("UPDATE users SET sessionid=null WHERE sessionid=?", req.Sessionid)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, resp)
		return
	}
	_, err = tx.Exec("DELETE FROM location WHERE sessionid=?", req.Sessionid)
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

	var username sql.NullString
	row := db.QueryRow(`SELECT username FROM users WHERE sessionid=$1`, req.UserID)
	err := row.Scan(&username)

	if err == sql.ErrNoRows {
		log.Printf("session id '%s' does not exist\n", req.UserID)
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
	var resp GetLocationResponse

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

	// Set Gin to release mode for production
    gin.SetMode(gin.ReleaseMode)

	db, err = sql.Open("sqlite3", fileName)

	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	// Configure trusted proxies
    router.SetTrustedProxies([]string{"127.0.0.1"}) 

	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowOrigins = []string{"http://localhost:5500", "http://127.0.0.1:5500" , "http://heron.cs.umanitoba.ca"}
	config.AllowMethods = []string{"POST", "GET", "OPTIONS"}
    config.AllowHeaders = []string{"Content-Type"}
	router.Use(cors.New(config))

	router.POST("/login/new", createAccount)
	router.POST("/login", login)
	router.POST("/logout", logout)
	router.POST("/location", postLocation)
	router.GET("/location", getLocation)

	router.Run(":8081")
}
