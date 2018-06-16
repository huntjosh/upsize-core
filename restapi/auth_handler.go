package restapi

import (
	"net/http"
	"upsizeAPI/models"
	"database/sql"
	"encoding/json"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"strconv"
)

func (a *Api) Authenticate(w http.ResponseWriter, r *http.Request) {
	var m map[string]string
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&m)
	password := m["password"]
	email := m["email"]
	if len(email) == 0 || len(password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please provide name and password to obtain the token"))
		return
	}

	u := models.User{Email: m["email"]}
	if err := u.GetUser(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "user not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if ComparePasswords([]byte(u.PasswordHash), []byte(password)) {
		// Authenticated, woohoo
		token, err := GetToken(u.Email, u.Role)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error generating JWT token: " + err.Error()))
		} else {
			w.Header().Set("Authorization", "Bearer "+token)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Token: " + token))
		}
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Name and password do not match"))
	return
}

func (a *Api) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing Authorization Header"))
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		claims, err := VerifyToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Error verifying JWT token: " + err.Error()))
			return
		}
		email := claims.(jwt.MapClaims)["authEmail"].(string)
		role := claims.(jwt.MapClaims)["authRole"].(string)
		r.Header.Set("authEmail", email)
		r.Header.Set("authRole", role)

		if role != "admin" {
			a.setUserAuthHeaders(r, role, email)
		}
		next.ServeHTTP(w, r)
	})
}

func (a *Api) setUserAuthHeaders(r *http.Request, role, email string) {
	companyId := strconv.Itoa(models.GetCompanyIDFromEmail(a.DB, email, role))
	r.Header.Set("authCompanyID", companyId)
	r.Header.Set("authId", strconv.Itoa(models.GetId(a.DB, email, role)))
}

func IsAdmin(authRole string) bool                      { return authRole == "admin" }
func IsManagerOrAdmin(authRole string) bool             { return authRole == "admin" || authRole == "manager" }