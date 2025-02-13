import { useContext, useState } from "react";
import { QuestionSet as QuestionSetType } from "../models";
import { CaseSensitive, Loader2 } from "lucide-react";
import { GetQuestionSet, GetQuestionSetArgs } from "../services/api";
import { ApiContext } from "@/common/api/apiContext";
import { SessionContext } from "@/subdomains/users/contexts/sessionContext";
import { UserAvatar } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Link } from "react-router-dom";

export const PreviewQuestionSet = ({ id }: { id: string }) => {
  const [questionSet, setQuestionSet] = useState<QuestionSetType>();
  const [loading, setLoading] = useState<boolean>(true)

  const sessionContext = useContext(SessionContext)
  const session = sessionContext.GetSession()!

  const getQuestionSetArgs: GetQuestionSetArgs = {
    Session: session,
    Id: id!
  }

  const api = useContext(ApiContext)

  const RefreshQuestionSet = () => GetQuestionSet(getQuestionSetArgs, api).then(res => {
    res.Ok && setQuestionSet(res.Model)
    setLoading(false)
  })

  loading && RefreshQuestionSet()

  return <>
    <main className="flex justify-center items-center w-screen h-screen p-2 fixed top-0 left-0 z-30 bg-background">
      <div className="bg-card text-card-foreground rounded-lg shadow-lg p-8 border w-full max-w-2xl h-full flex flex-col gap-4 overflow-y-auto">
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
            <Button asChild variant="link">
              <Link to="/question-set/search">return to searching</Link>
            </Button>
          </div>
        </>}



        {!loading && questionSet && <>
          {/* browsing */}
          <h2 className="text-2xl">{questionSet.Name.Value}</h2>
          <p className="text-md">{questionSet.Description.Value}</p>

          <div className="w-full flex flex-col justify-center items-center">
            <p className="self-start">created by</p>
            <UserAvatar user={questionSet.Owner} />
          </div>

          <ul className="list-decimal flex flex-col gap-6">
            {questionSet.Questions.map((question, i) => <li key={i}>
              <h3 className="text-lg">{question.Question.Value}</h3>
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
      </div>
    </main>
  </>
}