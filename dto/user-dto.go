package dto

type UserUpdateDTO struct {
	ID        uint64 `json:"id" form:"id"`
	FirstName string `json:"first_name" form:"first_name" binding:"required"`
	LastName  string `json:"last_name" form:"last_name" binding:"required"`
	Email     string `json:"email" form:"email" binding:"required,email"`
	Password  string `json:"password,omitempty" form:"password,omitempty"`
}

// type UserCreateDTO struct {
// 	FirstName string `json:"first_name" form:"first_name" binding:"required"`
// 	LastName  string `json:"last_name" form:"last_name" binding:"required"`
// 	Email     string `json:"email" form:"email" binding:"required" validate:"email"`
// 	Password  string `json:"password,omitempty" form:"password,omitempty" validate:"min:6" binding:"required"`
// }
