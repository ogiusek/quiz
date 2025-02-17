package usersmodule

import (
	"log"
	"quizapi/common"
	"quizapi/modules/modelmodule"

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
