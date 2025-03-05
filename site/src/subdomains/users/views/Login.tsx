import { useContext, useState } from "react";
import { ApiContext } from "../../../common/api/apiContext";
import { LogIn, LogInArgs } from "../services/api";
import { Login as LoginVO, Password } from "../valueobjects";
import { SessionStorage, SessionContext } from "../contexts/sessionContext";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { ApiDefinition } from "@/common/api/api";
import { ShowErrors, ShowVoErrors } from "@/components/ui/errors";
import { NotiesStorage, NotiesContext } from "@/common/noties/notiesContext";
import { Link } from "react-router-dom";
import { Nav } from "@/common/nav/nav";

const login = async (sessionContext: SessionStorage, api: ApiDefinition, noties: NotiesStorage, args: LogInArgs) => {
  const existingUserData = sessionContext.GetSession()?.Session()
  if (existingUserData)
    return noties.AddNoty({
      Message: `you are already logged in as '${existingUserData.UserName}'`,
      Type: 'noty',
    })

  const session = await LogIn(args, api)
  if (session.Ok)
    return sessionContext.SetSession(session.Model)
}

export default function Login() {
  const api = useContext(ApiContext)
  const sessionContext = useContext(SessionContext)
  const noties = useContext(NotiesContext)

  const [args, setArgs] = useState<LogInArgs>({
    Login: new LoginVO(""),
    Password: new Password("")
  })

  const errors = [...args.Login.Valid(), ...args.Password.Valid()]

  const onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    login(sessionContext, api, noties, args)
  }

  return <>
    <Nav />
    <main className="flex justify-center items-center min-h-screen">
      <form
        className="bg-card text-card-foreground rounded-lg shadow-lg p-8 border w-full max-w-md flex flex-col gap-4"
        onSubmit={onSubmit}>

        <h1 className="text-center text-3xl">Quiz</h1>
        <h2 className="text-2xl">Login</h2>

        <ShowVoErrors vo={args.Login}>
          <Label htmlFor="login">Login</Label>
          <Input id="login" type="text" placeholder="John" value={args.Login.Value}
            onChange={e => setArgs({ ...args, Login: new LoginVO(e.target.value) })} />
        </ShowVoErrors>

        <ShowVoErrors vo={args.Password}>
          <Label htmlFor="password">Password</Label>
          <Input id="password" type="password" placeholder="********" value={args.Password.Value}
            onChange={e => setArgs({ ...args, Password: new Password(e.target.value) })} />
        </ShowVoErrors>

        <ShowErrors allErrors={errors} errors={errors}>
          <Button aria-label="login link" className="w-full" type="submit" disabled={errors.length != 0}>Login</Button>
        </ShowErrors>

        <div className="flex flex-row justify-between">
          <Button aria-label="register link" asChild variant="link"><Link to={`/user/register?${window.location.href.split("?").filter((_, i) => i != 0).join("?")}`}>register</Link></Button>
          {/* <Button asChild variant="link"><Link to="/user/login">login</Link></Button> */}
        </div>
      </form>
    </main>
  </>
}