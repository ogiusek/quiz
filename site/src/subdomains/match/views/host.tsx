import { Nav } from "@/common/nav/nav"
import { WsContext } from "@/common/ws/wsContext"
import { Button } from "@/components/ui/button"
import { useContext } from "react"

export const Host = () => {
  const wsContext = useContext(WsContext)

  return <>
    <Nav />
    <main className="flex justify-center items-center h-screen p-2 pt-12">
      <div className="bg-card text-card-foreground rounded-lg shadow-lg p-8 border w-full max-w-md flex flex-col justify-center items-center gap-4">
        <Button aria-label="host" onClick={() => wsContext.SendMessage({ topic: "match/host", payload: {} })}>
          Host match
        </Button>
      </div>
    </main>
  </>
}