package matchmodule

import (
	"log"
	"quizapi/modules/eventsmodule"
	"quizapi/modules/timemodule"
	"quizapi/modules/wsmodule"
	"time"
)

// contents list ->
// event handlers

// created match handler
// changed match handler
// deleted match handler

// created match course handler
// changed match course handler
// deleted match course handler

// created answered question handler
// changed answered question handler
// deleted answered question handler

// created player handler
// changed player handler
// deleted player handler
// <- contents list

// event handlers

var eventHandlers map[string]func(any) = map[string]func(any){
	eventsmodule.NewEvent(CreatedMatchEvent{}).Topic: func(a any) { createdMatchHandler(a.(CreatedMatchEvent)) },
	eventsmodule.NewEvent(ChangedMatchEvent{}).Topic: func(a any) { changedMatchHandler(a.(ChangedMatchEvent)) },
	eventsmodule.NewEvent(DeletedMatchEvent{}).Topic: func(a any) { deletedMatchHandler(a.(DeletedMatchEvent)) },

	eventsmodule.NewEvent(CreatedMatchCourseEvent{}).Topic: func(a any) { createdMatchCourseHandler(a.(CreatedMatchCourseEvent)) },
	eventsmodule.NewEvent(ChangedMatchCourseEvent{}).Topic: func(a any) { changedMatchCourseHandler(a.(ChangedMatchCourseEvent)) },
	eventsmodule.NewEvent(DeletedMatchCourseEvent{}).Topic: func(a any) { deletedMatchCourseHandler(a.(DeletedMatchCourseEvent)) },

	eventsmodule.NewEvent(CreatedAnsweredQuestionEvent{}).Topic: func(a any) { createdAnsweredQuestionHandler(a.(CreatedAnsweredQuestionEvent)) },
	eventsmodule.NewEvent(ChangedAnsweredQuestionEvent{}).Topic: func(a any) { changedAnsweredQuestionHandler(a.(ChangedAnsweredQuestionEvent)) },
	eventsmodule.NewEvent(DeletedAnsweredQuestionEvent{}).Topic: func(a any) { deletedAnsweredQuestionHandler(a.(DeletedAnsweredQuestionEvent)) },

	eventsmodule.NewEvent(CreatedPlayerEvent{}).Topic: func(a any) { createdPlayerHandler(a.(CreatedPlayerEvent)) },
	eventsmodule.NewEvent(ChangedPlayerEvent{}).Topic: func(a any) { changedPlayerHandler(a.(ChangedPlayerEvent)) },
	eventsmodule.NewEvent(DeletedPlayerEvent{}).Topic: func(a any) { deletedPlayerHandler(a.(DeletedPlayerEvent)) },
}

// created match handler

func createdMatchHandler(event CreatedMatchEvent) {
	var repo MatchRepository
	event.Services.Inject(&repo)
	repo.Create(event.Services, event.Model)

	var manager wsmodule.SocketsMessager
	event.Services.Inject(&manager)
	for _, player := range event.Model.onlinePlayers() {
		manager.Send(wsmodule.SocketId(player.UserId), wsmodule.NewMessage("match/created_match", event.Model.FullDto(event.Services)))
	}
}

// changed match handler

func changedMatchHandler(event ChangedMatchEvent) {
	var repo MatchRepository
	event.Services.Inject(&repo)
	repo.Update(event.Services, event.Model)

	var manager wsmodule.SocketsMessager
	event.Services.Inject(&manager)
	for _, player := range event.Model.onlinePlayers() {
		manager.Send(wsmodule.SocketId(player.UserId), wsmodule.NewMessage("match/changed_match", event.Model.Dto()))
	}
}

// deleted match handler

func deletedMatchHandler(event DeletedMatchEvent) {
	var repo MatchRepository
	event.Services.Inject(&repo)
	err := repo.Delete(event.Services, event.Model.Id)
	if err != nil {
		var logger log.Logger
		event.Services.Inject(&logger)
		logger.Printf("error deleting match %s", err.Error())
	}

	var manager wsmodule.SocketsMessager
	event.Services.Inject(&manager)
	for _, player := range event.Model.onlinePlayers() {
		manager.Send(wsmodule.SocketId(player.UserId), wsmodule.NewMessage("match/deleted_match", event.Model.Dto()))
	}
}

//

// created match course handler

func createdMatchCourseHandler(event CreatedMatchCourseEvent) {
	var repo MatchCourseRepository
	event.Services.Inject(&repo)
	repo.Create(event.Services, event.Model)

	var manager wsmodule.SocketsMessager
	event.Services.Inject(&manager)
	for _, player := range event.Match.onlinePlayers() {
		manager.Send(wsmodule.SocketId(player.UserId), wsmodule.NewMessage("match/created_match_course", event.Model.FullDto(event.Services)))
	}

	var scheduler timemodule.Scheduler
	event.Services.Inject(&scheduler)
	scheduler.Schedule(event.Services, time.Time(event.Model.NextStep), func() {
		args := SyncArgs{
			MatchId: event.Model.MatchId,
		}
		args.Handle(event.Services)
	})
}

// changed match course handler

func changedMatchCourseHandler(event ChangedMatchCourseEvent) {
	var repo MatchCourseRepository
	event.Services.Inject(&repo)
	repo.Update(event.Services, event.Model)

	var manager wsmodule.SocketsMessager
	event.Services.Inject(&manager)
	for _, player := range event.Match.onlinePlayers() {
		manager.Send(wsmodule.SocketId(player.UserId), wsmodule.NewMessage("match/changed_match_course", event.Model.Dto(event.Services)))
	}

	var scheduler timemodule.Scheduler
	event.Services.Inject(&scheduler)
	scheduler.Schedule(event.Services, time.Time(event.Model.NextStep), func() {
		args := SyncArgs{
			MatchId: event.Model.MatchId,
		}
		args.Handle(event.Services)
	})
}

// deleted match course handler

func deletedMatchCourseHandler(event DeletedMatchCourseEvent) {
	var repo MatchCourseRepository
	event.Services.Inject(&repo)
	repo.Delete(event.Services, event.Model.Id)

	var manager wsmodule.SocketsMessager
	event.Services.Inject(&manager)
	for _, player := range event.Match.onlinePlayers() {
		manager.Send(wsmodule.SocketId(player.UserId), wsmodule.NewMessage("match/deleted_match_course", event.Model.Dto(event.Services)))
	}
}

//

// created answered question handler

func createdAnsweredQuestionHandler(event CreatedAnsweredQuestionEvent) {
	var repo AnsweredQuestionsRepository
	event.Services.Inject(&repo)
	repo.Create(event.Services, event.Model)

	var manager wsmodule.SocketsMessager
	event.Services.Inject(&manager)
	for _, player := range event.Match.onlinePlayers() {
		manager.Send(wsmodule.SocketId(player.UserId), wsmodule.NewMessage("match/created_answered_question", event.Model.Dto()))
	}
}

// changed answered question handler

func changedAnsweredQuestionHandler(event ChangedAnsweredQuestionEvent) {
	var repo AnsweredQuestionsRepository
	event.Services.Inject(&repo)
	repo.Update(event.Services, event.Model)

	var manager wsmodule.SocketsMessager
	event.Services.Inject(&manager)
	for _, player := range event.Match.onlinePlayers() {
		manager.Send(wsmodule.SocketId(player.UserId), wsmodule.NewMessage("match/changed_answered_question", event.Model.Dto()))
	}
}

// deleted answered question handler

func deletedAnsweredQuestionHandler(event DeletedAnsweredQuestionEvent) {
	var repo AnsweredQuestionsRepository
	event.Services.Inject(&repo)
	repo.Delete(event.Services, event.Model.Id)

	var manager wsmodule.SocketsMessager
	event.Services.Inject(&manager)
	for _, player := range event.Match.onlinePlayers() {
		manager.Send(wsmodule.SocketId(player.UserId), wsmodule.NewMessage("match/deleted_answered_question", event.Model.Dto()))
	}
}

//

// created player handler

func createdPlayerHandler(event CreatedPlayerEvent) {
	var repo PlayerRepository
	event.Services.Inject(&repo)
	repo.Create(event.Services, event.Model)

	var manager wsmodule.SocketsMessager
	event.Services.Inject(&manager)
	for _, player := range event.Match.onlinePlayers() {
		manager.Send(wsmodule.SocketId(player.UserId), wsmodule.NewMessage("match/created_player", event.Model.Dto()))
	}
}

// changed player handler

func changedPlayerHandler(event ChangedPlayerEvent) {
	var repo PlayerRepository
	event.Services.Inject(&repo)
	repo.Update(event.Services, event.Model)

	var manager wsmodule.SocketsMessager
	event.Services.Inject(&manager)
	for _, player := range event.Match.onlinePlayers() {
		manager.Send(wsmodule.SocketId(player.UserId), wsmodule.NewMessage("match/changed_player", event.Model.Dto()))
	}
}

// deleted player handler

func deletedPlayerHandler(event DeletedPlayerEvent) {
	var repo PlayerRepository
	event.Services.Inject(&repo)
	repo.Delete(event.Services, event.Model.Id)

	var manager wsmodule.SocketsMessager
	event.Services.Inject(&manager)
	for _, player := range event.Match.onlinePlayers() {
		manager.Send(wsmodule.SocketId(player.UserId), wsmodule.NewMessage("match/deleted_player", event.Model.Dto()))
	}
}
