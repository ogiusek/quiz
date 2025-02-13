import { useBottomScrollListener } from "react-bottom-scroll-listener";
import { QuestionSet } from "../models";
import { SearchQuestionSets } from "../services/api";
import { useContext, useState } from "react";
import { SessionContext } from "@/subdomains/users/contexts/sessionContext";
import { ApiContext } from "@/common/api/apiContext";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Loader2, SearchIcon } from "lucide-react";
import { UserAvatar } from "@/components/ui/avatar";

export const QuestionSetPicker = ({ onChoose: choose }: { onChoose: (_: QuestionSet) => void }) => {
  const api = useContext(ApiContext)
  const sessionContext = useContext(SessionContext)
  const session = sessionContext.GetSession()!

  const [sets, setSets] = useState<QuestionSet[]>([])
  const [currentSearch, setCurrentSearch] = useState("")
  const [search, setSearch] = useState("")
  const [date, setDate] = useState<string | undefined>()
  const [page, setPage] = useState(0)
  const [loading, setLoading] = useState(true)
  const [end, setEnd] = useState(false)

  const fetchSets = () => {
    if (end) return
    !loading && setLoading(true)
    SearchQuestionSets({
      Session: session!,
      Search: currentSearch,
      Page: page,
      LastUpdate: date == "" ? undefined : date
    }, api).then(res => {
      if (!res.Ok) return
      setLoading(false)
      setPage(page + 1)
      setSets([...sets, ...res.Model.Found])
      setEnd(res.Model.Found.length == 0)
      date != res.Model.Time && setDate(res.Model.Time)
    })
  }

  page == 0 && fetchSets()

  const onReachBottom = () => page != 0 && fetchSets()
  const ref = useBottomScrollListener(onReachBottom)

  return <>
    <main className="flex justify-center items-center w-screen h-screen p-2 fixed top-0 left-0 z-30 bg-background">
      <div ref={ref as React.RefObject<HTMLDivElement> | undefined} className="bg-card text-card-foreground rounded-lg shadow-lg p-8 border w-full max-w-2xl h-full flex flex-col gap-4 overflow-y-auto">
        <h1 className="text-3xl">Choose Set</h1>
        <form className="flex flex-row gap-2" onSubmit={(e) => {
          e.preventDefault()
          setSets([])
          setDate(undefined)
          setEnd(false)
          setCurrentSearch(search)
          setPage(0)
        }}>
          <Input value={search} placeholder="search" onChange={e => setSearch(e.target.value)} />
          <Button variant="ghost" type="submit">
            <SearchIcon />
          </Button>
        </form>

        <ul className="flex flex-col gap-2">
          {sets.length == 0 && end && <>
            <h2 className="text-2xl">No quiz found create new</h2>
          </>}
          {sets.map((set) => <li key={set.Id} className="border p-2 rounded-md">
            <button className="w-full h-full" onClick={() => choose(set)}>
              <div className="w-full flex flex-row justify-between">
                <h3 className="text-2xl">{set.Name.Value}</h3>
                <UserAvatar user={set.Owner} />
              </div>
              <p className="w-full text-left">{set.Description.Value}</p>
            </button>
          </li>)}
          {loading && <>
            <div className="flex justify-center items-center">
              <Loader2 className="animate-spin" size={240} />
            </div>
          </>}
        </ul>
      </div>
    </main>
  </>
}