package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/koki/randommatch/database"
	"github.com/koki/randommatch/entity"
	"github.com/koki/randommatch/utils/helper"
)

type SchedulingReq struct {
	Organization string          `json:"organization"`
	Schedule     entity.Schedule `json:"schedule"`
	Users        []entity.User   `json:"users"`
}
type SchedulingGroupReq struct {
	Organization string          `json:"organization"`
	Schedule     entity.Schedule `json:"schedule"`
	Groups       [][]entity.User `json:"group"`
}

type SchedulingTagReq struct {
	Organization string          `json:"organization"`
	Schedule     entity.Schedule `json:"schedule"`
	Tags         []string        `json:"tags"`
}

// Schedule for simple matching
func CreateSchedule(c *gin.Context) {
	/*
		  Steps:
			- Create schedule node in BD if not exist
			- Create a dummy tag for schedule
			- Assign a dummy tag to users
			- Assign a dummy tag to schedule
	*/
	defer helper.Duration(helper.Track("createShedule"))
	var req SchedulingReq

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json sent " + err.Error()})
		return
	}
	scheduleTag := "dummy_" + req.Schedule.Name
	err := database.CreateSchedule(req.Schedule, req.Organization)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	err = database.CreateTechTags(scheduleTag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	err = database.UserLinkTags(req.Users, scheduleTag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	err = database.ScheduleLinkTTags(scheduleTag, req.Schedule.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
}

// Schedule for group matching
func CreateScheduleGroup(c *gin.Context) {
	/* step :
	           - create a schedule if not exist
	   		- create dummy_tag1 and dummy_tag2  + uid of schedule
	   		- Make link between users from A and dummy_tag1 and users from B and dummy_tag2

	*/
	defer helper.Duration(helper.Track("createShedule"))
	var req SchedulingGroupReq

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json sent " + err.Error()})
		return
	}
	if len(req.Groups) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "you should send 2 groups"})
		return
	}
	err := database.CreateSchedule(req.Schedule, req.Organization)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	for i, users := range req.Groups {
		scheduleTagGroup := "dummy_group_" + strconv.Itoa(i) + "_" + req.Schedule.Name

		err = database.CreateTechTags(scheduleTagGroup)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		err = database.UserLinkTags(users, scheduleTagGroup)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		err = database.ScheduleLinkTTags(scheduleTagGroup, req.Schedule.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}
}

// Schedule for tag matching
func CreateScheduleTag(c *gin.Context) {
	/*
	      step :
	            - create a schedule if not exist
	   		 - make link between schedule and coresponding tags

	*/
	defer helper.Duration(helper.Track("createShedule"))
	var req SchedulingTagReq

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json sent " + err.Error()})
		return
	}
	if len(req.Tags) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "you should send 2 groups"})
		return
	}
	err := database.CreateSchedule(req.Schedule, req.Organization)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	for _, tag := range req.Tags {

		err = database.ScheduleLinkTags(tag, req.Schedule.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

}
