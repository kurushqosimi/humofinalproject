package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(h *Handler) *mux.Router {
	router := mux.NewRouter()
	auth := router.PathPrefix("/registration").Subrouter()
	auth.HandleFunc("", h.UserRegistration).Methods(http.MethodPost)
	router.HandleFunc("/sign_in", h.SignIn).Methods(http.MethodGet) //is it acceptable?
	jobs := router.PathPrefix("/jobs").Subrouter()
	jobs.Use(h.CheckUser)
	jobs.HandleFunc("/create_vacancy", h.VacancyCreator).Methods(http.MethodPost)
	jobs.HandleFunc("/change_vacancy/{id}", h.ChangeVacancy).Methods(http.MethodPatch)
	jobs.HandleFunc("/delete_vacancy/{id}", h.DeleteVacancy).Methods(http.MethodDelete)
	jobs.HandleFunc("/view_vacancy/{id}", h.ViewVacancy).Methods(http.MethodGet)
	jobs.HandleFunc("/view_all_vacancies/{page}", h.ViewAllVacancies).Methods(http.MethodGet)
	jobs.HandleFunc("/sign_out", h.SignOut).Methods(http.MethodGet)
	jobs.HandleFunc("/categories", h.Categories).Methods(http.MethodGet)
	jobs.HandleFunc("/apply/{id}", h.Apply).Methods(http.MethodPost)
	jobs.HandleFunc("/my_vacancies", h.MyVacancies).Methods(http.MethodGet)
	jobs.HandleFunc("/my_vacancies/{id}/applicants", h.VacancyApplicants).Methods(http.MethodGet)
	jobs.HandleFunc("/my_vacancies/{id}/applicants/{applicant_id}", h.VacancyApplicant).Methods(http.MethodGet)
	jobs.HandleFunc("/view_vacancies_by_category/{category_id}", h.ViewByCategory).Methods(http.MethodGet)
	jobs.HandleFunc("/profile", h.Profile).Methods(http.MethodGet)
	jobs.HandleFunc("/profile/change", h.ChangeProfile).Methods(http.MethodPatch)
	jobs.HandleFunc("/profile/delete", h.DeleteProfile).Methods(http.MethodDelete)
	jobs.HandleFunc("/notification/{vacancy_id}/{candidate_id}", h.SendNotification).Methods(http.MethodPost)
	jobs.HandleFunc("/my_notifications", h.Notifications).Methods(http.MethodGet)
	return router
}
