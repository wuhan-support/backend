package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

const (
	DotSeparator      = "、"
	ColonSeparator    = ":"
	OrSymbolSeparator = "|"
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

func NewSubmissionSupplyFromShimoDoc(doc string) *SubmissionSupply {
	submissionSupply := &SubmissionSupply{}
	columns := strings.Split(doc, OrSymbolSeparator)
	if len(columns) != reflect.TypeOf(submissionSupply).Elem().NumField()-1 {
		fmt.Printf("%d, %d, %+v\n", len(columns), reflect.TypeOf(submissionSupply).Elem().NumField(), columns)
		return nil
	}
	submissionSupply.Name = columns[0]
	submissionSupply.Need, submissionSupply.Unit, _ = ParseSubmissionSupplyHumanString(columns[1])
	submissionSupply.Daily, _, _ = ParseSubmissionSupplyHumanString(columns[2])
	submissionSupply.Have, _, _ = ParseSubmissionSupplyHumanString(columns[3])
	_, _, submissionSupply.Requirements = ParseSubmissionSupplyHumanString(columns[4])

	return submissionSupply
}

func ParseSubmissionSupplyHumanString(humanString string) (string, string, string) {
	fields := strings.Split(humanString, ColonSeparator)
	if len(fields) != 2 {
		return "", "", ""
	}
	part := fields[1]
	length := len(part)
	if length < 3 {
		return part, part, part
	} else {
		return part[:length-3], part[length-3:], part
	}
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

func RefactCommunitySubmissionFromShimoDoc(doc []byte) ([]byte, error) {
	slice := make([]map[string]interface{}, 0)
	err := json.Unmarshal(doc, &slice)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, m := range slice {
		if ms, exists := m["medicalsupplies"]; exists {
			if ms != nil && ms != "" {
				m["medicalsupplies"] = RefactSubmissionSupply(ms.(string))
			}
		}
		if ms, exists := m["livesupplies"]; exists {
			if ms != nil && ms != "" {
				m["livesupplies"] = RefactSubmissionSupply(ms.(string))
			}
		}
	}

	return json.Marshal(slice)
}

func RefactSubmissionSupply(humanString string) []*SubmissionSupply {
	supplySlice := make([]*SubmissionSupply, 0)
	elements := strings.Split(humanString, DotSeparator)
	for _, e := range elements {
		submissionSupply := NewSubmissionSupplyFromShimoDoc(e)
		if submissionSupply != nil {
			supplySlice = append(supplySlice, submissionSupply)
		}
	}
	return supplySlice
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
	return strings.Join(r, DotSeparator)
}

type GetSubmissionsRequest struct {
	Page  int `json:"page" validate:"required"`
	Limit int `json:"limit" validate:"required"`
}
