import { useContext, useEffect, useState } from "react"
import { MatchDto } from "../models"
import { NewMatchChangesListener } from "../services/match"
import { WsContext } from "@/common/ws/wsContext"
import { Play } from "../views/play"

const MatchChangesListener = NewMatchChangesListener()

export const MatchService: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const wsContext = useContext(WsContext)
  const [match, setMatch] = useState<MatchDto | undefined>()

  useEffect(() => {
    MatchChangesListener.SetWs(wsContext)
    MatchChangesListener.SetState(match)
    MatchChangesListener.SetSetter(setMatch)
  }, [match, setMatch])

  return <>
    {!match && children}
    {match && <Play match={match} />}
  </>
}