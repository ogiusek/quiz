package matchmodule

import (
	"log"
	"quizapi/common"
	"quizapi/modules/modelmodule"

	"gorm.io/gorm"
)

// match repository
// match course repository
// answered question repository
// player repository

// match repository

type MatchRepository interface {
	GetById(common.Ioc, modelmodule.ModelId) *MatchModel

	Create(common.Ioc, MatchModel) error
	Update(common.Ioc, MatchModel) error
	Delete(common.Ioc, modelmodule.ModelId) error
}

type matchRepositoryImpl struct{}

func NewMatchRepository() MatchRepository {
	return &matchRepositoryImpl{}
}

func (*matchRepositoryImpl) db(c common.Ioc) *gorm.DB {
	var dbStorage common.ServiceStorage[*gorm.DB]
	c.Inject(&dbStorage)
	return dbStorage.MustGet()
}

func (repo *matchRepositoryImpl) GetById(c common.Ioc, id modelmodule.ModelId) *MatchModel {
	var model MatchModel
	if tx := repo.db(c).
		// Preload("QuestionSet").
		Preload("QuestionSet.Questions").
		// Preload("QuestionSet.Owner").
		Preload("Players", func(db *gorm.DB) *gorm.DB { return db.Joins("User") }).
		Preload("Course", func(db *gorm.DB) *gorm.DB {
			return db.
				Preload("Questions", func(db *gorm.DB) *gorm.DB {
					return db.
						Joins("Question").
						Order("match_course_questions.question_index ASC")
				}).
				Preload("AnsweredQuestions.Question")
		}).
		Where("\"matches\".\"id\" = ?", id).
		First(&model); tx.Error != nil {
		return nil
	}

	return &model
}

func (repo *matchRepositoryImpl) Create(c common.Ioc, model MatchModel) error {
	if tx := repo.db(c).Create(&model); tx.Error != nil {
		var logger log.Logger
		c.Inject(&logger)
		logger.Printf("MatchRepository.Create error: %s", tx.Error.Error())
		return common.ErrRepositoryConflict
	}
	return nil
}

func (repo *matchRepositoryImpl) Update(c common.Ioc, model MatchModel) error {
	model.Changed()
	if tx := repo.db(c).Omit("QuestionSet", "Course", "Players").Where("changes = ? AND id = ?", model.Changes, model.Id).Save(&model); tx.Error != nil {
		return common.ErrRepositoryConflict
	} else if tx.RowsAffected == 0 {
		return common.ErrRepositoryParallelModification
	}
	return nil
}

func (repo *matchRepositoryImpl) Delete(c common.Ioc, id modelmodule.ModelId) error {
	if tx := repo.db(c).Where("id = ?", id).Delete(&MatchModel{}); tx.Error != nil {
		var logger log.Logger
		c.Inject(&logger)
		logger.Print(tx.Error.Error())
	} else if tx.RowsAffected == 0 {
		return common.ErrRepositoryNotFound
	}
	return nil
}

// match course repository

type MatchCourseRepository interface {
	Create(common.Ioc, MatchCourseModel) error
	Update(common.Ioc, MatchCourseModel) error
	Delete(common.Ioc, modelmodule.ModelId) error
}

type matchCourseRepositoryImpl struct{}

func NewMatchCourseRepository() MatchCourseRepository {
	return &matchCourseRepositoryImpl{}
}

func (repo *matchCourseRepositoryImpl) db(c common.Ioc) *gorm.DB {
	var dbStorage common.ServiceStorage[*gorm.DB]
	c.Inject(&dbStorage)
	return dbStorage.MustGet()
}

func (repo *matchCourseRepositoryImpl) Create(c common.Ioc, model MatchCourseModel) error {
	db := repo.db(c)
	if tx := db.Omit("Match", "Questions").Create(&model); tx.Error != nil {
		var logger log.Logger
		c.Inject(&logger)
		logger.Printf("MatchCourseRepository.Create error: %s", tx.Error.Error())
		return common.ErrRepositoryConflict
	}

	var questions []MatchCourseQuestionModel
	for _, question := range model.Questions {
		questions = append(questions, *question)
	}

	if tx := db.Omit("Question", "MatchCourse").Create(&questions); tx.Error != nil {
		return common.ErrRepositoryConflict
	}

	return nil
}

func (repo *matchCourseRepositoryImpl) Update(c common.Ioc, model MatchCourseModel) error {
	model.Changed()
	if tx := repo.db(c).Omit("Questions", "AnsweredQuestions").Where("changes = ? AND id = ?", model.Changes, model.Id).Save(&model); tx.Error != nil {
		return common.ErrRepositoryConflict
	} else if tx.RowsAffected == 0 {
		return common.ErrRepositoryParallelModification
	}
	return nil
}

func (repo *matchCourseRepositoryImpl) Delete(c common.Ioc, id modelmodule.ModelId) error {
	if tx := repo.db(c).Where("id = ?", id).Delete(&MatchCourseModel{}); tx.Error != nil {
		var logger log.Logger
		c.Inject(&logger)
		logger.Panic(tx.Error.Error())
	} else if tx.RowsAffected == 0 {
		return common.ErrRepositoryNotFound
	}
	return nil
}

// answered question repository

type AnsweredQuestionsRepository interface {
	Create(common.Ioc, AnsweredQuestionModel) error
	Update(common.Ioc, AnsweredQuestionModel) error
	Delete(common.Ioc, modelmodule.ModelId) error
}

type answeredQuestionsRepositoryImpl struct{}

func NewAnsweredQuestionsRepository() AnsweredQuestionsRepository {
	return &answeredQuestionsRepositoryImpl{}
}

func (repo *answeredQuestionsRepositoryImpl) db(c common.Ioc) *gorm.DB {
	var dbStorage common.ServiceStorage[*gorm.DB]
	c.Inject(&dbStorage)
	return dbStorage.MustGet()
}

func (repo *answeredQuestionsRepositoryImpl) Create(c common.Ioc, model AnsweredQuestionModel) error {
	if tx := repo.db(c).Save(&model); tx.Error != nil {
		var logger log.Logger
		c.Inject(&logger)
		logger.Printf("AnsweredQuestionsRepository.Create error: %s", tx.Error.Error())
		return common.ErrRepositoryConflict
	}
	return nil
}

func (repo *answeredQuestionsRepositoryImpl) Update(c common.Ioc, model AnsweredQuestionModel) error {
	model.Changed()
	if tx := repo.db(c).Where("changes = ? AND id = ?", model.Changes, model.Id).Save(&model); tx.Error != nil {
		return common.ErrRepositoryConflict
	} else if tx.RowsAffected == 0 {
		return common.ErrRepositoryParallelModification
	}
	return nil
}

func (repo *answeredQuestionsRepositoryImpl) Delete(c common.Ioc, id modelmodule.ModelId) error {
	if tx := repo.db(c).Where("id = ?", id).Delete(&AnsweredQuestionModel{}); tx.Error != nil {
		var logger log.Logger
		c.Inject(&logger)
		logger.Panic(tx.Error.Error())
	} else if tx.RowsAffected == 0 {
		return common.ErrRepositoryNotFound
	}
	return nil
}

// player repository

type PlayerRepository interface {
	GetByUserId(c common.Ioc, userId modelmodule.ModelId) *PlayerModel

	Create(common.Ioc, PlayerModel) error
	Update(common.Ioc, PlayerModel) error
	Delete(common.Ioc, modelmodule.ModelId) error
}

type playerRepositoryImpl struct{}

func NewPlayerRepository() PlayerRepository {
	return &playerRepositoryImpl{}
}

func (repo *playerRepositoryImpl) db(c common.Ioc) *gorm.DB {
	var dbStorage common.ServiceStorage[*gorm.DB]
	c.Inject(&dbStorage)
	return dbStorage.MustGet()
}

func (repo *playerRepositoryImpl) GetByUserId(c common.Ioc, userId modelmodule.ModelId) *PlayerModel {
	db := repo.db(c)
	var player PlayerModel
	if tx := db.
		Joins("User").
		Where("user_id = ?", userId).
		First(&player); tx.Error == gorm.ErrRecordNotFound {
		return nil
	} else if tx.Error != nil {
		var logger log.Logger
		c.Inject(&logger)
		logger.Printf("PlayerRepository.GetByUserId error: %s", tx.Error.Error())
		return nil
	}
	return &player
}

func (repo *playerRepositoryImpl) Create(c common.Ioc, model PlayerModel) error {
	if tx := repo.db(c).Save(&model); tx.Error != nil {
		var logger log.Logger
		c.Inject(&logger)
		logger.Printf("PlayerRepository.Create error: %s", tx.Error.Error())
		return common.ErrRepositoryConflict
	}
	return nil
}

func (repo *playerRepositoryImpl) Update(c common.Ioc, model PlayerModel) error {
	model.Changed()
	if tx := repo.db(c).Where("changes = ? AND id = ?", model.Changes, model.Id).Save(&model); tx.Error != nil {
		return common.ErrRepositoryConflict
	} else if tx.RowsAffected == 0 {
		return common.ErrRepositoryParallelModification
	}
	return nil
}

func (repo *playerRepositoryImpl) Delete(c common.Ioc, id modelmodule.ModelId) error {
	if tx := repo.db(c).Where("id = ?", id).Delete(&PlayerModel{}); tx.Error != nil {
		var logger log.Logger
		c.Inject(&logger)
		logger.Panic(tx.Error.Error())
	} else if tx.RowsAffected == 0 {
		return common.ErrRepositoryNotFound
	}
	return nil
}
