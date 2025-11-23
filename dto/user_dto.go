package dto

type UserDTO struct {
	IdUser     string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Alamat     string `json:"alamatUser"`
	Created_at string `json:"createdDate"`
}

type InsertUserDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Alamat   string `json:"alamat"`
}
