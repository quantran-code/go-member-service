package main

import "time"

// Member is a person who can borrow books. Email is the natural identity: two
// members with the same address would be the same person.
type Member struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
