import { WsContext } from "@/common/ws/wsContext"
import { useContext } from "react"

export const ReJoinService: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const wsContext = useContext(WsContext)
  wsContext.MessageListener("match/active_match", ({ match_id }: { match_id: string }) => {
    const usesHashRouter = window.location.href.includes('/#/')
    const link = `${location.protocol}//${location.host}${usesHashRouter ? '/#' : ''}/quiz/join/${match_id}`

    window.location.assign(link)
  })

  return <>
    {children}
  </>
}