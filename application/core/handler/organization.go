package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/koki/randommatch/database"
	"github.com/koki/randommatch/entity"
	"github.com/koki/randommatch/utils/helper"
)

func GetOrganization(c *gin.Context) {
	defer helper.Duration(helper.Track("getOrganization"))
	claims := c.MustGet("tokenClaims").(jwt.MapClaims)
	roles := claims["roles"].([]interface{})
	if !helper.Contains(roles, "Privilege.Approve") {
		c.JSON(http.StatusForbidden, gin.H{"message": "Operation denied permission missing"})
		return
	}
	id := c.Param("id")

	orga, err := database.GetOrganizationById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message":  err.Error()})
		return
	}
	if len(orga.Id) == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": orga})
}

func CreateOrganization(c *gin.Context) {
	defer helper.Duration(helper.Track("createOrganization"))
	claims := c.MustGet("tokenClaims").(jwt.MapClaims)
	roles := claims["roles"].([]interface{})
	if !helper.Contains(roles, "Privilege.Approve") {
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
