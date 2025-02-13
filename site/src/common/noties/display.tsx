import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert"
import { AlertCircle } from "lucide-react"
import { useMemo } from "react"
import { FullNoty } from "./dispatcher"

export function DisplayNoties({ noties, close, mouseEnter, mouseLeave }: {
  noties: { [key: string]: FullNoty },
  close: (id: string) => void,
  mouseEnter: (id: string) => void,
  mouseLeave: (id: string) => void
}) {
  const sorted = useMemo(() => Object.entries(noties)
    .sort(e => -e[1].CreatedAt.getUTCMilliseconds()),
    [noties]
  )

  return <>
    <div className="fixed right-0 top-5 flex flex-col gap-2 z-50">
      {sorted.map(([id, noty]) =>
        <Alert
          key={id}
          onClick={() => close(id)}
          onMouseEnter={() => mouseEnter(id)}
          onMouseLeave={() => mouseLeave(id)}
          style={{ maxWidth: "100vw" }}
          className={`w-max transition mr-5 ${!noty.Hide ? '' : 'opacity-0 translate-x-full'} hover:border-secondary`}
          variant={noty.Type == 'error' ? 'destructive' : noty.Type == 'success' ? 'success' : 'default'}
        >
          <AlertCircle className="h-6 w-6 mr-2" />
          <AlertTitle className="text-lg">{noty.Type == 'error' ? 'error' : 'success'}</AlertTitle>
          <AlertDescription className="text-md">{noty.Message}</AlertDescription>
        </Alert>
      )}
    </div >
  </>
}