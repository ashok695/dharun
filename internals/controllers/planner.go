package controllers

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/dharun/poc/database"
	"github.com/dharun/poc/internals/helpers"
	"github.com/dharun/poc/internals/models"
	"github.com/facette/natsort"
	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2"
)

var (
	newStatusID    = "5c552b391c9d440000904724"
	ActiveStatusID = "5c552b551c9d440000904725"
)

func GetPlannerDataV1(c *fiber.Ctx) ([]*models.TaskData, error) {
	var (
		tasksData        []models.TaskData
		userData         []models.UserData
		statusData       []models.StatusData
		rwtData          []models.TypesData
		DependenciesData []models.DependenciesData

		userMap   map[string]models.UserData
		statusMap map[string]models.StatusData
		typesMap  models.RWTMap
	)
	var wg sync.WaitGroup
	start := time.Now()
	errChan := make(chan error, 5)
	defer close(errChan)
	project := c.Query("projectID")
	wg.Add(5)
	go func() {
		start := time.Now()
		defer wg.Done()
		user, err := helpers.GetUserData(project)
		if err != nil {
			fmt.Println("Err in user chan")
			errChan <- err
			return
		}
		userData = user
		userMap = helpers.Usermap(userData)
		fmt.Println("time taken for user", time.Since(start))
	}()
	go func() {
		start := time.Now()
		defer wg.Done()
		status, err := helpers.GetStatusData(project)
		if err != nil {
			fmt.Println("Err in user chan")
			errChan <- err
			return
		}
		statusData = status
		statusMap = helpers.StatusMap(statusData)
		fmt.Println("time taken for status", time.Since(start))
	}()
	go func() {
		start := time.Now()
		defer wg.Done()
		types, err := helpers.GetTypesData(project)
		if err != nil {
			fmt.Println("Err in user chan")
			errChan <- err
			return
		}
		rwtData = types
		typesMap = helpers.TypesMap(rwtData)
		fmt.Println("time taken for types", time.Since(start))
	}()
	go func() {
		start := time.Now()
		defer wg.Done()
		tasks, err := GetPlannertasks(project)
		if err != nil {
			fmt.Println("err in task")
			errChan <- err
		}
		tasksData = tasks
		fmt.Println("time taken for tasklists", time.Since(start))
	}()
	go func() {
		start := time.Now()
		defer wg.Done()
		depData, err := GetDependenciesData(project)
		if err != nil {
			fmt.Println("err in task")
			errChan <- err
		}
		DependenciesData = depData
		fmt.Println("time taken for tasklists", time.Since(start))
	}()
	wg.Wait()
	select {
	case err := <-errChan:
		return nil, err
	default:
	}
	startb := time.Now()
	sort.Slice(tasksData, func(i, j int) bool {
		return natsort.Compare(tasksData[i].OrderID, tasksData[j].OrderID)
	})
	fmt.Println("time taken for sorting", time.Since(startb))
	startc := time.Now()
	FormulatingData(tasksData, userMap, statusMap, typesMap)
	fmt.Println("time taken for assign", time.Since(startc))
	startD := time.Now()
	a := FormulateBryntumStructure(tasksData)
	fmt.Println("time taken for formulating brytntum", time.Since(startD))
	fmt.Println("time taken for return ", time.Since(start))
	return a, nil
}

func GetPlannertasks(projectID string) ([]models.TaskData, error) {
	start := time.Now()
	var tasksData []models.TaskData
	var wg sync.WaitGroup
	errChan := make(chan error, 2)
	wg.Add(2)
	go func() {
		defer wg.Done()
		var id gocql.UUID
		var activePercentage, duration float64
		var project, actualend, actualstart, assignedto, constraintdate, constrainttype, delaylog, mongoid, orderid, phase,
			plannedfrom, plannedto, reftaskid, roletype, status, subphaseid, tasktype, title, types, workstream string

		iter := database.Session.Query(`SELECT project,id, activepercentage,actualend,actualstart, assignedto, constraintdate,constrainttype,delaylog,duration,mongoid, orderid, phase,
		plannedfrom, plannedto,reftaskid, roletype, status, subphaseid, tasktype, title, type, workstream
		FROM tasklists WHERE PROJECT =?`, projectID).Iter()

		// Scan the result into respective variables
		for iter.Scan(&project, &id, &activePercentage, &actualend, &actualstart, &assignedto, &constraintdate, &constrainttype, &delaylog, &duration, &mongoid, &orderid, &phase, &plannedfrom,
			&plannedto, &reftaskid, &roletype, &status, &subphaseid, &tasktype, &title, &types, &workstream) {

			// Create TaskData struct and append it to tasksData
			task := models.TaskData{
				Project:          project,
				ID:               id,
				ActivePercentage: activePercentage,
				ActualEnd:        actualend,
				ActualStart:      actualstart,
				AssignedToID:     assignedto,
				Constraintdate:   constraintdate,
				Constrainttype:   constrainttype,
				DelayLog:         delaylog,
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
				Type:             types,
				WorkstreamID:     workstream,
			}
			tasksData = append(tasksData, task)
		}

		err := iter.Close()
		if err != nil {
			fmt.Println("ERROR IN DATA TASKLISTS", err.Error())
			errChan <- err
		}
	}()
	go func() {
		defer wg.Done()
		var id gocql.UUID
		var duration float64
		var project, constraintdate, constrainttype, delaylog, mongoid, orderid, phase,
			role, title, tasktype, workstream string

		iter := database.Session.Query(`SELECT project,id,constraintdate,constrainttype,delaylog,duration,mongoid,orderid,phase,
	        role,title,type,workstream
	        from subphases WHERE project = ?`, projectID).Iter()

		for iter.Scan(&project, &id, &constraintdate, &constrainttype, &delaylog, &duration, &mongoid, &orderid, &phase, &role, &title, &tasktype, &workstream) {
			task := models.TaskData{
				Project:          project,
				ID:               id,
				ActivePercentage: 0,
				ActualEnd:        "",
				ActualStart:      "",
				AssignedToID:     "",
				Constraintdate:   constraintdate,
				Constrainttype:   constrainttype,
				DelayLog:         delaylog,
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
			tasksData = append(tasksData, task)
		}
		err := iter.Close()
		if err != nil {
			fmt.Println("ERROR IN DATA ASHOK", err.Error())
			errChan <- err
		}
	}()
	wg.Wait()
	select {
	case err := <-errChan:
		return nil, err
	default:
	}
	fmt.Println("planner time", time.Since(start))
	return tasksData, nil
}

func GetDependenciesData(projectID string) ([]models.DependenciesData, error) {
	var dependenciesData []models.DependenciesData
	var dbid gocql.UUID
	var project, gantttype, id, lagunit, mongoid, source, target string
	var lag, types int16

	iter := database.Session.Query(`SELECT project,dbid,gantttype,id,lag,
	lagunit,mongoid,source,target,type`).Iter()

	for iter.Scan(&project, &dbid, &gantttype, &id, &lag, &lagunit, &mongoid, &source, &target, &types) {
		depends := models.DependenciesData{
			Project:   project,
			DBID:      dbid,
			GanttType: gantttype,
			ID:        id,
			Lag:       int32(lag),
			LagUnit:   lagunit,
			MongoID:   mongoid,
			Source:    source,
			Target:    target,
			Type:      int(types),
		}
		dependenciesData = append(dependenciesData, depends)
	}
	return dependenciesData, nil
}

func FormulateBryntumStructure(tasksData []models.TaskData) []*models.TaskData {
	tasksMap := make(map[string]*models.TaskData)

	for i := range tasksData {
		tasksMap[tasksData[i].MongoID] = &tasksData[i]
	}

	var root []*models.TaskData

	for i := range tasksData {
		currentTask := tasksMap[tasksData[i].MongoID]
		if currentTask.ParentID != "" {
			if parent, exits := tasksMap[currentTask.ParentID]; exits {
				parent.Children = append(parent.Children, currentTask)
			} else {
				root = append(root, currentTask)
			}
		} else {
			root = append(root, currentTask)
		}
	}
	return root
}

func FormulatingData(tasksData []models.TaskData, userMap map[string]models.UserData, statusMap map[string]models.StatusData, typesMap models.RWTMap) {
	for i := range tasksData {
		if tasksData[i].AssignedToID != "" {
			if user, exists := userMap[tasksData[i].AssignedToID]; exists {
				tasksData[i].AssignedTo = []models.UserData{user}
			}
		}
		if tasksData[i].StatusID != "" {
			if status, exists := statusMap[tasksData[i].StatusID]; exists {
				tasksData[i].Status = []models.StatusData{status}
			}
		}
		if tasksData[i].RoleID != "" {
			if role, exists := typesMap.RoleMap[tasksData[i].RoleID]; exists {
				tasksData[i].Role = []models.TypesData{role}
			}
		}
		if tasksData[i].WorkstreamID != "" {
			if workstream, exists := typesMap.WorkstreamMap[tasksData[i].WorkstreamID]; exists {
				tasksData[i].Workstream = []models.TypesData{workstream}
			}
		}
		if tasksData[i].TaskTypeID != "" {
			if tasktype, exists := typesMap.TasktypeMap[tasksData[i].TaskTypeID]; exists {
				tasksData[i].TaskType = []models.TypesData{tasktype}
			}
		}
		if tasksData[i].RefTaskID != "" {
			tasksData[i].ParentID = tasksData[i].RefTaskID
		} else if tasksData[i].RefTaskID == "" && tasksData[i].SubPhaseID != "" {
			tasksData[i].ParentID = tasksData[i].SubPhaseID
		} else if tasksData[i].RefTaskID == "" && tasksData[i].SubPhaseID == "" {
			tasksData[i].ParentID = tasksData[i].Phase
		} else {
			tasksData[i].ParentID = ""
		}
		tasksData[i].PlannedStart = helpers.ParseDate(tasksData[i].PlannedFrom)
		tasksData[i].PlannedEnd = helpers.ParseDate(tasksData[i].PlannedTo)
		tasksData[i].ActualStartDate = helpers.ParseDate(tasksData[i].ActualStart)
		tasksData[i].ActualEndDate = helpers.ParseDate(tasksData[i].ActualEnd)
		if tasksData[i].PlannedEnd != nil {
			today := time.Now().UTC()
			if (*tasksData[i].PlannedEnd).Before(today) && (tasksData[i].StatusID == newStatusID || tasksData[i].StatusID == ActiveStatusID) {
				tasksData[i].OverDueDays = int16(today.Sub(*tasksData[i].PlannedEnd).Hours() / 24)
				if tasksData[i].OverDueDays > 0 {
					tasksData[i].IsOverdue = true
				}
			}
		}
		if tasksData[i].PlannedEnd != nil && tasksData[i].ActualEndDate != nil {
			// if tasksData[i].PlannedTo > tasksData[i].ActualEnd {
			tasksData[i].Variance = int8(tasksData[i].PlannedEnd.Sub(*tasksData[i].ActualEndDate))
			// }else{
			// 	tasksData[i].Variance = *tasksData[i].ActualEndDate.Sub(tasksData[i].PlannedEnd)
			// }
		}
	}
}
