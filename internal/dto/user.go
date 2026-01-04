package dto

type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email" example:"user@example.com"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginUserResponse struct {
	UserID       string `json:"user_id"`
	UserRole     string `json:"user_role"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UpdateUserRequest struct {
	Email     string `json:"email" validate:"omitempty,email" example:"user@example.com"`
	Password  string `json:"password" validate:"omitempty,min=8"`
	FirstName string `json:"first_name" validate:"omitempty"`
	LastName  string `json:"last_name" validate:"omitempty"`
	Role      string `json:"role" validate:"omitempty,oneof=user admin"`
}
