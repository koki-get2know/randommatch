package main

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/koki/randommatch/calendar"
	"github.com/koki/randommatch/convert"
	"github.com/koki/randommatch/database"
	"github.com/koki/randommatch/entity"
	"github.com/koki/randommatch/matcher"
	"github.com/koki/randommatch/middlewares"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type matchingReq struct {
	Size                 uint            `json:"size"`
	Users                []entity.User   `json:"users"`
	ForbiddenConnections [][]entity.User `json:"forbiddenConnections"`
}

type EmailReq struct {
	Matches []matcher.Match `json:"matches"`
}

type UsersFile struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// getAlbums responds with the list of all albums as JSON.
func getHealthCheck(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func generateMatchings(c *gin.Context) {
	defer duration(track("generateMatchings"))
	var req matchingReq

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json sent " + err.Error()})
		return
	}
	tuples := matcher.GenerateTuple(req.Users, [][]entity.User{}, matcher.Basic,
		req.ForbiddenConnections, req.Size, []entity.User{}, []entity.User{}, 0, 0)
	c.JSON(http.StatusCreated, gin.H{"data": tuples})
}

func uploadUsers(c *gin.Context) {
	defer duration(track("uploadUsers"))
	var usersFile UsersFile

	if err := c.ShouldBind(&usersFile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid file sent " + err.Error()})
		return
	}

	// https://stackoverflow.com/questions/45121457/how-to-get-file-posted-from-json-in-go-gin
	// https://github.com/gin-gonic/gin#model-binding-and-validation
	if usersFile.File.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "File exceeded max size"})
		return
	}

	users, err := convert.CsvToUsers(usersFile.File)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid file content " + err.Error()})
		return
	}

	jobId := uuid.New().String()
	err = database.CreateJobStatus(jobId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "job creation failed " + err.Error()})
		return
	}

	go database.CreateUsers(users, jobId)

	c.Header("Location", fmt.Sprintf("/users-creation-job/%v", jobId))
	c.JSON(http.StatusAccepted, gin.H{"message": "Job enqueued"})
}

func getUsers(c *gin.Context) {
	defer duration(track("getUsers"))
	users, err := database.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": users})
}

func getJobStatus(c *gin.Context) {
	defer duration(track("getJobStatus"))
	id := c.Param("id")
	status, err := database.GetJobStatus(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": status})
}

func emailMatches(c *gin.Context) {
	defer duration(track("emailMatches"))
	var req EmailReq

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json sent " + err.Error()})
		return
	}

	claims := c.MustGet("tokenClaims").(jwt.MapClaims)
	adminEmail := claims["preferred_username"].(string)

	go func() {
		for _, match := range req.Matches {
			match := match
			calendar.SendInvite(&match, adminEmail)
			break // send only once for now
		}
	}()
	c.JSON(http.StatusOK, gin.H{"message": "emails are being sent"})

}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	defer duration(track("getAlbums"))

	c.JSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	defer duration(track("postAlbums"))

	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	defer duration(track("getAlbumByID"))
	id := c.Param("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
	log.Printf("%v: %v\n", msg, time.Since(start))
}

func main() {
	//os.Setenv("NEO4J_AUTH", "***/***")
	_, exists := os.LookupEnv("NEO4J_AUTH")
	if exists {
		driver, err := database.Driver()
		if err != nil {
			fmt.Print(err)
			return
		}
		defer (*driver).Close()
	}

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.Use(middlewares.Cors())

	public := router.Group("")
	public.StaticFile("/api", "./api/swagger.yaml")
	public.GET("health-check", getHealthCheck)
	public.GET("/albums", getAlbums)
	public.GET("/albums/:id", getAlbumByID)
	public.POST("/albums", postAlbums)

	protected := router.Group("")
	protected.Use(middlewares.JwtAuth())
	protected.POST("/matchings", generateMatchings)
	protected.POST("/upload-users", uploadUsers)
	protected.GET("/users-creation-job/:id", getJobStatus)
	protected.GET("/users", getUsers)
	protected.POST("/email-matches", emailMatches)

	router.Run()

}
