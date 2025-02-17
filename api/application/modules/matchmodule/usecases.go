package matchmodule

import (
	"quizapi/common"
	"quizapi/modules/modelmodule"
	"quizapi/modules/questionsmodule"
	"quizapi/modules/usersmodule"
	"quizapi/modules/wsmodule"
)

// active game
// host game
// join game
// quit game
// change question set
// change questions amount
// start game
// reset game
// answer
// sync

// active game

type ActiveGameArgs struct{}

func (args *ActiveGameArgs) Valid() []error {
	return nil
}

func (args *ActiveGameArgs) Handle(c common.Ioc) error {
	var sessionStorage common.ServiceStorage[usersmodule.SessionDto]
	c.Inject(&sessionStorage)
	session := sessionStorage.Get()
	if session == nil {
		return usersmodule.ErrUnauthorized
	}

	var playerRepo PlayerRepository
	c.Inject(&playerRepo)

	existing_player := playerRepo.GetByUserId(c, session.UserId)
	if existing_player == nil {
		return nil
	}

	var manager wsmodule.SocketsMessager
	c.Inject(&manager)

	res := wsmodule.NewMessage("match/active_match", struct {
		MatchId modelmodule.ModelId `json:"match_id"`
	}{MatchId: existing_player.MatchId})

	var resStorage common.ServiceStorage[common.Response]
	c.Inject(&resStorage)
	resStorage.Set(res)

	return nil
}

// host game

type HostArgs struct{}

func (args *HostArgs) Valid() []error {
	return nil
}

func (args *HostArgs) Handle(c common.Ioc) error {
	var sessionStorage common.ServiceStorage[usersmodule.SessionDto]
	c.Inject(&sessionStorage)
	session := sessionStorage.Get()
	if session == nil {
		return usersmodule.ErrUnauthorized
	}

	var playerRepo PlayerRepository
	c.Inject(&playerRepo)

	existing_player := playerRepo.GetByUserId(c, session.UserId)
	if existing_player != nil {
		return errHttpAlreadyInMatch
	}

	var userRepo usersmodule.UserRepository
	c.Inject(&userRepo)
	user := userRepo.GetById(c, session.UserId)
	if user == nil {
		return usersmodule.ErrUserNotFound
	}

	var matchRepo MatchRepository
	c.Inject(&matchRepo)

	Host(c, *user)

	return nil
}

// join game

type JoinArgs struct {
	MatchId modelmodule.ModelId `json:"match_id"`
}

func (args *JoinArgs) Valid() []error {
	var errors []error
	for _, err := range args.MatchId.Valid() {
		errors = append(errors, common.ErrPath(err).Property("match_id"))
	}
	return errors
}

func (args *JoinArgs) Handle(c common.Ioc) error {
	var sessionStorage common.ServiceStorage[usersmodule.SessionDto]
	c.Inject(&sessionStorage)
	session := sessionStorage.Get()
	if session == nil {
		return usersmodule.ErrUnauthorized
	}

	var userRepo usersmodule.UserRepository
	c.Inject(&userRepo)
	user := userRepo.GetById(c, session.UserId)
	if user == nil {
		return usersmodule.ErrUserNotFound
	}

	var matchRepo MatchRepository
	c.Inject(&matchRepo)
	match := matchRepo.GetById(c, args.MatchId)
	if match == nil {
		return errHttpMatchDoNotExists
	}

	err := match.Join(c, *user)

	if err != nil {
		return err
	}

	var resStorage common.ServiceStorage[common.Response]
	c.Inject(&resStorage)

	resStorage.Set(wsmodule.NewMessage("match/created_match", match.FullDto(c)))

	return nil
}

// quit game

type QuitArgs struct{}

func (args *QuitArgs) Handle(c common.Ioc) error {
	var sessionStorage common.ServiceStorage[usersmodule.SessionDto]
	c.Inject(&sessionStorage)
	session := sessionStorage.Get()
	if session == nil {
		return usersmodule.ErrUnauthorized
	}

	var playerRepo PlayerRepository
	c.Inject(&playerRepo)

	player := playerRepo.GetByUserId(c, session.UserId)

	if player == nil {
		return errHttpPlayerNotFound
	}

	var matchRepo MatchRepository
	c.Inject(&matchRepo)

	match := matchRepo.GetById(c, player.MatchId)

	if match == nil {
		return errHttpMatchDoNotExists
	}

	match.Quit(c, *player.User)

	var socketStorage common.ServiceStorage[wsmodule.SocketId]
	c.Inject(&socketStorage)
	id := socketStorage.MustGet()

	var socketsMessager wsmodule.SocketsMessager
	c.Inject(&socketsMessager)
	socketsMessager.Send(id, wsmodule.NewMessage("match/deleted_match", match.Dto()))

	return nil
}

// change question set

type ChangeQuestionSetArgs struct {
	QuestionSetId modelmodule.ModelId `json:"question_set_id"`
}

func (args *ChangeQuestionSetArgs) Valid() []error {
	var errors []error
	for _, err := range args.QuestionSetId.Valid() {
		errors = append(errors, common.ErrPath(err).Property("question_set_id"))
	}
	return errors
}

func (args *ChangeQuestionSetArgs) Handle(c common.Ioc) error {
	var sessionStorage common.ServiceStorage[usersmodule.SessionDto]
	c.Inject(&sessionStorage)
	session := sessionStorage.Get()
	if session == nil {
		return usersmodule.ErrUnauthorized
	}

	var playerRepo PlayerRepository
	c.Inject(&playerRepo)
	player := playerRepo.GetByUserId(c, session.UserId)

	if player == nil {
		return errHttpPlayerNotFound
	}

	var questionSetRepo questionsmodule.QuestionSetRepository
	c.Inject(&questionSetRepo)
	questionSet := questionSetRepo.GetById(c, args.QuestionSetId)

	if questionSet == nil {
		return questionsmodule.ErrQuestionSetNotFound
	}

	var matchRepo MatchRepository
	c.Inject(&matchRepo)
	match := matchRepo.GetById(c, player.MatchId)

	if match == nil {
		return errHttpMatchDoNotExists
	}

	err := match.ChangeQuestionSet(c, player.UserId, *questionSet)

	return err
}

// change questions amount

type ChangeQuestionsAmountArgs struct {
	QuestionsAmount int `json:"questions_amount"`
}

func (args *ChangeQuestionsAmountArgs) Valid() []error {
	var errors []error
	if args.QuestionsAmount <= 0 {
		errors = append(errors, common.ErrPath(errQuestionsAmountHasToBePositive).Property("questions_amount"))
	}
	return errors
}

func (args *ChangeQuestionsAmountArgs) Handle(c common.Ioc) error {
	var sessionStorage common.ServiceStorage[usersmodule.SessionDto]
	c.Inject(&sessionStorage)
	session := sessionStorage.Get()
	if session == nil {
		return usersmodule.ErrUnauthorized
	}

	var playerRepo PlayerRepository
	c.Inject(&playerRepo)
	player := playerRepo.GetByUserId(c, session.UserId)

	if player == nil {
		return errHttpPlayerNotFound
	}

	var matchRepo MatchRepository
	c.Inject(&matchRepo)
	match := matchRepo.GetById(c, player.MatchId)

	if match == nil {
		return errHttpMatchDoNotExists
	}

	err := match.ChangeQuestionsAmount(c, session.UserId, args.QuestionsAmount)

	return err
}

// start game

type StartArgs struct{}

func (args *StartArgs) Handle(c common.Ioc) error {
	var sessionStorage common.ServiceStorage[usersmodule.SessionDto]
	c.Inject(&sessionStorage)
	session := sessionStorage.Get()
	if session == nil {
		return usersmodule.ErrUnauthorized
	}

	var playerRepo PlayerRepository
	c.Inject(&playerRepo)
	player := playerRepo.GetByUserId(c, session.UserId)

	if player == nil {
		return errHttpPlayerNotFound
	}

	var matchRepo MatchRepository
	c.Inject(&matchRepo)
	match := matchRepo.GetById(c, player.MatchId)

	if match == nil {
		return errHttpMatchDoNotExists
	}

	err := match.Start(c, session.UserId)

	return err
}

// reset game

type ResetArgs struct{}

func (args *ResetArgs) Handle(c common.Ioc) error {
	var sessionStorage common.ServiceStorage[usersmodule.SessionDto]
	c.Inject(&sessionStorage)
	session := sessionStorage.Get()
	if session == nil {
		return usersmodule.ErrUnauthorized
	}

	var playerRepo PlayerRepository
	c.Inject(&playerRepo)
	player := playerRepo.GetByUserId(c, session.UserId)

	if player == nil {
		return errHttpPlayerNotFound
	}

	var matchRepo MatchRepository
	c.Inject(&matchRepo)
	match := matchRepo.GetById(c, player.MatchId)

	if match == nil {
		return errHttpMatchDoNotExists
	}

	err := match.Reset(c, session.UserId)

	return err
}

// answer

type AnswerArgs struct {
	Answer questionsmodule.AnswerMessage `json:"answer"`
}

func (args *AnswerArgs) Valid() []error {
	var errors []error
	for _, err := range args.Answer.Valid() {
		errors = append(errors, common.ErrPath(err).Property("answer"))
	}
	return errors
}

func (args *AnswerArgs) Handle(c common.Ioc) error {
	var sessionStorage common.ServiceStorage[usersmodule.SessionDto]
	c.Inject(&sessionStorage)
	session := sessionStorage.Get()
	if session == nil {
		return usersmodule.ErrUnauthorized
	}

	var playerRepo PlayerRepository
	c.Inject(&playerRepo)
	player := playerRepo.GetByUserId(c, session.UserId)

	if player == nil {
		return errHttpPlayerNotFound
	}

	var matchRepo MatchRepository
	c.Inject(&matchRepo)
	match := matchRepo.GetById(c, player.MatchId)

	if match == nil {
		return errHttpMatchDoNotExists
	}

	if match.Course == nil {
		return errHttpMatchCourseDoNotExists
	}

	err := match.Course.Answer(c, match, session.UserId, args.Answer)

	return err
}

// sync

// sync is invoked by sheduler
type SyncArgs struct {
	MatchId modelmodule.ModelId `json:"match_id"`
}

func (args *SyncArgs) Valid() []error {
	var errors []error
	for _, err := range args.MatchId.Valid() {
		errors = append(errors, common.ErrPath(err).Property("match_id"))
	}
	return errors
}

func (args *SyncArgs) Handle(c common.Ioc) error {
	var matchRepo MatchRepository
	c.Inject(&matchRepo)
	match := matchRepo.GetById(c, args.MatchId)

	if match == nil {
		return errHttpMatchDoNotExists
	}

	if match.Course == nil {
		return errHttpMatchCourseDoNotExists
	}

	err := match.Course.Sync(c, match)

	return err
}
