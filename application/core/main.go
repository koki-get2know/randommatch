package main

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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

type groupMatchingReq struct {
	Size                 uint            `json:"size"`
	Groups               [][]entity.User `json:"groups"`
	ForbiddenConnections [][]entity.User `json:"forbiddenConnections"`
}

type EmailReq struct {
	Matches []matcher.Match `json:"matches"`
}

type UsersFile struct {
	Organization string                `form:"organization"`
	File         *multipart.FileHeader `form:"file" binding:"required"`
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

func linkfromMatching(tuples []matcher.Match) error {
	/*
	   Input : tuples of matchings
	   purpose: serialize the link in BD
	*/
	var connections [][]entity.User
	for _, match := range tuples {
		for _, u := range match.Users {
			for _, u1 := range match.Users {
				if u.Id != u1.Id {
					connections = append(connections, []entity.User{u, u1})
				}
			}
		}
	}

	return database.CreateLink(connections)
}

func generateMatchings(c *gin.Context) {
	defer duration(track("generateMatchings"))
	var req matchingReq

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json sent " + err.Error()})
		return
	}

	tuples := matcher.GenerateTuple(req.Users, [][]entity.User{}, matcher.Basic,
		req.ForbiddenConnections, req.Size, []entity.User{}, []entity.User{})
	c.JSON(http.StatusCreated, gin.H{"data": tuples})

}

func generateGroupMatchings(c *gin.Context) {
	defer duration(track("generateGroupMatchings"))
	var req groupMatchingReq

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json sent " + err.Error()})
		return
	}
	if len(req.Groups) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "you should send 2 groups"})
		return
	}

	tuples := matcher.GenerateTuple([]entity.User{}, [][]entity.User{}, matcher.Group,

		req.ForbiddenConnections, req.Size, req.Groups[0], req.Groups[1])

	c.JSON(http.StatusCreated, gin.H{"data": tuples})

}

func getTags(c *gin.Context) {
	defer duration(track("getTags"))
	tags, err := database.GetTags()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tags})
}

func getOrganization(c *gin.Context) {
	defer duration(track("getOrganization"))
	claims := c.MustGet("tokenClaims").(jwt.MapClaims)
	roles := claims["roles"].([]interface{})
	if !contains(roles, "Privilege.Approve") {
		c.JSON(http.StatusForbidden, gin.H{"message": "Operation denied permission missing"})
		return
	}
	id := c.Param("id")

	orga, err := database.GetOrganizationById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if len(orga.Id) == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": orga})
}

func createOrganization(c *gin.Context) {
	defer duration(track("createOrganization"))
	claims := c.MustGet("tokenClaims").(jwt.MapClaims)
	roles := claims["roles"].([]interface{})
	if !contains(roles, "Privilege.Approve") {
		c.JSON(http.StatusForbidden, gin.H{"message": "Operation denied permission missing"})
		return
	}

	var orga entity.Organization

	if err := c.BindJSON(&orga); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json sent " + err.Error()})
		return
	}
	orga.Name = strings.ToLower(orga.Name)
	orgaUid, err := database.CreateOrganization(orga)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "organization creation failed " + err.Error()})
		return
	}
	orga.Id = orgaUid
	c.Header("Location", fmt.Sprintf("/organizations/%v", orgaUid))
	c.JSON(http.StatusCreated, gin.H{"data": orga})
}

func uploadUsers(c *gin.Context) {
	defer duration(track("uploadUsers"))
	claims := c.MustGet("tokenClaims").(jwt.MapClaims)
	roles := claims["roles"].([]interface{})
	orgs := itemsWithPrefixInRole(roles, "Org.")

	var usersFile UsersFile

	if err := c.ShouldBind(&usersFile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid file sent " + err.Error()})
		return
	}
	usersFile.Organization = strings.ToLower(usersFile.Organization)
	if len(usersFile.Organization) > 0 && !containsString(orgs, usersFile.Organization) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Operation denied permission missing"})
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

	orgaUid, err := database.GetOrganizationByName(usersFile.Organization)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if len(orgaUid) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Organization " + usersFile.Organization + " not found"})
		return
	}

	jobId, err := database.CreateUsers(users, orgaUid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "users creation failed " + err.Error()})
		return
	}
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
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func deleteUser(c *gin.Context) {
	defer duration(track("deleteUser"))
	id := c.Param("id")
	claims := c.MustGet("tokenClaims").(jwt.MapClaims)
	roles := claims["roles"].([]interface{})
	if !contains(roles, "Privilege.Approve") {
		c.JSON(http.StatusForbidden, gin.H{"message": "Operation denied permission missing"})
		return
	}
	if err := database.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func deleteUsers(c *gin.Context) {
	defer duration(track("deleteUsers"))
	claims := c.MustGet("tokenClaims").(jwt.MapClaims)
	roles := claims["roles"].([]interface{})
	if !contains(roles, "Privilege.Approve") {
		c.JSON(http.StatusForbidden, gin.H{"message": "Operation denied permission missing"})
		return
	}
	if err := database.DeleteUsers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// End point to get all links in BD
func getLinks(c *gin.Context) {
	defer duration(track("getLinks"))
	links, err := database.GetLink()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": links})

}

func contains(s []any, e string) bool {
	for _, a := range s {
		if a.(string) == e {
			return true
		}
	}
	return false
}
func containsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func itemsWithPrefixInRole(s []any, prefix string) []string {
	orgs := []string{}
	for _, a := range s {
		if strings.HasPrefix(a.(string), prefix) {
			orgs = append(orgs, strings.ToLower(strings.TrimPrefix(a.(string), prefix)))
		}
	}
	return orgs
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

	// http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/

	jobId, err := calendar.SendInvite(req.Matches)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "mails sending failed " + err.Error()})
		return
	}
	c.Header("Location", fmt.Sprintf("/matching-email-job/%v", jobId))
	c.JSON(http.StatusOK, gin.H{"message": "emails are being sent"})

}

func getMatchingStats(c *gin.Context) {
	defer duration(track("getMatchings"))
	matchings, err := database.GetMatchingStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": matchings})
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
	_, exists := os.LookupEnv("NEO4J_AUTH")
	if exists {
		driver, err := database.Driver()
		if err != nil {
			fmt.Print(err)
		} else {
			defer (*driver).Close()
		}
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
	//protected.Use(middlewares.JwtAuth())
	protected.POST("/matchings", generateMatchings)
	protected.POST("/group-matchings", generateGroupMatchings)
	protected.POST("/upload-users", uploadUsers)
	protected.GET("/users-creation-job/:id", getJobStatus)
	protected.GET("/matching-email-job/:id", getJobStatus)
	protected.POST("/organizations", createOrganization)
	protected.GET("/organizations/:id", getOrganization)
	protected.GET("/users", getUsers)
	protected.DELETE("/users", deleteUsers)
	protected.DELETE("/users/:id", deleteUser)
	protected.GET("/tags", getTags)
	protected.POST("/email-matches", emailMatches)
	protected.GET("/matchings-stats", getMatchingStats)
	protected.GET("/links", getLinks)
	router.Run()
}
