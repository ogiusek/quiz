import { SessionAvatar } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { SessionContext } from "@/subdomains/users/contexts/sessionContext";
import { LogOut, MenuIcon, X } from "lucide-react";
import { useContext, useState } from "react";
import { Link } from "react-router-dom";

type LinkType = {
  user: "authorized" | "unauthorized" | "any"
  url: string
  message: string
}

type LinksType = {
  [key: string]: LinkType[]
}

const allLinks: LinksType = {
  "account": [
    { user: "authorized", url: "/user/profile", message: "profile" },
    { user: "unauthorized", url: "/user/login", message: "login" },
    { user: "unauthorized", url: "/user/reigster", message: "register" },
  ],
  "question sets": [
    { user: "authorized", url: "/question-set/my", message: "my question sets" },
    { user: "authorized", url: "/question-set/search", message: "search question sets" },
  ],
  "quiz": [
    { user: "authorized", url: "/quiz/host", message: "host" },
    { user: "authorized", url: "/quiz/join", message: "join" },
  ]
}

export const Nav = () => {
  const [show, setShow] = useState(false)
  const sessionContext = useContext(SessionContext)
  const session = sessionContext.GetSession()
  if (!session) {
    return <></>
  }

  const links: LinksType = Object.fromEntries(
    Object.entries(allLinks)
      .map(([section, links]) => [
        section,
        links.filter(e =>
          (!!session && e.user != 'unauthorized') ||
          (!session && e.user != 'authorized')
        )
      ]).filter(([_, links]) => links.length != 0)
  )

  return <>
    <Button
      variant="ghost"
      onClick={() => setShow(true)}
      className={`fixed top-4 transition right-4 ${!show ? "z-40" : "opacity-0 translate-x-full"}`}>
      <MenuIcon />
    </Button>
    <button className={`
      fixed top-0 left-0
      w-screen h-screen
      opacity-30 bg-black
      ${show ? "z-30" : "translate-x-full"}`}
      onClick={() => setShow(false)}></button>
    <nav style={{ maxWidth: "100vw" }} className={`
      fixed top-0 right-0
      w-96 h-screen p-4
      bg-card border rounded-l-md
      transition
      ${show ? "z-40" : "translate-x-full"}`}>
      <div className="flex flex-row justify-between gap-2 w-full">
        <div className="flex flex-row gap-2">
          {!!session && <>
            <SessionAvatar session={session.Session()!} />
            <Button
              variant="destructive"
              onClick={sessionContext.UnSetSession}>
              <LogOut />
            </Button>
          </>}
        </div>
        <Button
          variant="ghost"
          onClick={() => setShow(false)}>
          <X />
        </Button>
      </div>
      <br />
      {Object.entries(links).map(([section, links], i) => <div key={i}>
        <h3 className="text-2xl">{section}</h3>
        <ul>
          {links.map((link, i) => <li key={i} className="w-full">
            <Button variant="link" asChild className="w-full">
              <Link to={link.url}>{link.message}</Link>
            </Button>
          </li>)}
        </ul>
      </div>)}
    </nav>
  </>
}