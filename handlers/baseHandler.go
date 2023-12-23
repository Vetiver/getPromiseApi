package handlers

import (
	"getPromiseApi/db"
)

type BaseHandler struct {
	db *db.DB
}

func NewBaseHandler(db *db.DB) *BaseHandler {
	return &BaseHandler{
		db: db,
	}
}
