import React, { memo } from 'react'
import { Router } from './Router'
import { ThemeService } from './common/theme/context'
import { NotiesService } from './common/noties/notiesContext'
import { SessionService } from './subdomains/users/contexts/sessionContext'
import { ApiService } from './common/api/apiContext'
import { WsService } from './common/ws/wsContext'
import { WsNoties } from './common/ws/wsNoties'
import { ReJoinService } from './subdomains/match/contexts/rejoin'
import { MatchService } from './subdomains/match/contexts/match'
import { LeaderboardSerivce } from './subdomains/match/contexts/leaderboard'

const Url = import.meta.env.VITE_API_URL
const WsUrl = import.meta.env.VITE_API_WS_URL

const services: React.FC<{ children: React.ReactNode }>[] = [
  // universal
  ThemeService,
  NotiesService,
  WsService,
  WsNoties,
  SessionService,
  ApiService({ Url, WsUrl }),

  // match
  LeaderboardSerivce,
  MatchService,
  ReJoinService,
]

const ServicesComponent = services
  .map(s => memo(s))
  .reduce((ServiceComponent, NextComponent) =>
    memo(({ children }) => <ServiceComponent children={<NextComponent children={children} />} />)
  )

const App = () => (
  <ServicesComponent>
    <Router />
  </ServicesComponent>
)

export default App

