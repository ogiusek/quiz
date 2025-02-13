package usersmodule

import (
	"quizapi/common"
	"quizapi/modules/filesmodule"
	"quizapi/modules/modelmodule"
)

// user model
//

// user model

type UserModel struct {
	modelmodule.Model
	Name  UserName            `gorm:"column:name;uniqueIndex;not null"`
	Image filesmodule.ImageId `gorm:"column:image;not null"`
	Hash  string              `gorm:"column:hash;type:varchar(64);not null"`
}

func (UserModel) TableName() string { return "users" }

func NewUser(model modelmodule.Model, name UserName, image filesmodule.ImageId, password UserPassword, hasher common.Hasher) UserModel {
	return UserModel{
		Model: model,
		Name:  name,
		Image: image,
		Hash:  hasher.Hash([]byte(password)),
	}
}

func (user *UserModel) ChangeName(newName UserName) {
	user.Name = newName
}

func (user *UserModel) ChangeProfilePicture(newImage filesmodule.ImageId) {
	user.Image = newImage
}

func (user *UserModel) ChangePassword(password UserPassword, hasher common.Hasher) {
	user.Hash = hasher.Hash([]byte(password))
}
