package repositories

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"main/pkg/models"
	"time"
)

type Repository struct {
	DB     *gorm.DB
	logger *logrus.Logger
}

func GetConnection(config *models.Config) (*Repository, error) {
	dbUri := " host = " + config.DbSetting.Host + " port = " + config.DbSetting.Port + " user = " + config.DbSetting.Username + " password = " + config.DbSetting.Password + " dbname = " + config.DbSetting.Database
	fmt.Println(config.DbSetting)
	db, err := gorm.Open(postgres.Open(dbUri), &gorm.Config{}) //installs connection with database
	if err != nil {
		return nil, err
	}
	var rep Repository
	rep.DB = db
	return &rep, nil
}
func (r *Repository) UserRegistration(user *models.User, hash []byte) (error, int) {
	var id int
	sqlQuery := `select login from users where login = ? and active = true`
	tx := r.DB.Exec(sqlQuery, user.Login)
	err := tx.Error
	if err != nil {
		return err, 0
	}
	if tx.RowsAffected != 0 {
		return errors.New("this kind of account already exists"), 0
	}
	sqlQuery = `insert into personal_data (first_name, second_name, patronymic, phone, address, company, gender_id) 
				values (?,?,?,?,?,?,?) returning id `
	err = r.DB.Raw(sqlQuery, user.PersonalData.FirstName, user.PersonalData.SecondName, user.PersonalData.Patronymic, user.PersonalData.Phone, user.PersonalData.Address, user.PersonalData.Company, user.PersonalData.GenderId).Scan(&id).Error
	if err != nil {
		return err, 0
	}
	sqlQuery = `insert into users (login, password, personal_id) values (?, ?, ?) returning personal_id`
	err = r.DB.Raw(sqlQuery, user.Login, hash, id).Scan(&id).Error
	if err != nil {
		return err, 0
	}
	return nil, id
}

func (r *Repository) UserCheck(user *models.User) (*string, *int, error) {
	var newUser struct {
		Id         int       `json:"id"`
		Login      string    `json:"login"`
		Password   string    `json:"password"`
		PersonalId int       `json:"personal_id"`
		CreatedAt  time.Time `json:"created_at"`
		Active     bool      `json:"active"`
		UpdatedAt  time.Time `json:"updated_at"`
		DeletedAt  time.Time `json:"deleted_at"`
	}
	sqlQuery := `select * from users where login = ?  and active = true`
	err := r.DB.Raw(sqlQuery, user.Login).Scan(&newUser).Error
	if err != nil {
		return nil, nil, err
	}
	if newUser.Id == 0 {
		return nil, nil, errors.New("there is not such account")
	}
	return &newUser.Password, &newUser.Id, nil
}
func (r *Repository) AddToken(token *models.Token) error {
	sqlQuery := `insert into tokens (token, user_id) values (?, ?)`
	err := r.DB.Exec(sqlQuery, token.Token, token.UserID).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) TokenCheck(token *string) (int, error) {
	sqlQuery := `select * from tokens where token = ? and expiration_time > current_timestamp and active = true; `
	var check models.Token
	err := r.DB.Raw(sqlQuery, token).Scan(&check).Error
	log.Println(check)
	if err != nil {
		return 0, err
	}
	if check.Token == "" {
		return 0, errors.New("this token does not exists")

	}
	timeChecking := time.Now().After(check.ExpirationTime)
	if timeChecking == true {
		sqlQuery = `delete from tokens where id = ?`
		err := r.DB.Exec(sqlQuery, check.ID).Error
		if err != nil {
			return 0, err
		}
		return 0, errors.New("reenter to your account")
	}
	return check.UserID, nil
}
func (r *Repository) CreateVacancy(vacancy *models.Vacancies) error {
	sqlQuery := `insert into vacancies (title, terms, duration, fee, user_id, category) values (?,?,?,?,?,?)`
	err := r.DB.Exec(sqlQuery, vacancy.Title, vacancy.Terms, vacancy.Duration, vacancy.Fee, vacancy.UserID, vacancy.CategoryID).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) ChangeVacancy(vacancy *models.Vacancies, id int) error {
	sqlQuery := `update vacancies set title = ?, terms = ?, duration = ?, fee = ? where user_id = ? and active = true and id = ?`
	tx := r.DB.Exec(sqlQuery, vacancy.Title, vacancy.Terms, vacancy.Duration, vacancy.Fee, vacancy.UserID, id)
	err := tx.Error
	if err != nil {
		return err
	}
	if tx.RowsAffected != 0 {
		return errors.New("not your vacancy")
	}
	return nil
}
func (r *Repository) DeleteVacancy(id *int, userId *int) error {
	sqlQuery := `update vacancies set active = false where id = ? and user_id = ?;`
	err := r.DB.Exec(sqlQuery, id, userId).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) ViewAllVacancies(page *int) (*[]models.Vacancies, error) {
	slice := make([]models.Vacancies, 1, 5)
	sqlQuery := `select * from vacancies where expiration_time > current_timestamp and active= true limit 5  offset ?;`
	err := r.DB.Raw(sqlQuery, (*page-1)*5).Scan(&slice).Error
	if err != nil {
		return nil, err
	}
	return &slice, nil
}
func (r *Repository) ViewVacancy(Id *int) (*models.Vacancies, error) {
	var vacancy models.Vacancies
	sqlQuery := `select * from vacancies where id =? and active = true`
	tx := r.DB.Raw(sqlQuery, Id).Scan(&vacancy)
	err := tx.Error
	if err != nil {
		return nil, err
	}
	if vacancy.Title == "" {
		return nil, errors.New("does not exist such vacancy")
	}
	return &vacancy, nil
}
func (r *Repository) SignOut(Id *int) error {
	sqlQuery := `delete from tokens where user_id = ? and active = true`
	err := r.DB.Exec(sqlQuery, Id).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) Categories() (error, []models.Category) {
	categories := make([]models.Category, 7, 7)
	sqlQuery := `select * from categories`
	err := r.DB.Raw(sqlQuery).Scan(&categories).Error
	if err != nil {
		return err, nil
	}
	return nil, categories
}
func (r *Repository) Apply(id, userID *int) error {
	var creatorID int
	sqlQuery := `select user_id from vacancies where id = ? and active = true`
	err := r.DB.Raw(sqlQuery, id).Scan(&creatorID).Error
	if err != nil {
		return err
	}
	if creatorID == *userID {
		return errors.New("you cannot apply to your own vacancy")
	}
	sqlQuery = `insert into responses (creator_id, respondent_id, vacancy_id) values (?,?,?)`
	err = r.DB.Exec(sqlQuery, creatorID, userID, id).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) MyVacancies(userID *int) ([]models.Vacancies, error) {
	vacancies := make([]models.Vacancies, 1, 2)
	sqlQuery := `select * from vacancies where user_id = ? and active = true`
	err := r.DB.Raw(sqlQuery, userID).Scan(&vacancies).Error
	if err != nil {
		return nil, err
	}
	return vacancies, err
}
func (r *Repository) VacancyApplicants(vacancyID, userID *int) ([]int, error) {
	applicants := make([]int, 1, 2)
	sqlQuery := `select respondent_id from responses where creator_id = ? and vacancy_id = ?`
	err := r.DB.Raw(sqlQuery, userID, vacancyID).Scan(&applicants).Error
	if err != nil {
		return nil, err
	}
	return applicants, err
}
func (r *Repository) VacancyApplicant(vacancyID, userID, applicantID *int) error {
	sqlQuery := `select * from responses where creator_id=? and respondent_id=? and vacancy_id=?`
	tx := r.DB.Exec(sqlQuery, userID, applicantID, vacancyID)
	err := tx.Error
	if err != nil {
		return err
	}
	if tx.RowsAffected == 0 {
		return errors.New("not such applicant")
	}
	return nil
}
func (r *Repository) ViewByCategory(id int) ([]models.Vacancies, error) {
	vacancies := make([]models.Vacancies, 1, 2)
	sqlQuery := `select * from vacancies where category=? and active = true`
	err := r.DB.Raw(sqlQuery, id).Scan(&vacancies).Error
	if err != nil {
		return nil, err
	}
	return vacancies, err
}
func (r *Repository) ChangeProfile(user models.User) error {
	sqlQuery := `update personal_data set first_name = ?, second_name=? , patronymic=? , phone=? , address=? , company=? , gender_id=? , updated_at=current_timestamp where id = ? and active = true;`
	err := r.DB.Exec(sqlQuery, user.PersonalData.FirstName, user.PersonalData.SecondName, user.PersonalData.Patronymic, user.PersonalData.Phone, user.PersonalData.Address, user.PersonalData.Company, user.PersonalData.GenderId, user.ID).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) DeleteProfile(id int) error {
	sqlQuery := `update personal_data set active = false where id = ? `
	err := r.DB.Exec(sqlQuery, id).Error
	if err != nil {
		return err
	}
	sqlQuery = `update users set active = false where personal_id = ? `
	err = r.DB.Exec(sqlQuery, id).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) SendNotification(notification models.Notification) error {
	sqlQuery := `insert into notifications (comment, owner_id, sender_id, vacancy_id) values (?,?,?,?)`
	err := r.DB.Exec(sqlQuery, notification.Comment, notification.OwnerID, notification.SenderID, notification.VacancyID).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) NewNotification(userID int) ([]models.Notification, error) {
	notifications := make([]models.Notification, 1, 2)
	sqlQuery := `select * from notifications where owner_id = ? and status=true`
	err := r.DB.Raw(sqlQuery, userID).Scan(&notifications).Error
	if err != nil {
		return nil, err
	}
	return notifications, nil
}
