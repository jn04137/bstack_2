package team

import (
	"database/sql"
	"context"
	"log"

	"github.com/jmoiron/sqlx"
)

type TeamRepo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *TeamRepo {
	return &TeamRepo{
		db: db,
	}
}

func (r TeamRepo) createTeam(team *Team, ctx context.Context) error {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	res, err := tx.NamedExecContext(ctx, "INSERT INTO team (nano_id,team_name,team_desc) VALUES (:nano_id, :team_name, :team_desc)", team)
	if err != nil {
		tx.Rollback()
		return err
	}

	teamId, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO team_esea_division (team_id, esea_division) VALUES (?, ?)", teamId, team.ESEADivision.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}

func (r TeamRepo) getTeams() ([]Team, error) {
	teams := []Team{}
	err := r.db.Select(&teams, "SELECT id,nano_id,team_name,team_desc,created_at,updated_at FROM team")
	return teams, err
}

func (r TeamRepo) getTeam(nanoId string) (Team, error) {
	log.Printf("This is the nanoid: %v", nanoId)
	team := Team{}
	err := r.db.Get(&team, "select id,nano_id,team_name,team_desc,created_at,updated_at FROM team WHERE nano_id=?", nanoId)
	return team, err
}

func (r TeamRepo) getDivisions() ([]ESEADivision, error) {
	var divisions []ESEADivision
	err := r.db.Select(&divisions, "SELECT id,division_name from esea_division ORDER BY id DESC")

	return divisions, err
}

func (r TeamRepo) addTeamAchievement(achievement TeamAchievement) error {
	_, err := r.db.Exec("INSERT INTO team_achievements ((SELECT team.id WHERE team.nano_id=?),event,placement,details,date) VALUES (?,?,?,?)", 
		achievement.TeamNanoId, achievement.Event, achievement.Placement, achievement.Details, achievement.Date)
	if err != nil {
		return err
	}

	return err
}

