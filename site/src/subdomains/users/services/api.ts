import { ApiEndpoint } from "../../../common/api/api";
import { SessionDto, SessionDtoToModel } from "../dtos/dtos";
import { Session } from "../models/session";
import { Login, Name, Password } from "../valueobjects";

// register
// login
// refresh
// profile
// change name
// change password

// register

export type RegisterArgs = {
  Name: Name
  Password: Password
}

export const Register = ApiEndpoint<RegisterArgs, void>(async (args, api) => {
  const response = await fetch(`${api.Url}/api/user/register`, {
    method: "POST",
    body: JSON.stringify({ "name": args.Name.Value, "password": args.Password.Value })
  })

  if (!response.ok) {
    await api.InvalidResponseHandler(response)
    return { Ok: false }
  }

  return { Ok: true, Model: undefined }
})

// login

export type LogInArgs = {
  Login: Login
  Password: Password
}

export const LogIn = ApiEndpoint<LogInArgs, Session>(async (args, api) => {
  const response = await fetch(`${api.Url}/api/user/log-in`, {
    method: "POST",
    body: JSON.stringify({ "login": args.Login.Value, "password": args.Password.Value })
  })

  if (!response.ok) {
    await api.InvalidResponseHandler(response)
    return { Ok: false }
  }
  const result = await response.json()
  const session = SessionDto.parse(result)

  return { Ok: true, Model: SessionDtoToModel(session) }
})

// refresh

export type RefreshArgs = {
  Session: Session
}

export const Refresh = ApiEndpoint<RefreshArgs, Session>(async (args, api) => {
  const response = await fetch(`${api.Url}/api/user/refresh`, {
    method: "POST",
    body: JSON.stringify({ "session_token": args.Session.SessionToken(), "refresh_token": args.Session.RefreshToken() })
  })

  if (!response.ok) {
    await api.InvalidResponseHandler(response)
    return { Ok: false }
  }

  const result = await response.json()
  const session = SessionDto.parse(result)

  return { Ok: true, Model: SessionDtoToModel(session) }
})

// profile

export type ProfileArgs = {
  Session: Session
}

export const Profile = ApiEndpoint<ProfileArgs, void>(async (args, api) => {
  const response = await fetch(`${api.Url}/api/user/profile`, {
    headers: args.Session.Headers()
  })

  if (!response.ok) {
    await api.InvalidResponseHandler(response)
    return { Ok: false }
  }

  return { Ok: true, Model: undefined }
})

// change name

export type ChangeNameArgs = {
  Session: Session,
  NewName: Name
}

export const ChangeName = ApiEndpoint<ChangeNameArgs, void>(async (args, api) => {
  const response = await fetch(`${api.Url}/api/user/change-name`, {
    method: "POST",
    headers: { "Content-Type": "application/json", ...args.Session.Headers() },
    body: JSON.stringify({ "new_name": args.NewName.Value })
  })

  if (!response.ok) {
    await api.InvalidResponseHandler(response)
    return { Ok: false }
  }

  return { Ok: true, Model: undefined }
})

// change password

export type ChangePasswordArgs = {
  Session: Session
  NewPassword: Password
}

export const ChangePassword = ApiEndpoint<ChangePasswordArgs, void>(async (args, api) => {
  const response = await fetch(`${api.Url}/api/user/change-password`, {
    method: "POST",
    headers: { "Content-Type": "application/json", ...args.Session.Headers() },
    body: JSON.stringify({ "new_password": args.NewPassword.Value })
  })
  if (!response.ok) {
    await api.InvalidResponseHandler(response)
    return { Ok: false }
  }
  return { Ok: true, Model: undefined }
})
