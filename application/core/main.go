package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/koki/randommatch/matcher"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type matchingReq struct {
	Size                 uint             `json:"size"`
	Users                []matcher.User   `json:"users"`
	ForbiddenConnections [][]matcher.User `json:"forbiddenConnections"`
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
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json sent"})
	}
	tuples := matcher.GenerateTuple(req.Users, [][]matcher.User{}, matcher.Basic, req.ForbiddenConnections, req.Size, []matcher.User{}, []matcher.User{}, 0, 0)
	c.JSON(http.StatusCreated, gin.H{"data": tuples})
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

func helloFromNeo4j(c *gin.Context) {
	defer duration(track("helloFromNeo4j"))
	hello, err := helloNeo4j("bolt://match-db:7687", "neo4j", "test")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Exiting because of error" + err.Error()})
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": hello})
}

func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
	log.Printf("%v: %v\n", msg, time.Since(start))
}

func helloNeo4j(uri, username, password string) (string, error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return "", err
	}
	defer driver.Close()

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	defer session.Close()

	greeting, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"CREATE (a:Greeting) SET a.message = $message RETURN a.message + ', from node ' + id(a)",
			map[string]interface{}{"message": "hello, world from neo4j"})
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return "", err
	}

	return greeting.(string), nil
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.StaticFile("/api", "./api/swagger.yaml")
	router.GET("health-check", getHealthCheck)
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)
	router.GET("/neo4j", helloFromNeo4j)
	router.POST("/matchings", generateMatchings)

	router.Run()

}
