package helpers

import (
	"fmt"
	"github.com/dharun/poc/database"
	"github.com/dharun/poc/internals/models"
	"github.com/gocql/gocql"
	"time"
)

func GetUserData(projectID string) ([]models.UserData, error) {
	var UserData []models.UserData
	var Id gocql.UUID
	var fullName, mongoID, team string

	// ✅ Fix the query: Remove single quotes around '?'
	iter := database.Session.Query(`SELECT id, fullname, mongoid, team FROM users WHERE project = ?`, projectID).Iter()

	for iter.Scan(&Id, &fullName, &mongoID, &team) {
		user := models.UserData{
			ID:       Id,
			FullName: fullName,
			MongoID:  mongoID,
			Team:     team,
		}
		UserData = append(UserData, user)
	}

	// ✅ Handle errors properly
	err := iter.Close()
	if err != nil {
		fmt.Println("Error closing iterator:", err)
		return nil, err
	}
	fmt.Println("user done", len(UserData))
	return UserData, nil
}

func GetStatusData(projectID string) ([]models.StatusData, error) {
	var statusData []models.StatusData

	var id gocql.UUID
	var category, color, mongoid, status, workitem string

	iter := database.Session.Query(`SELECT id,category,color,mongoid,status,workitem from status WHERE project = ?`, projectID).Iter()
	for iter.Scan(&id, &category, &color, &mongoid, &status, &workitem) {
		status := models.StatusData{
			ID:       id,
			Category: category,
			Color:    color,
			MongoID:  mongoid,
			Status:   status,
			WorkItem: workitem,
		}
		statusData = append(statusData, status)
	}
	err := iter.Close()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Error in getting status data")
		return nil, err
	}
	fmt.Println("Status Done")
	return statusData, nil
}

func GetTypesData(projectID string) ([]models.TypesData, error) {
	var typesData []models.TypesData
	var id gocql.UUID
	var mongodid, name, types string

	iter := database.Session.Query(`SELECT id,mongoid,name,type FROM types WHERE project = ?`, projectID).Iter()
	for iter.Scan(&id, &mongodid, &name, &types) {
		types := models.TypesData{
			ID:      id,
			MongoID: mongodid,
			Name:    name,
			Type:    types,
		}
		typesData = append(typesData, types)
	}
	err := iter.Close()
	if err != nil {
		fmt.Println("ERROR in RWT data")
		return nil, err
	}
	fmt.Println("RWT done")
	return typesData, nil
}

// func GetTeamsData() ([]models.TeamData, error) {
// 	var teamsData []models.TeamData
// 	var id gocql.UUID
// 	var mongodid, title string

//		iter := database.Session.Query("SELECT id,mongoid,title FROM dharun.roles").Iter()
//		for iter.Scan(&id, &mongodid, &title) {
//			types := models.TeamData{
//				ID:      id,
//				MongoID: mongodid,
//				Title:   title,
//			}
//			teamsData = append(teamsData, types)
//		}
//		err := iter.Close()
//		if err != nil {
//			fmt.Println("ERROR in RWT data")
//			return nil, err
//		}
//		fmt.Println("RWT done")
//		return teamsData, nil
//	}
func ParseDate(date string) *time.Time {
	layout := "2006-01-02T15:04:05Z"
	if date == "" {
		return nil
	}
	pdate, err := time.Parse(layout, date)
	if err != nil {
		fmt.Println("Error in parsing the date", err.Error())
	}
	return &pdate
}
