package main

import "time"

// User represents a user in the system
type User struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"` // Do not expose hashed password in JSON
	CreatedAt      time.Time `json:"created_at"`
	Profile        Profile   `json:"profile"`
}

// Profile represents a user's profile
type Profile struct {
	UserID    int       `json:"user_id"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
}
