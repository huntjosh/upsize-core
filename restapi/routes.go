package restapi

import (
	_ "github.com/lib/pq"
	"net/http"
)

func (a *Api) initializeCompanyRoutes() {
	a.Router.Handle("/companies", a.AuthMiddleware(http.HandlerFunc(a.getCompanies))).Methods("GET")
	a.Router.Handle("/company", a.AuthMiddleware(http.HandlerFunc(a.createCompany))).Methods("PUT")
	a.Router.Handle("/company/{id:[0-9]+}", a.AuthMiddleware(http.HandlerFunc(a.getCompany))).Methods("GET")
	a.Router.Handle("/company/{id:[0-9]+}", a.AuthMiddleware(http.HandlerFunc(a.updateCompany))).Methods("POST")
	a.Router.Handle("/company/{id:[0-9]+}", a.AuthMiddleware(http.HandlerFunc(a.deleteCompany))).Methods("DELETE")

	a.Router.Handle("/company/{id:[0-9]+}/skills", a.AuthMiddleware(http.HandlerFunc(a.getCompanySkills))).Methods("GET")
	a.Router.Handle("/company/{id:[0-9]+}/skill", a.AuthMiddleware(http.HandlerFunc(a.createCompanySkill))).Methods("PUT")
	a.Router.Handle("/company/{company_id:[0-9]+}/skill/{skill_id:[0-9]+}", a.AuthMiddleware(http.HandlerFunc(a.getCompanySkill))).Methods("GET")
	a.Router.Handle("/company/{company_id:[0-9]+}/skill/{skill_id:[0-9]+}", a.AuthMiddleware(http.HandlerFunc(a.deleteCompanySkill))).Methods("DELETE")

	a.Router.Handle("/company/{id:[0-9]+}/contractors", a.AuthMiddleware(http.HandlerFunc(a.getCompanyContractors))).Methods("GET")
	a.Router.Handle("/company/{id:[0-9]+}/jobs", a.AuthMiddleware(http.HandlerFunc(a.getCompanyJobs))).Methods("GET")
}

func (a *Api) initializeContractorRoutes() {
	a.Router.Handle("/contractors", a.AuthMiddleware(http.HandlerFunc(a.getContractors))).Methods("GET")
	a.Router.Handle("/contractor", a.AuthMiddleware(http.HandlerFunc(a.createContractor))).Methods("PUT")
	a.Router.Handle("/contractor/{id:[0-9]+}", a.AuthMiddleware(http.HandlerFunc(a.getContractor))).Methods("GET")
	a.Router.Handle("/contractor/{id:[0-9]+}", a.AuthMiddleware(http.HandlerFunc(a.updateContractor))).Methods("POST")
	a.Router.Handle("/contractor/{id:[0-9]+}", a.AuthMiddleware(http.HandlerFunc(a.deleteContractor))).Methods("DELETE")
	a.Router.Handle("/contractor/{id:[0-9]+}/company", a.AuthMiddleware(http.HandlerFunc(a.getContractorCompany))).Methods("GET")

	a.Router.Handle("/contractor/{contractor_id:[0-9]+}/jobs/unseenCounts", a.AuthMiddleware(http.HandlerFunc(a.getContractorJobUnseenCounts))).Methods("GET")
	a.Router.Handle("/contractor/{contractor_id:[0-9]+}/jobs", a.AuthMiddleware(http.HandlerFunc(a.getContractorJobs))).Methods("GET")
	a.Router.Handle("/contractor/{contractor_id:[0-9]+}/job", a.AuthMiddleware(http.HandlerFunc(a.createContractorJob))).Methods("PUT")
	a.Router.Handle("/contractor/{contractor_id:[0-9]+}/job/{job_id:[0-9]+}", a.AuthMiddleware(http.HandlerFunc(a.getContractorJob))).Methods("GET")
	a.Router.Handle("/contractor/{contractor_id:[0-9]+}/job/{job_id:[0-9]+}", a.AuthMiddleware(http.HandlerFunc(a.updateContractorJob))).Methods("POST")
	a.Router.Handle("/contractor/{contractor_id:[0-9]+}/job/{job_id:[0-9]+}", a.AuthMiddleware(http.HandlerFunc(a.deleteContractorJob))).Methods("DELETE")
}

func (a *Api) initializeManagerRoutes() {
	a.Router.Handle("/managers", a.AuthMiddleware(http.HandlerFunc(a.getManagers))).Methods("GET")
	a.Router.Handle("/manager", a.AuthMiddleware(http.HandlerFunc(a.createManager))).Methods("PUT")
	a.Router.Handle("/manager/{id:[0-9]+}", a.AuthMiddleware(http.HandlerFunc(a.getManager))).Methods("GET")
	a.Router.Handle("/manager/{id:[0-9]+}", a.AuthMiddleware(http.HandlerFunc(a.updateManager))).Methods("POST")
	a.Router.Handle("/manager/{id:[0-9]+}", a.AuthMiddleware(http.HandlerFunc(a.deleteManager))).Methods("DELETE")
	a.Router.Handle("/manager/{id:[0-9]+}/company", a.AuthMiddleware(http.HandlerFunc(a.getManagerCompany))).Methods("GET")
	a.Router.Handle("/manager/{id:[0-9]+}/jobs", a.AuthMiddleware(http.HandlerFunc(a.getManagerJobs))).Methods("GET")
}

func (a *Api) initializeJobRoutes() {
	a.Router.Handle("/jobs", a.AuthMiddleware(http.HandlerFunc(a.getJobs))).Methods("GET")
	a.Router.Handle("/job", a.AuthMiddleware(http.HandlerFunc(a.createJob))).Methods("PUT")
	a.Router.Handle("/job/{id:[0-9]+}", a.AuthMiddleware(http.HandlerFunc(a.getJob))).Methods("GET")
	a.Router.Handle("/job/{id:[0-9]+}", a.AuthMiddleware(http.HandlerFunc(a.updateJob))).Methods("POST")
	a.Router.Handle("/job/{id:[0-9]+}", a.AuthMiddleware(http.HandlerFunc(a.deleteJob))).Methods("DELETE")
	a.Router.Handle("/job/{id:[0-9]+}/contractors", a.AuthMiddleware(http.HandlerFunc(a.getJobContractors))).Methods("GET")
}

func (a *Api) initializeSkillRoutes() {
	a.Router.Handle("/skills", a.AuthMiddleware(http.HandlerFunc(a.getSkills))).Methods("GET")
	a.Router.Handle("/skill", a.AuthMiddleware(http.HandlerFunc(a.createSkill))).Methods("PUT")
	a.Router.Handle("/skill/{id:[0-9]+}", a.AuthMiddleware(http.HandlerFunc(a.getSkill))).Methods("GET")
	a.Router.Handle("/skill/{id:[0-9]+}", a.AuthMiddleware(http.HandlerFunc(a.updateSkill))).Methods("POST")
	a.Router.Handle("/skill/{id:[0-9]+}", a.AuthMiddleware(http.HandlerFunc(a.deleteSkill))).Methods("DELETE")
}

func (a *Api) initializeUserRoutes() {
	a.Router.Handle("/user", a.AuthMiddleware(http.HandlerFunc(a.createUser))).Methods("PUT")
	a.Router.Handle("/user", a.AuthMiddleware(http.HandlerFunc(a.getUser))).Methods("GET")
	a.Router.Handle("/user/role", a.AuthMiddleware(http.HandlerFunc(a.getUserRole))).Methods("GET")
	a.Router.Handle("/user", a.AuthMiddleware(http.HandlerFunc(a.updateUser))).Methods("POST")
	a.Router.Handle("/user", a.AuthMiddleware(http.HandlerFunc(a.deleteUser))).Methods("DELETE")
}

func (a *Api) initializeAuthRoutes() {
	a.Router.HandleFunc("/authorize", a.Authenticate).Methods("POST")
}
