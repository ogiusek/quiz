package usersmodule

import (
	"quizapi/common"
	"quizapi/modules/filesmodule"
	"quizapi/modules/modelmodule"
	"quizapi/modules/wsmodule"
)

// user model
// user sockets

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

// user sockets

type UserSocket struct {
	SocketId wsmodule.SocketId   `gorm:"column:socket_id;type:varchar(36);not null;primaryKey"`
	UserId   modelmodule.ModelId `gorm:"column:user_id;type:varchar(36);not null"`
}

func NewUserSocket(socketId wsmodule.SocketId, userId modelmodule.ModelId) UserSocket {
	return UserSocket{
		SocketId: socketId,
		UserId:   userId,
	}
}
