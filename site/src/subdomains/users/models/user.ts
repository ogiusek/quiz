import { Name } from "../valueobjects"

export type User = {
  Id: string
  UserName: Name
  Image: string
}

export const ToUser = (payload: any): User => {
  return {
    Id: payload["id"],
    UserName: new Name(payload["name"]),
    Image: payload["image"]
  }
}