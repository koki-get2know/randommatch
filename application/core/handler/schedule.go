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
	Groups       [][]entity.User `json:"group"`
	Tags         []string        `json:"tags"`
}

func CreateScheduleType(c *gin.Context) {

	defer helper.Duration(helper.Track("createShedule"))
	var req SchedulingReq

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json sent " + err.Error()})
		return
	}

	switch req.Schedule.MatchingType {
	case entity.Simple:
		/*
			  Steps:
				- Create schedule node in BD if not exist
				- Create a dummy tag for schedule
				- Assign a dummy tag to users
				- Assign a dummy tag to schedule
		*/
		scheduleTag := "dummy_" + req.Schedule.Name
		uid, err := database.CreateSchedule(req.Schedule, req.Organization)
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

		err = database.ScheduleLinkTTags(scheduleTag, uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	case entity.Groups:
		/* step :
		           - create a schedule if not exist
		   		- create dummy_tag1 and dummy_tag2  + uid of schedule
		   		- Make link between users from A and dummy_tag1 and users from B and dummy_tag2

		*/

		if len(req.Groups) < 2 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "you should send 2 groups"})
			return
		}
		uid, err := database.CreateSchedule(req.Schedule, req.Organization)
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

			err = database.ScheduleLinkTTags(scheduleTagGroup, uid)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
		}
	case entity.Tags:
		/*
		      step :
		            - create a schedule if not exist
		   		 - make link between schedule and coresponding tags

		*/

		if len(req.Tags) < 2 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "you should send 2 groups"})
			return
		}
		uid, err := database.CreateSchedule(req.Schedule, req.Organization)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		for _, tag := range req.Tags {

			err = database.ScheduleLinkTags(tag, uid)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
		}
	}

}
