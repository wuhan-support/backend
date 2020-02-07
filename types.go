package main

import (
	"gopkg.in/go-playground/validator.v9"
)

type GeneralResponse struct {
	Message string `json:"message"`
}

type ReportRequest struct {
	Type    string `json:"type" validate:"required"`
	Cause   string `json:"cause" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type CustomValidator struct {
	validator *validator.Validate
}

type SubmissionSupply struct {
	//gorm.Model
	Name string `json:"name,omitempty" gorm:"" validate:"required"`
	Need string `json:"need,omitempty" gorm:"" validate:"required"`
	Daily string `json:"daily,omitempty" gorm:"" validate:"required"`
	Have string `json:"have,omitempty" gorm:"" validate:"required"`
	Requirements string `json:"requirements,omitempty" gorm:"" validate:"required"`
}

type Submission struct {
	//ID              int    `json:"id" gorm:"primary_key;not null;auto_increment"`
	//gorm.Model
	Name            string             `json:"name,omitempty" gorm:"" validate:"required,max=100"`
	Province        string             `json:"province,omitempty" gorm:"" validate:"required"`
	City            string             `json:"city,omitempty" gorm:"" validate:"required"`
	Suburb          string             `json:"suburb,omitempty" gorm:"" validate:"required"`
	Address         string             `json:"address,omitempty" gorm:"" validate:"required"`
	Patients        string                `json:"patients,omitempty" gorm:"" validate:""`
	Beds            string                `json:"beds,omitempty" gorm:"" validate:""`
	ContactName     string             `json:"contactName,omitempty" gorm:"" validate:""`
	ContactOrg      string             `json:"contactOrg,omitempty" gorm:"" validate:"required"`
	ContactPhone    string             `json:"contactPhone,omitempty" gorm:"" validate:"required"`
	Supplies        []SubmissionSupply `json:"supplies,omitempty" gorm:"" validate:"required"`
	Pathways        string             `json:"pathways,omitempty" gorm:"" validate:"required"`
	LogisticStatus  string             `json:"logisticStatus,omitempty" gorm:"" validate:""`
	Source          string             `json:"source,omitempty" gorm:"" validate:""`
	Proof           string             `json:"proof,omitempty" gorm:"" validate:""`
	Notes           string             `json:"notes,omitempty" gorm:"type:text" validate:""`
}

type GetSubmissionsRequest struct {
	Page int `json:"page" validate:"required"`
	Limit int `json:"limit" validate:"required"`
}
