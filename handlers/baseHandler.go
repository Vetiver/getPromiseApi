package handlers

import (
	"getPromiseApi/db"
)

type UserGet struct {
	Parce []db.User `json:"parce"`
 }


 type BaseHandler struct {
	db *db.DB
 }

func NewBaseHandler(db *db.DB) *BaseHandler {
	return &BaseHandler{
		db: db,
	}
}
