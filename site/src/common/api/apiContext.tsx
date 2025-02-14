import { createContext, useContext, useEffect } from "react";
import { ApiDefinition } from "./api";
import { SessionContext } from "@/subdomains/users/contexts/sessionContext";
import { Refresh } from "@/subdomains/users/services/api";
import { NotiesContext } from "../noties/notiesContext";
import { WsContext } from "../ws/wsContext";

export const ApiContext = createContext<ApiDefinition>({
  Url: "",
  WsUrl: "",
  ErrorHandler() { throw new Error("not implemented") },
  InvalidResponseHandler() { throw new Error("not implemented") },
})

type ApiEnv = {
  Url: string
  WsUrl: string
}

// this shouldn't be in common because of its dependencies but is for convienience
// this makes this service easier to locate
export const ApiService: (_: ApiEnv) => React.FC<{ children: React.ReactNode }> = ({ Url, WsUrl }) => ({ children }) => {
  const wsContext = useContext(WsContext)
  const sessionStorage = useContext(SessionContext)
  const notiesStorage = useContext(NotiesContext)

  const api: ApiDefinition = {
    Url: Url,
    WsUrl: WsUrl,

    ErrorHandler(error) {
      notiesStorage.AddNoty({
        Type: "error",
        Message: error.Message
      })
      console.error('cannot fetch', error)
    },

    async InvalidResponseHandler(response) {
      const session = sessionStorage.GetSession()

      if (response.status == 401 && session && session) {
        const response = await Refresh({ Session: session }, {
          ...api,
          async InvalidResponseHandler(_) {
            sessionStorage.UnSetSession()
            return undefined
          },
        })
        if (response.Ok) {
          sessionStorage.SetSession(response.Model)
          return { additional: { Session: response.Model } }
        }
      }

      const errorText = await response.text()
      notiesStorage.AddNoty({
        Type: "error",
        Message: errorText
      })
      console.error('invalid response', errorText)
    },
  }

  useEffect(() => {
    if (!sessionStorage.GetSession() || sessionStorage.GetSession()?.Valid() == 'invalid') {
      wsContext.Close()
      return
    }

    wsContext.Connector(async () => {
      if (sessionStorage.GetSession()?.Valid() == 'expired') {
        const res = await Refresh({ Session: sessionStorage.GetSession()! }, api)

        if (res.Ok) sessionStorage.SetSession(res.Model)
      }

      if (sessionStorage.GetSession()?.Valid() != 'valid') return

      const ws = new WebSocket(`${api.WsUrl}/ws?authorization=${sessionStorage.GetSession()!.SessionToken()}`)
      return ws
    })
    wsContext.Connect()
  }, [sessionStorage])

  return <ApiContext.Provider value={api}>
    {children}
  </ApiContext.Provider>
}