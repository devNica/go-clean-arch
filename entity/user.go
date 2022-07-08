package entity

type User struct {
	ID        uint64  `json:"id" gorm:"primaryKey:auto_increment"`
	FirstName string  `json:"first_name" gorm:"type:varchar(255);not null"`
	LastName  string  `json:"last_name" gorm:"type:varchar(255);not null"`
	Email     string  `json:"email" gorm:"uniqueIndex;type:varchar(255);not null"`
	Password  string  `json:"-" gorm:"->;<-;not null"`
	Token     string  `json:"token,omitempty" gorm:"-"`
	Books     *[]Book `json:"books,omitempty"`
}
