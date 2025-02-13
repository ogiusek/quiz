import { AnsweredQuestionDto, MatchCourseDto, PlayerDto, ToAnsweredQuestionDto, ToMatchCourseDto, ToPlayerDto } from './../models';
import { WsDefinition } from './../../../common/ws/wsContext';
import { MatchDto, ToMatchDto } from "../models";
import { HrefAfterMatch } from '@/Router';

export interface IMatchChangesListener {
  SetWs(_: WsDefinition): void
  SetState(_: MatchDto | undefined): void
  SetSetter(_: ((_: MatchDto | undefined) => void) | undefined): void
}

class MatchChangesListener implements IMatchChangesListener {
  private _ws: WsDefinition | undefined
  private _state: MatchDto | undefined
  private _setter: ((_: MatchDto | undefined) => void) | undefined

  SetState(state: MatchDto | undefined): void { this._state = state }
  SetSetter(setter: ((_: MatchDto | undefined) => void) | undefined): void { this._setter = setter }

  private _SetState(state: MatchDto | undefined) {
    this._state = state
    this._setter && this._setter(state)
  }

  SetWs(ws: WsDefinition): void {
    if (this._ws) return
    this._ws = ws
    this._ws.CloseListener(() => (this._ws = undefined))

    this._ws.MessageListener("match/created_match", (payload) =>
      this._SetState(ToMatchDto(payload)))
    this._ws.MessageListener("match/changed_match", (payload) => ((payload: MatchDto) =>
      // this._SetState({ ...this._state!, ...ToMatchDto(payload) }))
      this._SetState({
        Id: this._state!.Id,
        State: payload.State,
        QuestionSetId: payload.QuestionSetId,
        QuestionsAmount: payload.QuestionsAmount,
        HostUserId: payload.HostUserId,
        Course: this._state!.Course,
        Players: this._state!.Players,
      }))(ToMatchDto(payload)))
    this._ws.MessageListener("match/deleted_match", (_) => {
      window.location.href = HrefAfterMatch
      this._SetState(undefined)
    })

    this._ws.MessageListener("match/created_match_course", (payload) =>
      this._SetState({ ...this._state!, Course: ToMatchCourseDto(payload) }))
    this._ws.MessageListener("match/changed_match_course", (payload) => ((payload: MatchCourseDto) => this._SetState({
      ...this._state!, Course: {
        MatchId: this._state!.Course!.MatchId,
        CurrentQuestionIndex: payload.CurrentQuestionIndex,
        CurrentQuestion: payload.CurrentQuestion,
        Step: payload.Step,
        LastStep: payload.LastStep,
        NextStep: payload.NextStep,
        AnsweredQuestions: this._state!.Course!.AnsweredQuestions,
      }
    }))(ToMatchCourseDto(payload)))
    this._ws.MessageListener("match/deleted_match_course", (_) =>
      this._SetState({ ...this._state!, Course: undefined }))

    this._ws.MessageListener("match/created_answered_question", (payload) => ((payload: AnsweredQuestionDto) =>
      this._SetState({ ...this._state!, Course: { ...this._state!.Course!, AnsweredQuestions: [...this._state!.Course!.AnsweredQuestions, payload] } })
    )(ToAnsweredQuestionDto(payload)))
    this._ws.MessageListener("match/changed_answered_question", (payload) => ((payload: AnsweredQuestionDto) =>
      this._SetState({ ...this._state!, Course: { ...this._state!.Course!, AnsweredQuestions: this._state!.Course!.AnsweredQuestions.map(q => q.Id == payload.Id ? payload : q) } })
    )(ToAnsweredQuestionDto(payload)))
    this._ws.MessageListener("match/deleted_answered_question", (payload) => ((payload: AnsweredQuestionDto) =>
      this._SetState({ ...this._state!, Course: { ...this._state!.Course!, AnsweredQuestions: this._state!.Course!.AnsweredQuestions.filter(q => q.Id != payload.Id) } })
    )(ToAnsweredQuestionDto(payload)))

    this._ws.MessageListener("match/created_player", (payload) => ((payload: PlayerDto) =>
      this._state && this._SetState({ ...this._state!, Players: [...this._state!.Players, payload] })
    )(ToPlayerDto(payload)))
    this._ws.MessageListener("match/changed_player", (payload) => ((payload: PlayerDto) =>
      this._state && this._SetState({ ...this._state!, Players: [...this._state!.Players].map(p => p.Id == payload.Id ? payload : p) })
    )(ToPlayerDto(payload)))
    this._ws.MessageListener("match/deleted_player", (payload) => ((payload: PlayerDto) =>
      this._SetState({ ...this._state!, Players: [...this._state!.Players].filter(p => p.Id != payload.Id) })
    )(ToPlayerDto(payload)))

    // this._ws.MessageListener("match/active_match", (payload: { match_id: string }) => this._reJoinListener && this._reJoinListener(payload.match_id))
  }
}

export const NewMatchChangesListener = (): IMatchChangesListener => new MatchChangesListener