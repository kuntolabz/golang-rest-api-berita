package controllers

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/kunto/golang-rest-api-berita/config"
	"github.com/kunto/golang-rest-api-berita/dto"
	"github.com/kunto/golang-rest-api-berita/models"
	"github.com/kunto/golang-rest-api-berita/utils"
)

func GetListUsers(c *gin.Context) {
	search := strings.TrimSpace(c.Query("search"))
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limit := 10
	offset := 0

	// Validasi dan parse limit
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err != nil || parsed < 1 || parsed > 100 {
			utils.ErrorResponse(c, 400, "Limit harus angka antara 1-100")
			return
		} else {
			limit = parsed
		}
	}

	// Validasi dan parse offset
	if offsetStr != "" {
		if parsed, err := strconv.Atoi(offsetStr); err != nil || parsed < 0 {
			utils.ErrorResponse(c, 400, "Offset harus angka non-negatif")
			return
		} else {
			offset = parsed
		}
	}

	// Validasi search
	if len(search) > 100 {
		utils.ErrorResponse(c, 400, "Search terlalu panjang (maksimal 100 karakter)")
		return
	}

	// Bangun WHERE dan params (tetap raw SQL)
	where := "WHERE 1=1"
	params := []interface{}{}
	if search != "" {
		where += " AND (LOWER(name) LIKE LOWER(?) OR LOWER(email) LIKE LOWER(?) OR LOWER(username) LIKE LOWER(?))"
		like := "%" + search + "%"
		params = append(params, like, like, like)
	}

	// Hitung total row
	var total int64
	totalQuery := "SELECT COUNT(*) FROM users " + where
	if err := config.DB.Raw(totalQuery, params...).Scan(&total).Error; err != nil {
		log.Printf("Error counting users: %v", err) // Logging internal
		utils.ErrorResponse(c, 500, "Gagal menghitung data")
		return
	}

	// Data list (tetap raw SQL)
	var users []dto.UserDTO
	dataQuery := `
        SELECT id_user, name, email, username, alamat, created_at,status
        FROM users
    ` + where + ` ORDER BY name ASC LIMIT ? OFFSET ?`
	params = append(params, limit, offset)

	if err := config.DB.Raw(dataQuery, params...).Scan(&users).Error; err != nil {
		log.Printf("Error fetching users: %v", err)
		utils.ErrorResponse(c, 500, "Gagal mengambil data")
		return
	}

	// Success
	utils.SuccessResponse(c, users, total, offset, limit)
}

func CreateUser(c *gin.Context) {
	var input dto.InsertUserDTO

	// Trim input
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)
	input.Username = strings.TrimSpace(input.Username)
	input.Password = strings.TrimSpace(input.Password)
	input.Alamat = strings.TrimSpace(input.Alamat)
	input.IdRole = strings.TrimSpace(input.IdRole)

	// Validasi input dasar
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, 400, "Input tidak valid", err.Error())
		return
	}
	if input.Name == "" || input.Email == "" || input.Username == "" || input.Password == "" {
		utils.ErrorResponse(c, 400, "Name, email, username, dan password wajib diisi")
		return
	}

	// Validasi email dengan regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(input.Email) {
		utils.ErrorResponse(c, 400, "Format email tidak valid")
		return
	}

	// Validasi password
	if len(input.Password) < 8 || !regexp.MustCompile(`[A-Z]`).MatchString(input.Password) || !regexp.MustCompile(`[a-z]`).MatchString(input.Password) || !regexp.MustCompile(`[0-9]`).MatchString(input.Password) {
		utils.ErrorResponse(c, 400, "Password minimal 8 karakter, harus mengandung huruf besar, kecil, dan angka")
		return
	}

	// Cek email dengan error handling dan case-insensitive
	var exists bool
	if err := config.DB.Raw("SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(email) = LOWER(?))", input.Email).Scan(&exists).Error; err != nil {
		utils.ErrorResponse(c, 500, "Gagal memeriksa email", err.Error())
		return
	}
	if exists {
		utils.ErrorResponse(c, 400, "Email sudah digunakan")
		return
	}

	// Cek username (mirip)
	if err := config.DB.Raw("SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(username) = LOWER(?))", input.Username).Scan(&exists).Error; err != nil {
		utils.ErrorResponse(c, 500, "Gagal memeriksa username", err.Error())
		return
	}
	if exists {
		utils.ErrorResponse(c, 400, "Username sudah digunakan")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.ErrorResponse(c, 500, "Gagal memproses password", err.Error())
		return
	}

	// Insert dengan transaksi dan error handling konsisten
	var lastID string
	var createdDate string
	tx := config.DB.Begin()
	query := `INSERT INTO users (name, email, username, password, alamat, created_at, status, id_role) VALUES (?, ?, ?, ?, ?, NOW(), 1,?) RETURNING id_user, created_at;`
	err = tx.Raw(query, input.Name, input.Email, input.Username, string(hashedPassword), input.Alamat, input.IdRole).Row().Scan(&lastID, &createdDate)
	if err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, 500, "Gagal insert data", err.Error())
		return
	}
	tx.Commit()

	// Response
	utils.Success(c, gin.H{
		"id":          lastID,
		"name":        input.Name,
		"email":       input.Email,
		"username":    input.Username,
		"alamat":      input.Alamat,
		"createdDate": createdDate,
		"status":      1,
		"id_role":     input.IdRole,
	})
}

func GetUserDetail(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		utils.ErrorResponse(c, 400, "ID user tidak boleh kosong")
		return
	}

	// Query detail user
	var user dto.UserDTO
	query := `SELECT id_user,name,email,username,alamat,created_at,status FROM users WHERE id_user = ?`
	if err := config.DB.Raw(query, id).Scan(&user).Error; err != nil {
		log.Printf("Error fetching user detail: %v", err)
		utils.ErrorResponse(c, 500, "Gagal mengambil data")
		return
	}

	// Jika user tidak ditemukan
	if user.IdUser == "" {
		utils.ErrorResponse(c, 404, "User tidak ditemukan")
		return
	}

	utils.SuccessResponse(c, user, 1, 0, 1)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		utils.ErrorResponse(c, 400, "ID user tidak boleh kosong")
		return
	}

	// Bind data input JSON ke DTO
	var input dto.InsertUserDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, 400, "Input tidak valid", err.Error())
		return
	}

	// Cek user ada atau tidak
	var user models.User
	if err := config.DB.Where("id_user = ?", id).First(&user).Error; err != nil {
		utils.ErrorResponse(c, 404, "User tidak ditemukan")
		return
	}

	// Bangun map update
	data := map[string]interface{}{}

	if input.Name != "" {
		data["name"] = input.Name
	}
	if input.Email != "" {
		data["email"] = input.Email
	}
	if input.Username != "" {
		data["username"] = input.Username
	}
	if input.Alamat != "" {
		data["alamat"] = input.Alamat
	}
	if input.IdRole != "" {
		data["id_role"] = input.IdRole
	}

	// NOTE: jika tidak ada field yang dikirim akan error silent
	if len(data) == 0 {
		utils.ErrorResponse(c, 400, "Tidak ada field yang diupdate")
		return
	}

	// Update DB
	if err := config.DB.Model(&models.User{}).
		Where("id_user = ?", id).
		Updates(data).Error; err != nil {

		log.Printf("Error update user: %v", err)
		utils.ErrorResponse(c, 500, "Gagal memperbarui user")
		return
	}

	// Ambil kembali untuk response (data terbaru)
	config.DB.Where("id_user = ?", id).First(&user)

	// Response hanya field aman (DTO response)
	resp := dto.UserDTO{
		IdUser:   id,
		Name:     user.Name,
		Email:    user.Email,
		Username: user.Username,
		Alamat:   user.Alamat,
		IdRole:   user.IdRole,
	}

	utils.SuccessResponse(c, resp, 1, 0, 1)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		utils.ErrorResponse(c, 400, "ID user tidak boleh kosong")
		return
	}

	// 👇 Cek apakah user ada
	var user models.User
	if err := config.DB.Where("id_user = ?", id).First(&user).Error; err != nil {
		utils.ErrorResponse(c, 404, "User tidak ditemukan")
		return
	}

	// 👇 Delete permanen
	if err := config.DB.Delete(&user).Error; err != nil {
		utils.ErrorResponse(c, 500, "Gagal menghapus user")
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message": "User berhasil dihapus",
		"id":      id,
	}, 1, 0, 0)
}
