import { SessionContext } from "@/subdomains/users/contexts/sessionContext";
import React, { useContext } from "react";
import { Navigate, useLocation } from "react-router";

export const goToParam = 'on_authorize_go_to'

export function Authorized(Component: (() => React.JSX.Element) | React.MemoExoticComponent<() => React.JSX.Element>, defaultUrl: string): () => React.JSX.Element {
  return () => {
    const sessionStorage = useContext(SessionContext)
    const session = sessionStorage.GetSession()

    if (session && session.Valid() != 'invalid')
      return <Component />

    const location = useLocation()

    const hash = window.location.hash
    const queryString = !hash ?
      window.location.search : hash.includes('?') ?
        hash.split('?')[1] : ""
    const queryParams = new URLSearchParams(queryString)
    queryParams.set(goToParam, location.pathname)

    return <Navigate to={`${defaultUrl}?${queryParams.toString()}`} />
  }
}