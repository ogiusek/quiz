import { ApiContext } from "@/common/api/apiContext"
import { useContext, useState } from "react"
import { SessionContext } from "../contexts/sessionContext"
import { NotiesContext } from "@/common/noties/notiesContext"
import { Button } from "@/components/ui/button"
import { SessionAvatar } from "@/components/ui/avatar"
import { Label } from "@/components/ui/label"
import { Input } from "@/components/ui/input"
import { Edit, LogOutIcon } from "lucide-react"
import { ShowErrors, ShowVoErrors } from "@/components/ui/errors"
import { Name, Password } from "../valueobjects"
import { ChangeName, ChangePassword } from "../services/api"
import { Nav } from "@/common/nav/nav"

export function Profile() {
  const api = useContext(ApiContext)
  const sessionContext = useContext(SessionContext)
  const noties = useContext(NotiesContext)

  const session = sessionContext.GetSession()
  const sessionData = session?.Session()
  if (!sessionData) return <></>

  const [name, setName] = useState(new Name(sessionData.UserName))
  const [password, setPassword] = useState(new Password(''))
  const [repeatPassword, setRepeatPassword] = useState(new Password(''))
  const repeatPasswordErrors = [new Error('passwords must match')]
  const repeatPasswordFoundErrors: Error[] = password.Value == repeatPassword.Value ? [] : repeatPasswordErrors

  const SaveName = async () => {
    const response = await ChangeName({
      Session: session!,
      NewName: name
    }, api)
    if (!response.Ok) return

    noties.AddNoty({
      Type: "success",
      Message: `succesfuly changed name to ${name.Value}`
    })
  }

  const SavePassword = async () => {
    const response = await ChangePassword({
      Session: session!,
      NewPassword: password
    }, api)
    if (!response.Ok) return

    noties.AddNoty({
      Type: "success",
      Message: `succesfuly changed password`
    })
  }

  const LogOut = sessionContext.UnSetSession

  return <>
    <Nav />
    <main className="flex justify-center items-center min-h-screen">
      <div className="bg-card text-card-foreground rounded-lg shadow-lg p-8 border w-full max-w-md flex flex-col gap-4">
        <h1 className="text-center text-3xl">Quiz</h1>
        <div className="flex flex-row justify-between">
          <h2 className="text-2xl">Profile</h2>
          <Button aria-label="log out" variant="destructive" className="w-max ml-auto" onClick={LogOut}><LogOutIcon /></Button>
        </div>

        <div className="w-full flex justify-center items-center">
          <SessionAvatar session={sessionData} />
        </div>

        <ShowVoErrors vo={name}>
          <Label htmlFor="name">Name</Label>
          <div className="flex flex-row gap-2 justify-center items-center">
            <Input id="name" type="text" placeholder="John" value={name.Value} onChange={e => setName(new Name(e.target.value))} />
            <Button aria-label="edit name " variant="outline" disabled={name.Valid().length != 0} onClick={SaveName}>
              <Edit />
            </Button>
          </div>
        </ShowVoErrors>

        {/* <ShowErrors errors={[...password.Errors(), ...repeatPasswordErrors]} allErrors={[...password.Valid(), ...repeatPasswordFoundErrors]}> */}
        <ShowErrors allErrors={[...password.Errors(), ...repeatPasswordErrors]} errors={[...password.Valid(), ...repeatPasswordFoundErrors]}>
          <Label htmlFor="password">Change Password</Label>
          <div className="flex flex-row gap-2 justify-center items-center">
            <Input id="password" type="password" placeholder="********" value={password.Value} onChange={e => setPassword(new Password(e.target.value))} />
          </div>
          <Label htmlFor="repeat_password">Repeat Password</Label>
          <div className="flex flex-row gap-2 justify-center items-center">
            <Input id="repeat_password" type="password" placeholder="********" value={repeatPassword.Value} onChange={e => setRepeatPassword(new Password(e.target.value))} />
            <Button aria-label="edit password" variant="outline" disabled={[...password.Valid(), ...repeatPasswordFoundErrors].length != 0} onClick={SavePassword}>
              <Edit />
            </Button>
          </div>
        </ShowErrors>
      </div>
    </main>
  </>
} 