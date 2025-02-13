export type NotyType = 'error' | 'noty' | 'success'

export type Noty = {
  Message: string
  Type?: NotyType
  LifeTimeInSeconds?: 'forever' | number
}

export type FullNoty = Noty & {
  CreatedAt: Date
  ChangedAt: Date | 'now'
  Type: NotyType
  LifeTimeInSeconds: 'forever' | number
  Shown?: true
  Hide: boolean
}

export type NotyDispatcherState = { [key: string]: FullNoty }

export type NotyDispatcherAction = {
  action: 'add_noty',
  payload: {
    id: string
    noty: FullNoty
  }
} | {
  action: 'show_noties'
} | {
  action: 'mouse_enter',
  payload: { id: string }
} | {
  action: 'mouse_leave',
  payload: { id: string }
} | {
  action: 'close',
  payload: { id: string }
} | {
  action: 'remove',
  payload: { id: string }
}


export function NotyDispatcher(state: NotyDispatcherState, action: NotyDispatcherAction): NotyDispatcherState {
  let newState = { ...state }

  switch (action.action) {
    case 'add_noty':
      newState[action.payload.id] = action.payload.noty
      break
    case 'show_noties':
      Object.keys(newState).forEach(k => {
        if (newState[k].Shown) return
        newState[k].Shown = true
        newState[k].Hide = false
      })
      break
    case 'mouse_enter':
      if (newState[action.payload.id]) {
        newState[action.payload.id].ChangedAt = 'now'
      }
      break
    case 'mouse_leave':
      if (newState[action.payload.id]) {
        newState[action.payload.id].ChangedAt = new Date()
      }
      break
    case 'close':
      if (newState[action.payload.id]) {
        newState[action.payload.id].Hide = true
      }
      break
    case 'remove':
      if (newState[action.payload.id]) {
        delete newState[action.payload.id]
      }
      break
  }

  return newState
}