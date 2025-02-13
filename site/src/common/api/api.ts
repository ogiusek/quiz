type Redo = {
  additional: { [key: string]: any }
} | void

export type ApiDefinition = {
  Url: string
  InvalidResponseHandler: (response: Response) => Promise<Redo>
  ErrorHandler: (error: any) => void
}

type ApiResponse<T> = Promise<{ Ok: false } | { Ok: true, Model: T }>

type ApiEndpoint<Args, Res> = (args: Args, api: ApiDefinition) => ApiResponse<Res>

/**
 * Handles partly ApiDefinition
 * - handles fully  'ErrorHandler' by using try catch.
 * - handles partly 'InvalidReponseHandler' by automaticaly re-doing request but has to be invoked first
 * 
 * IMPORTANT
 * there are properies in args which can be refreshed after redo.
 * things like Session for example, so follow specific naming pattern otherwise user will have to manually redo action.
 * 
 * naming pattern: Name properties with the same name as their type
 */
export const ApiEndpoint = <TArgs, TRes>(handler: ApiEndpoint<TArgs, TRes>): ApiEndpoint<TArgs, TRes> => async (args, api) => {
  try {
    let redo = true
    while (redo) {
      redo = false
      const res = await handler(args, {
        Url: api.Url,
        ErrorHandler: api.ErrorHandler,
        InvalidResponseHandler: async (response) => {
          const redoRes = await api.InvalidResponseHandler(response)
          if (!redoRes) return
          redo = true
          args = {
            ...args,
            ...redoRes.additional,
          }
        }
      })

      if (redo) continue

      return res
    }
  } catch (err) {
    api.ErrorHandler(err)
  }
  return { Ok: false }
}
