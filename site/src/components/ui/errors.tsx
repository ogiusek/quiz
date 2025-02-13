import { VO } from "@/common/vo/vo"
import React from "react"

export function ShowErrors({ allErrors: allErrors, errors, children }: { allErrors: Error[], errors: Error[], children?: React.ReactNode }): React.JSX.Element {
  let notifications: { [key: string]: boolean } = {}
  allErrors.map(e => notifications[e.message] = false)
  errors.map(e => notifications[e.message] = true)

  return <div className={`relative group`}>
    {children}
    {errors.length == 0 ? '' :
      <div className="absolute left-0 w-full hidden z-10 mt-2 group-hover:block group-focus-within:block bg-card text-card-foreground border rounded-sm p-4 gap-2">
        {Object.entries(notifications)
          .map((value, i) =>
            <p key={i} className={value[1] ? "text-red-500" : "text-green-500"}>{value[0]}</p>
          )}
      </div>}
  </div>
}

export function ShowVoErrors({ vo, children }: { vo: VO, children?: React.ReactNode }): React.JSX.Element {
  return <ShowErrors allErrors={vo.Errors()} errors={vo.Valid()}>
    {children}
  </ShowErrors>
}