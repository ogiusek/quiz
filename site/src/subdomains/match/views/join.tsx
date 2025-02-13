import { Nav } from "@/common/nav/nav"
import { Loader2 } from "lucide-react"
import { useContext, useEffect, useState } from "react"
import { Input } from "@/components/ui/input"
import { Link, useParams } from "react-router-dom"
import { Button } from "@/components/ui/button"
import { WsContext } from "@/common/ws/wsContext"

export const JoinId = () => {
  const { id } = useParams()
  const wsContext = useContext(WsContext)
  const [found, setFound] = useState<'not sent' | 'loading' | 'found' | 'not found'>('not sent')

  useEffect(() => {
    if (found == 'not sent') {
      wsContext.ConnectListener(() => {
        wsContext.SendMessage({ topic: "match/join", payload: { match_id: id } })
        setFound('loading')
      })
    }

    if (found == 'loading') {
      const joinListener = () => found == 'loading' && setFound('found')
      const errListener = () => found == 'loading' && setFound('not found')

      wsContext.TextListener(errListener)
      wsContext.MessageListener('match/created_match', joinListener)
      return () => {
        wsContext.DelTextListener(errListener)
        wsContext.DelMessageListener('match/created_match', joinListener)
      }
    }
  }, [found])

  useEffect(() => setFound('not sent'), [id])

  return <>
    <Nav />
    <main className="flex justify-center items-center h-screen p-2">
      <div className="bg-card text-card-foreground rounded-lg shadow-lg p-8 border w-full max-w-2xl h-full flex flex-col justify-center items-center gap-4">
        {(found == 'loading' || found == 'not sent') && <Loader2 className="animate-spin" size={240} />}
        {found == 'not found' && <>
          <h1 className="text-3xl">Match with this id do not exist</h1>
        </>}
      </div>
    </main>
  </>
}

export const Join = () => {
  const [matchId, setMatchId] = useState<string>("")
  return <>
    <Nav />
    <main className="flex justify-center items-center h-screen p-2">
      <div className="bg-card text-card-foreground w-full max-w-lg rounded-lg shadow-lg p-8 border flex flex-col gap-4">
        <h1 className="text-2xl">Join match</h1>
        <Input value={matchId} onChange={e => setMatchId(e.target.value)} placeholder="eecca86c-5c89-4eb0-b6e6-cc611f4c8992" />
        <Button asChild variant="outline">
          <Link to={`/quiz/join/${matchId}`}>Join</Link>
        </Button>
      </div>
    </main>
  </>
}