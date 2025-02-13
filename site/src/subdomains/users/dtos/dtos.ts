import { z } from "zod";
import { Session } from "../models/session";

// session
// user

// session

export const SessionDto = z.object({
  session_token: z.string(),
  refresh_token: z.string()
});

type SessionDtoType = z.infer<typeof SessionDto>

export const SessionDtoToModel = (dto: SessionDtoType): Session => {
  return new Session(dto.session_token, dto.refresh_token);
}

// user
// TODO