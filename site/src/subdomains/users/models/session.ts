import { jwtDecode } from "jwt-decode";

// i only made that a separate directory is because of not being able to clasify this
// and i'm having dilema should models in other sub-domains be a file or a directory

export type SessionUserData = {
  UserId: string
  UserName: string
  UserImage: string
}

export class Session {
  private _sessionToken: string;
  private _refreshToken: string;

  constructor(sessionToken: string, refreshToken: string) {
    this._sessionToken = sessionToken;
    this._refreshToken = refreshToken;
  }

  SessionToken(): string { return this._sessionToken }
  RefreshToken(): string { return this._refreshToken }

  Session(): SessionUserData | undefined {
    try {
      const session: { [key: string]: string } = jwtDecode(this._sessionToken)
      return {
        UserId: session["user_id"],
        UserName: session["user_name"],
        UserImage: session["user_image"],
      }
    } catch (_) {
      return undefined
    }
  }

  Valid(): 'valid' | 'expired' | 'invalid' {
    try {
      const currentTime = Math.floor(Date.now() / 1000)
      const session: { [key: string]: any } = jwtDecode(this._sessionToken)
      const refresh: { [key: string]: any } = jwtDecode(this._refreshToken)
      const sessionExpired = session.exp < currentTime
      const refreshExpired = refresh.exp < currentTime

      if (refreshExpired && sessionExpired)
        return 'invalid'

      if (sessionExpired)
        return 'expired'

      return 'valid'
    } catch (err) {
      return 'invalid'
    }
  }

  Headers(): Record<string, string> {
    return { "authorization": this._sessionToken }
  }

  Params(): URLSearchParams {
    return new URLSearchParams({ "authorization": this._sessionToken })
  }
}