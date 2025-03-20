package models

import (
	"time"

	"github.com/gocql/gocql"
)

type Workbook struct {
	Project           string       `json:"project"`
	ID                gocql.UUID   `json:"id"`
	ActualEnd         string       `json:"actualEnd"`
	ActualStart       string       `json:"actualStart"`
	AssignedToID      string       `json:"-"`
	AssignedTo        []UserData   `json:"assignedTo"`
	Duration          float64      `json:"duration"`
	MongoID           string       `json:"mongoid"`
	OrderID           string       `json:"orderID"`
	Phase             string       `json:"phaseID"`
	PlannedFrom       string       `json:"plannedFrom"`
	PlannedTo         string       `json:"plannedTo"`
	RefTaskID         string       `json:"refTaskID"`
	RoleID            string       `json:"-"`
	Role              []TypesData  `json:"role"`
	StatusID          string       `json:"-"`
	Status            []StatusData `json:"status"`
	SubPhaseID        string       `json:"subphaseID"`
	TaskTypeID        string       `json:"-"`
	TaskType          []TypesData  `json:"TaskType"`
	Title             string       `json:"name"`
	Type              string       `json:"type"`
	WorkstreamID      string       `json:"-"`
	Workstream        []TypesData  `json:"Workstream"`
	TaskLevels        string       `json:"taskLevels"`
	PlannedStart      *time.Time   `json:"startDate"`
	PlannedEnd        *time.Time   `json:"endDate"`
	ActualStartDate   *time.Time   `json:"actualStartDate"`
	ActualEndDate     *time.Time   `json:"actualEndDate"`
	ActivePercentage  float64      `json:"activePercentage"`
	Startvariance     string       `json:"startvariance"`
	ParentPhase       string       `json:"parentPhase"`
	ParentWorkpackage string       `json:"parentWorkpackage"`
	ParenttaskL1      string       `json:"parentTaskL1"`
	ParenttaskL2      string       `json:"parentTaskL2"`
}
