import { useContext, useState } from "react";
import { ApiContext } from "@/common/api/apiContext";
import { NotiesContext } from "@/common/noties/notiesContext";
import { SessionContext } from "@/subdomains/users/contexts/sessionContext";
import { Link, useParams } from "react-router-dom";
import type { AnswerWithType, NewQuestion, QuestionSet as QuestionSetType } from "../models";
import { ChangeQuestion, ChangeQuestionAnswer, ChangeQuestionSetDescription, ChangeQuestionSetName, CreateQuestion, DeleteQuestion, DeleteQuestionSet, GetQuestionSet, GetQuestionSetArgs } from "../services/api";
import { CaseSensitive, Loader2, Plus, Save, Trash } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Toggle } from "@/components/ui/toggle";
import { UserAvatar } from "@/components/ui/avatar";
import { ShowVoErrors } from "@/components/ui/errors";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectGroup, SelectItem, SelectLabel, SelectTrigger, SelectValue } from "@/components/ui/select";
import { AnswerInput, AnswerInputs, AnswerMessage, AnswerOptions, AnswerType, Question as QuestionVO } from "../valueobjects";
import { Checkbox } from "@/components/ui/checkbox";
import { Nav } from "@/common/nav/nav";

function AnswerInputsComponent({ value, onChange }: { value: AnswerInputs, onChange: (_: AnswerInputs) => void }) {
  return <>
    <ul>
      {value.CorrectAnswers.map((input, inputIndex) => <li key={inputIndex} className="flex flex-row justify-between">
        <ShowVoErrors vo={input.Answer}>
          <Label htmlFor={`input-answer-${inputIndex}`}>answer</Label>
          <Input name={`input-answer-${inputIndex}`} type="text" placeholder="physics quiz" value={input.Answer.Value} onChange={e => {
            value.CorrectAnswers[inputIndex].Answer.Value = e.target.value
            onChange(value)
          }} />
        </ShowVoErrors>
        <div className="flex flex-row mt-auto gap-2 items-center">
          <Checkbox checked={input.CaseSensitive} onCheckedChange={e => {
            value.CorrectAnswers[inputIndex].CaseSensitive = e == true
            onChange(value)
          }} name={`input-answer-case-${inputIndex}`} />
          <Label htmlFor={`input-answer-case-${inputIndex}`}>case sentive</Label>

          <Button aria-label="remove" variant="outline" onClick={() => {
            value.CorrectAnswers.splice(inputIndex, 1)
            onChange(value)
          }} >
            <Trash />
          </Button>
        </div>
      </li>)}
    </ul>
    <br />
    <div className="w-full flex flex-col items-end">
      <Button aria-label="add" variant="outline" onClick={() => {
        value.CorrectAnswers.push(new AnswerInput(new AnswerMessage(""), false))
        onChange(value)
      }}>
        <Plus />
      </Button>
    </div>
  </>
}

function AnswerOptionsComponent({ value, onChange }: { value: AnswerOptions, onChange: (_: AnswerOptions) => void }) {
  return <>
    <ul>
      <li className="text-green-500">
        <ShowVoErrors vo={value.CorrectAnswer}>
          <Label htmlFor={`option-answer`}>correct option</Label>
          <Input name={`options-answer`} type="text" placeholder="physics quiz" value={value.CorrectAnswer.Value} onChange={e => {
            value.CorrectAnswer.Value = e.target.value
            onChange(value)
          }} />
        </ShowVoErrors>
      </li>
      {value.Answers.map((option, optionIndex) => <li key={optionIndex}>
        <ShowVoErrors vo={option}>
          <Label htmlFor={`option-answer-${optionIndex}`}>incorrect option</Label>
          <div className="flex flex-row gap-2">
            <Input name={`option-answer-${optionIndex}`} type="text" placeholder="physics quiz" value={option.Value} onChange={e => {
              value.Answers[optionIndex].Value = e.target.value;
              onChange(value)
            }} />
            <Button aria-label="remove" variant="outline" onClick={() => {
              value.Answers.splice(optionIndex, 1)
              onChange(value)
            }} >
              <Trash />
            </Button>
          </div>
        </ShowVoErrors>
      </li>)}
    </ul>
    <br />
    <div className="w-full flex flex-col items-end">
      <Button aria-label="add" variant="outline" onClick={() => {
        value.Answers.push(new AnswerMessage(""))
        onChange(value)
      }}>
        <Plus />
      </Button>
    </div>
  </>
}

function AnswerTypeComponent({ value, onChange }: { value: AnswerType, onChange: (_: AnswerType) => void }) {
  return <div className="flex items-center justify-center">
    <Select value={value} onValueChange={(e: AnswerType) => onChange(e)}>
      <SelectTrigger className="w-[180px]">
        <SelectValue placeholder="Select answer type" />
        <SelectValue />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectLabel>Select answer type</SelectLabel>
          <SelectItem value="i">input</SelectItem>
          <SelectItem value="o">option</SelectItem>
        </SelectGroup>
      </SelectContent>
    </Select>
  </div>
}

function AnswerWithTypeComponent({ value, onChange }: { value: AnswerWithType, onChange: (_: AnswerWithType) => void }) {
  return <>
    <AnswerTypeComponent value={value.AnswerType} onChange={(answerType) => {
      value.AnswerType = answerType
      answerType == "i" && (value.Answer = new AnswerInputs([new AnswerInput(new AnswerMessage(""), true)]))
      answerType == "o" && (value.Answer = new AnswerOptions([], new AnswerMessage("")))
      onChange(value)
    }} />
    {value.AnswerType == 'i' && <>
      <AnswerInputsComponent value={value.Answer} onChange={(v) => {
        (value.Answer as AnswerInputs) = v
        onChange(value)
      }} />
    </>}
    {value.AnswerType == 'o' && <>
      <AnswerOptionsComponent value={value.Answer} onChange={(v) => {
        (value.Answer as AnswerOptions) = v
        onChange(value)
      }} />
    </>}
  </>
}

const getDefaultNewQuestion: (_: string) => NewQuestion = (questionSetId: string) => ({
  QuestionSetId: questionSetId,
  Question: new QuestionVO(""),
  AnswerType: "i",
  Answer: new AnswerInputs([]),
})

// i hoped memoizing this would remove site refresh on session refresh but react does not support handling state (uses fp)
// in my opinion framework should use oop and components should be like value objects (identical component for the same state)
// i do not see here space to manage state because fp lacks it 
export const QuestionSet = () => {
  const { id } = useParams();
  const [questionSet, setQuestionSet] = useState<QuestionSetType>();
  const [newQuestion, setNewQuestion] = useState<NewQuestion>(getDefaultNewQuestion(id!));
  const [loading, setLoading] = useState<boolean>(true)
  const [modifies, setModifies] = useState<boolean>(false)

  const sessionContext = useContext(SessionContext)
  const session = sessionContext.GetSession()!
  const owns = questionSet && questionSet.Owner.Id == session.Session()?.UserId

  const getQuestionSetArgs: GetQuestionSetArgs = {
    Session: session,
    Id: id!
  }

  const api = useContext(ApiContext)
  const noties = useContext(NotiesContext)

  const RefreshQuestionSet = () => GetQuestionSet(getQuestionSetArgs, api).then(res => {
    res.Ok && setQuestionSet(res.Model)
    setLoading(false)
  })

  loading && RefreshQuestionSet()

  // a lot of js is nested because there are many components combined and if whole js would be up here this could quicly become mess (more then it already is)
  // if you think why there are many components in one, answer is that they aren't used anywhere else so extraction is not necessary
  return <>
    <Nav />
    <main style={{ maxWidth: "100vw" }} className="flex justify-center items-center h-screen p-2 pt-12">
      <div className="bg-card text-card-foreground rounded-lg shadow-lg p-2 border w-full max-w-2xl h-full flex flex-col gap-4 overflow-y-auto">
        {loading && <>
          {/* loader */}
          <div className="flex justify-center items-center w-full h-full">
            <Loader2 className="animate-spin" size={240} />
          </div>
        </>}



        {!loading && !questionSet && <>
          {/* 404 */}
          <h2 className="text-6xl text-center">404</h2>
          <h1 className="text-3xl text-center">Not found</h1>
          <div className="w-full">
            <Button aria-label="return to searching link" asChild variant="link">
              <Link to="/question-set/search">return to searching</Link>
            </Button>
          </div>
        </>}



        {!loading && questionSet && owns && <>
          {/* toggler */}
          <div className="flex flex-row w-full justify-between">
            <Button aria-label="delete question set" variant="destructive" onClick={() => {
              DeleteQuestionSet({
                Session: session,
                Id: id!
              }, api).then(res => {
                if (!res.Ok) return
                noties.AddNoty({
                  Type: "noty",
                  Message: "succesfuly deleted this page do not exist anymore",
                })
              })
            }}>
              <Trash />
            </Button>
            <Toggle className="w-max ml-auto" onClick={_ => setModifies(!modifies)}>modify</Toggle>
          </div>
        </>}



        {!loading && questionSet && !modifies && <>
          {/* browsing */}
          <h2 className="text-2xl">{questionSet.Name.Value}</h2>
          <p className="text-md">{questionSet.Description.Value}</p>

          <div className="w-full flex flex-col justify-center items-center">
            <p className="self-start">created by</p>
            <UserAvatar user={questionSet.Owner} />
          </div>

          <ul className="flex flex-col gap-6">
            {questionSet.Questions.map((question, i) => <li key={i}>
              <h3 className="text-lg">{i + 1}. {question.Question.Value}</h3>
              {question.AnswerType == 'i' && <>
                <ul>
                  {question.Answer.CorrectAnswers.map((input, i) => <li key={i} className="flex flex-row justify-between">
                    <h4>{input.Answer.Value}</h4>
                    {input.CaseSensitive && <CaseSensitive />}
                  </li>)}
                </ul>
              </>}
              {question.AnswerType == 'o' && <>
                <ul>
                  <li className="text-green-500">{question.Answer.CorrectAnswer.Value}</li>
                  {question.Answer.Answers.map((option, i) => <li key={i} className="text-red-500">
                    {option.Value}
                  </li>)}
                </ul>
              </>}
            </li>)}
          </ul>
        </>}



        {!loading && questionSet && modifies && <>
          {/* modifying */}

          <ShowVoErrors vo={questionSet.Name}>
            <Label htmlFor="name">Name</Label>
            <div className="flex flex-row gap-2 justify-center items-center">
              <Input name="name" type="text" placeholder="John" value={questionSet.Name.Value} onChange={e => {
                questionSet.Name.Value = e.target.value
                setQuestionSet({ ...questionSet })
              }} />
              <Button aria-label="change question set name" variant="outline" disabled={questionSet.Name.Valid().length != 0}
                onClick={async () => {
                  const res = await ChangeQuestionSetName({
                    Session: session,
                    Id: questionSet.Id,
                    NewName: questionSet.Name,
                  }, api)
                  if (!res.Ok) return
                  noties.AddNoty({ Type: "success", Message: "Saved" })
                }}
              >
                <Save />
              </Button>
            </div>
          </ShowVoErrors>

          <ShowVoErrors vo={questionSet.Description}>
            <Label htmlFor="desc">Description</Label>
            <div className="flex flex-row gap-2 justify-center items-center">
              <Input name="desc" type="text" placeholder="physics quiz" value={questionSet.Description.Value} onChange={e => {
                questionSet.Description.Value = e.target.value
                setQuestionSet({ ...questionSet })
              }} />
              <Button aria-label="Change question set description" variant="outline" disabled={questionSet.Name.Valid().length != 0}
                onClick={async () => {
                  const res = await ChangeQuestionSetDescription({
                    Session: session,
                    Id: questionSet.Id,
                    NewDescription: questionSet.Description,
                  }, api)
                  if (!res.Ok) return
                  noties.AddNoty({ Type: "success", Message: "Saved" })
                }}
              >
                <Save />
              </Button>
            </div>
          </ShowVoErrors>

          <ul className="flex flex-col gap-6">
            {questionSet.Questions.map((question, questionIndex) => <li key={questionIndex} className="border p-2 rounded-sm">
              <div className="flex flex-row justify-between">
                <h2 className="text-xl">{questionIndex + 1}</h2>
                <Button aria-label="remove question" variant="destructive" onClick={async () => {
                  const res = await DeleteQuestion({
                    Session: session,
                    Id: question.Id,
                  }, api)
                  if (res.Ok) {
                    setQuestionSet({ ...questionSet, Questions: questionSet.Questions.filter(q => q.Id != question.Id) })
                  }
                }}>
                  <Trash />
                </Button>
              </div>
              <ShowVoErrors vo={question.Question}>
                <Label htmlFor="question">Question</Label>
                <br />
                <div className="flex flex-row gap-2 justify-center items-center">
                  <Input name="question" type="text" placeholder="How many atoms are in the universe" value={question.Question.Value} onChange={e => {
                    questionSet.Questions[questionIndex].Question.Value = e.target.value
                    setQuestionSet({ ...questionSet })
                  }} />
                  <Button aria-label="change question question" variant="secondary" disabled={questionSet.Name.Valid().length != 0}
                    onClick={async () => {
                      const res = await ChangeQuestion({
                        Session: session,
                        Id: question.Id,
                        NewQuestion: question.Question,
                      }, api)
                      if (!res.Ok) return
                      noties.AddNoty({ Type: "success", Message: "Saved" })
                    }}
                  >
                    <Save />
                  </Button>
                </div>
              </ShowVoErrors>
              <br />
              <ShowVoErrors vo={question.Answer}>
                <AnswerWithTypeComponent value={question} onChange={(answer) => {
                  questionSet.Questions[questionIndex].AnswerType = answer.AnswerType
                  questionSet.Questions[questionIndex].Answer = answer.Answer
                  setQuestionSet({ ...questionSet })
                }} />
                <br />
                <div className="w-full flex flex-col items-end">
                  <Button aria-label="save" variant="secondary" disabled={question.Answer.Valid().length != 0}
                    onClick={async () => { // @ts-ignore
                      const res = await ChangeQuestionAnswer({
                        Session: session,
                        Id: question.Id,
                        NewAnswerType: question.AnswerType,
                        NewAnswer: question.Answer,
                      }, api)
                      if (!res.Ok) return
                      noties.AddNoty({ Type: "success", Message: "Saved" })
                    }}
                  >
                    <Save />
                  </Button>
                </div>
              </ShowVoErrors>
            </li>)}
          </ul>
          <div className="w-full flex flex-col gap-2 border p-2 rounded-md">
            <ShowVoErrors vo={newQuestion.Question}>
              <Label htmlFor="question">Question</Label>
              {/* <div className="flex flex-row gap-2 justify-center items-center"> */}
              <Input name="question" type="text" placeholder="How many atoms are in the universe" value={newQuestion.Question.Value} onChange={e => {
                newQuestion.Question.Value = e.target.value
                setNewQuestion({ ...newQuestion })
              }} />
              {/* </div> */}
            </ShowVoErrors>
            <ShowVoErrors vo={newQuestion.Answer}>
              <AnswerWithTypeComponent value={newQuestion} onChange={(value) => {
                newQuestion.Answer = value.Answer
                newQuestion.AnswerType = value.AnswerType
                setNewQuestion({ ...newQuestion })
              }} />
            </ShowVoErrors>
            <Button aria-label="create question" variant="secondary" disabled={[...newQuestion.Question.Valid(), ...newQuestion.Answer.Valid()].length != 0} onClick={() => {
              CreateQuestion({
                Session: session,
                ...newQuestion
              }, api).then(res => {
                if (!res.Ok) return

                questionSet.Questions.push({
                  Id: res.Model.Id,
                  ...newQuestion
                })
                setQuestionSet({ ...questionSet })
                setNewQuestion(getDefaultNewQuestion(id!))
              })
            }}>
              <Plus />
            </Button>
          </div>
        </>}
      </div>
    </main>
  </>
}