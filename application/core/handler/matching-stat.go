package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koki/randommatch/database"
	"github.com/koki/randommatch/utils/helper"
)

func GetMatchingStats(c *gin.Context) {
	defer helper.Duration(helper.Track("getMatchings"))
	matchings, err := database.GetMatchingStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": matchings})
}
