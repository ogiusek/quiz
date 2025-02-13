package matchmodule

import (
	"log"
	"quizapi/common"
	"quizapi/modules/eventsmodule"
	"quizapi/modules/modelmodule"
	"quizapi/modules/questionsmodule"
	"quizapi/modules/timemodule"
	"quizapi/modules/usersmodule"
	"quizapi/modules/wsmodule"
	"time"
)

// match model
// match course model
// answered question model
// player model
// match course questions

// match model

type MatchModel struct {
	modelmodule.Model
	State           MatchState                        `gorm:"column:state;not null"`
	QuestionSetId   *modelmodule.ModelId              `gorm:"column:question_set_id;null"`
	QuestionSet     *questionsmodule.QuestionSetModel `gorm:"foreignKey:QuestionSetId"`
	QuestionsAmount int                               `gorm:"column:questions_amount;type:integer;not null"`
	Course          *MatchCourseModel                 `gorm:"foreignKey:MatchId;constraint:OnDelete:CASCADE"`
	Players         []*PlayerModel                    `gorm:"foreignKey:MatchId;constraint:OnDelete:CASCADE"`
	HostUserId      modelmodule.ModelId               `gorm:"column:host_user_id;not null"`
}

var (
	// errAlreadyJoinedMatch             error = common.NewHttpError("you are already in match", 409)
	errMatchAlreadyStarted            error = common.NewHttpError("match already started", 409)
	errNotHost                        error = common.NewHttpError("you are not a host", 403)
	errToManyQuestions                error = common.NewHttpError("to many questions", 400)
	errQuestionsAmountHasToBePositive error = common.NewHttpError("questions amount has to be positive", 400)
	errMissingQuestionSet             error = common.NewHttpError("missing question set", 400)
)

func Host(c common.Ioc, user usersmodule.UserModel) {
	match := MatchModel{
		Model:           modelmodule.NewModel(c),
		State:           StatePrepare,
		QuestionSet:     nil,
		QuestionSetId:   nil,
		QuestionsAmount: 0,
		Course:          nil,
		Players:         nil,
		HostUserId:      user.Id,
	}
	player := NewPlayer(c, match.Id, user)
	match.Players = append(match.Players, &player)
	var manager eventsmodule.EventManager
	c.Inject(&manager)
	manager.Dispach(NewCreatedMatchEvent(c, match))
}

func (match *MatchModel) onlinePlayers() []*PlayerModel {
	var players []*PlayerModel
	for _, player := range match.Players {
		if player.Online {
			players = append(players, player)
		}
	}
	return players
}

func (match *MatchModel) Join(c common.Ioc, user usersmodule.UserModel) error {
	var eventManager eventsmodule.EventManager
	c.Inject(&eventManager)

	for _, player := range match.Players {
		if player.User.Id != user.Id {
			continue
		}

		if match.State == StatePrepare || player.Online {
			// return errAlreadyJoinedMatch
			var resStorage common.ServiceStorage[common.Response]
			c.Inject(&resStorage)
			resStorage.Set(wsmodule.NewMessage("match/created_match", match.FullDto(c)))
			return nil
		}

		player.Online = true
		eventManager.Dispach(NewChangedPlayerEvent(c, *match, *player))
		return nil
	}

	if match.State != StatePrepare {
		return errMatchAlreadyStarted
	}

	player := NewPlayer(c, match.Id, user)
	match.Players = append(match.Players, &player)
	eventManager.Dispach(NewCreatedPlayerEvent(c, *match, player))

	return nil
}

func (match *MatchModel) Quit(c common.Ioc, user usersmodule.UserModel) {
	var eventManager eventsmodule.EventManager
	c.Inject(&eventManager)
	var players []*PlayerModel
	for _, player := range match.onlinePlayers() {
		if player.User.Id != user.Id {
			players = append(players, player)
			continue
		}

		if match.State == StatePrepare {
			eventManager.Dispach(NewDeletedPlayerEvent(c, *match, *player))
			continue
		}

		player.Online = false
		eventManager.Dispach(NewChangedPlayerEvent(c, *match, *player))
		players = append(players, player)
	}

	match.Players = players

	containsPlayers := false
	for _, player := range match.Players {
		if !player.Online {
			continue
		}
		containsPlayers = true
		break
	}

	if !containsPlayers {
		eventManager.Dispach(NewDeletedMatchEvent(c, *match))
		return
	}
}

func (match *MatchModel) ChangeQuestionSet(c common.Ioc, userId modelmodule.ModelId, questionSet questionsmodule.QuestionSetModel) error {
	var eventManager eventsmodule.EventManager
	c.Inject(&eventManager)
	if match.State != StatePrepare {
		return errMatchAlreadyStarted
	}

	if match.HostUserId != userId {
		return errNotHost
	}

	match.QuestionSetId = &questionSet.Id
	match.QuestionSet = &questionSet

	if len(match.QuestionSet.Questions) < match.QuestionsAmount {
		match.QuestionsAmount = len(match.QuestionSet.Questions)
	}

	eventManager.Dispach(NewChangedMatchEvent(c, *match))

	return nil
}

func (match *MatchModel) ChangeQuestionsAmount(c common.Ioc, userId modelmodule.ModelId, questionsAmount int) error {
	var eventManager eventsmodule.EventManager
	c.Inject(&eventManager)

	if match.State != StatePrepare {
		return errMatchAlreadyStarted
	}

	if match.HostUserId != userId {
		return errNotHost
	}

	if match.QuestionSet == nil {
		return errMissingQuestionSet
	}

	if questionsAmount < 1 {
		return errQuestionsAmountHasToBePositive
	}

	if len(match.QuestionSet.Questions) < questionsAmount {
		return errToManyQuestions
	}

	match.QuestionsAmount = questionsAmount
	eventManager.Dispach(NewChangedMatchEvent(c, *match))

	return nil
}

func (match *MatchModel) Start(c common.Ioc, userId modelmodule.ModelId) error {
	var eventManager eventsmodule.EventManager
	c.Inject(&eventManager)

	if match.State != StatePrepare {
		return errMatchAlreadyStarted
	}

	if match.HostUserId != userId {
		return errNotHost
	}

	if match.QuestionSet == nil {
		return errMissingQuestionSet
	}

	if match.QuestionsAmount < 1 {
		return errQuestionsAmountHasToBePositive
	}

	match.State = StatePlaying

	questions := match.QuestionSet.GetRandomQuestions(c, match.QuestionsAmount)

	course := NewMatchCourse(c, match, questions)
	match.Course = &course

	eventManager.Dispach(NewChangedMatchEvent(c, *match))
	eventManager.Dispach(NewCreatedMatchCourseEvent(c, *match, course))

	return nil
}

func (match *MatchModel) Reset(c common.Ioc, userId modelmodule.ModelId) error {
	if userId != match.HostUserId {
		return errNotHost
	}
	match.reset(c)
	return nil
}

func (match *MatchModel) reset(c common.Ioc) {
	var eventManager eventsmodule.EventManager
	c.Inject(&eventManager)
	var players []*PlayerModel
	for _, player := range match.Players {
		if !player.Online {
			eventManager.Dispach(NewDeletedPlayerEvent(c, *match, *player))
			continue
		}
		if player.Score != 0 {
			player.Score = 0
			players = append(players, player)
			eventManager.Dispach(NewChangedPlayerEvent(c, *match, *player))
			continue
		}
		players = append(players, player)
	}
	match.Players = players
	course := match.Course
	match.Course = nil
	match.State = StatePrepare
	eventManager.Dispach(NewDeletedMatchCourseEvent(c, *match, *course))
	eventManager.Dispach(NewChangedMatchEvent(c, *match))
}

func (MatchModel) TableName() string {
	return "matches"
}

// match course model

type MatchCourseModel struct {
	modelmodule.Model
	MatchId           modelmodule.ModelId         `gorm:"column:match_id;not null"`
	Questions         []*MatchCourseQuestionModel `gorm:"foreignKey:MatchCourseId;constraint:OnDelete:CASCADE"`
	AnsweredQuestions []*AnsweredQuestionModel    `gorm:"foreignKey:MatchCourseId;constraint:OnDelete:CASCADE"`
	CurrentQuestion   int                         `gorm:"column:current_question;type:integer;not null"`
	Step              MatchCourseStep             `gorm:"column:step;not null"`
	LastStep          MatchCourseStepTime         `gorm:"column:last_step;not null"`
	NextStep          MatchCourseStepTime         `gorm:"column:next_step;not null"`
}

var (
	errNotQuestioning error = common.NewHttpError("this is not question time", 400)
	errUserNotInMatch error = common.NewHttpError("user does not play in match", 404)
)

const ( // of course all of this can be later extracted this is made this way for now
	scorePerCorrectQuestion   int           = 1000
	scorePerIncorrectQuestion int           = -250
	breakDuration             time.Duration = time.Second * 3
	questionDuration          time.Duration = time.Second * 30
)

func (course *MatchCourseModel) changeStep(step MatchCourseStep, now time.Time, duration time.Duration) {
	course.Step = step
	course.LastStep = MatchCourseStepTime(now)
	course.NextStep = MatchCourseStepTime(now.Add(duration))
}

func (course *MatchCourseModel) finish(clock timemodule.Clock) {
	course.Step = MatchCourseFinished
	now := clock.Now()
	course.LastStep = MatchCourseStepTime(now)
	course.NextStep = MatchCourseStepTime(now.Add(time.Second))
}

func (course *MatchCourseModel) hasMoreQuestions() bool {
	return course.CurrentQuestion+1 < len(course.questions())
}

func (course *MatchCourseModel) questions() []*questionsmodule.QuestionModel {
	var questions []*questionsmodule.QuestionModel
	for _, question := range course.Questions {
		questions = append(questions, question.Question)
	}
	return questions
}

func (course *MatchCourseModel) Answer(c common.Ioc, match *MatchModel, userId modelmodule.ModelId, userAnswer questionsmodule.AnswerMessage) error {
	var eventManager eventsmodule.EventManager
	c.Inject(&eventManager)

	var clock timemodule.Clock
	c.Inject(&clock)

	if course.Step != MatchCourseQuestion {
		return errNotQuestioning
	}

	var player *PlayerModel
	for _, comparedPlayer := range match.onlinePlayers() {
		if comparedPlayer.User.Id != userId {
			continue
		}
		player = comparedPlayer
		break
	}

	if player == nil {
		return errUserNotInMatch
	}

	question := course.questions()[course.CurrentQuestion]
	now := clock.Now()
	answerTime := now.Sub(time.Time(course.LastStep))
	answered := NewAnsweredQuestion(c, course.Id, *question, userAnswer, AnswerTime(answerTime), AnsweredAt(clock.Now()), userId)
	course.AnsweredQuestions = append(course.AnsweredQuestions, &answered)
	eventManager.Dispach(NewCreatedAnsweredQuestionEvent(c, *match, answered))

	if answered.AnsweredCorrectly {
		player.Score += scorePerCorrectQuestion
	} else {
		player.Score += scorePerIncorrectQuestion
	}

	eventManager.Dispach(NewChangedPlayerEvent(c, *match, *player))

	if course.hasMoreQuestions() {
		course.changeStep(MatchCourseBreak, now, breakDuration)
	} else {
		course.finish(clock)
	}
	eventManager.Dispach(NewChangedMatchCourseEvent(c, *match, *course))

	return nil
}

func (course *MatchCourseModel) Sync(c common.Ioc, match *MatchModel) error {
	var eventManager eventsmodule.EventManager
	c.Inject(&eventManager)

	var clock timemodule.Clock
	c.Inject(&clock)

	now := clock.Now()
	if now.UnixMilli() < time.Time(course.NextStep).UnixMilli() {
		return nil
	}

	switch course.Step {
	case MatchCourseBreak:
		course.changeStep(MatchCourseQuestion, now, questionDuration)
		course.CurrentQuestion += 1
		eventManager.Dispach(NewChangedMatchCourseEvent(c, *match, *course))
		return nil
	case MatchCourseQuestion:
		question := course.questions()[course.CurrentQuestion]
		now := clock.Now()
		answerTime := now.Sub(time.Time(course.LastStep))
		answered := NewNotAnsweredQuestion(c, course.Id, *question, AnswerTime(answerTime), AnsweredAt(clock.Now()))
		course.AnsweredQuestions = append(course.AnsweredQuestions, &answered)
		eventManager.Dispach(NewCreatedAnsweredQuestionEvent(c, *match, answered))

		if course.hasMoreQuestions() {
			course.changeStep(MatchCourseBreak, now, breakDuration)
		} else {
			course.finish(clock)
		}
		eventManager.Dispach(NewChangedMatchCourseEvent(c, *match, *course))
		return nil
	case MatchCourseFinished:
		match.reset(c)
		return nil
	}

	log.Panic("missing step handler")
	return nil
}

func NewMatchCourse(c common.Ioc, match *MatchModel, questions []*questionsmodule.QuestionModel) MatchCourseModel {
	var clock timemodule.Clock
	c.Inject(&clock)
	now := clock.Now()

	course := MatchCourseModel{
		Model:             modelmodule.NewModel(c),
		MatchId:           match.Id,
		AnsweredQuestions: make([]*AnsweredQuestionModel, 0),
		CurrentQuestion:   -1,
		Step:              MatchCourseBreak,
		LastStep:          MatchCourseStepTime(now),
		NextStep:          MatchCourseStepTime(now.Add(breakDuration)),
	}

	var relations []*MatchCourseQuestionModel
	for i, question := range questions {
		model := NewMatchCourseQuestion(course.Id, question.Id, i)
		relations = append(relations, &model)
	}
	course.Questions = relations
	return course
}

func (MatchCourseModel) TableName() string {
	return "match_courses"
}

// answered question model

type AnsweredQuestionModel struct {
	modelmodule.Model
	MatchCourseId     modelmodule.ModelId            `gorm:"column:match_course_id;not null;constraint:OnDelete:CASCADE"`
	QuestionId        modelmodule.ModelId            `gorm:"column:question_id;not null"`
	Question          *questionsmodule.QuestionModel `gorm:"foreignKey:QuestionId"`
	AnsweredCorrectly bool                           `gorm:"column:answered_correctly;type:boolean;not null"`
	AnswerTime        AnswerTime                     `gorm:"column:answer_time;not null"`
	AnswredAt         AnsweredAt                     `gorm:"column:answered_at;not null"`
	UserId            *modelmodule.ModelId           `gorm:"column:user_id;null"`
}

func NewAnsweredQuestion(c common.Ioc, matchCourseId modelmodule.ModelId, question questionsmodule.QuestionModel, answer questionsmodule.AnswerMessage, answerTime AnswerTime, answeredAt AnsweredAt, userId modelmodule.ModelId) AnsweredQuestionModel {
	return AnsweredQuestionModel{
		Model:             modelmodule.NewModel(c),
		MatchCourseId:     matchCourseId,
		QuestionId:        question.Id,
		Question:          &question,
		AnsweredCorrectly: question.GetAnswer().IsCorrectAnswer(answer),
		AnswerTime:        answerTime,
		AnswredAt:         answeredAt,
		UserId:            &userId,
	}
}

func NewNotAnsweredQuestion(c common.Ioc, matchCourseId modelmodule.ModelId, question questionsmodule.QuestionModel, answerTime AnswerTime, answeredAt AnsweredAt) AnsweredQuestionModel {
	return AnsweredQuestionModel{
		Model:             modelmodule.NewModel(c),
		MatchCourseId:     matchCourseId,
		QuestionId:        question.Id,
		Question:          &question,
		AnsweredCorrectly: false,
		AnswerTime:        answerTime,
		AnswredAt:         answeredAt,
		UserId:            nil,
	}
}

func (AnsweredQuestionModel) TableName() string {
	return "answered_questions"
}

// player model

type PlayerModel struct {
	modelmodule.Model
	MatchId modelmodule.ModelId    `gorm:"column:match_id;not null"`
	UserId  modelmodule.ModelId    `gorm:"column:user_id;not null"`
	User    *usersmodule.UserModel `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	Online  bool                   `gorm:"column:online;type:boolean;not null"`
	Score   int                    `gorm:"column:score;type:integer;not null"`
}

func NewPlayer(c common.Ioc, matchId modelmodule.ModelId, user usersmodule.UserModel) PlayerModel {
	return PlayerModel{
		Model:   modelmodule.NewModel(c),
		MatchId: matchId,
		UserId:  user.Id,
		User:    &user,
		Online:  true,
		Score:   0,
	}
}

func (PlayerModel) TableName() string {
	return "players"
}

// match course questions

type MatchCourseQuestionModel struct {
	MatchCourseId modelmodule.ModelId            `gorm:"column:match_course_id;not null"`
	QuestionId    modelmodule.ModelId            `gorm:"column:question_id;not null"`
	Index         int                            `gorm:"column:question_index;type:integer;not null"`
	MatchCourse   *MatchCourseModel              `gorm:"foreignKey:MatchCourseId;constraint:OnDelete:CASCADE"`
	Question      *questionsmodule.QuestionModel `gorm:"foreignKey:QuestionId"`
}

func (MatchCourseQuestionModel) TableName() string {
	return "match_course_questions"
}

func NewMatchCourseQuestion(matchCourseId modelmodule.ModelId, questionId modelmodule.ModelId, index int) MatchCourseQuestionModel {
	return MatchCourseQuestionModel{
		MatchCourseId: matchCourseId,
		QuestionId:    questionId,
		Index:         index,
	}
}
