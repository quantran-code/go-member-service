package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type memberRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (r memberRequest) validate() string {
	if strings.TrimSpace(r.Name) == "" {
		return "name is required"
	}
	email := strings.TrimSpace(r.Email)
	if email == "" {
		return "email is required"
	}
	if !strings.Contains(email, "@") || strings.HasPrefix(email, "@") || strings.HasSuffix(email, "@") {
		return "email must be a valid address"
	}
	return ""
}

// NewRouter wires the member API. Every error uses the project shape
// {"error": "..."}, and a duplicate email is a conflict rather than a 500.
func NewRouter(db *gorm.DB) http.Handler {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/members", func(c *gin.Context) {
		var req memberRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
			return
		}
		if message := req.validate(); message != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": message})
			return
		}
		member := Member{Name: strings.TrimSpace(req.Name), Email: strings.ToLower(strings.TrimSpace(req.Email))}
		if err := db.Create(&member).Error; err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}
		c.JSON(http.StatusCreated, member)
	})

	r.GET("/members", func(c *gin.Context) {
		var members []Member
		if err := db.Order("id").Find(&members).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "listing members"})
			return
		}
		if members == nil {
			members = []Member{}
		}
		c.JSON(http.StatusOK, members)
	})

	r.GET("/members/:id", func(c *gin.Context) {
		id, ok := memberID(c)
		if !ok {
			return
		}
		var member Member
		err := db.First(&member, id).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "loading member"})
			return
		}
		c.JSON(http.StatusOK, member)
	})

	r.PUT("/members/:id", func(c *gin.Context) {
		id, ok := memberID(c)
		if !ok {
			return
		}
		var req memberRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
			return
		}
		if message := req.validate(); message != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": message})
			return
		}
		var member Member
		if err := db.First(&member, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
			return
		}
		member.Name = strings.TrimSpace(req.Name)
		member.Email = strings.ToLower(strings.TrimSpace(req.Email))
		if err := db.Save(&member).Error; err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}
		c.JSON(http.StatusOK, member)
	})

	r.DELETE("/members/:id", func(c *gin.Context) {
		id, ok := memberID(c)
		if !ok {
			return
		}
		deleted := db.Delete(&Member{}, id)
		if deleted.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "deleting member"})
			return
		}
		if deleted.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
			return
		}
		c.Status(http.StatusNoContent)
	})

	return r
}

func memberID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid member id"})
		return 0, false
	}
	return uint(id), true
}
