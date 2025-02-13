import { ApiContext } from "@/common/api/apiContext"
import { SessionContext } from "@/subdomains/users/contexts/sessionContext"
import { useContext, useState } from "react"
import { QuestionSet } from "../models"
import { CreateQuestionSet, CreateQuestionSetArgs, DeleteQuestionSet, SearchQuestionSets } from "../services/api"
import { useBottomScrollListener } from 'react-bottom-scroll-listener'
import { Link } from "react-router-dom"
import { UserAvatar } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import { Plus, Trash } from "lucide-react"
import { ShowVoErrors } from "@/components/ui/errors"
import { QuestionSetDescription, QuestionSetName } from "../valueobjects"
import { Label } from "@/components/ui/label"
import { Input } from "@/components/ui/input"
import { Name } from "@/subdomains/users/valueobjects"
import { NotiesContext } from "@/common/noties/notiesContext"
import { Nav } from "@/common/nav/nav"

export function MyQuestionSets() {
  const api = useContext(ApiContext)
  const sessionContext = useContext(SessionContext)
  const noties = useContext(NotiesContext)
  const session = sessionContext.GetSession()!

  const [sets, setSets] = useState<QuestionSet[]>([])
  const [date, setDate] = useState<string>()
  const [page, setPage] = useState(0)
  const [end, setEnd] = useState(false)

  const fetchSets = () => {
    if (end) return
    SearchQuestionSets({
      Session: session!,
      Page: page,
      OwnerId: session.Session()?.UserId,
      LastUpdate: date == "" ? undefined : date
    }, api).then(res => {
      if (!res.Ok) return
      setPage(page + 1)
      setSets([...sets, ...res.Model.Found])
      setEnd(res.Model.Found.length == 0)
      date != res.Model.Time && setDate(res.Model.Time)
    })
  }

  page == 0 && fetchSets()

  const onReachBottom = () => page != 0 && fetchSets()
  const ref = useBottomScrollListener(onReachBottom)

  const [newSet, setNewSet] = useState<CreateQuestionSetArgs>({
    Session: session,
    Name: new QuestionSetName(""),
    Description: new QuestionSetDescription("")
  })

  // fetch
  // allow to add new
  // redirect to selected sets

  return <>
    <Nav />
    <main className="flex justify-center items-center h-screen p-2">
      <div ref={ref as React.RefObject<HTMLDivElement> | undefined} className="bg-card text-card-foreground rounded-lg shadow-lg p-8 border w-full max-w-2xl h-full flex flex-col gap-4 overflow-y-auto">
        <h1 className="text-3xl">My Question sets</h1>
        <div className="w-full flex flex-col gap-2 p-2 border rounded-md">
          <h2 className="text-xl">New set</h2>
          <ShowVoErrors vo={newSet.Name}>
            <Label htmlFor="name">Name</Label>
            <Input name="name" type="text" placeholder="zodiac signs" value={newSet.Name.Value} onChange={e => {
              newSet.Name.Value = e.target.value
              setNewSet({ ...newSet })
            }} />
          </ShowVoErrors>

          <ShowVoErrors vo={newSet.Description}>
            <Label htmlFor="desc">Description</Label>
            <Input name="desc" type="text" placeholder="are you sure you know everything about zodiac signs ?" value={newSet.Description.Value} onChange={e => {
              newSet.Description.Value = e.target.value
              setNewSet({ ...newSet })
            }} />
          </ShowVoErrors>

          <Button variant="secondary" onClick={() => {
            CreateQuestionSet({ ...newSet, Session: session, }, api).then(res => {
              if (!res.Ok) return
              const userSession = session.Session()!
              setSets([
                {
                  Id: res.Model.Id,
                  Name: newSet.Name,
                  Description: newSet.Description,
                  Questions: [],
                  Owner: {
                    Id: userSession.UserId,
                    UserName: new Name(userSession.UserName),
                    Image: userSession.UserImage,
                  }
                },
                ...sets,
              ])
              setNewSet({
                Session: session,
                Name: new QuestionSetName(""),
                Description: new QuestionSetDescription(""),
              })
            })
          }}>
            <Plus />
          </Button>
        </div>
        <ul className="flex flex-col gap-2">
          {sets.length == 0 && <>
            <h2 className="text-2xl">No quiz found create new</h2>
          </>}
          {sets.map((set, setIndex) => <li key={set.Id} className="border p-2 rounded-md">
            <Link to={`/question-set/get/${set.Id}`}>
              <div className="w-full flex flex-row justify-between">
                <h3 className="text-2xl">{set.Name.Value}</h3>

                <UserAvatar user={set.Owner} />
              </div>
              <p>{set.Description.Value}</p>
            </Link>
            <div className="flex flex-col items-end">
              <Button variant="destructive" onClick={() => {
                DeleteQuestionSet({
                  Session: session,
                  Id: set.Id
                }, api).then(res => {
                  if (!res.Ok) return

                  sets.splice(setIndex, 1)
                  setSets([...sets])
                  noties.AddNoty({
                    Type: "noty",
                    Message: "deleted set",
                  })
                })
              }}>
                <Trash />
              </Button>
            </div>
          </li>)}
        </ul>
      </div>
    </main>
  </>
}