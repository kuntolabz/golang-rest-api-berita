package repositories

import (
	"errors"

	dto "github.com/kunto/golang-rest-api-berita/dto/cms"
	"github.com/kunto/golang-rest-api-berita/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetList(search string, limit int, offset int) ([]dto.UserDTO, int64, error)
	CheckEmailExists(email string) (bool, error)
	CheckUsernameExists(username string) (bool, error)
	CreateUser(user models.Ms_user) (string, string, error)
	GetByID(id string) (dto.UserDTO, error)
	UpdateUser(id string, data map[string]interface{}) error
	DeleteUser(id string) error
	GetModelByID(id string) (models.Ms_user, error)
	GetByEmail(email string) (models.Ms_user, error)
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
	totalQuery := "SELECT COUNT(*) FROM ms_users " + where
	if err := r.db.Raw(totalQuery, params...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	// DATA RESULT
	var users []dto.UserDTO
	dataQuery := `
		SELECT id_user, name, email, username, alamat, created_at, status
		FROM ms_users ` + where + `
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
	err := r.db.Raw("SELECT EXISTS(SELECT 1 FROM ms_users WHERE LOWER(email)=LOWER(?))", email).
		Scan(&exists).Error
	return exists, err
}

func (r *userRepository) CheckUsernameExists(username string) (bool, error) {
	var exists bool
	err := r.db.Raw("SELECT EXISTS(SELECT 1 FROM ms_users WHERE LOWER(username)=LOWER(?))", username).
		Scan(&exists).Error
	return exists, err
}

func (r *userRepository) CreateUser(user models.Ms_user) (string, string, error) {
	var lastID string
	var createdAt string

	query := `
		INSERT INTO ms_users (name,email,username,password,alamat,created_at,status,id_role,created_by)
		VALUES (?,?,?,?,?,NOW(),1,?,?)
		RETURNING id_user, created_at;
	`

	err := r.db.Raw(query,
		user.Name,
		user.Email,
		user.Username,
		user.Password,
		user.Alamat,
		user.IdRole,
		user.CreatedBy,
	).Row().Scan(&lastID, &createdAt)

	return lastID, createdAt, err
}

func (r *userRepository) GetByID(id string) (dto.UserDTO, error) {
	var user dto.UserDTO
	q := `SELECT id_user,name,email,username,alamat,created_at,status,id_role FROM ms_users WHERE id_user=?`

	if err := r.db.Raw(q, id).Scan(&user).Error; err != nil {
		return dto.UserDTO{}, err
	}

	if user.IdUser == "" {
		return dto.UserDTO{}, errors.New("user tidak ditemukan")
	}

	return user, nil
}

func (r *userRepository) GetModelByID(id string) (models.Ms_user, error) {
	var user models.Ms_user
	err := r.db.Where("id_user=?", id).First(&user).Error
	return user, err
}

func (r *userRepository) DeleteUser(id string) error {
	return r.db.Delete(&models.Ms_user{}, "id_user=?", id).Error
}

func (r *userRepository) UpdateUser(id string, data map[string]interface{}) error {
	data["updated_by"] = data["updated_by"]

	return r.db.
		Model(&models.Ms_user{}).
		Where("id_user = ?", id).
		Updates(data).Error
}

func (r *userRepository) GetByEmail(email string) (models.Ms_user, error) {
	var user models.Ms_user
	err := r.db.Where("email=?", email).First(&user).Error
	return user, err
}
