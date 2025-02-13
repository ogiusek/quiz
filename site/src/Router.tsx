import { HashRouter as Router, Route, Navigate, Routes } from 'react-router-dom'
import Login from "./subdomains/users/views/Login"
import Register from './subdomains/users/views/Register'
import { memo, useContext } from 'react'
import { SessionContext } from './subdomains/users/contexts/sessionContext'
import { Authorized } from './common/routeWrappers/authorized'
import { Profile } from './subdomains/users/views/Profile'
import { UnAuthorized } from './common/routeWrappers/unauthorized'
import { Search } from './subdomains/questions/views/Search'
import { QuestionSet } from './subdomains/questions/views/QuestionSet'
import { MyQuestionSets } from './subdomains/questions/views/MyQuestionSets'
import { Host } from './subdomains/match/views/host'
import { Join, JoinId } from './subdomains/match/views/join'

export const HrefAfterMatch = "/#/quiz/join"

const RouterWrapper = memo(() => {
  const sessionStorage = useContext(SessionContext);

  const defaultAuthorizedEndpoint = "/user/profile"
  const defaultUnAuthorizedEndpoint = "/user/login"

  return <Router>
    <Routes>
      {/* unauthorized */}
      <Route path='/user/login' Component={UnAuthorized(Login, defaultAuthorizedEndpoint)} />
      <Route path='/user/register' Component={UnAuthorized(Register, defaultAuthorizedEndpoint)} />

      {/*  authorized */}
      <Route path='/user/profile' Component={Authorized(Profile, defaultUnAuthorizedEndpoint)} />
      <Route path='/question-set/search' Component={Authorized(Search, defaultUnAuthorizedEndpoint)} />
      <Route path='/question-set/my' Component={Authorized(MyQuestionSets, defaultUnAuthorizedEndpoint)} />
      <Route path='/question-set/get/:id' Component={Authorized(QuestionSet, defaultUnAuthorizedEndpoint)} />
      <Route path='/quiz/host' Component={Authorized(Host, defaultUnAuthorizedEndpoint)} />
      <Route path='/quiz/join/:id' Component={Authorized(JoinId, defaultUnAuthorizedEndpoint)} />
      <Route path='/quiz/join' Component={Authorized(Join, defaultUnAuthorizedEndpoint)} />

      {/* other */}
      <Route path='*' element={<>
        {sessionStorage.GetSession()?.Valid() == 'invalid' ?
          <Navigate to={defaultUnAuthorizedEndpoint} /> :
          <Navigate to={defaultAuthorizedEndpoint} />}
      </>} />
    </Routes>
  </Router>
})

export { RouterWrapper as Router }