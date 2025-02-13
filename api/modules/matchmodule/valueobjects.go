package matchmodule

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

// match state
// match course step
// answered at
// anser time
// match course step time

// match state

type MatchState string

const (
	StatePrepare MatchState = "prepare"
	StatePlaying MatchState = "playing"
)

func (MatchState) GormDataType() string { return "varchar(12)" }

var (
	errUknownMatchState error = errors.New("this is not a match state")
)

func (state *MatchState) Valid() []error {
	switch *state {
	case StatePrepare:
		return nil
	case StatePlaying:
		return nil
	}
	return []error{errUknownMatchState}
}

// match course step

type MatchCourseStep string

const (
	MatchCourseQuestion MatchCourseStep = "question"
	MatchCourseBreak    MatchCourseStep = "break"
	MatchCourseFinished MatchCourseStep = "finished"
)

func (MatchCourseStep) GormDataType() string { return "varchar(12)" }

var (
	errUknownMatchCourseStep error = errors.New("this is not a match course step")
)

func (step *MatchCourseStep) Valid() []error {
	switch *step {
	case MatchCourseQuestion:
		return nil
	case MatchCourseBreak:
		return nil
	case MatchCourseFinished:
		return nil
	}
	return []error{errUknownMatchCourseStep}
}

// answered at

type AnsweredAt time.Time

func (AnsweredAt) GormDataType() string { return "TIMESTAMP" }

// Implement Valuer interface
func (e AnsweredAt) Value() (driver.Value, error) {
	return time.Time(e).Format(time.RFC3339), nil
}

// Implement Scanner interface
func (e *AnsweredAt) Scan(value interface{}) error {
	switch x := value.(type) {
	case time.Time:
		*e = AnsweredAt(x)
	case string:
		// Parse string to time if needed
		parsedTime, err := time.Parse(time.RFC3339, x)
		if err != nil {
			return fmt.Errorf("cannot scan timestamp: %v", err)
		}
		*e = AnsweredAt(parsedTime)
	case []byte:
		// Handle byte slice representation
		parsedTime, err := time.Parse(time.RFC3339, string(x))
		if err != nil {
			return fmt.Errorf("cannot scan timestamp: %v", err)
		}
		*e = AnsweredAt(parsedTime)
	case nil:
		*e = AnsweredAt(time.Time{})
	default:
		return fmt.Errorf("unsupported scan type %T", value)
	}
	return nil
}

// anser time

type AnswerTime time.Duration

func (AnswerTime) GormDataType() string { return "INTEGER" }

// Implement Valuer interface
func (e AnswerTime) Value() (driver.Value, error) {
	return time.Duration(e).Milliseconds(), nil
}

// Implement Scanner interface
func (e *AnswerTime) Scan(value interface{}) error {
	count, ok := value.(int64)
	if !ok {
		return errors.New("NaN")
	}
	*e = AnswerTime(time.Millisecond * time.Duration(count))
	return nil
}

// match course step time

type MatchCourseStepTime time.Time

func (MatchCourseStepTime) GormDataType() string { return "TIMESTAMP" }

// Implement Valuer interface
func (e MatchCourseStepTime) Value() (driver.Value, error) {
	return time.Time(e).Format(time.RFC3339), nil
}

// Implement Scanner interface
func (e *MatchCourseStepTime) Scan(value interface{}) error {
	switch x := value.(type) {
	case time.Time:
		*e = MatchCourseStepTime(x)
	case string:
		// Parse string to time if needed
		parsedTime, err := time.Parse(time.RFC3339, x)
		if err != nil {
			return fmt.Errorf("cannot scan timestamp: %v", err)
		}
		*e = MatchCourseStepTime(parsedTime)
	case []byte:
		// Handle byte slice representation
		parsedTime, err := time.Parse(time.RFC3339, string(x))
		if err != nil {
			return fmt.Errorf("cannot scan timestamp: %v", err)
		}
		*e = MatchCourseStepTime(parsedTime)
	case nil:
		*e = MatchCourseStepTime(time.Time{})
	default:
		return fmt.Errorf("unsupported scan type %T", value)
	}
	return nil
}
