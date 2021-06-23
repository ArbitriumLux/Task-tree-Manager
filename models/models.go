package models

import "github.com/jinzhu/gorm"

type Tasks struct {
	gorm.Model
	TaskName  string `gorm:"unique"`
	MiniTasks []MiniTasks
}

type MiniTasks struct {
	gorm.Model
	MiniTaskName string `gorm:"unique"`
	LaborCosts   []LaborCosts
	TasksID      uint
}

type LaborCosts struct {
	gorm.Model
	Cost        int
	MiniTasksID uint
}
