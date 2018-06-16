package tests

import (
	"testing"
	"net/http"
	"bytes"
	"encoding/json"
	"strconv"
	"upsizeAPI/restapi"
	"upsizeAPI/models"
)

func TestGetCompanyID(t *testing.T) {
	FreshDatabase()
	pwd := restapi.HashAndSalt([]byte("111kkk"))
	user := models.User{Email: "josh@test.com", PasswordHash: pwd, Role: "manager"}
	b, _ := json.Marshal(user)
	req, _ := http.NewRequest("PUT", "/user", bytes.NewBuffer(b))
	response := executeRequest(req, "manager")
	checkResponseCode(t, http.StatusCreated, response.Code)

	manager := models.Manager{Name: "Josh", Email: "josh@test.com", Phone: "02040492034", CompanyID: 1}
	b, _ = json.Marshal(manager)
	req, _ = http.NewRequest("PUT", "/manager", bytes.NewBuffer(b))
	response = executeRequest(req, "manager")
	checkResponseCode(t, http.StatusCreated, response.Code)

	companyId := models.GetCompanyIDFromEmail(a.DB, "josh@test.com", "manager")
	if companyId != 1.0 {
		t.Errorf("Weird company id")
	}
}

func TestCreateUser(t *testing.T) {
	FreshDatabase()
	pwd := restapi.HashAndSalt([]byte("111kkk"))
	user := models.User{Email: "josh@test.com", PasswordHash: pwd, Role: "manager"}
	b, _ := json.Marshal(user)
	req, _ := http.NewRequest("PUT", "/user", bytes.NewBuffer(b))
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusCreated, response.Code)
}

func TestGetNonExistentUser(t *testing.T) {
	FreshDatabase()
	payload := []byte(`{"email":"1blahblah@gmail.com"}`)
	req, _ := http.NewRequest("GET", "/user", bytes.NewBuffer(payload))
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

func TestGetUser(t *testing.T) {
	FreshDatabase()
	addUsers(1)
	payload := []byte(`{"email":"1blahblah@gmail.com"}`)
	req, _ := http.NewRequest("GET", "/user", bytes.NewBuffer(payload))
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetUserRole(t *testing.T) {
	FreshDatabase()
	req, _ := http.NewRequest("GET", "/user/role", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addUsers(count int) {
	if count < 1 {
		count = 1
	}
	pwd := restapi.HashAndSalt([]byte("123456"))

	for i := 0; i < count; i++ {
		_, err := a.DB.Exec("INSERT INTO users(email, password_hash, role) VALUES($1, $2, $3)", strconv.Itoa(i+1)+"blahblah@gmail.com", pwd, "manager")
		if err != nil {
			panic(err.Error())
		}
	}

}

func TestUpdateUser(t *testing.T) {
	FreshDatabase()
	addUsers(1)

	payload := []byte(`{"email":"manager@test.com"}`)
	req, _ := http.NewRequest("GET", "/user", bytes.NewBuffer(payload))
	response := executeRequest(req, "manager")
	var originalUser map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalUser)

	payload = []byte(`{"password_hash":"99"}`)

	req, _ = http.NewRequest("POST", "/user", bytes.NewBuffer(payload))
	response = executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["password_hash"] == originalUser["password_hash"] {
		t.Errorf("Expected the password_hash to change from '%v' to '%v'. Got '%v'", originalUser["password_hash"], m["password_hash"], m["password_hash"])
	}
}

func TestValidatePassword(t *testing.T) {
	pwd := restapi.HashAndSalt([]byte("111kkk"))
	plainText := "111kkk"
	if !restapi.ComparePasswords([]byte(pwd), []byte(plainText)) {
		t.Errorf("PWD didnt match!!")
	}
}

func TestDeleteUser(t *testing.T) {
	FreshDatabase()
	addUsers(1)
	payload := []byte(`{"email":"1blahblah@gmail.com"}`)

	req, _ := http.NewRequest("GET", "/user", bytes.NewBuffer(payload))
	response := executeRequest(req, "admin")
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/user", bytes.NewBuffer(payload))
	response = executeRequest(req, "admin")
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/user", bytes.NewBuffer(payload))
	response = executeRequest(req, "manager")
	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}
