package modelmodule

import (
	"database/sql/driver"
	"quizapi/common"
	"quizapi/modules/timemodule"
	"time"

	"github.com/google/uuid"
)

type ModelId string

func (ModelId) GormDataType() string { return "varchar(36)" }
func (id *ModelId) Valid() []error {
	var errors []error
	if *id == "" {
		errors = append(errors, ErrIdCannotBeEmpty)
	}
	return errors
}

//

type ModelChanges int

func (ModelChanges) GormDataType() string    { return "integer" }
func (changes *ModelChanges) Valid() []error { return nil }

func (changes *ModelChanges) changed() {
	*changes = ModelChanges(*changes + 1)
}

//

type ModelCreatedAt time.Time

func (ModelCreatedAt) GormDataType() string      { return "TIMESTAMP" }
func (createdAt *ModelCreatedAt) Valid() []error { return nil }

func (createdAt ModelCreatedAt) Value() (driver.Value, error) {
	return createdAt.Format(), nil
}

func (createdAt ModelCreatedAt) Format() string {
	return timemodule.FormatDate(time.Time(createdAt))
}

//

type Model struct {
	Id        ModelId        `gorm:"column:id;primaryKey"`
	Changes   ModelChanges   `gorm:"column:changes;not null"`
	CreatedAt ModelCreatedAt `gorm:"column:created_at;not null"`
}

func (m *Model) Valid() []error {
	var errors []error
	for _, err := range m.Id.Valid() {
		errors = append(errors, common.ErrPath(err).Property("id"))
	}
	for _, err := range m.Changes.Valid() {
		errors = append(errors, common.ErrPath(err).Property("changes"))
	}
	for _, err := range m.CreatedAt.Valid() {
		errors = append(errors, common.ErrPath(err).Property("created_at"))
	}
	return errors
}

func (model *Model) Changed() {
	model.Changes.changed()
}

func NewModel(c common.Ioc) Model {
	var clock timemodule.Clock
	c.Inject(&clock)
	return Model{
		Id:        ModelId(uuid.NewString()), // space for making it a service
		Changes:   ModelChanges(0),
		CreatedAt: ModelCreatedAt(clock.Now()),
	}
}
