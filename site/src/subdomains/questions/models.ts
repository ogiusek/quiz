import { ToUser, User } from "../users/models/user"
import { AnswerInputs, AnswerInputsAnswerType, AnswerOptions, AnswerOptionsAnswerType, Question as QuestionVO, QuestionSetDescription, QuestionSetName } from "./valueobjects"

// question set
// answer with type
// question

// question set

export type QuestionSet = {
  Id: string
  Name: QuestionSetName
  Description: QuestionSetDescription
  Owner: User
  Questions: Question[]
}

export const ToQuestionSet = (payload: any): QuestionSet => {
  return {
    Id: payload['id'],
    Name: new QuestionSetName(payload['name']),
    Description: new QuestionSetDescription(payload['description']),
    Owner: ToUser(payload['owner']),
    Questions: [...(payload['questions'] ?? [])].map(payload => ToQuestion(payload))
  }
}

// answer with type

export type AnswerWithType = {
  AnswerType: AnswerOptionsAnswerType
  Answer: AnswerOptions
} | {
  AnswerType: AnswerInputsAnswerType
  Answer: AnswerInputs
}

// question

export type NewQuestion = {
  QuestionSetId: string
  Question: QuestionVO
} & AnswerWithType

export type Question = {
  Id: string
} & NewQuestion

export const ToQuestion = (payload: any): Question => {
  const panic = (): any => {
    throw new Error("not implemented")
  }
  return {
    Id: payload["id"],
    QuestionSetId: payload['question_set_id'],
    Question: new QuestionVO(payload['question']),
    AnswerType: payload["answer_type"],
    Answer: payload["answer_type"] == "o" ? AnswerOptions.Decode(payload["answer"]) :
      payload["answer_type"] == "i" ? AnswerInputs.Decode(payload["answer"]) :
        panic()
  }
}