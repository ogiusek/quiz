import { createContext, useEffect, useReducer } from "react"
import { DisplayNoties } from "./display"
import { FullNoty, Noty, NotyDispatcher, NotyDispatcherAction, NotyDispatcherState } from "./dispatcher"

export type { Noty }

export interface NotiesStorage {
  AddNoty(noty: Noty): void
}

export const NotiesContext = createContext<NotiesStorage>({
  AddNoty(_) { throw new Error("not implemented") },
});

type plan = (state: NotyDispatcherState, dispatcher: React.Dispatch<NotyDispatcherAction>) => void

const DefaultNotyLifeTimeInSeconds = 30

export const NotiesService: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  // 'plans' exist to because otherwise setTimeout does not work (setTimeout(() => {/* this would not work */}))
  const [plans, dispachPlans] = useReducer(
    (
      state: plan[],
      action: { action: 'add', payload: plan } | { action: 'clean' }
    ) => {
      if (action.action == 'clean') return []
      return [...state, action.payload]
    }, [])

  const [noties, dispatchNoties] = useReducer(NotyDispatcher, {})

  const close = (id: string, _: NotyDispatcherState, dispatchNoties: React.Dispatch<NotyDispatcherAction>) => {
    dispatchNoties({
      action: 'close',
      payload: { id: id }
    })

    setTimeout(() => {
      dispatchNoties({
        action: 'remove',
        payload: { id: id }
      })
    }, 150)
  }

  const planClose = (id: string, noties: NotyDispatcherState, dispatchNoties: React.Dispatch<NotyDispatcherAction>) => {
    const noty = noties[id]
    if (!noty || noty.LifeTimeInSeconds == "forever") return

    setTimeout(() => {
      const now = new Date()
      const noty = noties[id]
      if (noty.ChangedAt == 'now' || noty.LifeTimeInSeconds == 'forever') return

      const then = new Date(noty.ChangedAt)
      then.setSeconds(noty.ChangedAt.getSeconds() + noty.LifeTimeInSeconds)
      if (then > now) return

      close(id, noties, dispatchNoties)
    }, noty.LifeTimeInSeconds * 1000);
  }

  const mouseEnter = (id: string, _: NotyDispatcherState, dispatchNoties: React.Dispatch<NotyDispatcherAction>) => {
    dispatchNoties({
      action: 'mouse_enter',
      payload: { id: id }
    })
  }

  const mouseLeave = (id: string, _: NotyDispatcherState, dispatchNoties: React.Dispatch<NotyDispatcherAction>) => {
    dispatchNoties({
      action: 'mouse_leave',
      payload: { id: id }

    })
    dispachPlans({
      action: 'add',
      payload(noties, dispatchNoties) {
        setTimeout(() => planClose(id, noties, dispatchNoties), 100)
      },
    })
  }

  useEffect(() => {
    if (plans.length != 0) {
      plans.map(plan => plan(noties, dispatchNoties))
      dispachPlans({ action: 'clean' })
    }
  })


  return <NotiesContext.Provider value={{
    AddNoty(noty) {
      const now = new Date()
      const id = `${now.toUTCString()}-${Math.floor(Math.random() * 10000)}`
      const newNoty: FullNoty = {
        Type: 'noty',
        CreatedAt: now,
        ChangedAt: now,
        LifeTimeInSeconds: DefaultNotyLifeTimeInSeconds,
        Hide: true,
        ...noty,
      }
      dispatchNoties({
        action: 'add_noty', payload: {
          id: id,
          noty: newNoty
        }
      })
      dispachPlans({
        action: 'add',
        payload(noties, dispatchNoties) {
          dispatchNoties({ action: 'show_noties' })
          planClose(id, noties, dispatchNoties)
        },
      })
    },
  }}>
    <DisplayNoties noties={noties}
      close={id => close(id, noties, dispatchNoties)}
      mouseEnter={id => mouseEnter(id, noties, dispatchNoties)}
      mouseLeave={id => mouseLeave(id, noties, dispatchNoties)} />
    {children}
  </NotiesContext.Provider >
}