import { useContext, useEffect, useMemo, useState } from "react"
import { MatchDto } from "../models"
import { SessionContext } from "@/subdomains/users/contexts/sessionContext"
import { Button } from "@/components/ui/button"
import { Contact, Copy, Eye, LogOut, Power, PowerOff, Save, Send, X } from "lucide-react"
import { WsContext } from "@/common/ws/wsContext"
import { QuestionSetPicker } from "@/subdomains/questions/components/QuestionSetPicker"
import { Input } from "@/components/ui/input"
import { NotiesContext } from "@/common/noties/notiesContext"
import QRCode from "react-qr-code"
import { PreviewQuestionSet } from "@/subdomains/questions/components/PreviewQuestionSet"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { AnswerInputs, AnswerOptions, AnswerType } from "@/subdomains/questions/valueobjects"
import { LeaderboardContext } from "../contexts/leaderboard"

const RenderOptionsAnswer = ({ answer, onAnswer }: { answer: AnswerOptions, onAnswer: (_: string) => void }) => {
  return <>
    <ul className="w-full h-full flex flex-row gap-2 flex-wrap">
      {answer.Answers.map((answer, i) => <li key={i} className=" w-2/5 flex-grow">
        <Button className="w-full h-full" variant="outline" onClick={() => onAnswer(answer.Value)}>{answer.Value}</Button>
      </li>)}
    </ul>
  </>
}

const RenderInputsAnswer = ({ onAnswer }: { answer: AnswerInputs, onAnswer: (_: string) => void }) => {
  const [val, setVal] = useState<string>('')
  return <>
    <form className="w-full h-full flex flex-row gap-2" onSubmit={e => {
      e.preventDefault()
      val && onAnswer(val)
    }}>
      <Input placeholder="answer" value={val} onChange={e => setVal(e.target.value)} />
      <Button type="submit" disabled={!val}>
        <Send />
      </Button>
    </form>
  </>
}

const RenderAnswer = ({ type, answer, onAnswer }: { type: AnswerType, answer: AnswerInputs | AnswerOptions, onAnswer: (_: string) => void }) => {
  switch (type) {
    case 'i': return <RenderInputsAnswer answer={answer as AnswerInputs} onAnswer={onAnswer} />
    case 'o': return <RenderOptionsAnswer answer={answer as AnswerOptions} onAnswer={onAnswer} />
  }
  throw new Error('not implemented this answer type')
}

const CountTo = ({ date }: { date: Date }) => {
  const [changes, setChanges] = useState(0)
  const now = new Date().getTime()
  const compared = useMemo(() => date.getTime() + 1000, [])
  const diff = Math.floor((compared - now) / 1000)
  setTimeout(() => setChanges(changes + 1), 250);
  return <>{diff < 0 ? 0 : diff}</>
}

export const Play = ({ match }: { match: MatchDto }) => {
  // admin
  const [questionsAmount, setQuestionsAmount] = useState<number>(match.QuestionsAmount ?? 0)
  const [choosingQuestionSet, setChoosingQuestionSet] = useState<boolean>(false)
  const [previewingQuesitonSet, setPreviewingQuestionSet] = useState<boolean>(false)

  // universal
  const [showPlayers, setShowPlayers] = useState<boolean>(false)
  const [updatedLeaderboard, setUpdatedLeaderboard] = useState<boolean>(false)

  // contexts
  const wsContext = useContext(WsContext)
  const notiesContext = useContext(NotiesContext)
  const sessionContext = useContext(SessionContext)
  const leaderboardContext = useContext(LeaderboardContext)

  const userId = sessionContext.GetSession()!.Session()!.UserId
  const usesHashRouter = window.location.href.includes('/#/')
  const link = `${location.protocol}//${location.host}${usesHashRouter ? '/#' : ''}/quiz/join/${match.Id}`

  const sortedPlayers = match.Players.sort((p1, p2) => p2.Score - p1.Score)

  useEffect(() => {
    if (match.Course?.Step == 'finished' && !updatedLeaderboard) {
      leaderboardContext.FinishedMatch(match)
      setUpdatedLeaderboard(true)
    } else if (match.Course?.Step != 'finished') {
      setUpdatedLeaderboard(false)
    }
  }, [match, updatedLeaderboard])

  useEffect(() => {
    setQuestionsAmount(match.QuestionsAmount)
  }, [match.QuestionsAmount])

  return <>
    <main className="flex justify-center items-center h-screen">
      <div className="bg-card text-card-foreground rounded-lg shadow-lg p-8 border w-full max-w-3xl h-5/6 flex flex-col gap-4">
        <h1 className="text-3xl">You're in match</h1>

        <div className={`transition fixed top-0 right-0 w-full max-w-sm h-full z-20 p-2 bg-card rounded-l-md border ${showPlayers ? '' : 'translate-x-full'}`}>
          <div className="w-full flex flex-row justify-start">
            <Button variant="destructive" onClick={() => setShowPlayers(false)}>
              <X />
            </Button>
          </div>
          <ul className="flex flex-col gap-2 mt-4">
            <li className="flex flex-row justify-between">
              <p>avatar</p>
              <p>score</p>
              <p>active</p>
            </li>
            {match.Players.map(player => <li key={player.Id} className="flex flex-row justify-between">
              <Avatar className="border flex items-center justify-center">
                <AvatarImage src={player.User.Image} />
                <AvatarFallback>{player.User.UserName.Value.slice(0, 2).toUpperCase()}</AvatarFallback>
              </Avatar>
              <p>{player.Score}</p>
              {player.Online ? <Power /> : <PowerOff />}
            </li>)}
          </ul>
        </div>

        <div className="flex flex-row justify-between">
          <Button variant="outline" onClick={() => setShowPlayers(!showPlayers)}>
            <Contact />
          </Button>
          <Button variant="destructive" disabled={match.State != 'prepare'} onClick={() => {
            wsContext.SendMessage({ topic: "match/quit", payload: {} })
          }}>
            <LogOut />
          </Button>
        </div>

        {(match.State == "prepare" || match.Course?.Step == 'finished') && <>
          <div className="flex flex-row gap-2">
            <Input value={link} onChange={() => { }} disabled />
            <Button onClick={async () => {
              try {
                await navigator.clipboard.writeText(link);
                notiesContext.AddNoty({ Type: "success", Message: "copied!" })
              } catch (error) {
                notiesContext.AddNoty({ Type: "error", Message: "failed to copy" })
              }
            }}>
              <Copy />
            </Button>
          </div>
          <div className="flex flex-row justify-center">
            <QRCode value={link} className="border-8 border-primary" />
          </div>
          {match.HostUserId == userId ? <>
            {/* admin */}
            <div className="w-full flex flex-row gap-2">
              <Button variant="outline" className="grow" onClick={() => setChoosingQuestionSet(true)}>
                Change question set
              </Button>
              <Button
                variant="outline"
                disabled={!match.QuestionSetId}
                onClick={() => setPreviewingQuestionSet(true)}>
                <Eye />
              </Button>
            </div>
            {choosingQuestionSet && <>
              <Button
                className="fixed top-5 left-5 z-20"
                variant="destructive"
                onClick={() => setChoosingQuestionSet(false)}>
                <X />
              </Button>
              <QuestionSetPicker onChoose={set => {
                setChoosingQuestionSet(false)
                wsContext.SendMessage({ topic: "match/change-question-set", payload: { question_set_id: set.Id } })
              }} />
            </>}

            {previewingQuesitonSet && <>
              <Button
                className="fixed top-5 left-5 z-20"
                variant="destructive"
                onClick={() => setPreviewingQuestionSet(false)}>
                <X />
              </Button>
              <PreviewQuestionSet id={match.QuestionSetId} />
            </>}

            <form className="flex flex-row gap-2" onSubmit={e => {
              e.preventDefault()
              wsContext.SendMessage({ topic: "match/change-questions-amount", payload: { questions_amount: questionsAmount } })
            }}>
              <Input placeholder="11" value={questionsAmount} onChange={e => {
                const amount = Number(e.target.value)
                if (!isNaN(amount)) setQuestionsAmount(amount)
                else notiesContext.AddNoty({
                  Type: "error",
                  Message: "invalid number"
                })
              }} onBlur={() => {
                wsContext.SendMessage({ topic: "match/change-questions-amount", payload: { questions_amount: questionsAmount } })
              }} />
              <Button variant={questionsAmount != match.QuestionsAmount ? "default" : "ghost"} type="submit">
                <Save />
              </Button>
            </form>

            <Button variant="default" onClick={() => wsContext.SendMessage({ topic: "match/start", payload: {} })}>
              Start game
            </Button>
          </> : <>
            {/* player */}
            <p>Host is preparing match for you</p>
          </>}
        </>}

        {match.State == 'playing' && match.Course && <>
          <p className="w-full text-end">Question: {match.Course!.CurrentQuestionIndex + 1}/{match.QuestionsAmount}</p>


          {match.Course.Step == 'question' && <>
            <h2 className="text-3xl w-full text-center">{match.Course!.CurrentQuestion!.Question.Value}</h2>
            <RenderAnswer
              type={match.Course!.CurrentQuestion!.AnswerType}
              answer={match.Course!.CurrentQuestion!.Answer}
              onAnswer={answer => {
                wsContext.SendMessage({
                  topic: "match/answer",
                  payload: { answer: answer }
                })
              }}
            />
          </>}


          {match.Course.Step == 'break' && <>
            <p className="text-3xl text-center">
              <CountTo date={new Date(match.Course.NextStep)} />
            </p>

            {match.Course.AnsweredQuestions.length !== 0 && (() => {
              const answered = match.Course!.AnsweredQuestions[match.Course.AnsweredQuestions.length - 1]
              const player = match.Players.filter(p => p.UserId == answered.UserId).at(0)
              return <>
                {player == null ? <>
                  <p className="text-3xl">nobody answered</p>
                </> : <>
                  {answered.AnsweredCorrectly ? <>
                    <p className="text-3xl text-green-500">{player.User.UserName.Value} answered correctly</p>
                  </> : <>
                    <p className="text-3xl text-red-500">{player.User.UserName.Value} answered incorrectly</p>
                  </>}
                </>}
              </>
            })()}

            <ul className="flex flex-col gap-2 mt-4">
              <li className="flex flex-row justify-between">
                <p>avatar</p>
                <p>score</p>
                <p>active</p>
              </li>
              {sortedPlayers.map(player => <li key={player.Id} className="flex flex-row justify-between">
                <Avatar className="border flex items-center justify-center">
                  <AvatarImage src={player.User.Image} />
                  <AvatarFallback>{player.User.UserName.Value.slice(0, 2).toUpperCase()}</AvatarFallback>
                </Avatar>
                <p>{player.Score}</p>
                {player.Online ? <Power /> : <PowerOff />}
              </li>)}
            </ul>
          </>}


        </>}
      </div>
    </main>
  </>
}