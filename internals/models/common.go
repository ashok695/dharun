package models

import "github.com/gocql/gocql"

type UserData struct {
	ID       gocql.UUID `json:"id"`
	FullName string     `json:"fullName"`
	MongoID  string     `json:"mongoID"`
	Role     string     `json:"role"`
	Team     string     `json:"team"`
	Project  string     `json:"project"`
}

type StatusData struct {
	ID       gocql.UUID `json:"id"`
	Category string     `json:"category"`
	Color    string     `json:"color"`
	MongoID  string     `json:"mongoID"`
	Status   string     `json:"status"`
	WorkItem string     `json:"workitem"`
}

type TypesData struct {
	ID      gocql.UUID `json:"id"`
	Type    string     `json:"type"`
	MongoID string     `json:"mongoID"`
	Name    string     `json:"name"`
}
