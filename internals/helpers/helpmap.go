package helpers

import "github.com/dharun/poc/internals/models"

func Usermap(userdata []models.UserData) map[string]models.UserData {
	usermap := make(map[string]models.UserData, len(userdata))

	for _, user := range userdata {
		usermap[user.MongoID] = user
	}
	return usermap
}

func StatusMap(statusData []models.StatusData) map[string]models.StatusData {
	statusMap := make(map[string]models.StatusData, len(statusData))

	for _, status := range statusData {
		if status.WorkItem == "Task" {
			statusMap[status.MongoID] = status
		}
	}
	return statusMap
}

func TypesMap(typesData []models.TypesData) models.RWTMap {
	workStreamMap := make(map[string]models.TypesData, len(typesData))
	roleMap := make(map[string]models.TypesData, len(typesData))
	tasktypeMap := make(map[string]models.TypesData, len(typesData))
	for _, types := range typesData {
		if types.Type == "role" {
			roleMap[types.MongoID] = types
		} else if types.Type == "workstream" {
			workStreamMap[types.MongoID] = types
		} else if types.Type == "tasktype" {
			tasktypeMap[types.MongoID] = types
		}
	}
	rwtmap := models.RWTMap{
		RoleMap:       roleMap,
		WorkstreamMap: workStreamMap,
		TasktypeMap:   tasktypeMap,
	}
	return rwtmap
}
