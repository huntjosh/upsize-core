package tests

import (
	"testing"
	"net/http"
	"bytes"
	"encoding/json"
	"upsizeAPI/restapi"
	"github.com/dgrijalva/jwt-go"
)

func TestCreateJWT(t *testing.T) {
	token, err := restapi.GetToken("joshhhunt@gmail.com", "manager")
	if err != nil {
		t.Errorf("Token gen error!!")
	}

	claims, err := restapi.VerifyToken(token)

	if err != nil {
		t.Errorf("Token verify error!!")
	}

	if claims.(jwt.MapClaims)["authEmail"].(string) != "joshhhunt@gmail.com" {
		t.Errorf("Invalid email from jwt!!")
	}

	if claims.(jwt.MapClaims)["authRole"].(string) != "manager" {
		t.Errorf("Invalid role from jwt!!")
	}
}

func TestAuthorizeValid(t *testing.T) {
	EmptyAuthTables()
	addUsers(1)

	payload := []byte(`{"email":"1blahblah@gmail.com","password":"123456"}`)

	req, _ := http.NewRequest("POST", "/authorize", bytes.NewBuffer(payload))
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if response.Header()["Authorization"] == nil {
		t.Errorf("Weird bearer token")
	}
	FillAuthTables()
}

func TestAuthorizeEmailInvalid(t *testing.T) {
	EmptyAuthTables()
	addUsers(1)

	payload := []byte(`{"email":"blahhha@test.com","password":"123456"}`)

	req, _ := http.NewRequest("POST", "/authorize", bytes.NewBuffer(payload))
	response := executeRequest(req, "contractor")

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
	FillAuthTables()
}

func TestAuthorizePwdInvalid(t *testing.T) {
	EmptyAuthTables()
	addUsers(1)

	payload := []byte(`{"email":"1blahblah@gmail.com","password":"111"}`)

	req, _ := http.NewRequest("POST", "/authorize", bytes.NewBuffer(payload))
	response := executeRequest(req, "contractor")

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
	FillAuthTables()
}
