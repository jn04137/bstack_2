package team

import (
	"database/sql"
	_ "encoding/json"
)

type Team struct {
	Id int `json:"id" db:"id"`
	NanoId string `json:"nanoId" db:"nano_id"`
	TeamName string `json:"teamName" db:"team_name"`
	TeamDesc string `json:"teamDesc,omitempty" db:"team_desc"`
	ESEADivision *ESEADivision `json:"eseaDivision,omitempty"`
	CreatedAt sql.NullTime `json:"createdAt" db:"created_at"`
	UpdatedAt sql.NullTime `json:"updatedAt" db:"updated_at"`
}

type ESEADivision struct {
	Id int `json:"id" db:"id"`
	DivisionName string `json:"divisionName" db:"division_name"`
}
