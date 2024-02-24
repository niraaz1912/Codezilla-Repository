package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
	var login LoginRequest
	var resp Empty

	if err := c.BindJSON(&login); err != nil {
		c.IndentedJSON(http.StatusBadRequest, resp)
		return
	}

	sessionid := uuid.New()

	cookie, err := c.Cookie("session")

	if err != nil {
		cookie = sessionid.String()
		c.SetCookie("session", cookie, 3500, "/", "localhost", false, true)
	}

	c.IndentedJSON(http.StatusOK, resp)
}

func login(c *gin.Context) {
	var login LoginRequest
	var resp Empty

	if err := c.BindJSON(&login); err != nil {
		c.IndentedJSON(http.StatusBadRequest, resp)
		return
	}

	sessionid := uuid.New()

	cookie, err := c.Cookie("session")

	if err != nil {
		cookie = sessionid.String()
		c.SetCookie("session", cookie, 3500, "/", "localhost", false, true)
	}

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
	fmt.Println("hello, world")

	router := gin.Default()
	router.POST("/login", login)
	router.POST("/location", postLocation)
	router.GET("/location", getLocation)

	router.Run("localhost:8081")
}
