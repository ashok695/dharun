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

		statusMap[status.MongoID] = status

	}
	return statusMap
}

func TypesMap(typesData []models.TypesData) (models.RWTMap, models.RWTData) {
	workStreamMap := make(map[string]models.TypesData, len(typesData))
	roleMap := make(map[string]models.TypesData, len(typesData))
	tasktypeMap := make(map[string]models.TypesData, len(typesData))
	workStreamData := make([]models.TypesData, 0, len(typesData))
	roleData := make([]models.TypesData, 0, len(typesData))
	tasktypeData := make([]models.TypesData, 0, len(typesData))
	for _, types := range typesData {
		if types.Type == "role" {
			roleData = append(roleData, types)
			roleMap[types.MongoID] = types
		} else if types.Type == "workstream" {
			workStreamData = append(workStreamData, types)
			workStreamMap[types.MongoID] = types
		} else if types.Type == "tasktype" {
			tasktypeData = append(tasktypeData, types)
			tasktypeMap[types.MongoID] = types
		}
	}
	rwtmap := models.RWTMap{
		RoleMap:       roleMap,
		WorkstreamMap: workStreamMap,
		TasktypeMap:   tasktypeMap,
	}
	rwtData := models.RWTData{
		RoleData:       roleData,
		WorkStreamData: workStreamData,
		TasktypeData:   tasktypeData,
	}
	return rwtmap, rwtData
}
