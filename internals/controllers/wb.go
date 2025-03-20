package controllers

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dharun/poc/database"
	"github.com/dharun/poc/internals/helpers"
	"github.com/dharun/poc/internals/models"
	"github.com/facette/natsort"
	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
)

func GetWorkbook(c *fiber.Ctx) ([]map[string]interface{}, error) {
	start := time.Now()
	var tasksData []models.Workbook
	var subphaseData []models.Workbook
	statusMap := make(map[string]models.StatusData)
	userMap := make(map[string]models.UserData)
	var typesMap models.RWTMap
	errChan := make(chan error, 4)

	var wg sync.WaitGroup
	wg.Add(5)
	projectID := c.Query("projectID")
	go func() {
		defer wg.Done()
		td, err := GetTaskWorkbookData(projectID)
		if err != nil {
			errChan <- err
		}
		tasksData = td
	}()
	go func() {
		defer wg.Done()
		sd, err := helpers.GetStatusData(projectID)
		if err != nil {
			errChan <- err
		}
		statusMap = helpers.StatusMap(sd)
	}()
	go func() {
		defer wg.Done()
		ud, err := helpers.GetUserData(projectID)
		if err != nil {
			errChan <- err
		}
		userMap = helpers.Usermap(ud)
	}()
	go func() {
		defer wg.Done()
		td, err := helpers.GetTypesData(projectID)
		if err != nil {
			errChan <- err
		}
		typesMap, _ = helpers.TypesMap(td)
	}()
	go func() {
		defer wg.Done()
		sd, err := GetSubphaseWorkbookData(projectID)
		if err != nil {
			errChan <- err
		}
		subphaseData = sd
	}()
	wg.Wait()
	select {
	case err := <-errChan:
		return nil, err
	default:
	}
	fmt.Println(len(statusMap), len(userMap), len(typesMap.RoleMap), len(subphaseData))
	highestlevelTasks, firstleveltasks := AssignMetaDataFortasks(tasksData, subphaseData, statusMap, userMap, typesMap)
	AssignMetaDataForSubphases(subphaseData, firstleveltasks)
	f := append(subphaseData, tasksData...)
	fmt.Println("length of overall", len(f))
	ss := time.Now()
	sort.Slice(f, func(i, j int) bool {
		return natsort.Compare(f[i].OrderID, f[j].OrderID)
	})
	fmt.Println("time talem for Sorting", time.Since(ss))
	a := FormulateWorkbookStructure(f, highestlevelTasks)
	fmt.Println("time taken for workbook api", time.Since(start))
	return a, nil
}

func GetTaskWorkbookData(projectID string) ([]models.Workbook, error) {
	var tasksData []models.Workbook
	var id gocql.UUID
	var activePercentage, duration float64
	var project, actualend, actualstart, assignedto, mongoid, orderid, phase,
		plannedfrom, plannedto, reftaskid, roletype, status, subphaseid, tasktype, title, workstream string

	iter := database.Session.Query(`SELECT project,id, activepercentage,actualend,actualstart, assignedto,duration,mongoid, orderid, phase,
		plannedfrom, plannedto,reftaskid, roletype, status, subphaseid, tasktype, title,workstream
		FROM tasklists WHERE PROJECT =?`, projectID).Iter()

	// Scan the result into respective variables
	for iter.Scan(&project, &id, &activePercentage, &actualend, &actualstart, &assignedto, &duration, &mongoid, &orderid, &phase, &plannedfrom,
		&plannedto, &reftaskid, &roletype, &status, &subphaseid, &tasktype, &title, &workstream) {

		// Create TaskData struct and append it to tasksData
		task := models.Workbook{
			Project:          project,
			ID:               id,
			ActivePercentage: activePercentage,
			ActualEnd:        actualend,
			ActualStart:      actualstart,
			AssignedToID:     assignedto,
			Duration:         duration,
			MongoID:          mongoid,
			OrderID:          orderid,
			Phase:            phase,
			PlannedFrom:      plannedfrom,
			PlannedTo:        plannedto,
			RefTaskID:        reftaskid,
			RoleID:           roletype,
			StatusID:         status,
			SubPhaseID:       subphaseid,
			TaskTypeID:       tasktype,
			Title:            title,
			WorkstreamID:     workstream,
		}
		tasksData = append(tasksData, task)
	}

	err := iter.Close()
	if err != nil {
		fmt.Println("ERROR IN DATA TASKLISTS", err.Error())
		return nil, err
	}
	fmt.Println("len", len(tasksData))
	return tasksData, nil
}

func GetSubphaseWorkbookData(projectID string) ([]models.Workbook, error) {
	var subphaseData []models.Workbook
	var id gocql.UUID
	var duration float64
	var project, mongoid, orderid, phase,
		role, title, tasktype, workstream string

	iter := database.Session.Query(`SELECT project,id,duration,mongoid,orderid,phase,
	        role,title,type,workstream
	        from subphases WHERE project = ?`, projectID).Iter()

	for iter.Scan(&project, &id, &duration, &mongoid, &orderid, &phase, &role, &title, &tasktype, &workstream) {
		task := models.Workbook{
			Project:          project,
			ID:               id,
			ActivePercentage: 0,
			ActualEnd:        "",
			ActualStart:      "",
			AssignedToID:     "",
			Duration:         duration,
			MongoID:          mongoid,
			OrderID:          orderid,
			Phase:            phase,
			PlannedFrom:      "",
			PlannedTo:        "",
			RefTaskID:        "",
			RoleID:           role,
			StatusID:         "",
			SubPhaseID:       "",
			Title:            title,
			TaskTypeID:       tasktype,
			Type:             "",
			WorkstreamID:     workstream,
		}
		subphaseData = append(subphaseData, task)
	}
	err := iter.Close()
	if err != nil {
		fmt.Println("ERROR IN DATA ASHOK", err.Error())
		return nil, err
	}
	return subphaseData, nil
}

func AssignMetaDataFortasks(tasksData []models.Workbook, subphaseData []models.Workbook, statusMap map[string]models.StatusData, userMap map[string]models.UserData, typesMap models.RWTMap) (int, []models.Workbook) {
	mapfortask := make(map[string]models.Workbook)
	mapforsubphase := make(map[string]models.Workbook)
	for i := range tasksData {
		mapfortask[tasksData[i].OrderID] = tasksData[i]
	}
	for i := range subphaseData {
		mapforsubphase[subphaseData[i].OrderID] = subphaseData[i]
	}
	highestLevel := 1
	var firstLeveltasks []models.Workbook
	for i := range tasksData {
		wbsSplit := strings.Split(tasksData[i].OrderID, ".")
		level := len(wbsSplit) - 2
		if level > highestLevel {
			highestLevel = level
		}
		tasksData[i].TaskLevels = "taskL" + strconv.Itoa(level)
		//Assigning owner
		if owner, exists := userMap[tasksData[i].AssignedToID]; exists {
			tasksData[i].AssignedTo = []models.UserData{owner}
		}
		//Assigning Status
		if status, exists := statusMap[tasksData[i].StatusID]; exists {
			tasksData[i].Status = []models.StatusData{status}
		}
		if role, exists := typesMap.RoleMap[tasksData[i].RoleID]; exists {
			tasksData[i].Role = []models.TypesData{role}
		}
		if ws, exists := typesMap.WorkstreamMap[tasksData[i].WorkstreamID]; exists {
			tasksData[i].Workstream = []models.TypesData{ws}
		}
		if tt, exists := typesMap.TasktypeMap[tasksData[i].TaskTypeID]; exists {
			tasksData[i].TaskType = []models.TypesData{tt}
		}
		tasksData[i].PlannedStart = helpers.ParseDate(tasksData[i].PlannedFrom)
		tasksData[i].PlannedEnd = helpers.ParseDate(tasksData[i].PlannedTo)
		tasksData[i].ActualStartDate = helpers.ParseDate(tasksData[i].ActualStart)
		tasksData[i].ActualEndDate = helpers.ParseDate(tasksData[i].ActualEnd)
		tasksData[i].Startvariance = Findvariance(tasksData[i].PlannedStart, tasksData[i].ActualStartDate)
		if tasksData[i].TaskLevels == "taskL1" {
			firstLeveltasks = append(firstLeveltasks, tasksData[i])
		}

		splittedwbslength := strings.Split(tasksData[i].OrderID, ".")
		switch {
		case len(splittedwbslength) == 3:
			parentwbsForPhase := strings.Join(splittedwbslength[:1], ".")
			parentwbsForWorkpackage := strings.Join(splittedwbslength[:2], ".")
			if pp, exists := mapforsubphase[parentwbsForPhase]; exists {
				tasksData[i].ParentPhase = pp.Title
			}
			if wp, exists := mapforsubphase[parentwbsForWorkpackage]; exists {
				tasksData[i].ParentWorkpackage = wp.Title
			}
		case len(splittedwbslength) == 4:
			parentwbsForPhase := strings.Join(splittedwbslength[:1], ".")
			parentwbsForWorkpackage := strings.Join(splittedwbslength[:2], ".")
			parentwbsForTaskL1 := strings.Join(splittedwbslength[:3], ".")
			if pp, exists := mapforsubphase[parentwbsForPhase]; exists {
				tasksData[i].ParentPhase = pp.Title
			}
			if wp, exists := mapforsubphase[parentwbsForWorkpackage]; exists {
				tasksData[i].ParentWorkpackage = wp.Title
			}
			if pt, exists := mapfortask[parentwbsForTaskL1]; exists {
				tasksData[i].ParenttaskL1 = pt.Title
			}
		case len(splittedwbslength) > 4:
			parentwbsForTaskL1 := strings.Join(splittedwbslength[:3], ".")
			parentwbsForTaskL2 := strings.Join(splittedwbslength[:4], ".")
			parentwbsForPhase := strings.Join(splittedwbslength[:1], ".")
			parentwbsForWorkpackage := strings.Join(splittedwbslength[:2], ".")
			if pt, exists := mapfortask[parentwbsForTaskL1]; exists {
				tasksData[i].ParenttaskL1 = pt.Title
			}
			if pt2, exists := mapfortask[parentwbsForTaskL2]; exists {
				tasksData[i].ParenttaskL2 = pt2.Title
			}
			if pp, exists := mapforsubphase[parentwbsForPhase]; exists {
				tasksData[i].ParentPhase = pp.Title
			}
			if wp, exists := mapforsubphase[parentwbsForWorkpackage]; exists {
				tasksData[i].ParentWorkpackage = wp.Title
			}
		}
	}
	return highestLevel, firstLeveltasks
}

func AssignMetaDataForSubphases(subphaseData []models.Workbook, firstLevelTasks []models.Workbook) {
	mapforsubphase := make(map[string]models.Workbook)
	for i := range subphaseData {
		mapforsubphase[subphaseData[i].OrderID] = subphaseData[i]
	}
	fmt.Println("len of ft", len(firstLevelTasks))
	firstleveltasksMap := make(map[string][]models.Workbook, len(subphaseData))
	for i := range firstLevelTasks {
		firstleveltasksMap[firstLevelTasks[i].SubPhaseID] = append(firstleveltasksMap[firstLevelTasks[i].SubPhaseID], firstLevelTasks[i])
	}
	for i := range subphaseData {
		val := strings.Split(subphaseData[i].OrderID, ".")
		if len(val) == 1 {
			subphaseData[i].TaskLevels = "phase"
		} else {
			subphaseData[i].TaskLevels = "workpackage"
		}
		if subphaseData[i].TaskLevels == "workpackage" {
			// dates roll-up
			if subphase, exists := firstleveltasksMap[subphaseData[i].MongoID]; exists {
				for j := range subphase {
					if subphase[j].PlannedStart != nil {
						if subphaseData[i].PlannedStart == nil || subphase[j].PlannedStart.Before(*subphaseData[i].PlannedStart) {
							subphaseData[i].PlannedStart = subphase[j].PlannedStart
						}
					}
					if subphase[j].PlannedEnd != nil {
						if subphaseData[i].PlannedEnd == nil || subphase[j].PlannedEnd.After(*subphaseData[i].PlannedEnd) {
							subphaseData[i].PlannedEnd = subphase[j].PlannedEnd
						}
					}
					if subphase[j].ActualStartDate != nil {
						if subphaseData[i].ActualStartDate == nil || subphase[j].ActualStartDate.Before(*subphaseData[i].ActualStartDate) {
							subphaseData[i].ActualStartDate = subphase[j].ActualStartDate
						}
					}
					if subphase[j].ActualEndDate != nil {
						if subphaseData[i].ActualEndDate == nil || subphase[j].ActualEndDate.After(*subphaseData[i].ActualEndDate) {
							subphaseData[i].ActualEndDate = subphase[j].ActualEndDate
						}
					}
				}
			}
		}

		// Assigning parent phase & workpackage
		splittedwbslength := strings.Split(subphaseData[i].OrderID, ".")
		switch {
		case len(splittedwbslength) > 1:
			parentwbsForPhase := strings.Join(splittedwbslength[:1], ".")
			parentwbsForWorkpackage := strings.Join(splittedwbslength[:2], ".")

			if pp, exists := mapforsubphase[parentwbsForPhase]; exists {
				subphaseData[i].ParentPhase = pp.Title
			}
			if pw, exists := mapforsubphase[parentwbsForWorkpackage]; exists {
				subphaseData[i].ParentWorkpackage = pw.Title
			}
		}
	}
}

func FormulateWorkbookStructure(tasksData []models.Workbook, ht int) []map[string]interface{} {
	start := time.Now()
	// wdr := make([]map[string]interface{}, len(tasksData))
	var wdr []map[string]interface{}
	for _, task := range tasksData {
		taskres := map[string]interface{}{
			"assignedTo":        task.AssignedTo,
			"taskLevels":        task.TaskLevels,
			"status":            task.Status,
			"role":              task.Role,
			"workstream":        task.Workstream,
			"tasktype":          task.TaskType,
			"wbs":               task.OrderID,
			"plannedStart":      task.PlannedStart,
			"plannnedEnd":       task.PlannedEnd,
			"startvariance":     task.Startvariance,
			"actualStart":       task.ActualStartDate,
			"actualEnd":         task.ActualEndDate,
			"title":             "",
			"workpackage":       "",
			"parentaskL1":       task.ParenttaskL1,
			"parenttaskL2":      task.ParenttaskL2,
			"parentPhase":       task.ParentPhase,
			"parentWorkpackage": task.ParentWorkpackage,
		}
		for i := 1; i <= ht; i++ {
			val := "taskL" + strconv.Itoa(i)
			taskres[val] = ""
		}
		if _, exists := taskres[task.TaskLevels]; exists {
			taskres[task.TaskLevels] = task.Title
		} else {
			if task.TaskLevels == "phase" {
				taskres["title"] = task.Title
			} else {
				taskres["workpackage"] = task.Title
			}
		}

		wdr = append(wdr, taskres)
	}
	fmt.Println("time taken for formulating workbook structure", time.Since(start))
	return wdr
}

func Findvariance(plannedDate *time.Time, actualDate *time.Time) string {
	if plannedDate == nil || actualDate == nil {
		return "Not yet Started"
	}
	days := actualDate.Sub(*plannedDate).Hours() / 24
	if days == 0 {
		return "0 Days ahead"
	} else if days > 0 {
		return strconv.FormatFloat(days, 'f', 0, 64) + " days ahead"
	} else {
		return strconv.FormatFloat(days, 'f', 0, 64) + " days delayed"
	}
}
