package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/koki/randommatch/calendar"
	"github.com/koki/randommatch/database"
	"github.com/koki/randommatch/entity"
	"github.com/koki/randommatch/handler"
	"github.com/koki/randommatch/matcher"
	"github.com/koki/randommatch/middlewares"
	"github.com/koki/randommatch/utils/helper"
)

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

type tagMatchingReq struct {
	Size                 uint            `json:"size"`
	Tags                 []string        `json:"tags"`
	ForbiddenConnections [][]entity.User `json:"forbiddenConnections"`
	Exclude              []entity.User   `json:"excludeUsers"`
	Organization         string          `json:"organization"`
}

type EmailReq struct {
	Matches []matcher.Match `json:"matches"`
}

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
	defer helper.Duration(helper.Track("generateMatchings"))
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
	defer helper.Duration(helper.Track("generateGroupMatchings"))
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

func generateMatchingByTag(c *gin.Context) {
	defer helper.Duration(helper.Track("generateGroupMatchings"))
	var req tagMatchingReq

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json sent " + err.Error()})
		return
	}
	if len(req.Tags) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "you should send 2 tags"})
		return
	}
	usersTag1, err := database.GetUsersByTag(req.Organization, req.Tags[0])

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	usersTag2, err := database.GetUsersByTag(req.Organization, req.Tags[1])

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	for _, u := range req.Exclude {
		usersTag1 = u.RmUser(usersTag1)
	}
	for _, u := range req.Exclude {
		usersTag2 = u.RmUser(usersTag2)
	}

	tuples := matcher.GenerateTuple([]entity.User{}, [][]entity.User{}, matcher.Group,

		req.ForbiddenConnections, req.Size, usersTag1, usersTag2)

	c.JSON(http.StatusCreated, gin.H{"data": tuples})
}

func getTags(c *gin.Context) {
	defer helper.Duration(helper.Track("getTags"))
	tags, err := database.GetTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tags})
}

// End point to get all links in BD
func getLinks(c *gin.Context) {
	defer helper.Duration(helper.Track("getLinks"))
	links, err := database.GetLink()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": links})

}

func emailMatches(c *gin.Context) {
	defer helper.Duration(helper.Track("emailMatches"))
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

func main() {
	_, exists := os.LookupEnv("NEO4J_AUTH")
	if exists {
		driver, err := database.Driver()
		if err != nil {
			log.Print(err)
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
	public.GET("/albums", handler.GetAlbums)
	public.GET("/albums/:id", handler.GetAlbumByID)
	public.POST("/albums", handler.PostAlbums)

	protected := router.Group("")
	protected.Use(middlewares.JwtAuth())
	protected.POST("/matchings", generateMatchings)
	protected.POST("/group-matchings", generateGroupMatchings)
	protected.POST("/tag-matchings", generateMatchingByTag)
	protected.POST("/upload-users", handler.UploadUsers)
	protected.POST("/organizations", handler.CreateOrganization)
	protected.POST("/email-matches", emailMatches)

	protected.GET("/users-creation-job/:id", handler.GetJobStatus)
	protected.GET("/matching-email-job/:id", handler.GetJobStatus)
	protected.GET("/organizations/:id", handler.GetOrganization)
	protected.GET("/users", handler.GetUsers)

	protected.DELETE("/users", handler.DeleteUsers)
	protected.DELETE("/users/:id", handler.DeleteUser)
	protected.GET("/tags", getTags)
	protected.GET("/matchings-stats", handler.GetMatchingStats)
	protected.GET("/links", getLinks)

	router.Run()
}
