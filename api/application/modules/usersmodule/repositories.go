package usersmodule

import (
	"log"
	"quizapi/common"
	"quizapi/modules/modelmodule"
	"quizapi/modules/wsmodule"

	"gorm.io/gorm"
)

// user repository
//

// user repository
type UserRepository interface {
	GetById(common.Ioc, modelmodule.ModelId) *UserModel
	GetByName(common.Ioc, UserName) *UserModel
	// Search() []UserModel

	Create(common.Ioc, UserModel) error
	Update(common.Ioc, UserModel) error
	Delete(common.Ioc, modelmodule.ModelId) error
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (repo *userRepository) db(c common.Ioc) *gorm.DB {
	var dbStorage common.ServiceStorage[*gorm.DB]
	c.Inject(&dbStorage)
	return dbStorage.MustGet()
}

func (repo *userRepository) GetById(c common.Ioc, id modelmodule.ModelId) *UserModel {
	var user UserModel
	if tx := repo.db(c).Where("id = ?", id).First(&user); tx.Error != nil {
		return nil
	}
	return &user
}

func (repo *userRepository) GetByName(c common.Ioc, name UserName) *UserModel {
	var user UserModel
	if tx := repo.db(c).Where("name = ?", name).First(&user); tx.Error != nil {
		return nil
	}
	return &user
}

func (repo *userRepository) Create(c common.Ioc, user UserModel) error {
	if tx := repo.db(c).Create(&user); tx.Error != nil {
		return common.ErrRepositoryConflict
	}
	return nil
}

func (repo *userRepository) Update(c common.Ioc, user UserModel) error {
	user.Changed()
	if tx := repo.db(c).Where("changes = ?", user.Changes).Save(&user); tx.Error != nil {
		return common.ErrRepositoryConflict // we assume all errors are duplicate key
	} else if tx.RowsAffected == 0 {
		return common.ErrRepositoryParallelModification
	}
	return nil
}

func (repo *userRepository) Delete(c common.Ioc, id modelmodule.ModelId) error {
	if tx := repo.db(c).Where("id = ?", id).Delete(&UserModel{}); tx.Error != nil {
		var logger log.Logger
		c.Inject(&logger)
		logger.Panic(tx.Error.Error())
	} else if tx.RowsAffected == 0 {
		return common.ErrRepositoryNotFound
	}
	return nil
}

// user socket repo

type UserSocketRepo interface {
	GetAll(c common.Ioc) []UserSocket
	GetByUser(c common.Ioc, userId modelmodule.ModelId) []UserSocket
	GetBySocket(c common.Ioc, socketId wsmodule.SocketId) *UserSocket

	Delete(c common.Ioc, socket UserSocket)
	DeleteBySocketId(c common.Ioc, id wsmodule.SocketId)
	Create(c common.Ioc, socket UserSocket)
}

type socketRepoImpl struct{}

func NewUserSocketRepository() UserSocketRepo {
	return &socketRepoImpl{}
}

func (repo *socketRepoImpl) db(c common.Ioc) *gorm.DB {
	var dbStorage common.ServiceStorage[*gorm.DB]
	c.Inject(&dbStorage)
	return dbStorage.MustGet()
}

func (repo *socketRepoImpl) GetAll(c common.Ioc) []UserSocket {
	var userSockets []UserSocket
	repo.db(c).Find(&userSockets)
	return userSockets
}

func (repo *socketRepoImpl) GetByUser(c common.Ioc, userId modelmodule.ModelId) []UserSocket {
	var sockets []UserSocket
	repo.db(c).Where("user_id = ?", userId).Find(&sockets)
	return sockets
}

func (repo *socketRepoImpl) GetBySocket(c common.Ioc, socketId wsmodule.SocketId) *UserSocket {
	var socket UserSocket
	if tx := repo.db(c).Where("socket_id = ?", socketId).First(&socket); tx.Error != nil {
		return nil
	}
	return &socket
}

func (repo *socketRepoImpl) Delete(c common.Ioc, socket UserSocket) {
	repo.db(c).Where("user_id = ? AND socket_id = ?", socket.UserId, socket.SocketId).Delete(&UserSocket{})
}

func (repo *socketRepoImpl) DeleteBySocketId(c common.Ioc, id wsmodule.SocketId) {
	repo.db(c).Where("socket_id = ?", id).Delete(&UserSocket{})
}

func (repo *socketRepoImpl) Create(c common.Ioc, socket UserSocket) {
	if tx := repo.db(c).Create(&socket); tx.Error != nil {
		log.Print(tx.Error.Error())
	}
}
