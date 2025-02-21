import { ApiDefinition } from "@/common/api/api"
import { NotiesContext, NotiesStorage } from "@/common/noties/notiesContext"
import { Register as ApiRegister, RegisterArgs } from "../services/api"
import { useContext, useState } from "react"
import { ApiContext } from "@/common/api/apiContext"
import { Name, Password } from "../valueobjects"
import { ShowErrors, ShowVoErrors } from "@/components/ui/errors"
import { Label } from "@radix-ui/react-label"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Link } from "react-router-dom"
import { Nav } from "@/common/nav/nav"

const register = async (api: ApiDefinition, noties: NotiesStorage, args: RegisterArgs) => {
  const response = await ApiRegister(args, api)
  if (response.Ok)
    noties.AddNoty({
      Type: "noty",
      Message: "Success you can login "
    })
}

export default function Register() {
  const api = useContext(ApiContext)
  const noties = useContext(NotiesContext)

  const [args, setArgs] = useState<RegisterArgs>({
    Name: new Name(""),
    Password: new Password("")
  })
  const [repeatPassword, setNewPassword] = useState(new Password(""))

  let error = args.Password.Value == repeatPassword.Value ? null : new Error("passwords must equal")
  const errors = [...args.Name.Valid(), ...args.Password.Valid()]
  error && errors.push(error)

  const onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    register(api, noties, args)
  }

  return <>
    <Nav />
    <main className="flex justify-center items-center min-h-screen">
      <form
        className="bg-card text-card-foreground rounded-lg shadow-lg p-8 border w-full max-w-md flex flex-col gap-4"
        onSubmit={onSubmit}>

        <h1 className="text-center text-3xl">Quiz</h1>
        <h2 className="text-2xl">Register</h2>

        <ShowVoErrors vo={args.Name}>
          <Label htmlFor="name">Name</Label>
          <Input id="name" type="text" placeholder="John" value={args.Name.Value}
            onChange={e => setArgs({ ...args, Name: new Name(e.target.value) })} />
        </ShowVoErrors>

        <ShowVoErrors vo={args.Password}>
          <Label htmlFor="password">Password</Label>
          <Input id="password" type="password" placeholder="********" value={args.Password.Value}
            onChange={e => setArgs({ ...args, Password: new Password(e.target.value) })} />
        </ShowVoErrors>

        <ShowErrors allErrors={[]} errors={[]}>
          <Label htmlFor="repeat_password">Repeat Password</Label>
          <Input id="repeat_password" type="password" placeholder="********" value={repeatPassword.Value}
            onChange={e => setNewPassword(new Password(e.target.value))} />
        </ShowErrors>

        <ShowErrors allErrors={errors} errors={errors}>
          <Button aria-label="login link" className="w-full" type="submit" disabled={errors.length != 0}>Register</Button>
        </ShowErrors>

        <div className="flex flex-row justify-between">
          {/* <Button asChild variant="link"><Link to="/user/register">register</Link></Button> */}
          <Button aria-label="login link" asChild variant="link"><Link to={`/user/login?${window.location.href.split("?").filter((_, i) => i != 0).join("?")}`}>login</Link></Button>
        </div>
      </form>
    </main>
  </>
}