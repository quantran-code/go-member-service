// Command member-service owns library members.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("member service listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, NewRouter()))
}

// NewRouter wires the member API.
func NewRouter() http.Handler {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(RequestIDMiddleware())
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	members := NewMemberStore()
	r.GET("/members/:id", handleGetMember(members))
	return r
}
