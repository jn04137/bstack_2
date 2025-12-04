package team

import (
	"database/sql"
	"context"

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
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO team (nano_id,team_name,team_desc) VALUES (:nano_id, :team_name, :team_desc)", team)
	if err != nil {
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

func (r TeamRepo) getTeam() {

}

func (r TeamRepo) getDivisions() ([]ESEADivision, error) {
	var divisions []ESEADivision
	err := r.db.Select(&divisions, "SELECT id,division_name from esea_division ORDER BY id ASC")

	return divisions, err
}
