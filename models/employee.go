package models

import "time"

type Employee struct {
	Id         uint    `gorm:"primaryKey"`
	EmpName    string  `json:"empName"`
	EmpAddress string  `json:"empAddress"`
	EmpTel     string  `json:"empTel"`
	EmpSalary  float32 `json:"empSalary"`
	EmpPhoto   string  `json:"empPhoto"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
