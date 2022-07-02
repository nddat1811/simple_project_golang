package dto


//LoginDTO is used when client get from /login url
type LoginDTO struct{
	Email string `json:"email" form:"email" binding:"required" validate:"email"`
	Password string `json:"password" form:"password" validate:"min:6" binding:"required"`
}