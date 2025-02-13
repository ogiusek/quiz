import { Session } from './../../users/models/session';
import { ApiEndpoint } from "@/common/api/api"
import { AnswerInputs, AnswerInputsAnswerType, AnswerOptions, AnswerOptionsAnswerType, Question, QuestionSetDescription, QuestionSetName } from "../valueobjects"
import { QuestionSet, ToQuestionSet } from '../models';

// create question set
// change question set name
// change question set description
// search question set
// get question set
// delete question set
// create question
// change question question
// change question answer
// delete question

// create question set

export type CreateQuestionSetArgs = {
  Session: Session
  Name: QuestionSetName
  Description: QuestionSetDescription
}

export const CreateQuestionSet = ApiEndpoint<CreateQuestionSetArgs, { Id: string }>(async (args, api) => {
  const response = await fetch(`${api.Url}/api/question-set/create`, {
    method: "POST",
    headers: { "Content-Type": "application-json", ...args.Session.Headers() },
    body: JSON.stringify({ "name": args.Name.Value, "description": args.Description.Value })
  })

  if (!response.ok) {
    await api.InvalidResponseHandler(response)
    return { Ok: false }
  }

  const payload = await response.json()

  return { Ok: true, Model: { Id: payload['id'] } }
})

// change question set name

export type ChangeQuestionSetNameArgs = {
  Session: Session
  Id: string
  NewName: QuestionSetName
}

export const ChangeQuestionSetName = ApiEndpoint<ChangeQuestionSetNameArgs, void>(async (args, api) => {
  const response = await fetch(`${api.Url}/api/question-set/change-name`, {
    method: "POST",
    headers: { "Content-Type": "application-json", ...args.Session.Headers() },
    body: JSON.stringify({ "id": args.Id, "new_name": args.NewName.Value })
  })

  if (!response.ok) {
    await api.InvalidResponseHandler(response)
    return { Ok: false }
  }

  return { Ok: true, Model: undefined }
})

// change question set description

export type ChangeQuestionSetDescriptionArgs = {
  Session: Session
  Id: string
  NewDescription: QuestionSetDescription
}

export const ChangeQuestionSetDescription = ApiEndpoint<ChangeQuestionSetDescriptionArgs, void>(async (args, api) => {
  const response = await fetch(`${api.Url}/api/question-set/change-description`, {
    method: "POST",
    headers: { "Content-Type": "application-json", ...args.Session.Headers() },
    body: JSON.stringify({ "id": args.Id, "new_description": args.NewDescription.Value })
  })

  if (!response.ok) {
    await api.InvalidResponseHandler(response)
    return { Ok: false }
  }

  return { Ok: true, Model: undefined }
})

// search question set

export type SearchQuestionSetsArgs = {
  Session: Session
  Search?: string
  Page?: number
  OwnerId?: string
  LastUpdate?: string
}

export type SearchQuestionSetsRes = {
  Found: QuestionSet[],
  Time: string,
}

export const SearchQuestionSets = ApiEndpoint<SearchQuestionSetsArgs, SearchQuestionSetsRes>(async (args, api) => {
  const argsObject: any = {}
  args.Search && (argsObject["search"] = args.Search)
  args.Page && (argsObject["page"] = args.Page)
  args.OwnerId && (argsObject["owner_id"] = args.OwnerId)
  args.LastUpdate && (argsObject["last_update"] = args.LastUpdate)

  const searchParams = new URLSearchParams({ args: JSON.stringify(argsObject) })

  const response = await fetch(`${api.Url}/api/question-set/search?${searchParams.toString()}`, {
    method: "GET",
    headers: { ...args.Session.Headers() },
  })

  if (!response.ok) {
    await api.InvalidResponseHandler(response)
    return { Ok: false }
  }

  const payload = await response.json()
  return {
    Ok: true,
    Model: {
      Found: [...payload['found']].map(payload => ToQuestionSet(payload)),
      Time: payload['when'],
    }
  }
})

// get question set

export type GetQuestionSetArgs = {
  Session: Session
  Id: string
}

// @ts-ignore
export const GetQuestionSet = ApiEndpoint<GetQuestionSetArgs, QuestionSet>(async (args, api) => {
  const searchParams = new URLSearchParams({
    args: JSON.stringify({
      'id': args.Id
    })
  })

  const response = await fetch(`${api.Url}/api/question-set/get?${searchParams.toString()}`, {
    method: "GET",
    headers: { ...args.Session.Headers() },
  })

  if (!response.ok) {
    await api.InvalidResponseHandler(response)
    return { Ok: false }
  }

  const payload = await response.json()

  return { Ok: true, Model: ToQuestionSet(payload.model) }
})

// delete question set

export type DeleteQuestionSetArgs = {
  Session: Session
  Id: string
}

export const DeleteQuestionSet = ApiEndpoint<DeleteQuestionSetArgs, void>(async (args, api) => {
  const response = await fetch(`${api.Url}/api/question-set/delete`, {
    method: "DELETE",
    headers: { "Content-Type": "application-json", ...args.Session.Headers() },
    body: JSON.stringify({ "id": args.Id })
  })

  if (!response.ok) {
    await api.InvalidResponseHandler(response)
    return { Ok: false }
  }

  return { Ok: true, Model: undefined }
})

// create question

export type CreateQuestionArgs = {
  Session: Session
  QuestionSetId: string
  Question: Question
} & ({
  AnswerType: AnswerInputsAnswerType
  Answer: AnswerInputs
} | {
  AnswerType: AnswerOptionsAnswerType
  Answer: AnswerOptions
})

export const CreateQuestion = ApiEndpoint<CreateQuestionArgs, { Id: string }>(async (args, api) => {
  const response = await fetch(`${api.Url}/api/question/create`, {
    method: "POST",
    headers: { "Content-Type": "application-json", ...args.Session.Headers() },
    body: JSON.stringify({
      "question_set_id": args.QuestionSetId,
      "question": args.Question.Value,
      "answer_type": args.AnswerType,
      "answer": args.Answer.Dto()
    })
  })

  if (!response.ok) {
    await api.InvalidResponseHandler(response)
    return { Ok: false }
  }

  const payload = await response.json()

  return { Ok: true, Model: { Id: payload['id'] } }
})

// change question question

export type ChangeQuestionArgs = {
  Session: Session
  Id: string
  NewQuestion: Question
}

export const ChangeQuestion = ApiEndpoint<ChangeQuestionArgs, void>(async (args, api) => {
  const response = await fetch(`${api.Url}/api/question/change-question`, {
    method: "POST",
    headers: { "Content-Type": "application-json", ...args.Session.Headers() },
    body: JSON.stringify({
      "id": args.Id,
      "new_question": args.NewQuestion.Value,
    })
  })

  if (!response.ok) {
    await api.InvalidResponseHandler(response)
    return { Ok: false }
  }

  return { Ok: true, Model: undefined }
})

// change question answer

export type ChangeQuestionAnswerArgs = {
  Session: Session
  Id: string
} & ({
  NewAnswerType: AnswerOptionsAnswerType
  NewAnswer: AnswerOptions
} | {
  NewAnswerType: AnswerInputsAnswerType
  NewAnswer: AnswerInputs
})

export const ChangeQuestionAnswer = ApiEndpoint<ChangeQuestionAnswerArgs, void>(async (args, api) => {
  const response = await fetch(`${api.Url}/api/question/change-answer`, {
    method: "POST",
    headers: { "Content-Type": "application-json", ...args.Session.Headers() },
    body: JSON.stringify({
      "id": args.Id,
      "new_answer_type": args.NewAnswerType,
      "new_answer": args.NewAnswer.Dto(),
    })
  })

  if (!response.ok) {
    await api.InvalidResponseHandler(response)
    return { Ok: false }
  }

  return { Ok: true, Model: undefined }
})


// delete question

export type DeleteQuestionArgs = {
  Session: Session
  Id: string
}

export const DeleteQuestion = ApiEndpoint<DeleteQuestionArgs, void>(async (args, api) => {
  const response = await fetch(`${api.Url}/api/question/delete`, {
    method: "DELETE",
    headers: { "Content-Type": "application-json", ...args.Session.Headers() },
    body: JSON.stringify({ "id": args.Id })
  })

  if (!response.ok) {
    await api.InvalidResponseHandler(response)
    return { Ok: false }
  }

  return { Ok: true, Model: undefined }
})