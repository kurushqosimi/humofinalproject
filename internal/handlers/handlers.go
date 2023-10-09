package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"main/internal/services"
	"main/pkg/models"
	"net/http"
	"os"
	"strconv"
)

type Handler struct {
	Service *services.Service
	logger  *logrus.Logger
}

func NewHandler(service *services.Service, logger *logrus.Logger) *Handler {
	return &Handler{
		Service: service,
		logger:  logger,
	}
}

func (h *Handler) UserRegistration(w http.ResponseWriter, r *http.Request) {
	var user models.User

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(bytes, &user); err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "", 14)
	pdf.SetX(10)
	pdf.SetY(10)
	pdf.Cell(0, 10, fmt.Sprintf("First Name: %s", user.PersonalData.FirstName))
	pdf.SetX(10)
	pdf.SetY(20)
	pdf.Cell(0, 10, fmt.Sprintf("Second Name: %s", user.PersonalData.SecondName))
	pdf.SetX(10)
	pdf.SetY(30)
	pdf.Cell(0, 10, fmt.Sprintf("Patronymic: %s", user.PersonalData.Patronymic))
	pdf.SetX(10)
	pdf.SetY(40)
	pdf.Cell(0, 10, fmt.Sprintf("Phone: %s", user.PersonalData.Phone))
	pdf.SetX(10)
	pdf.SetY(50)
	pdf.Cell(0, 10, fmt.Sprintf("Address: %s", user.PersonalData.Address))
	pdf.SetX(10)
	pdf.SetY(60)
	pdf.Cell(0, 10, fmt.Sprintf("Company: %s", user.PersonalData.Company))

	err, id := h.Service.UserRegistration(&user)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = pdf.OutputFileAndClose(strconv.Itoa(id) + ".pdf"); err != nil {
		h.logger.Error("Failed to save PDF:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var user models.User

	user.Login = r.Header.Get("login")
	user.Password = r.Header.Get("password")

	id, err := h.Service.SignIn(user)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := h.Service.GetToken(id)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Authentication", token.Token)
}

func (h *Handler) VacancyCreator(w http.ResponseWriter, r *http.Request) {
	var vacancy models.Vacancies

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(bytes, &vacancy); err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	UserID, ok := r.Context().Value("id").(int)
	if !ok {
		h.logger.Error("Couldn't find an ID!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vacancy.UserID = UserID

	if err = h.Service.CreateVacancy(&vacancy); err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *Handler) ChangeVacancy(w http.ResponseWriter, r *http.Request) {
	var (
		id      int
		vacancy models.Vacancies
	)

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(bytes, &vacancy); err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	UserID, ok := r.Context().Value("id").(int)
	if !ok {
		h.logger.Error("Couldn't find an ID!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vacancy.UserID = UserID
	if err = h.Service.ChangeVacancy(&vacancy, id); err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *Handler) DeleteVacancy(w http.ResponseWriter, r *http.Request) {
	var (
		id int
	)

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	UserID, ok := r.Context().Value("id").(int)
	if !ok {
		h.logger.Error("Couldn't find an ID!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = h.Service.DeleteVacancy(&id, &UserID); err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *Handler) ViewVacancy(w http.ResponseWriter, r *http.Request) {
	var (
		id         int
		newVacancy struct {
			ID       int     `json:"id"`
			Title    string  `json:"title"`
			Terms    string  `json:"terms"`
			Duration int     `json:"duration"`
			Fee      float64 `json:"fee"`
		}
	)

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vacancy, err := h.Service.ViewVacancy(&id)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newVacancy.ID = vacancy.ID
	newVacancy.Title = vacancy.Title
	newVacancy.Terms = vacancy.Terms
	newVacancy.Duration = vacancy.Duration
	newVacancy.Fee = vacancy.Fee

	bytes, err := json.Marshal(newVacancy)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func (h *Handler) ViewAllVacancies(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageStr := vars["page"]
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slice, err := h.Service.ViewAllVacancies(&page)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newSlice := make([]struct {
		ID       int
		Title    string
		Terms    string
		Duration int
		Fee      float64
		Category int
	}, len(*slice))

	for i := 0; i < len(*slice); i++ {
		newSlice[i].ID = (*slice)[i].ID
		newSlice[i].Title = (*slice)[i].Title
		newSlice[i].Terms = (*slice)[i].Terms
		newSlice[i].Duration = (*slice)[i].Duration
		newSlice[i].Fee = (*slice)[i].Fee
		newSlice[i].Category = (*slice)[i].CategoryID
	}

	if len(newSlice) == 1 && newSlice[0].ID == 0 {
		return
	}

	bytes, err := json.Marshal(newSlice)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func (h *Handler) SignOut(w http.ResponseWriter, r *http.Request) {
	UserID, ok := r.Context().Value("id").(int)
	if !ok {
		h.logger.Error("Couldn't find an ID!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.Service.SignOut(&UserID); err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Categories(w http.ResponseWriter, r *http.Request) {
	err, categories := h.Service.Categories()
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(categories)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func (h *Handler) Apply(w http.ResponseWriter, r *http.Request) {
	var (
		id int
	)

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	UserID, ok := r.Context().Value("id").(int)
	if !ok {
		h.logger.Error("Couldn't find an ID!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = h.Service.Apply(&id, &UserID); err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *Handler) MyVacancies(w http.ResponseWriter, r *http.Request) {
	UserID, ok := r.Context().Value("id").(int)
	if !ok {
		h.logger.Error("Couldn't find an ID!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vacancies, err := h.Service.MyVacancies(&UserID)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newSlice := make([]struct {
		ID       int
		Title    string
		Terms    string
		Duration int
		Fee      float64
		Category int
	}, len(vacancies))

	for i := 0; i < len(vacancies); i++ {
		newSlice[i].ID = (vacancies)[i].ID
		newSlice[i].Title = (vacancies)[i].Title
		newSlice[i].Terms = (vacancies)[i].Terms
		newSlice[i].Duration = (vacancies)[i].Duration
		newSlice[i].Fee = (vacancies)[i].Fee
		newSlice[i].Category = (vacancies)[i].CategoryID
	}

	if len(newSlice) == 1 && newSlice[0].ID == 0 {
		return
	}

	bytes, err := json.Marshal(newSlice)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func (h *Handler) VacancyApplicants(w http.ResponseWriter, r *http.Request) {
	var (
		id int
	)

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	UserID, ok := r.Context().Value("id").(int)
	if !ok {
		h.logger.Error("Couldn't find an ID!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	applicants, err := h.Service.VacancyApplicants(&id, &UserID)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(applicants)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func (h *Handler) VacancyApplicant(w http.ResponseWriter, r *http.Request) {
	var (
		id int
	)

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	UserID, ok := r.Context().Value("id").(int)
	if !ok {
		h.logger.Error("Couldn't find an ID!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	applicantIDSTR := vars["applicant_id"]
	applicantID, err := strconv.Atoi(applicantIDSTR)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = h.Service.VacancyApplicant(&id, &UserID, &applicantID); err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pdfBytes, err := ioutil.ReadFile(fmt.Sprintf("%d.pdf", applicantID))
	if err != nil {
		http.Error(w, "Failed to read PDF file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=file.pdf")
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Length", strconv.Itoa(len(pdfBytes)))

	_, err = w.Write(pdfBytes)
	if err != nil {
		http.Error(w, "Failed to send PDF file", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ViewByCategory(w http.ResponseWriter, r *http.Request) {
	var (
		id int
	)

	vars := mux.Vars(r)
	idStr := vars["category_id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vacancies, err := h.Service.ViewByCategory(id)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newSlice := make([]struct {
		ID       int
		Title    string
		Terms    string
		Duration int
		Fee      float64
		Category int
	}, len(vacancies))

	for i := 0; i < len(vacancies); i++ {
		newSlice[i].ID = (vacancies)[i].ID
		newSlice[i].Title = (vacancies)[i].Title
		newSlice[i].Terms = (vacancies)[i].Terms
		newSlice[i].Duration = (vacancies)[i].Duration
		newSlice[i].Fee = (vacancies)[i].Fee
	}

	if len(newSlice) == 1 && newSlice[0].ID == 0 {
		return
	}

	bytes, err := json.Marshal(newSlice)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	UserID, ok := r.Context().Value("id").(int)
	if !ok {
		h.logger.Error("Couldn't find an ID!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pdfBytes, err := ioutil.ReadFile(fmt.Sprintf("%d.pdf", UserID))
	if err != nil {
		http.Error(w, "Failed to read PDF file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=file.pdf")
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Length", strconv.Itoa(len(pdfBytes)))

	_, err = w.Write(pdfBytes)
	if err != nil {
		http.Error(w, "Failed to send PDF file", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ChangeProfile(w http.ResponseWriter, r *http.Request) {
	var (
		user models.User
	)

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(bytes, &user.PersonalData); err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	UserID, ok := r.Context().Value("id").(int)
	if !ok {
		h.logger.Error("Couldn't find an ID!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = os.Remove(fmt.Sprintf("%d.pdf", UserID))
	if err != nil {
		h.logger.Error("Ошибка при удалении файла:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.ID = UserID

	if err = h.Service.ChangeProfile(user); err != nil {
		h.logger.Error("Couldn't find an ID!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "", 14)
	pdf.SetX(10)
	pdf.SetY(10)
	pdf.Cell(0, 10, fmt.Sprintf("First Name: %s", user.PersonalData.FirstName))
	pdf.SetX(10)
	pdf.SetY(20)
	pdf.Cell(0, 10, fmt.Sprintf("Second Name: %s", user.PersonalData.SecondName))
	pdf.SetX(10)
	pdf.SetY(30)
	pdf.Cell(0, 10, fmt.Sprintf("Patronymic: %s", user.PersonalData.Patronymic))
	pdf.SetX(10)
	pdf.SetY(40)
	pdf.Cell(0, 10, fmt.Sprintf("Phone: %s", user.PersonalData.Phone))
	pdf.SetX(10)
	pdf.SetY(50)
	pdf.Cell(0, 10, fmt.Sprintf("Address: %s", user.PersonalData.Address))
	pdf.SetX(10)
	pdf.SetY(60)
	pdf.Cell(0, 10, fmt.Sprintf("Company: %s", user.PersonalData.Company))

	if err = pdf.OutputFileAndClose(strconv.Itoa(UserID) + ".pdf"); err != nil {
		h.logger.Error("Failed to save PDF:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	UserID, ok := r.Context().Value("id").(int)
	if !ok {
		h.logger.Error("Couldn't find an ID!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Service.DeleteProfile(UserID)
}

func (h *Handler) SendNotification(w http.ResponseWriter, r *http.Request) {
	var (
		notification models.Notification
	)

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(bytes, &notification); err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	SenderID, ok := r.Context().Value("id").(int)
	if !ok {
		h.logger.Error("Couldn't find an ID!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	vacancyIDStr := vars["vacancy_id"]
	vacancyID, err := strconv.Atoi(vacancyIDStr)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	candidateIDStr := vars["candidate_id"]
	candidateID, err := strconv.Atoi(candidateIDStr)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	notification.OwnerID = candidateID
	notification.SenderID = SenderID
	notification.VacancyID = vacancyID

	if err = h.Service.SendNotification(notification); err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Notifications(w http.ResponseWriter, r *http.Request) {
	UserID, ok := r.Context().Value("id").(int)
	if !ok {
		h.logger.Error("Couldn't find an ID!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	notifications, err := h.Service.NewNotifications(UserID)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(notifications)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}
