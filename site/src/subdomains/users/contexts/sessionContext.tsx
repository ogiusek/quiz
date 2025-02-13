import { createContext, useEffect, useState } from "react";
import { Session } from "../models/session";

export interface SessionStorage {
  GetSession: () => (Session | void)
  SetSession: (session: Session) => any
  UnSetSession: () => void
}

export const SessionContext = createContext<SessionStorage>({
  GetSession() { throw new Error("not implemented") },
  SetSession() { throw new Error("not implemented") },
  UnSetSession() { throw new Error("not implemented") },
})

function SaveSession(session: Session) {
  const sessionToken = session?.SessionToken()
  sessionToken && localStorage.setItem("session_token", sessionToken)

  const refreshToken = session?.RefreshToken()
  refreshToken && localStorage.setItem("refresh_token", refreshToken)
}

function RetrieveSession(): Session | undefined {
  const sessionToken = localStorage.getItem('session_token')
  const refreshToken = localStorage.getItem('refresh_token')

  if (!sessionToken || !refreshToken)
    return

  const session = new Session(
    sessionToken,
    refreshToken
  )

  if (session.Valid() == 'invalid')
    return

  return session
}

function DestroySession() {
  localStorage.removeItem('session_token')
  localStorage.removeItem('refresh_token')
}

export const SessionService: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [session, setSession] = useState<Session | undefined>(RetrieveSession());

  useEffect(() => {
    session && SaveSession(session)
  }, [session])

  return <SessionContext.Provider value={{
    // TODO check do using redux help here
    GetSession: () => session,
    SetSession: (session) => { setSession(session) },
    UnSetSession: () => {
      DestroySession()
      setSession(undefined)
    }
  }}>
    {children}
  </SessionContext.Provider>
}