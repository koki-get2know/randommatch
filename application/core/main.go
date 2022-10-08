package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
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
	Matches      []matcher.Match `json:"matches"`
	Organization string          `json:"organization"`
	Duration     int64           `json:"duration"`
	Date         time.Time       `json:"date"`
}

func getHealthCheck(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

type scheduleMatchingReq struct {
	Uid          string `json:"code"`
	Organization string `json:"organization"`
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

	tuples := matcher.GenerateTuple(req.Users, [][]entity.User{}, entity.Basic,
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

	tuples := matcher.GenerateTuple([]entity.User{}, [][]entity.User{}, entity.Group,

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

	tuples := matcher.GenerateTuple([]entity.User{}, [][]entity.User{}, entity.Group,

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
	claims := c.MustGet("tokenClaims").(jwt.MapClaims)
	roles := claims["roles"].([]interface{})
	orgs := helper.ItemsWithPrefixInRole(roles, "Org.")

	orgaName := strings.ToLower(req.Organization)
	if len(orgaName) > 0 && !helper.ContainsString(orgs, orgaName) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Operation denied permission missing"})
		return
	}
	orgaUid, err := database.GetOrganizationByName(orgaName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/

	if req.Duration == 0 {
		req.Duration = 15
	}

	jobId, err := calendar.SendInvite(req.Matches, orgaUid, req.Duration, req.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "mails sending failed " + err.Error()})
		return
	}
	c.Header("Location", fmt.Sprintf("/matching-email-job/%v", jobId))
	c.JSON(http.StatusOK, gin.H{"message": "emails are being sent"})

}

func matchingBySchedule(schedule entity.Schedule, organization string) ([]matcher.Match, database.JobStatus, error) {

	/*steps
	   - Scan schedule
	   - Load matchingtype
	      -if simple search users connected to dummy+ uid of schedule
		  -if group getUsersByTags(dummy_tag1+uid) , getUsersByTags(dummy_tag2+uid)
		  -if tag  getUsersByTags(tag1), getUsersByTags(tag2)
	   - Return
	*/
	var tuples []matcher.Match
	var out database.JobStatus
	switch schedule.MatchingType {
	case entity.Simple:
		techTag := "dummy_" + schedule.Name
		users, err := database.GetUsersByTechTag(schedule.Id, organization, techTag)

		if err != nil {
			out = database.Failed
			return nil, out, err
		}

		tuples = matcher.GenerateTuple(users, [][]entity.User{}, entity.Basic,
			[][]entity.User{}, uint(schedule.Size), []entity.User{}, []entity.User{})

	case entity.Groups:

		techTag1 := "dummy_group_" + strconv.Itoa(0) + "_" + schedule.Name
		techTag2 := "dummy_group_" + strconv.Itoa(1) + "_" + schedule.Name
		userGroup1, err := database.GetUsersByTechTag(schedule.Id, organization, techTag1)
		if err != nil {
			out = database.Failed
			return nil, out, err
		}
		userGroup2, err := database.GetUsersByTechTag(schedule.Id, organization, techTag2)
		if err != nil {
			out = database.Failed
			return nil, out, err
		}

		tuples = matcher.GenerateTuple([]entity.User{}, [][]entity.User{}, entity.Group,
			[][]entity.User{}, uint(schedule.Size), userGroup1, userGroup2)

	case entity.Tags:

		tags, err := database.GetTagBySchedule(schedule.Id)
		if err != nil {
			out = database.Failed

			return nil, out, err
		}
		userGroup1, err := database.GetUsersByTag(organization, tags[0].Name)
		if err != nil {
			out = database.Failed
			return nil, out, err
		}
		userGroup2, err := database.GetUsersByTag(organization, tags[1].Name)
		if err != nil {
			out = database.Failed
			return nil, out, err
		}
		tuples = matcher.GenerateTuple([]entity.User{}, [][]entity.User{}, entity.Group,
			[][]entity.User{}, uint(schedule.Size), userGroup1, userGroup2)

	}

	err := database.UpdateSchedule(schedule, organization)
	if err != nil {
		out = database.Failed
		return nil, out, err
	}
	out = database.Done
	return tuples, out, nil
}
func scheduleJob(organization string) {

	orgaUid, err := database.GetOrganizationByName(organization)
	log.Println(orgaUid)
	if err != nil {
		log.Println("Error while get uid of organization", err.Error())
		return
	}
	schedules, err := database.GetScheduleJob(organization)

	if err != nil {
		log.Println("Error while get schedules nodes in database", err)
		return
	}

	for _, schedule := range schedules {

		jobId := uuid.New().String()

		if err := database.CreateJobStatus(jobId); err != nil {
			log.Println("Error while creating job", jobId, err)
		}

		var status database.JobStatus

		tuples, status, err := matchingBySchedule(schedule, organization)
		log.Println(tuples)
		if err != nil {
			if er := database.UpdateJobErrors(jobId, []string{err.Error()}); er != nil {
				log.Println("Error while updating job", jobId, er)
			}
		} else {
			// sending mails
			// TODO:add code for authenticating
			//_, err := calendar.SendInvite(tuples, orgaUid)
			//if err != nil {
			//	log.Println("mails sending failed " + err.Error()) 681474948 kemayou fabrice , omer disgn

			//}
		}

		if err := database.UpdateJobStatus(jobId, status); err != nil {
			log.Println("Error while updating job", jobId, err)

		}
		if err := database.ScheduleLinkJobs(jobId, schedule.Id); err != nil {
			log.Println("Error while build a relationship schedule and job state", jobId, schedule.Id, err)
		}
		log.Println("execution of schedule is a success", schedule.Id)
	}

}
func processScheduledJob(c *gin.Context) {
	defer helper.Duration(helper.Track("processScheduledJob"))
	organisations, err := database.GetOrganizations()

	if err != nil {
		log.Println("Error while get all orga in database", err)
		return
	}
	for _, orga := range organisations {
		go scheduleJob(orga.Name)
	}
	c.JSON(http.StatusOK, gin.H{"message": "schedule are being excecute with successfull"})
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

	protected.POST("/schedule", handler.CreateScheduleType)

	protected.GET("/users-creation-job/:id", handler.GetJobStatus)
	protected.GET("/matching-email-job/:id", handler.GetJobStatus)
	protected.GET("/organizations/:id", handler.GetOrganization)
	protected.GET("/users", handler.GetUsers)
	protected.GET("/tags", getTags)
	protected.GET("/matchings-stats", handler.GetMatchingStats)
	protected.GET("/links", getLinks)
	protected.GET("/execute-schedule", processScheduledJob)

	protected.DELETE("/users", handler.DeleteUsers)
	protected.DELETE("/users/:id", handler.DeleteUser)

	router.Run()
}
