package services

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"main/internal/repositories"
	"main/pkg/models"
	"time"
)

type Service struct {
	Repository *repositories.Repository
	logger     *logrus.Logger
}

func NewService(repository *repositories.Repository, logger *logrus.Logger) *Service {
	return &Service{
		Repository: repository,
		logger:     logger,
	}
}
func (s *Service) UserRegistration(user *models.User) (error, int) {
	err := s.Validation(&user.Password)
	if err != nil {
		return err, 0
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err, 0
	}
	err, id := s.Repository.UserRegistration(user, hash)
	if err != nil {
		return err, 0
	}
	return nil, id
}
func (s *Service) SignIn(user models.User) (int, error) {
	password, id, err := s.Repository.UserCheck(&user)
	if err != nil {
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(*password), []byte(user.Password))
	if err != nil {
		return 0, err
	}
	return *id, nil
}
func (s *Service) GetToken(userID int) (*models.Token, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	calms := token.Claims.(jwt.MapClaims)
	calms["userId"] = userID
	calms["time"] = time.Now()
	secretKey := []byte("secret")
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return nil, err
	}
	var NewToken models.Token
	NewToken.Token = tokenString
	NewToken.UserID = userID
	err = s.Repository.AddToken(&NewToken)
	if err != nil {
		return nil, err
	}
	return &NewToken, nil
}

func (s *Service) Validation(password *string) error {
	runes := []rune(*password)
	if len(runes) < 8 {
		return errors.New("your password is too short")
	}
	return nil
}
func (s *Service) TokenCheck(token *string) (int, error) {
	return s.Repository.TokenCheck(token)
}
func (s *Service) CreateVacancy(vacancy *models.Vacancies) error {
	return s.Repository.CreateVacancy(vacancy)
}
func (s *Service) ChangeVacancy(vacancy *models.Vacancies, id int) error {
	return s.Repository.ChangeVacancy(vacancy, id)
}
func (s *Service) DeleteVacancy(id *int, userId *int) error {
	return s.Repository.DeleteVacancy(id, userId)
}
func (s *Service) ViewAllVacancies(page *int) (*[]models.Vacancies, error) {
	return s.Repository.ViewAllVacancies(page)
}
func (s *Service) ViewVacancy(id *int) (*models.Vacancies, error) {
	return s.Repository.ViewVacancy(id)
}
func (s *Service) SignOut(id *int) error {
	return s.Repository.SignOut(id)
}

func (s *Service) Categories() (error, []models.Category) {
	return s.Repository.Categories()
}
func (s *Service) Apply(id, userID *int) error {
	return s.Repository.Apply(id, userID)
}
func (s *Service) MyVacancies(userID *int) ([]models.Vacancies, error) {
	return s.Repository.MyVacancies(userID)
}
func (s *Service) VacancyApplicants(vacancyID, userID *int) ([]int, error) {
	return s.Repository.VacancyApplicants(vacancyID, userID)
}
func (s *Service) VacancyApplicant(vacancyID, userID, applicantID *int) error {
	return s.Repository.VacancyApplicant(vacancyID, userID, applicantID)
}
func (s *Service) ViewByCategory(id int) ([]models.Vacancies, error) {
	return s.Repository.ViewByCategory(id)
}
func (s *Service) ChangeProfile(user models.User) error {
	return s.Repository.ChangeProfile(user)
}
func (s *Service) DeleteProfile(id int) error {
	return s.Repository.DeleteProfile(id)
}
func (s *Service) SendNotification(notification models.Notification) error {
	return s.Repository.SendNotification(notification)
}
func (s *Service) NewNotifications(userID int) ([]models.Notification, error) {
	return s.Repository.NewNotification(userID)
}
