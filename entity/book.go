package entity

// Book struct represents books table in database
type Book struct {
	ID          uint64 `json:"id" gorm:"primaryKey:auto_increment"`
	Title       string `json:"title" gorm:"type:varchar(255);not null"`
	Description string `json:"description" gorm:"type:text"`
	UserID      uint64 `json:"user_id" gorm:"not null"`
	User        User   `json:"user" gorm:"foreignKey:UserID;constraint:onUpdate:CASCADE,onDelete:CASCADE"`
}
