package handler

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/koki/randommatch/convert"
	"github.com/koki/randommatch/database"
	"github.com/koki/randommatch/entity"
	"github.com/koki/randommatch/utils/helper"
)

type UsersFile struct {
	Organization string `form:"organization"`
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func UploadUsers(c *gin.Context) {
	defer helper.Duration(helper.Track("uploadUsers"))
	claims := c.MustGet("tokenClaims").(jwt.MapClaims)
	roles := claims["roles"].([]interface{})
	orgs := helper.ItemsWithPrefixInRole(roles, "Org.")


	var usersFile UsersFile

	if err := c.ShouldBind(&usersFile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid file sent " + err.Error()})
		return
	}
	usersFile.Organization = strings.ToLower(usersFile.Organization)
	if len(usersFile.Organization) > 0 && !helper.ContainsString(orgs, usersFile.Organization) {
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

func GetUsers(c *gin.Context) {
	defer helper.Duration(helper.Track("getUsers"))
	var err error
	var users []entity.User
	org, ok := c.GetQuery("organization");
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "organization is mandatory"})
		return
	}
	if tag, ok := c.GetQuery("tag"); !ok {
		users, err = database.GetUsers(org)
	}	else {
		//Get users having the tag specified in the query
		users, err = database.GetUsersByTag(org, tag)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func DeleteUser(c *gin.Context) {
	defer helper.Duration(helper.Track("deleteUser"))
	id := c.Param("id")
	claims := c.MustGet("tokenClaims").(jwt.MapClaims)
	roles := claims["roles"].([]interface{})
	if !helper.Contains(roles, "Privilege.Approve") {
		c.JSON(http.StatusForbidden, gin.H{"message": "Operation denied permission missing"})
		return
	}
	if err := database.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func DeleteUsers(c *gin.Context) {
	defer helper.Duration(helper.Track("deleteUsers"))
	claims := c.MustGet("tokenClaims").(jwt.MapClaims)
	roles := claims["roles"].([]interface{})
	if !helper.Contains(roles, "Privilege.Approve") {
		c.JSON(http.StatusForbidden, gin.H{"message": "Operation denied permission missing"})
		return
	}
	if err := database.DeleteUsers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}