package main

import "gopkg.in/go-playground/validator.v9"

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

type CollectForm struct {
	ID              int    `json:"id" gorm:"primary_key;not null;auto_increment"`
	Name            string `json:"name,omitempty" gorm:"" validate:"required,max=100"`
	Province        string `json:"province,omitempty" gorm:"" validate:"required"`
	City            string `json:"city,omitempty" gorm:"" validate:"required"`
	Suburb          string `json:"suburb,omitempty" gorm:"" validate:"required"`
	Address         string `json:"address,omitempty" gorm:"" validate:"required"`
	Patients        int    `json:"patients,omitempty" gorm:"" validate:"required"`
	Beds            int    `json:"beds,omitempty" gorm:"" validate:"required"`
	ContactName     string `json:"contactName,omitempty" gorm:"" validate:""`
	ContactOrg      string `json:"contactOrg,omitempty" gorm:"" validate:"required"`
	ContactPhone    string `json:"contactPhone,omitempty" gorm:"" validate:"required"`
	Supplies        string `json:"supplies,omitempty" gorm:"" validate:"required"`
	Pathways        string `json:"pathways,omitempty" gorm:"" validate:"personal|redcross"`
	LogisticsStatus string `json:"logisticsStatus,omitempty" gorm:"" validate:""`
	Source          string `json:"source,omitempty" gorm:"" validate:""`
	Proof           string `json:"proof,omitempty" gorm:"" validate:""`
	Notes           string `json:"notes,omitempty" gorm:"type:text" validate:""`
	CreateTime      string `json:"createTime" gorm:"type:timestamp;default:current_timestamp" validate:""`
}
