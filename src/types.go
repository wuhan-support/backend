package main

import "gopkg.in/go-playground/validator.v9"

type GeneralResponse struct {
	Message string `json:"message"`
}

type ReportRequest struct {
	Type string `json:"type" validate:"required"`
	Cause string `json:"cause" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type CustomValidator struct {
	validator *validator.Validate
}
