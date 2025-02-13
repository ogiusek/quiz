import { AnsweredAt, AnswerTime, MatchCourseStep, MatchState } from "./valueobjects"
import { ToUser, User } from "../users/models/user"
import { Question, ToQuestion } from "../questions/models"

// match dto
// match course dto
// answered question dto
// player dto

// match dto

export type MatchDto = {
  Id: string
  State: MatchState
  QuestionSetId: string
  QuestionsAmount: number
  HostUserId: string
  Course: MatchCourseDto | undefined
  Players: PlayerDto[]
}

export const ToMatchDto = (payload: any): MatchDto => ({
  Id: payload.id,
  State: payload.state,
  QuestionSetId: payload.question_set_id,
  QuestionsAmount: payload.questions_amount,
  HostUserId: payload.host_user_id,
  Course: payload.course ? ToMatchCourseDto(payload.course) : undefined,
  Players: [...(payload.players ?? [])].map(player => ToPlayerDto(player)),
})

// match course dto

export type MatchCourseDto = {
  MatchId: string
  CurrentQuestionIndex: number
  CurrentQuestion: Question | undefined
  Step: MatchCourseStep
  LastStep: string
  NextStep: string
  AnsweredQuestions: AnsweredQuestionDto[]
}

export const ToMatchCourseDto = (payload: any): MatchCourseDto => ({
  MatchId: payload.match_id,
  CurrentQuestionIndex: payload.current_question_index,
  CurrentQuestion: payload.current_question ? ToQuestion(payload.current_question) : undefined,
  Step: payload.step,
  LastStep: payload.last_step,
  NextStep: payload.next_step,
  AnsweredQuestions: [...(payload.answered_questions ?? [])].map(aq => ToAnsweredQuestionDto(aq)),
})

// answered question dto

export type AnsweredQuestionDto = {
  Id: string
  MatchCourseId: string
  QuestionId: string
  Question: Question | undefined
  AnsweredCorrectly: boolean
  AnswerTime: AnswerTime
  AnswredAt: AnsweredAt
  UserId: string | null
}

export const ToAnsweredQuestionDto = (payload: any): AnsweredQuestionDto => ({
  Id: payload.id,
  MatchCourseId: payload.match_course_id,
  QuestionId: payload.question_id,
  Question: payload.question ? ToQuestion(payload.question) : undefined,
  AnsweredCorrectly: payload.answered_correctly,
  AnswerTime: payload.answer_time,
  AnswredAt: payload.answered_at,
  UserId: payload.user_id
})


// player dto

export type PlayerDto = {
  Id: string
  MatchId: string
  UserId: string
  User: User
  Online: boolean
  Score: number
}

export const ToPlayerDto = (payload: any): PlayerDto => ({
  Id: payload.id,
  MatchId: payload.match_id,
  UserId: payload.user_id,
  User: ToUser(payload.user),
  Online: payload.online,
  Score: payload.score,
})