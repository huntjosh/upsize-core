package restapi

import (
	"database/sql"
	"net/http"
	"strconv"
	"upsizeAPI/models"
	"github.com/gorilla/mux"
	"time"
)

func (a *Api) createSkill(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("create skill", startTime)
	if !IsManagerOrAdmin(r.Header.Get("authRole")) {
		respondWithError(w, http.StatusUnauthorized, "Admin or manager access required")
	}
	var s models.Skill
	if !validPayload(w, r, &s) {
		return
	}
	defer r.Body.Close()

	if err := s.CreateSkill(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, s)
}

func (a *Api) updateSkill(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("update Skill", startTime)
	if !IsAdmin(r.Header.Get("authRole")) {
		respondWithError(w, http.StatusUnauthorized, "Admin access required")
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid skill ID")
		return
	}

	var s models.Skill
	if !validPayload(w, r, &s) {
		return
	}
	defer r.Body.Close()
	s.ID = id

	if err := s.UpdateSkill(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, s)
}

func (a *Api) deleteSkill(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("delete Skill", startTime)
	if !IsAdmin(r.Header.Get("authRole")) {
		respondWithError(w, http.StatusUnauthorized, "Admin access required")
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Skill ID")
		return
	}

	s := models.Skill{ID: id}
	if err := s.DeleteSkill(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *Api) getSkill(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("get Skill", startTime)
	if !IsManagerOrAdmin(r.Header.Get("authRole")) {
		respondWithError(w, http.StatusUnauthorized, "Admin or manager access required")
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid skill ID")
		return
	}

	s := models.Skill{ID: id}
	if err := s.GetSkill(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Skill not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, s)
}

func (a *Api) getSkills(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("get skills", startTime)
	if !IsManagerOrAdmin(r.Header.Get("authRole")) {
		respondWithError(w, http.StatusUnauthorized, "Admin or manager access required")
	}
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	skills, err := models.GetSkills(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, skills)
}
