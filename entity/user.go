package entity

import (
	"go-blogrpl/utils"

	"gorm.io/gorm"
)

type User struct {
	utils.Model
	Name         string        `json:"name" binding:"required"`
	Username     string        `json:"username" binding:"required"`
	Email        string        `json:"email" binding:"required"`
	Password     string        `json:"password" binding:"required"`
	Role         string        `json:"role" binding:"required"`
	Blogs        []Blog        `json:"blog,omitempty"`
	BlogLikes    []BlogLike    `json:"bloglikes,omitempty"`
	CommentLikes []CommentLike `json:"commentlikes,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	var err error
	u.Password, err = utils.PasswordHash(u.Password)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if u.Password != "" && tx.Statement.Changed("Password") {
		var err error
		u.Password, err = utils.PasswordHash(u.Password)
		if err != nil {
			return err
		}
	}
	return nil
}
