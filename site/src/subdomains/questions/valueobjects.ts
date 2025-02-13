import { VO } from "@/common/vo/vo"

// answer message
// answer
// answer options
// answer input
// answer inputs
// answer type
// question content
// question set name
// question set description

// react would be much better if components just were classes

// answer message

const answerMessageCannotBeEmptyError = new Error("answer cannot be empty");
const answerMessageToLongError = new Error("name cannot have more than 64 characters");

export class AnswerMessage implements VO {
  Value: string

  constructor(value: string) {
    this.Value = value
  }

  Valid(): Error[] {
    const errors: Error[] = []
    this.Value == "" && errors.push(answerMessageCannotBeEmptyError)
    this.Value.length > 64 && errors.push(answerMessageToLongError)
    return errors
  }

  Errors(): Error[] {
    return [
      answerMessageCannotBeEmptyError,
      answerMessageToLongError,
    ]
  }

}

// answer type

export type AnswerType = "o" | "i"

// answer options

export type AnswerOptionsAnswerType = "o"

const ToFewOptionsAnswersError: Error = new Error("there has to be 1 or 3 other answers")

export class AnswerOptions implements VO {
  Answers: AnswerMessage[];
  CorrectAnswer: AnswerMessage;

  constructor(answers: AnswerMessage[], correctAnswer: AnswerMessage) {
    this.Answers = answers
    this.CorrectAnswer = correctAnswer
  }

  Valid(): Error[] {
    const errors: Error[] = []
    ![1, 3].includes(this.Answers.length) && errors.push(ToFewOptionsAnswersError)
    return errors
  }

  Errors(): Error[] {
    return [
      ToFewOptionsAnswersError,
    ]
  }

  Dto(): any {
    return {
      answers: this.Answers.map(e => e.Value),
      correct_answer: this.CorrectAnswer.Value,
    }
  }

  static Decode(payload: any): AnswerOptions {
    return new AnswerOptions(
      [...payload["answers"]].map(payload => new AnswerMessage(payload)),
      new AnswerMessage(payload["correct_answer"])
    )
  }
}

// answer input

export class AnswerInput implements VO {
  Answer: AnswerMessage;
  CaseSensitive: boolean;

  constructor(answer: AnswerMessage, caseSensitive: boolean) {
    this.Answer = answer
    this.CaseSensitive = caseSensitive
  }

  Valid(): Error[] { return [] }
  Errors(): Error[] { return [] }

  Dto(): any {
    return {
      "answer": this.Answer.Value,
      "case_sensitive": this.CaseSensitive,
    }
  }
}

// answer inputs

export type AnswerInputsAnswerType = "i"

const
  ToFewInputAnswersError: Error = new Error("there have to be at least 1 correct input answer")

export class AnswerInputs implements VO {
  CorrectAnswers: AnswerInput[]

  constructor(correctAnswers: AnswerInput[]) {
    this.CorrectAnswers = correctAnswers
  }

  Valid(): Error[] {
    const errors: Error[] = []
    this.CorrectAnswers.length == 0 && errors.push(ToFewInputAnswersError)
    return errors
  }

  Errors(): Error[] {
    return [
      ToFewInputAnswersError,
    ]
  }

  Dto(): any {
    return {
      "correct_answers": this.CorrectAnswers.map(e => e.Dto())
    }
  }

  static Decode(payload: any): AnswerInputs {
    return new AnswerInputs(
      [...(payload['correct_answers'] ?? [])].map(payload => new AnswerInput(new AnswerMessage(payload["answer"]), payload['case_sensitive']))
    )
  }
}

// question content

const
  QuestionCannotBeEmptyError: Error = new Error("question cannot be empty"),
  QuestionToLongError: Error = new Error("question cannot exceed 64 characters")


export class Question implements VO {
  Value: string

  constructor(value: string) {
    this.Value = value
  }

  Valid(): Error[] {
    const errors: Error[] = []
    this.Value == "" && errors.push(QuestionCannotBeEmptyError)
    this.Value.length > 64 && errors.push(QuestionToLongError)
    return errors
  }

  Errors(): Error[] {
    return [
      QuestionCannotBeEmptyError,
      QuestionToLongError,
    ]
  }
}

// question set name

const
  QuestionSetNameCannotBeEmptyError: Error = new Error("question set name cannot be empty"),
  QuestionSetNameToLongError: Error = new Error("question set name exceeded 64 characters")

export class QuestionSetName implements VO {
  Value: string;

  constructor(value: string) {
    this.Value = value
  }

  Valid(): Error[] {
    const errors = []
    this.Value == "" && errors.push(QuestionSetNameCannotBeEmptyError)
    this.Value.length > 64 && errors.push(QuestionSetNameToLongError)
    return errors
  }

  Errors(): Error[] {
    return [
      QuestionSetNameCannotBeEmptyError,
      QuestionSetNameToLongError
    ]
  }
}

// question set description

const
  QuestionSetDescriptionCannotBeEmptyError: Error = new Error("question set description cannot be empty"),
  QuestionSetDescriptiontoLongError: Error = new Error("question set description exceeded 512 characters")

export class QuestionSetDescription implements VO {
  Value: string;
  constructor(value: string) {
    this.Value = value
  }

  Valid(): Error[] {
    const errors: Error[] = []
    this.Value == "" && errors.push(QuestionSetDescriptionCannotBeEmptyError)
    this.Value.length > 512 && errors.push(QuestionSetDescriptiontoLongError)
    return errors
  }
  Errors(): Error[] {
    return [
      QuestionSetDescriptionCannotBeEmptyError,
      QuestionSetDescriptiontoLongError,
    ]
  }
}
