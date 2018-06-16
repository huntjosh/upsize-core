package restapi

import (
	"database/sql"
	"net/http"
	"upsizeAPI/models"
	"time"
	"strconv"
)

func (a *Api) createUser(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("create user", startTime)
	if !IsManagerOrAdmin(r.Header.Get("authRole")) {
		respondWithError(w, http.StatusUnauthorized, "Admin or manager access required")
	}
	var u models.User
	if !validPayload(w, r, &u) {
		return
	}
	defer r.Body.Close()

	if err := u.CreateUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, u)
}

func (a *Api) updateUser(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("update user", startTime)

	var u models.User
	if !validPayload(w, r, &u) {
		return
	}

	role := r.Header.Get("authRole")
	if role != "admin" {
		u.Email = r.Header.Get("authEmail")
	} else if u.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid email payload")
		return
	}

	defer r.Body.Close()
	if u.PasswordHash != "" {
		u.PasswordHash = HashAndSalt([]byte(u.PasswordHash))
	}

	if err := u.UpdateUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (a *Api) deleteUser(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("delete user", startTime)
	if !IsAdmin(r.Header.Get("authRole")) {
		respondWithError(w, http.StatusUnauthorized, "Admin access required")
	}

	var u models.User
	if !validPayload(w, r, &u) {
		return
	}

	defer r.Body.Close()

	if err := u.DeleteUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *Api) getUser(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("get user", startTime)

	u := models.User{Email: r.Header.Get("authEmail")}

	defer r.Body.Close()
	if err := u.GetUserNoPassword(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (a *Api) getUserRole(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("get user", startTime)
	id, _ := strconv.Atoi(r.Header.Get("authId"))

	if r.Header.Get("authRole") == "manager" {
		m := models.Manager{ID: id}
		if err := m.GetManager(a.DB); err != nil {
			switch err {
			case sql.ErrNoRows:
				respondWithError(w, http.StatusNotFound, "Manager not found")
			default:
				respondWithError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}
		respondWithJSON(w, http.StatusOK, m)
	} else {
		c := models.Contractor{ID: id}
		if err := c.GetContractor(a.DB); err != nil {
			switch err {
			case sql.ErrNoRows:
				respondWithError(w, http.StatusNotFound, "Contractor not found")
			default:
				respondWithError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}
		respondWithJSON(w, http.StatusOK, c)
	}

}
