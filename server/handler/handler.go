package handler

import "github.com/scch94/ins_log"

type Handler struct {
	Utfi string
}

func NewHandler() *Handler {
	handler := &Handler{}
	handler.Utfi = ins_log.GenerateUTFI()
	return handler
}
