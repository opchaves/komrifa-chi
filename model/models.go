// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package model

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SchemaMigration struct {
	Version int64 `json:"version"`
	Dirty   bool  `json:"dirty"`
}

type User struct {
	ID                uuid.UUID        `json:"id"`
	Name              string           `json:"name"`
	Email             string           `json:"email"`
	Password          string           `json:"password"`
	Verified          bool             `json:"verified"`
	Active            bool             `json:"active"`
	VerificationToken *string          `json:"verification_token"`
	Avatar            string           `json:"avatar"`
	CreatedAt         pgtype.Timestamp `json:"created_at"`
	UpdatedAt         pgtype.Timestamp `json:"updated_at"`
}
