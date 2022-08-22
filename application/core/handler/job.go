package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koki/randommatch/database"
	"github.com/koki/randommatch/utils/helper"
)
func GetJobStatus(c *gin.Context) {
	defer helper.Duration(helper.Track("getJobStatus"))
	id := c.Param("id")
	status, err := database.GetJobStatus(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": status})
}