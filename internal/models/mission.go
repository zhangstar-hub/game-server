package models

import "my_app/internal/db"

type MissionModel struct {
	ID   uint   `gorm:"primary_key"`
	Data string `gorm:"type:json"`
}

func GetMission(id uint) *MissionModel {
	mission := &MissionModel{ID: id}
	db.DB.First(mission)
	return mission
}

func CreateMission(id uint, data string) *MissionModel {
	mission := &MissionModel{ID: id, Data: data}
	tx := db.DB.Create(mission)
	if tx.RowsAffected == 0 {
		return nil
	}
	return mission
}
