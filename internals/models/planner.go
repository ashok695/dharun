package models

import (
	"time"

	"github.com/gocql/gocql"
)

type TaskData struct {
	Project          string       `json:"project"`
	ID               gocql.UUID   `json:"id"`
	ActivePercentage float64      `json:"activePercentage"`
	ActualEnd        string       `json:"actualEnd"`
	ActualStart      string       `json:"actualStart"`
	AssignedToID     string       `json:"-"`
	AssignedTo       []UserData   `json:"assignedTo"`
	Constraintdate   string       `json:"constraintdate"`
	Constrainttype   string       `json:"constrainttype"`
	DelayLog         string       `json:"delayLog"`
	Duration         float64      `json:"duration"`
	MongoID          string       `json:"mongoID"`
	OrderID          string       `json:"orderID"`
	Phase            string       `json:"phaseID"`
	PlannedFrom      string       `json:"plannedFrom"`
	PlannedTo        string       `json:"plannedTo"`
	RefTaskID        string       `json:"refTaskID"`
	RoleID           string       `json:"-"`
	Role             []TypesData  `json:"role"`
	StatusID         string       `json:"-"`
	Status           []StatusData `json:"status"`
	SubPhaseID       string       `json:"subphaseID"`
	TaskTypeID       string       `json:"-"`
	TaskType         []TypesData  `json:"TaskType"`
	Title            string       `json:"name"`
	Type             string       `json:"type"`
	WorkstreamID     string       `json:"-"`
	Workstream       []TypesData  `json:"Workstream"`
	Children         []*TaskData  `json:"children"`
	ParentID         string       `json:"parentID"`
	OverDueDays      int16        `json:"overdueDays"`
	PlannedStart     *time.Time   `json:"plannedStart"`
	PlannedEnd       *time.Time   `json:"plannedEnd"`
	IsOverdue        bool         `json:"isOverdue"`
	Variance         int8         `json:"variance"`
	ActualStartDate  *time.Time   `json:"actualStartDate"`
	ActualEndDate    *time.Time   `json:"actualEndDate"`
	Level            string       `json:"level"`
}

type PlannerMap struct {
	UserMap   map[string]UserData
	StatusMap map[string]StatusData
	RwtMap    map[string]TypesData
}

type RWTMap struct {
	RoleMap       map[string]TypesData
	WorkstreamMap map[string]TypesData
	TasktypeMap   map[string]TypesData
}

type RWTData struct {
	RoleData       []TypesData
	WorkStreamData []TypesData
	TasktypeData   []TypesData
}

type DependenciesData struct {
	Project   string     `json:"project"`
	DBID      gocql.UUID `json:"dbid"`
	GanttType string     `json:"ganttype"`
	ID        string     `json:"id"`
	Lag       int32      `json:"lag"`
	LagUnit   string     `json:"lagunit"`
	MongoID   string     `json:"mongoid"`
	Source    string     `json:"source"`
	Target    string     `json:"target"`
	Type      int        `json:"type"`
}
type Project struct {
	Calendar     string `json:"calendar"`
	DaysPerMonth int8   `json:"daysPerMonth"`
	DaysPerWeek  int8   `json:"daysPerWeek"`
	EndDate      string `json:"endDate"`
	HoursPerDay  int8   `json:"hoursPerDay"`
	StartDate    string `json:"startDate"`
}
type PlannerStruct struct {
	Dependencies    []DependenciesData `json:"dependencies"`
	Project         Project            `json:"project"`
	Resources       []UserData         `json:"resources"`
	RoleTypeOptions []TypesData        `json:"roletypeoptions"`
	StatusType      []StatusData       `json:"statusType"`
	TaskType        []TypesData        `json:"taskType"`
	Tasks           []*TaskData        `json:"tasks"`
	WorkStream      []TypesData        `json:"workstreamtype"`
}
