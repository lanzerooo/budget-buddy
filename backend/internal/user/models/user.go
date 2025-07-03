package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

//ДОБАВЛЕНЫ СТРУКТУРЫ ДЛЯ НОВЫХ ЗАПРОСОВ

type UserProfileResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type UpdateProfileRequest struct {
	Name string `json:"name"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
