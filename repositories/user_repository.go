package repositories

import (
	"errors"

	"github.com/kunto/golang-rest-api-berita/dto"
	"github.com/kunto/golang-rest-api-berita/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetList(search string, limit int, offset int) ([]dto.UserDTO, int64, error)
	CheckEmailExists(email string) (bool, error)
	CheckUsernameExists(username string) (bool, error)
	CreateUser(user models.User) (string, string, error)
	GetByID(id string) (dto.UserDTO, error)
	UpdateUser(id string, data map[string]interface{}) error
	DeleteUser(id string) error
	GetModelByID(id string) (models.User, error)
	GetByEmail(email string) (models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetList(search string, limit int, offset int) ([]dto.UserDTO, int64, error) {
	where := "WHERE 1=1"
	params := []interface{}{}

	if search != "" {
		like := "%" + search + "%"
		where += " AND (LOWER(name) LIKE LOWER(?) OR LOWER(email) LIKE LOWER(?) OR LOWER(username) LIKE LOWER(?))"
		params = append(params, like, like, like)
	}

	// COUNT
	var total int64
	totalQuery := "SELECT COUNT(*) FROM users " + where
	if err := r.db.Raw(totalQuery, params...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	// DATA RESULT
	var users []dto.UserDTO
	dataQuery := `
		SELECT id_user, name, email, username, alamat, created_at, status
		FROM users ` + where + `
		ORDER BY name ASC
		LIMIT ? OFFSET ?
	`
	params = append(params, limit, offset)

	if err := r.db.Raw(dataQuery, params...).Scan(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) CheckEmailExists(email string) (bool, error) {
	var exists bool
	err := r.db.Raw("SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(email)=LOWER(?))", email).
		Scan(&exists).Error
	return exists, err
}

func (r *userRepository) CheckUsernameExists(username string) (bool, error) {
	var exists bool
	err := r.db.Raw("SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(username)=LOWER(?))", username).
		Scan(&exists).Error
	return exists, err
}

func (r *userRepository) CreateUser(user models.User) (string, string, error) {
	var lastID string
	var createdAt string

	query := `
		INSERT INTO users (name,email,username,password,alamat,created_at,status,id_role)
		VALUES (?,?,?,?,?,NOW(),1,?)
		RETURNING id_user, created_at;
	`

	err := r.db.Raw(query,
		user.Name,
		user.Email,
		user.Username,
		user.Password,
		user.Alamat,
		user.IdRole,
	).Row().Scan(&lastID, &createdAt)

	return lastID, createdAt, err
}

func (r *userRepository) GetByID(id string) (dto.UserDTO, error) {
	var user dto.UserDTO
	q := `SELECT id_user,name,email,username,alamat,created_at,status,id_role FROM users WHERE id_user=?`

	if err := r.db.Raw(q, id).Scan(&user).Error; err != nil {
		return dto.UserDTO{}, err
	}

	if user.IdUser == "" {
		return dto.UserDTO{}, errors.New("user tidak ditemukan")
	}

	return user, nil
}

func (r *userRepository) GetModelByID(id string) (models.User, error) {
	var user models.User
	err := r.db.Where("id_user=?", id).First(&user).Error
	return user, err
}

func (r *userRepository) DeleteUser(id string) error {
	return r.db.Delete(&models.User{}, "id_user=?", id).Error
}

func (r *userRepository) UpdateUser(id string, data map[string]interface{}) error {
	return r.db.Model(&models.User{}).Where("id_user=?", id).Updates(data).Error
}

func (r *userRepository) GetByEmail(email string) (models.User, error) {
	var user models.User
	err := r.db.Where("email=?", email).First(&user).Error
	return user, err
}
