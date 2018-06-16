package restapi

import (
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"time"
	"strconv"
	"github.com/rs/cors"
)

type Api struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *Api) Run(addr string) {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})

	handler := c.Handler(a.Router)
	log.Fatal(http.ListenAndServe(":8000", handler))
}

func (a *Api) Initialize(user, password, dbName string) {
	fmt.Println("Booting up UpsizeCore")
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbName)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
	fmt.Println("UpsizeCore is online")
}

func (a *Api) initializeRoutes() {
	a.initializeCompanyRoutes()
	a.initializeContractorRoutes()
	a.initializeJobRoutes()
	a.initializeManagerRoutes()
	a.initializeSkillRoutes()
	a.initializeUserRoutes()
	a.initializeAuthRoutes()
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func logFinished(resourceName string, startTime time.Time) {
	executionTime := strconv.FormatInt((time.Now().UnixNano()-startTime.UnixNano())/int64(time.Millisecond), 10)
	fmt.Println("Executed " + resourceName + " request in " + executionTime + "ms")
}

func validPayload(w http.ResponseWriter, r *http.Request, model interface{}) bool {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&model); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return false
	}

	return true
}