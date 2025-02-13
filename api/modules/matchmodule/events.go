package matchmodule

import (
	"quizapi/common"
	"quizapi/modules/eventsmodule"
)

// contents list ->
// events

// created match
// changed match
// deleted match

// created match course
// changed match course
// deleted match course

// created answered question
// changed answered question
// deleted answered question

// created player
// changed player
// deleted player
// <- contents list

// events

var events []eventsmodule.Event = []eventsmodule.Event{
	eventsmodule.NewEvent(CreatedMatchEvent{}),
	eventsmodule.NewEvent(ChangedMatchEvent{}),
	eventsmodule.NewEvent(DeletedMatchEvent{}),
	eventsmodule.NewEvent(CreatedMatchCourseEvent{}),
	eventsmodule.NewEvent(ChangedMatchCourseEvent{}),
	eventsmodule.NewEvent(DeletedMatchCourseEvent{}),
	eventsmodule.NewEvent(CreatedAnsweredQuestionEvent{}),
	eventsmodule.NewEvent(ChangedAnsweredQuestionEvent{}),
	eventsmodule.NewEvent(DeletedAnsweredQuestionEvent{}),
	eventsmodule.NewEvent(CreatedPlayerEvent{}),
	eventsmodule.NewEvent(ChangedPlayerEvent{}),
	eventsmodule.NewEvent(DeletedPlayerEvent{}),
}

//

// created match
type CreatedMatchEvent struct {
	Services common.Ioc
	Model    MatchModel
}

func NewCreatedMatchEvent(c common.Ioc, model MatchModel) eventsmodule.Event {
	return eventsmodule.NewEvent(CreatedMatchEvent{Services: c, Model: model})
}

// changed match
type ChangedMatchEvent struct {
	Services common.Ioc
	Model    MatchModel
}

func NewChangedMatchEvent(c common.Ioc, model MatchModel) eventsmodule.Event {
	return eventsmodule.NewEvent(ChangedMatchEvent{Services: c, Model: model})
}

// deleted match
type DeletedMatchEvent struct {
	Services common.Ioc
	Model    MatchModel
}

func NewDeletedMatchEvent(c common.Ioc, model MatchModel) eventsmodule.Event {
	return eventsmodule.NewEvent(DeletedMatchEvent{Services: c, Model: model})
}

//

// created match course
type CreatedMatchCourseEvent struct {
	Services common.Ioc
	Match    MatchModel
	Model    MatchCourseModel
}

func NewCreatedMatchCourseEvent(c common.Ioc, match MatchModel, model MatchCourseModel) eventsmodule.Event {
	return eventsmodule.NewEvent(CreatedMatchCourseEvent{Services: c, Match: match, Model: model})
}

// changed match course
type ChangedMatchCourseEvent struct {
	Services common.Ioc
	Match    MatchModel
	Model    MatchCourseModel
}

func NewChangedMatchCourseEvent(c common.Ioc, match MatchModel, model MatchCourseModel) eventsmodule.Event {
	return eventsmodule.NewEvent(ChangedMatchCourseEvent{Services: c, Match: match, Model: model})
}

// deleted match course
type DeletedMatchCourseEvent struct {
	Services common.Ioc
	Match    MatchModel
	Model    MatchCourseModel
}

func NewDeletedMatchCourseEvent(c common.Ioc, match MatchModel, model MatchCourseModel) eventsmodule.Event {
	return eventsmodule.NewEvent(DeletedMatchCourseEvent{Services: c, Match: match, Model: model})
}

//

// created answered question
type CreatedAnsweredQuestionEvent struct {
	Services common.Ioc
	Match    MatchModel
	Model    AnsweredQuestionModel
}

func NewCreatedAnsweredQuestionEvent(c common.Ioc, match MatchModel, model AnsweredQuestionModel) eventsmodule.Event {
	return eventsmodule.NewEvent(CreatedAnsweredQuestionEvent{Services: c, Match: match, Model: model})
}

// changed answered question
type ChangedAnsweredQuestionEvent struct {
	Services common.Ioc
	Match    MatchModel
	Model    AnsweredQuestionModel
}

func NewChangedAnsweredQuestionEvent(c common.Ioc, match MatchModel, model AnsweredQuestionModel) eventsmodule.Event {
	return eventsmodule.NewEvent(ChangedAnsweredQuestionEvent{Services: c, Match: match, Model: model})
}

// deleted answered question
type DeletedAnsweredQuestionEvent struct {
	Services common.Ioc
	Match    MatchModel
	Model    AnsweredQuestionModel
}

func NewDeletedAnsweredQuestionEvent(c common.Ioc, match MatchModel, model AnsweredQuestionModel) eventsmodule.Event {
	return eventsmodule.NewEvent(DeletedAnsweredQuestionEvent{Services: c, Match: match, Model: model})
}

//

// created player
type CreatedPlayerEvent struct {
	Services common.Ioc
	Match    MatchModel
	Model    PlayerModel
}

func NewCreatedPlayerEvent(c common.Ioc, match MatchModel, model PlayerModel) eventsmodule.Event {
	return eventsmodule.NewEvent(CreatedPlayerEvent{Services: c, Match: match, Model: model})
}

// changed player
type ChangedPlayerEvent struct {
	Services common.Ioc
	Match    MatchModel
	Model    PlayerModel
}

func NewChangedPlayerEvent(c common.Ioc, match MatchModel, model PlayerModel) eventsmodule.Event {
	return eventsmodule.NewEvent(ChangedPlayerEvent{Services: c, Match: match, Model: model})
}

// deleted player
type DeletedPlayerEvent struct {
	Services common.Ioc
	Match    MatchModel
	Model    PlayerModel
}

func NewDeletedPlayerEvent(c common.Ioc, match MatchModel, model PlayerModel) eventsmodule.Event {
	return eventsmodule.NewEvent(DeletedPlayerEvent{Services: c, Match: match, Model: model})
}
