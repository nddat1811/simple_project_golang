package entity

//Book struct represents book table in DB
type Book struct {
	ID         uint64 `gorm:"primary_key:auto_increment" json:"id"`
	Title      string `gorm:"type:varchar(255)" json:"tittle"`
	Desciption string `gorm:"type:text" json:"description"`
	UserId     uint64 `gorm:"not null" json:"-"`
	User       User   `gorm:"foreignkey:UserID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"user"`
}
