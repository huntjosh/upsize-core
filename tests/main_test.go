package tests

import (
	"os"
	"testing"
	"net/http"
	"net/http/httptest"
	"bytes"
	"upsizeAPI/restapi"
)

var a restapi.Api

func TestMain(m *testing.M) {
	restapi.SetupEnv()
	a = restapi.Api{}
	a.Initialize(
		os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"))
	FreshDatabase()
	FillDatabase()
	code := m.Run()

	FreshDatabase()

	os.Exit(code)
}

func FreshDatabase() {
	tables := []string{"skills", "jobs", "contractor_skills", "companies", "company_skills", "contractor_jobs"}
	_, err := a.DB.Exec(`
TRUNCATE skills, jobs, contractor_skills, companies, company_skills, contractor_jobs;
`)

	if err != nil {
		panic(err)
	}

	var queryBuffer bytes.Buffer
	for _, tableName := range tables {
		queryBuffer.WriteString("ALTER SEQUENCE ")
		queryBuffer.WriteString(tableName)
		queryBuffer.WriteString("_id_seq RESTART WITH 1;")
	}

	_, err = a.DB.Exec(queryBuffer.String())

	if err != nil {
		panic(err)
	}
	EmptyAuthTables()
	FillAuthTables()
}

func EmptyAuthTables() {
	_, err := a.DB.Exec(`
TRUNCATE users, managers, contractors;
ALTER SEQUENCE users_id_seq RESTART WITH 1;
ALTER SEQUENCE managers_id_seq RESTART WITH 1;
ALTER SEQUENCE contractors_id_seq RESTART WITH 1;
`)
	if err != nil {
		panic(err)
	}
}

func FillAuthTables() {
	pwd := restapi.HashAndSalt([]byte("123456"))
	_, err := a.DB.Exec("INSERT INTO users(email, password_hash, role) VALUES($1, $2, $3)", "manager@test.com", pwd, "manager")
	if err != nil {
		panic(err.Error())
	}

	_, err = a.DB.Exec("INSERT INTO users(email, password_hash, role) VALUES($1, $2, $3)", "contractor@test.com", pwd, "contractor")
	if err != nil {
		panic(err.Error())
	}

	_, err = a.DB.Exec("INSERT INTO managers(name, email, phone, company_id) VALUES($1, $2, $3, $4)", "bob", "manager@test.com", "02040490234", 1)
	if err != nil {
		panic(err.Error())
	}

	_, err = a.DB.Exec("INSERT INTO contractors(name, charge_rate, email, enabled, phone, company_id, available) VALUES ($1, $2, $3, $4, $5, $6, $7)", "bob", "25.1", "contractor@test.com", true, "1234", 1, true)
	if err != nil {
		panic(err.Error())
	}
}


func FillDatabase() {
	_, err := a.DB.Exec(`


INSERT INTO skills VALUES (1, 'Python'), (2, 'Java');

INSERT INTO skills VALUES (3, '3D Modelling'), (4, 'MySQL'), (5, 'Front End Development'), (6, 'Back End Development'),
(7, 'Web Design'), (8, 'GoLang'), (9, 'PostgreSQL');

INSERT INTO jobs VALUES (1, 'Write python script', '2 Days','2018-01-08 04:05:06', '2018-01-08 04:05:06', 'filling', 'Need a simple script written asap', 1),
(2, 'Java development', '1 Week','2018-01-08 04:05:06', '2018-01-08 04:05:06', 'filling', 'Intermediate level java development for microservice', 2),
(3, 'Front End Development', '1 Hour','2018-01-08 04:05:06', '2018-01-08 04:05:06', 'underway', 'UI/UX Expert required', 2);

INSERT INTO contractor_skills VALUES (1, 1, 3), (2, 1, 2), (3, 1, 1),(4, 2, 3), (5, 2, 4), (6, 2, 5),(7, 3, 3), (8, 3, 5), (9, 3, 6);

INSERT INTO companies VALUES (1, 'Maelstrom Software'), (2, 'Maelstrom 3D Animation');

Insert into contractor_jobs values (1, 1, 'invited', false, 1),(2, 2, 'invited', false, 1), (3, 2, 'invited', false, 1);

insert into company_skills values (1, 1, 1), (2, 1, 2), (3, 1, 3), (4, 1, 4),(5, 1, 5), (6, 1, 6), (7, 1, 7), (8, 1, 8),
(9, 2, 1), (10, 2, 2), (11, 2, 3), (12, 2, 4),(13, 2, 5), (14, 2, 6), (15, 2, 7), (16, 2, 8);
`)

	if err != nil {
		panic(err)
	}
}

func TestEmptyTable(t *testing.T) {
	FreshDatabase()

	req, _ := http.NewRequest("GET", "/companies", nil)
	response := executeRequest(req, "admin")

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request, roleType string) *httptest.ResponseRecorder {
	token, err := restapi.GetToken(roleType+"@test.com", roleType)
	if err != nil {
		panic(err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}