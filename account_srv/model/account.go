package model

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Mobile   string `gorm:"index:idx_mobile;unique;size:11;not null"`
	Password string `grom:"type:size:64;not null"`
	NickName string `gorm:"size:32"`
	Salt     string `gorm:"size:16"`
	Gender   string `gorm:"size:6;default:male"`
	Role     int    `gorm:"type:int;default:1;comment'1-普通用户，2-管理员'"`
}
