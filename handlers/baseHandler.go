package handlers

import (
	"getPromiseApi/db"
	"os"
)
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
type UserGet struct {
	Parce []db.User `json:"parce"`
}

type BaseHandler struct {
	db   *db.DB
	Code map[string]*db.User
}

func NewBaseHandler(pool *db.DB) *BaseHandler {
	return &BaseHandler{
		db:   pool,
		Code: make(map[string]*db.User),
	}
}
