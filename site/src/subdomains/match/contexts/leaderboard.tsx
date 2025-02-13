import React, { createContext, useState } from "react";
import { MatchDto } from "../models";
import { UserAvatar } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { X } from "lucide-react";

export interface Leaderboard {
  FinishedMatch(match: MatchDto): void
}

export const LeaderboardContext = createContext<Leaderboard>({
  FinishedMatch(_) { throw new Error("not implemented") },
})

export const LeaderboardSerivce: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [match, setMatch] = useState<MatchDto | undefined>(undefined)

  const sortedPlayers = match?.Players.sort((p1, p2) => p2.Score - p1.Score)
  const winner = sortedPlayers?.at(0)
  return <>
    <LeaderboardContext.Provider value={{
      FinishedMatch(match) {
        setMatch(JSON.parse(JSON.stringify(match)))
      },
    }}>
      {match && <>
        {match?.Course?.Step == 'finished' && <>
          <main style={{ maxWidth: "100vw" }} className="flex justify-center items-center h-screen fixed top-0 left-0 z-30 bg-background">
            <div className="bg-card text-card-foreground rounded-lg shadow-lg p-8 border w-full max-w-3xl h-5/6 flex flex-col gap-4">
              <div className="flex flex-row justify-between">
                <h2>Leader board</h2>
                <Button onClick={() => setMatch(undefined)}><X /></Button>
              </div>
              <h3 className="text-2xl">Winner is <span className="text-3xl">{winner?.User.UserName.Value}</span></h3>


              <ul className="flex flex-col gap-2 mt-4">
                <li className="flex flex-row justify-between">
                  <p>avatar</p>
                  <p>rank</p>
                  <p>score</p>
                </li>
                {sortedPlayers?.map((player, i) => <li key={player.Id} className="flex flex-row justify-between">
                  <UserAvatar user={player.User} />
                  <p>{i + 1}</p>
                  <p>{player.Score}</p>
                </li>)}
              </ul>
            </div>
          </main>
        </>}
      </>}
      {children}
    </LeaderboardContext.Provider>
  </>
}