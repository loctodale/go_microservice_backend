// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"
	"time"
)

type PreGoKeyToken9999 struct {
	TokenID      uint64
	ShopID       uint64
	PublicKey    string
	PrivateKey   string
	RefreshToken string
	KeyCreatedAt time.Time
	KeyUpdatedAt time.Time
	KeyDeletedAt time.Time
}

type PreGoShopBase9999 struct {
	ShopID        uint64
	ShopAccount   string
	ShopPassword  string
	ShopStatus    sql.NullInt32
	ShopCreatedAt time.Time
	ShopUpdatedAt time.Time
	ShopDeletedAt time.Time
}

type PreGoTokenUsed9999 struct {
	TokenUsedID      uint64
	TokenID          uint64
	RefreshTokenUsed string
}