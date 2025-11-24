package dto

type UserDTO struct {
	IdUser     string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Alamat     string `json:"alamatUser"`
	Created_at string `json:"createdDate"`
	Status     string `json:"status"`
	IdRole     string `json:"id_role"`
}

type InsertUserDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Alamat   string `json:"alamat"`
	IdRole   string `json:"id_role" binding:"required"`
}
