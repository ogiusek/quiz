import React, { useContext, useEffect } from "react";
import { WsContext } from "./wsContext";
import { NotiesContext } from "../noties/notiesContext";

export const WsNoties: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const wsContext = useContext(WsContext)
  const notiesContext = useContext(NotiesContext)
  useEffect(() => {
    const listener = (message: string) => notiesContext.AddNoty({
      Type: "error",
      Message: message
    })
    wsContext.TextListener(listener)
    return () => {
      wsContext.DelTextListener(listener)
    }
  })
  return <>
    {children}
  </>
}