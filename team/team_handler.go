package team

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/matoous/go-nanoid/v2"
)

type TeamHandler struct {
	TeamRepo *TeamRepo
}

func NewHandler(teamRepo *TeamRepo) *TeamHandler {
	return &TeamHandler{
		teamRepo,
	}
}

func (h TeamHandler) GetHandler() *chi.Mux {
	r := chi.NewMux()

	r.Get("/", h.getTeam)
	r.Get("/all", h.getTeams)
	r.Get("/getTeam/{teamNanoId}", h.getTeam)
	r.Get("/eseaDivisions", h.getESEADivisions)
	r.Post("/create", h.createTeam)
	r.Post("/addAchievement", h.addAchievement)

	return r
}

func (h TeamHandler) createTeam(w http.ResponseWriter, r *http.Request) {
	var team Team
	err := json.NewDecoder(r.Body).Decode(&team)
	if err != nil {
		log.Printf("Encountered error trying to decode body: %v", err)
		return
	}
	team.NanoId, err = gonanoid.New()
	if err != nil {
		log.Printf("Encountered problems trying to generate new NanoId: %v", err)
	}

	err = h.TeamRepo.createTeam(&team, r.Context())
	if err != nil {
		log.Printf("Encountered error trying to create team: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h TeamHandler) getTeam(w http.ResponseWriter, r *http.Request) {
	teamNanoId := chi.URLParam(r, "teamNanoId")

	team, err := h.TeamRepo.getTeam(teamNanoId)
	if err != nil {
		log.Printf("Failed to get team: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(team)
	if err != nil {
		log.Printf("Failed to encode team")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h TeamHandler) getTeams(w http.ResponseWriter, r *http.Request) {
	teams, err := h.TeamRepo.getTeams()
	if err != nil {
		log.Printf("failed to retrieve teams: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&teams)
	if err != nil {
		log.Printf("Failed to encode team")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h TeamHandler) getESEADivisions(w http.ResponseWriter, r *http.Request) {
	divisions, err := h.TeamRepo.getDivisions()
	if err != nil {
		log.Printf("Failed to retrieve divisions from db: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&divisions)
	w.WriteHeader(http.StatusOK)
}

func (h TeamHandler) addAchievement(w http.ResponseWriter, r *http.Request) {
	var achievement TeamAchievement
	err := json.NewDecoder(r.Body).Decode(&achievement)
	if err != nil {
		log.Printf("Failed to decode body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.TeamRepo.addTeamAchievement(achievement)
}

