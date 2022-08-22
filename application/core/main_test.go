package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/koki/randommatch/handler"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestGetAlbums(t *testing.T) {
	r := SetUpRouter()
	r.GET("/albums", handler.GetAlbums)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/albums", nil)
	r.ServeHTTP(w, req)

	var albums []handler.Album
	json.Unmarshal(w.Body.Bytes(), &albums)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, albums)
}

func TestGetAlbumByID(t *testing.T) {
	r := SetUpRouter()
	r.GET("/albums/:id", handler.GetAlbumByID)

	req, _ := http.NewRequest("GET", "/albums/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var album handler.Album
	json.Unmarshal(w.Body.Bytes(), &album)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, album.ID, "1")
	assert.Equal(t, album.Title, "Blue Train")
	assert.Equal(t, album.Artist, "John Coltrane")
	assert.Equal(t, album.Price, 56.99)

	req, _ = http.NewRequest("GET", "/albums/10", nil)
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPostAlbums(t *testing.T) {
	r := SetUpRouter()
	r.POST("/albums", handler.PostAlbums)
	albumId := xid.New().String()
	album := handler.Album{
		ID:     albumId,
		Title:  "Cryto is BTC the GOAT",
		Artist: "The GOAT",
		Price:  40.80,
	}
	jsonValue, _ := json.Marshal(album)
	req, _ := http.NewRequest("POST", "/albums", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestHealthCheck(t *testing.T) {
	r := SetUpRouter()
	r.GET("/health-check", getHealthCheck)
	req, _ := http.NewRequest("GET", "/health-check", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}
