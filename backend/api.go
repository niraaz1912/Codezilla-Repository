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
    _, err = db.Exec("INSERT OR IGNORE INTO users (username, passhash, role) VALUES (?, ?, ?)", req.Username, hashedPassword, "user")
    if err != nil {
        log.Println("Error inserting new user:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
        return
    }

    log.Println("Account created successfully for", req.Username)
    c.JSON(http.StatusOK, gin.H{"message": "Account created successfully"})
}




type LoginResponse struct {
    Sessionid *uuid.UUID `json:"sessionid"`
    Role      string      `json:"role"`
}

func login(c *gin.Context) {
    var req LoginRequest
    if err := c.BindJSON(&req); err != nil {
        log.Println("Error parsing request body:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    var storedHashedPassword, role string
    row := db.QueryRow("SELECT passhash, role FROM users WHERE username = ?", req.Username)
    if err := row.Scan(&storedHashedPassword, &role); err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
            return
        }
        log.Println("Error retrieving user role:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(req.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    sessionID := uuid.New()
    startTime := time.Now()
    _, err := db.Exec("UPDATE users SET sessionid = ?, start_time = ?, end_time = NULL WHERE username = ?", sessionID, startTime, req.Username)
    if err != nil {
        log.Println("Error updating session:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
        return
    }

    c.JSON(http.StatusOK, LoginResponse{Sessionid: &sessionID, Role: role})
}


type Session struct {
    Username   string `json:"username"`
    Role       string `json:"role"`
    StartTime  string `json:"start_time"`
    EndTime    string `json:"end_time,omitempty"` // Optional field for active sessions
}

func getSessions(c *gin.Context) {
    usernameFilter := c.Query("username") // Get `username` query parameter
    roleFilter := c.Query("role")        // Get `role` query parameter

    // Build the query dynamically based on filters
    query := "SELECT username, role, start_time, end_time FROM users WHERE 1=1"
    args := []interface{}{}

    if usernameFilter != "" {
        query += " AND username LIKE ?"
        args = append(args, "%"+usernameFilter+"%")
    }
    if roleFilter != "" {
        query += " AND role = ?"
        args = append(args, roleFilter)
    }

    rows, err := db.Query(query, args...)
    if err != nil {
        log.Println("Error fetching sessions:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sessions"})
        return
    }
    defer rows.Close()

    var sessions []Session
    for rows.Next() {
        var session Session
        var endTime sql.NullString
        if err := rows.Scan(&session.Username, &session.Role, &session.StartTime, &endTime); err != nil {
            log.Println("Error scanning session row:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sessions"})
            return
        }

        session.EndTime = "Active"
        if endTime.Valid {
            session.EndTime = endTime.String
        }

        sessions = append(sessions, session)
    }

    c.JSON(http.StatusOK, sessions)
}






func logout(c *gin.Context) {
    var req LogoutRequest
    var resp Empty

    log.Println("Logout called!")
    if err := c.BindJSON(&req); err != nil {
        log.Println("Error parsing request body:", err)
        c.IndentedJSON(http.StatusBadRequest, resp)
        return
    }

    log.Printf("Received session ID: %v", req.Sessionid) // Log the received session ID

    if req.Sessionid == nil {
        log.Println("Session ID is missing or invalid")
        c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
        return
    }

    var cookiestr sql.NullString
    row := db.QueryRow(`SELECT sessionid FROM users WHERE sessionid = ?`, req.Sessionid)
    err := row.Scan(&cookiestr)

    if err == sql.ErrNoRows {
        log.Printf("Session ID '%s' does not exist\n", req.Sessionid)
        c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Session not found"})
        return
    }

    log.Println("Session found, proceeding to logout")
    // Update database and invalidate session ID
    tx, _ := db.Begin()
    endTime := time.Now()

    _, err = tx.Exec("UPDATE users SET sessionid = NULL, end_time = ? WHERE sessionid = ?", endTime, req.Sessionid)
    if err != nil {
        tx.Rollback()
        log.Println("Error updating session end time:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
        return
    }
    tx.Commit()

    c.SetCookie("sessionid", "", -1, "/", "localhost", false, true)

    log.Println("Logout successful for session ID:", req.Sessionid)
    c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}


func postLocation(c *gin.Context) {
    var req PostUserLocation

    if err := c.BindJSON(&req); err != nil {
        log.Println("Error parsing location data:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    if req.Latitude == nil || req.Longitude == nil || req.UserID == nil {
        log.Println("Invalid latitude, longitude, or session ID")
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location data"})
        return
    }

    // Validate session ID
    var username sql.NullString
    row := db.QueryRow("SELECT username FROM users WHERE sessionid = ?", req.UserID)
    err := row.Scan(&username)

    if err == sql.ErrNoRows {
        log.Printf("Session ID '%s' not found", req.UserID)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
        return
    }

    // Insert location data into database
    tx, err := db.Begin()
    if err != nil {
        log.Println("Error starting transaction:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }

    _, err = tx.Exec("INSERT OR REPLACE INTO locationhistory (username, longitude, latitude, time) VALUES (?, ?, ?, ?)", username.String, req.Longitude, req.Latitude, time.Now().Unix())
    if err != nil {
        tx.Rollback()
        log.Println("Error inserting location data:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log location"})
        return
    }

    tx.Commit()
    log.Println("Location data logged for user:", username.String)
    c.JSON(http.StatusCreated, gin.H{"message": "Location logged successfully"})
}


func getLocation(c *gin.Context) {
    var resp GetLocationResponse

    // Validate session ID and get user role
    sessionID := c.Request.Header.Get("SessionID")
    if sessionID == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Session ID missing"})
        return
    }

    var username string
    var role string
    row := db.QueryRow("SELECT username, role FROM users WHERE sessionid = ?", sessionID)
    if err := row.Scan(&username, &role); err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
            return
        }
        log.Println("Error retrieving user role:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
        return
    }

    // Check role
    if role != "admin" {
        log.Printf("User %s does not have admin privileges", username)
        c.JSON(http.StatusOK, gin.H{"locations": map[string][]LocationInstance{}})
        return
    }

    // Admin: Get all user locations
    rows, err := db.Query("SELECT users.username, locationhistory.longitude, locationhistory.latitude, locationhistory.time FROM locationhistory LEFT JOIN users ON users.username = locationhistory.username")
    if err != nil {
        log.Println(err)
        c.JSON(http.StatusInternalServerError, resp)
        return
    }

    locationHistoriesMap := map[string][]LocationInstance{}
    for rows.Next() {
        var key string
        location := LocationInstance{}
        if err := rows.Scan(&key, &location.Longitude, &location.Latitude, &location.Time); err != nil {
            log.Println(err)
            c.JSON(http.StatusInternalServerError, resp)
            return
        }

        locationHistoriesMap[key] = append(locationHistoriesMap[key], location)
    }

    resp.Locations = locationHistoriesMap
    c.JSON(http.StatusOK, resp)
}

func getUserRole(c *gin.Context) {
    sessionID := c.Request.Header.Get("SessionID")
    if sessionID == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Session ID missing"})
        return
    }

    var role string
    row := db.QueryRow("SELECT role FROM users WHERE sessionid = ?", sessionID)
    if err := row.Scan(&role); err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
            return
        }
        log.Println("Error retrieving user role:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"role": role})
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
    config.AllowHeaders = []string{"Content-Type", "SessionID"}
	router.Use(cors.New(config))

	router.POST("/login/new", createAccount)
	router.POST("/login", login)
	router.POST("/logout", logout)
	router.POST("/location", postLocation)
	router.GET("/location", getLocation)
    router.GET("/user/role", getUserRole)
    router.GET("/sessions", getSessions)


	router.Run(":8081")
}
