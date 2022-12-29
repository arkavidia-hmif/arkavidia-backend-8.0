package repository

type SignInAdminRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,ascii"`
}
