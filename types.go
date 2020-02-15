package main

import (
	"strings"

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
	Name         string `json:"name,omitempty" gorm:"" validate:"required"`
	Unit         string `json:"unit,omitempty" gorm:"" validate:"required"`
	Need         string `json:"need,omitempty" gorm:"" validate:"required"`
	Daily        string `json:"daily,omitempty" gorm:"" validate:"required"`
	Have         string `json:"have,omitempty" gorm:"" validate:"required"`
	Requirements string `json:"requirements,omitempty" gorm:"" validate:"required"`
}

// 口罩|需求数量:10个|每日消耗1:个|库存:0个|要求:
func (s *SubmissionSupply) HumanString() string {
	var sb strings.Builder
	sb.WriteString(s.Name)
	sb.WriteString("|需求数量:")
	sb.WriteString(s.Need)
	sb.WriteString(s.Unit)
	sb.WriteString("|每日消耗:")
	sb.WriteString(s.Daily)
	sb.WriteString(s.Unit)
	sb.WriteString("|库存:")
	sb.WriteString(s.Have)
	sb.WriteString(s.Unit)
	sb.WriteString("|要求:")
	sb.WriteString(s.Requirements)
	return sb.String()
}

type Submission struct {
	//ID              int    `json:"id" gorm:"primary_key;not null;auto_increment"`
	//gorm.Model
	Name           string             `json:"name,omitempty" gorm:"" validate:"required,max=100"`
	Province       string             `json:"province,omitempty" gorm:"" validate:"required"`
	City           string             `json:"city,omitempty" gorm:"" validate:"required"`
	Suburb         string             `json:"suburb,omitempty" gorm:"" validate:"required"`
	Address        string             `json:"address,omitempty" gorm:"" validate:"required"`
	Patients       string             `json:"patients,omitempty" gorm:"" validate:""`
	Beds           string             `json:"beds,omitempty" gorm:"" validate:""`
	ContactName    string             `json:"contactName,omitempty" gorm:"" validate:""`
	ContactOrg     string             `json:"contactOrg,omitempty" gorm:"" validate:"required"`
	ContactPhone   string             `json:"contactPhone,omitempty" gorm:"" validate:"required"`
	Supplies       []SubmissionSupply `json:"supplies,omitempty" gorm:"" validate:"required"`
	Pathways       string             `json:"pathways,omitempty" gorm:"" validate:"required"`
	LogisticStatus string             `json:"logisticStatus,omitempty" gorm:"" validate:""`
	Source         string             `json:"source,omitempty" gorm:"" validate:""`
	Proof          string             `json:"proof,omitempty" gorm:"" validate:""`
	Notes          string             `json:"notes,omitempty" gorm:"type:text" validate:""`
}

// 社区物资需求
type CommunitySubmission struct {
	Name              string              `json:"name,omitempty" gorm:"" validate:"required,max=100"`
	Age               int                 `json:"age,omitempty" gorm:"" validate:"required,min=1,max=200"`
	Province          string              `json:"province,omitempty" gorm:"" validate:"required"`
	City              string              `json:"city,omitempty" gorm:"" validate:"required"`
	Suburb            string              `json:"suburb,omitempty" gorm:"" validate:"required"`
	Address           string              `json:"address,omitempty" gorm:"" validate:"required"`
	ContactPhone      string              `json:"contactPhone,omitempty" gorm:"" validate:"omitempty"`
	AgentName         string              `json:"agentName,omitempty" gorm:"" validate:"omitempty"`
	AgentContactPhone string              `json:"agentContactPhone,omitempty" gorm:"" validate:"omitempty"`
	MedicalSupplies   []*SubmissionSupply `json:"medicalSupplies,omitempty" gorm:"" validate:"omitempty"`
	LiveSupplies      []*SubmissionSupply `json:"liveSupplies,omitempty" gorm:"" validate:"omitempty"`
	NeedsVehicle      bool                `json:"needsVehicle,omitempty" gorm:"" validate:"omitempty"`
	Notes             string              `json:"notes,omitempty" gorm:"type:text" validate:"omitempty"`
}

func (c *CommunitySubmission) Values() []interface{} {
	r := []interface{}{
		c.Name,
		c.Age,
		c.Province,
		c.City,
		c.Suburb,
		c.Address,
		c.ContactPhone,
		c.AgentName,
		c.AgentContactPhone,
		JoinSubmissionSupplySlice(c.MedicalSupplies),
		JoinSubmissionSupplySlice(c.LiveSupplies),
		c.NeedsVehicle,
		c.Notes,
	}
	return r
}

func JoinSubmissionSupplySlice(supplySlice []*SubmissionSupply) string {
	r := make([]string, 0)
	for _, s := range supplySlice {
		r = append(r, s.HumanString())
	}
	return strings.Join(r, "、")
}

type GetSubmissionsRequest struct {
	Page  int `json:"page" validate:"required"`
	Limit int `json:"limit" validate:"required"`
}
