package datacontroller

import (
	dbmanager "hanyoung/logi-tracker/internal/database"
)

type DataController struct {
	db     *dbmanager.DataBaseManager
	admins map[string]bool
}

func GetInstance() *DataController {
	db := dbmanager.GetInstance()
	admins := make(map[string]bool)
	return &DataController{db, admins}
}
