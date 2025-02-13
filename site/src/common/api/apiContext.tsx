import { createContext, useContext, useEffect } from "react";
import { ApiDefinition } from "./api";
import { SessionContext } from "@/subdomains/users/contexts/sessionContext";
import { Refresh } from "@/subdomains/users/services/api";
import { NotiesContext } from "../noties/notiesContext";
import { WsContext } from "../ws/wsContext";

export const ApiContext = createContext<ApiDefinition>({
  Url: "",
  ErrorHandler() { throw new Error("not implemented") },
  InvalidResponseHandler() { throw new Error("not implemented") },
})

// this shouldn't be in common because of its dependencies but is for convienience
// this makes this service easier to locate
export const ApiService: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const wsContext = useContext(WsContext)
  const sessionStorage = useContext(SessionContext)
  const notiesStorage = useContext(NotiesContext)

  const api: ApiDefinition = {
    Url: `http://${window.location.hostname}:5050`,

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
        const response = await Refresh({ Session: session }, api)
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

      const ws = new WebSocket(`ws://${window.location.hostname}:5050/ws?authorization=${sessionStorage.GetSession()!.SessionToken()}`)
      return ws
    })
    wsContext.Connect()
  })

  return <ApiContext.Provider value={api}>
    {children}
  </ApiContext.Provider>
}