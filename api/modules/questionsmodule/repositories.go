package questionsmodule

import (
	"fmt"
	"log"
	"quizapi/common"
	"quizapi/modules/modelmodule"
	"time"

	"gorm.io/gorm"
)

// question set repository
// question repository

// question set repository

type SearchQuestionSets struct {
	Phrase     string              `json:"search"`
	Page       uint                `json:"page"`
	OwnerId    modelmodule.ModelId `json:"owner_id"`
	LastUpdate time.Time           `json:"last_update"`
}

func (args *SearchQuestionSets) Valid() []error {
	var errors []error
	return errors
}

type QuestionSetRepository interface {
	GetById(common.Ioc, modelmodule.ModelId) *QuestionSetModel
	GetByOwnerId(common.Ioc, modelmodule.ModelId) []QuestionSetModel
	Search(common.Ioc, SearchQuestionSets) []QuestionSetModel

	Create(common.Ioc, QuestionSetModel) error
	Update(common.Ioc, QuestionSetModel) error
	Delete(common.Ioc, modelmodule.ModelId) error
}

type questionSetRepository struct{}

func NewQuestionSetRepository() QuestionSetRepository {
	return &questionSetRepository{}
}

func (repo *questionSetRepository) db(c common.Ioc) *gorm.DB {
	var dbStorage common.ServiceStorage[*gorm.DB]
	c.Inject(&dbStorage)
	return dbStorage.MustGet()
}

func (repo *questionSetRepository) GetById(c common.Ioc, id modelmodule.ModelId) *QuestionSetModel {
	var model QuestionSetModel
	if tx := repo.db(c).
		Preload("Questions", func(db *gorm.DB) *gorm.DB {
			return db.Order("questions.created_at ASC")
		}).
		Joins("Owner").
		Where("question_sets.id = ?", id).
		Order("question_sets.created_at DESC").
		First(&model); tx.Error != nil {
		return nil
	}
	return &model
}

func (repo *questionSetRepository) GetByOwnerId(c common.Ioc, id modelmodule.ModelId) []QuestionSetModel {
	var models []QuestionSetModel
	if tx := repo.db(c).
		Preload("Questions", func(db *gorm.DB) *gorm.DB {
			return db.Order("questions.created_at ASC")
		}).
		Joins("Owner").
		Where("question_sets.owner_id = ?", id).
		// Order("question_sets.created_at DESC").
		Find(&models); tx.Error != nil {
		log.Panic(tx.Error.Error())
	}
	return models
}

const pageSize uint = 10

func (repo *questionSetRepository) Search(c common.Ioc, args SearchQuestionSets) []QuestionSetModel {
	var models []QuestionSetModel
	tx := repo.db(c).
		Debug().
		Table("question_sets").
		Preload("Questions", func(db *gorm.DB) *gorm.DB {
			return db.Order("questions.created_at ASC")
		}).
		Joins("Owner").
		Order("question_sets.created_at DESC").
		Offset(int(pageSize) * int(args.Page)).
		Limit(int(pageSize))

	var defaultArgs SearchQuestionSets
	if args.Phrase != defaultArgs.Phrase {
		tx = tx.Where(
			"question_sets.name LIKE ? OR question_sets.description LIKE ?",
			fmt.Sprintf("%%%s%%", args.Phrase),
			fmt.Sprintf("%%%s%%", args.Phrase))
	}

	if args.LastUpdate != defaultArgs.LastUpdate {
		tx = tx.Where("question_sets.created_at <= ?", args.LastUpdate)
	}

	if args.OwnerId != defaultArgs.OwnerId {
		tx = tx.Where("question_sets.owner_id = ?", args.OwnerId)
	}

	tx = tx.Find(&models)
	if tx.Error != nil {
		log.Panic(tx.Error.Error())
	}
	return models
}

func (repo *questionSetRepository) Create(c common.Ioc, model QuestionSetModel) error {
	if tx := repo.db(c).Save(&model); tx.Error != nil {
		return common.ErrRepositoryConflict
	}
	return nil
}

func (repo *questionSetRepository) Update(c common.Ioc, model QuestionSetModel) error {
	model.Changed()
	if tx := repo.db(c).Where("changes = ? AND id = ?", model.Changes, model.Id).Save(&model); tx.Error != nil {
		return common.ErrRepositoryConflict
	} else if tx.RowsAffected == 0 {
		return common.ErrRepositoryParallelModification
	}
	return nil
}

func (repo *questionSetRepository) Delete(c common.Ioc, id modelmodule.ModelId) error {
	if tx := repo.db(c).Where("id = ?", id).Delete(&QuestionSetModel{}); tx.Error != nil {
		log.Panic(tx.Error.Error())
	} else if tx.RowsAffected == 0 {
		return common.ErrRepositoryNotFound
	}
	return nil
}

// question repository

type QuestionRepository interface {
	GetById(common.Ioc, modelmodule.ModelId) *QuestionModel

	Create(common.Ioc, QuestionModel) error
	Update(common.Ioc, QuestionModel) error
	Delete(common.Ioc, modelmodule.ModelId) error
}

type questionRepository struct{}

func NewQuestionRepository() QuestionRepository {
	return &questionRepository{}
}

func (repo *questionRepository) db(c common.Ioc) *gorm.DB {
	var dbStorage common.ServiceStorage[*gorm.DB]
	c.Inject(&dbStorage)
	return dbStorage.MustGet()
}

func (repo *questionRepository) GetById(c common.Ioc, id modelmodule.ModelId) *QuestionModel {
	var model QuestionModel
	if tx := repo.db(c).
		Joins("QuestionSet").
		Joins("QuestionSet.Owner").
		Where("questions.id = ?", id).
		First(&model); tx.Error != nil {
		return nil
	}
	return &model
}

func (repo *questionRepository) Create(c common.Ioc, model QuestionModel) error {
	if tx := repo.db(c).Omit("QuestionSet").Create(&model); tx.Error != nil {
		return common.ErrRepositoryConflict
	}
	return nil
}

func (repo *questionRepository) Update(c common.Ioc, model QuestionModel) error {
	model.Changed()
	if tx := repo.db(c).Where("changes = ? AND id = ?", model.Changes, model.Id).Save(&model); tx.Error != nil {
		return common.ErrRepositoryConflict
	} else if tx.RowsAffected == 0 {
		return common.ErrRepositoryNotFound
	}
	return nil
}

func (repo *questionRepository) Delete(c common.Ioc, id modelmodule.ModelId) error {
	if tx := repo.db(c).Where("id = ?", id).Delete(&QuestionModel{}); tx.Error != nil {
		log.Panic(tx.Error.Error())
	} else if tx.RowsAffected == 0 {
		return common.ErrRepositoryNotFound
	}
	return nil
}
