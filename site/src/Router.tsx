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
      <Route path='/user/login' element={<UnAuthorized children={<Login />} defaultUrl={defaultAuthorizedEndpoint} />} />
      <Route path='/user/register' element={<UnAuthorized children={<Register />} defaultUrl={defaultAuthorizedEndpoint} />} />

      {/*  authorized */}
      <Route path='/user/profile' element={<Authorized children={<Profile />} defaultUrl={defaultUnAuthorizedEndpoint} />} />
      <Route path='/question-set/search' element={<Authorized children={<Search />} defaultUrl={defaultUnAuthorizedEndpoint} />} />
      <Route path='/question-set/my' element={<Authorized children={<MyQuestionSets />} defaultUrl={defaultUnAuthorizedEndpoint} />} />
      <Route path='/question-set/get/:id' element={<Authorized children={<QuestionSet />} defaultUrl={defaultUnAuthorizedEndpoint} />} />
      <Route path='/quiz/host' element={<Authorized children={<Host />} defaultUrl={defaultUnAuthorizedEndpoint} />} />
      <Route path='/quiz/join/:id' element={<Authorized children={<JoinId />} defaultUrl={defaultUnAuthorizedEndpoint} />} />
      <Route path='/quiz/join' element={<Authorized children={<Join />} defaultUrl={defaultUnAuthorizedEndpoint} />} />

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