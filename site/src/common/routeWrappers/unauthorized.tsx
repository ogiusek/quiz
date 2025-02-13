import { SessionContext } from "@/subdomains/users/contexts/sessionContext"
import React, { useContext } from "react"
import { Navigate } from "react-router"
import { goToParam } from "./authorized"

export function UnAuthorized(Component: (() => React.JSX.Element) | React.MemoExoticComponent<() => React.JSX.Element>, defaultUrl: string): () => React.JSX.Element {
  return () => {
    const sessionStorage = useContext(SessionContext)
    const session = sessionStorage.GetSession()

    if (!session || session.Valid() == 'invalid')
      return <Component />

    const hash = window.location.hash
    const queryString = !hash ?
      window.location.search : hash.includes('?') ?
        hash.split('?')[1] : ""
    const queryParams = new URLSearchParams(queryString)
    const onAuthorizeGoTo = queryParams.get(goToParam)
    queryParams.delete(goToParam)
    const goTo = onAuthorizeGoTo ? onAuthorizeGoTo : defaultUrl

    return <Navigate to={`${goTo}?${queryParams.toString()}`} />
  }
}