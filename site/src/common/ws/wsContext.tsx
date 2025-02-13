import { createContext } from "react"

type ConnectListener = () => void
type CloseListener = () => void
type MessageListener = (_: any) => void
type TextListener = (_: any) => void

type MessageTopic = string
type MessagePayload = any
type Message = {
  topic: MessageTopic
  payload: MessagePayload
}

export interface WsDefinition {
  IsConnected(): boolean
  Connector(ws: () => Promise<WebSocket | void>): void
  Connect(): void
  ConnectListener(_: ConnectListener): void

  MessageListener(topic: string, listener: MessageListener): void
  DelMessageListener(topic: string, listener: MessageListener): void
  TextListener(listener: TextListener): void
  DelTextListener(listener: TextListener): void
  SendMessage(message: Message): void

  CloseListener(listener: CloseListener): void
  DelCloseListener(listener: CloseListener): void
  Close(): void
}

export const WsContext = createContext<WsDefinition>({
  IsConnected() { throw new Error("not implemented") },
  Connector: function (_: () => Promise<WebSocket>): void { throw new Error("Function not implemented.") },
  Connect: function () { throw new Error("Function not implemented.") },
  ConnectListener(_) { throw new Error("Function not implemented.") },
  MessageListener: function (_: string, __: MessageListener): void { throw new Error("Function not implemented.") },
  DelMessageListener: function (_: string, __: MessageListener): void { throw new Error("Function not implemented.") },
  TextListener: function (_: TextListener): void { throw new Error("Function not implemented.") },
  DelTextListener: function (_: TextListener): void { throw new Error("Function not implemented.") },
  SendMessage: function (_: Message): void { throw new Error("Function not implemented.") },
  CloseListener: function (_: CloseListener): void { throw new Error("Function not implemented.") },
  DelCloseListener: function (_: CloseListener): void { throw new Error("Function not implemented.") },
  Close: function (): void { throw new Error("Function not implemented.") },
})

let ws: WebSocket | void
let connector: (() => Promise<WebSocket>) | void
let messageListeners: { [key: string]: MessageListener[] } = {}
let textListeners: TextListener[] = []
let connectListeners: CloseListener[] = []
let closeListeners: CloseListener[] = []

export const WsService: React.FC<{ children: React.ReactNode }> = ({ children }) => {

  const wsDef: WsDefinition = {
    IsConnected: () => (!!ws && ws.readyState == ws.OPEN),
    Connector: (newConnector) => {
      connector = async () => {
        ws?.close()
        ws = await newConnector()
        if (!ws) {
          throw new Error("cannot connect to websocket")
        }
        wsDef.IsConnected() ?
          connectListeners = connectListeners.filter(listener => { listener() }) :
          ws.onopen = () => connectListeners = connectListeners.filter(listener => { listener() })

        ws.onclose = () => closeListeners.map(listener => listener())
        ws.onmessage = e => {
          const raw = e.data
          try {
            const json: Message = JSON.parse(raw)
            setTimeout(() => messageListeners[json.topic]?.map(listener => listener(json.payload)))
          } catch (_) {
            textListeners.map(listener => listener(raw))
          }
        }
        return ws
      }
    },
    Connect: () => {
      if (wsDef.IsConnected()) return
      if (!connector)
        throw new Error("missing connector")
      connector()
    },
    ConnectListener: (listener) => { wsDef.IsConnected() ? listener() : connectListeners.push(listener) },
    MessageListener: (topic, listener) => (messageListeners[topic] = [...(messageListeners[topic] ?? []), listener]),
    DelMessageListener: (topic, listener) => (messageListeners[topic] = messageListeners[topic].filter(l => l != listener)),
    TextListener: (listener) => textListeners.push(listener),
    DelTextListener: (listener) => (textListeners = textListeners.filter(l => l != listener)),
    SendMessage: (message) => {
      if (wsDef.IsConnected()) {
        ws?.send(JSON.stringify(message))
      } else if (connector) {
        wsDef.ConnectListener(() => {
          ws?.send(JSON.stringify(message))
        })
        connector()
      } else {
        throw new Error("connector is missing")
      }
    },
    CloseListener: (listener) => closeListeners.push(listener),
    DelCloseListener: (listener) => (closeListeners = closeListeners.filter(l => l != listener)),
    Close: () => ws?.close()
  }

  return <>
    <WsContext.Provider value={wsDef}>
      {children}
    </WsContext.Provider>
  </>
}