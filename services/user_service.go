package services

import (
	"errors"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/kunto/golang-rest-api-berita/dto"
	"github.com/kunto/golang-rest-api-berita/models"
	"github.com/kunto/golang-rest-api-berita/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetList(search string, limit, offset int) ([]dto.UserDTO, int64, error)
	Create(req dto.InsertUserDTO) (interface{}, error)
	GetDetail(id string) (dto.UserDTO, error)
	Update(id string, req dto.InsertUserDTO) (dto.UserDTO, error)
	Delete(id string) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(r repositories.UserRepository) UserService {
	return &userService{r}
}

func (s *userService) GetList(search string, limit, offset int) ([]dto.UserDTO, int64, error) {
	if limit < 1 || limit > 100 {
		return nil, 0, errors.New("limit harus 1-100")
	}
	if offset < 0 {
		return nil, 0, errors.New("offset tidak valid")
	}

	return s.repo.GetList(search, limit, offset)
}

func (s *userService) Create(req dto.InsertUserDTO) (interface{}, error) {

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	if req.Name == "" || req.Email == "" || req.Username == "" || req.Password == "" {
		return nil, errors.New("data wajib diisi")
	}

	// Email regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return nil, errors.New("format email salah")
	}

	// Password rule
	if len(req.Password) < 8 ||
		!regexp.MustCompile(`[A-Z]`).MatchString(req.Password) ||
		!regexp.MustCompile(`[a-z]`).MatchString(req.Password) ||
		!regexp.MustCompile(`[0-9]`).MatchString(req.Password) {
		return nil, errors.New("password minimal 8 karakter, kombinasi huruf & angka")
	}

	// Cek unique
	if exists, _ := s.repo.CheckEmailExists(req.Email); exists {
		return nil, errors.New("email sudah digunakan")
	}
	if exists, _ := s.repo.CheckUsernameExists(req.Username); exists {
		return nil, errors.New("username sudah digunakan")
	}

	// Hash
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Username: req.Username,
		Password: string(hash),
		Alamat:   req.Alamat,
		IdRole:   req.IdRole,
	}

	id, createdAt, err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":         id,
		"name":       user.Name,
		"email":      user.Email,
		"username":   user.Username,
		"alamat":     user.Alamat,
		"created_at": createdAt,
		"status":     1,
		"id_role":    user.IdRole,
	}, nil
}

func (s *userService) GetDetail(id string) (dto.UserDTO, error) {

	if _, err := uuid.Parse(id); err != nil {
		return dto.UserDTO{}, errors.New("id user tidak valid")
	}

	user, err := s.repo.GetByID(id)
	if err != nil {
		return dto.UserDTO{}, errors.New("user tidak ditemukan")
	}

	return user, nil
}

func (s *userService) Delete(id string) error {
	_, err := s.repo.GetModelByID(id)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}
	return s.repo.DeleteUser(id)
}

func (s *userService) Update(id string, input dto.InsertUserDTO) (dto.UserDTO, error) {
	// cek user exist
	_, err := s.repo.GetByID(id)
	if err != nil {
		return dto.UserDTO{}, errors.New("user tidak ditemukan")
	}

	// bangun map update
	updateData := map[string]interface{}{}

	if input.Name != "" {
		updateData["name"] = input.Name
	}
	if input.Email != "" {
		updateData["email"] = input.Email
	}
	if input.Username != "" {
		updateData["username"] = input.Username
	}
	if input.Alamat != "" {
		updateData["alamat"] = input.Alamat
	}
	if input.IdRole != "" {
		updateData["id_role"] = input.IdRole
	}

	if len(updateData) == 0 {
		return dto.UserDTO{}, errors.New("tidak ada field yang diupdate")
	}

	// save
	if err := s.repo.UpdateUser(id, updateData); err != nil {
		return dto.UserDTO{}, errors.New("gagal memperbarui user")
	}

	// reload updated data
	updated, _ := s.repo.GetByID(id)

	return dto.UserDTO{
		IdUser:   updated.IdUser,
		Name:     updated.Name,
		Email:    updated.Email,
		Username: updated.Username,
		Alamat:   updated.Alamat,
		IdRole:   updated.IdRole,
	}, nil
}
