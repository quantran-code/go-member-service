package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Member is a library member as the API exposes it.
type Member struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// MemberStore holds the members this service knows about. It is in memory only.
type MemberStore struct {
	members map[string]Member
}

// NewMemberStore returns a store seeded with a few members.
func NewMemberStore() *MemberStore {
	return &MemberStore{
		members: map[string]Member{
			"1": {ID: "1", Name: "Ada Lovelace", Email: "ada.lovelace@example.com"},
			"2": {ID: "2", Name: "Alan Turing", Email: "alan.turing@example.com"},
			"3": {ID: "3", Name: "Grace Hopper", Email: "grace.hopper@example.com"},
		},
	}
}

// Get returns the member with the given id; the second result is false when there is no such member.
func (s *MemberStore) Get(id string) (Member, bool) {
	member, ok := s.members[id]
	return member, ok
}

// errorBody is the JSON shape of an error response.
type errorBody struct {
	Error string `json:"error"`
}

// handleGetMember serves GET /members/:id.
func handleGetMember(store *MemberStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		log.Printf("request_id=%s method=%s path=%s member_id=%s", RequestID(c), c.Request.Method, c.Request.URL.Path, id)

		member, ok := store.Get(id)
		if !ok {
			c.JSON(http.StatusNotFound, errorBody{Error: "member not found"})
			return
		}
		c.JSON(http.StatusOK, member)
	}
}
