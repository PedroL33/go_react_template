import { useContext } from 'react';
// import { LoginPage } from './LoginPage';
import {
  BrowserRouter as Router,
  Route,
  Routes,
  Navigate
} from "react-router-dom";
import { LoginPage } from '../pages/Login';
import { Settings } from '../pages/Settings';
import { AuthContext } from '../contexts/AuthContext';
import { Dashboard } from '../pages/Dashboard';
import { VerifyLoginPage } from 'pages/VerifyLogin';

export const RouteWrapper = () => {

  const { authToken } = useContext(AuthContext);

  if(authToken === "") {
    return (
      <Router>
        <Routes>
          <Route path="/login" element={<LoginPage/>}/>
          <Route path="/two_factor_authentication" element={<VerifyLoginPage />}/>
          <Route path="*" element={<Navigate to="/login"/>}/>
        </Routes>
      </Router>
    )
  }

  return (
    <Router>
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route path="/settings" element={<Settings />} />
        <Route path="*" element={<Navigate to="/"/>}/>
      </Routes>
    </Router>
  )
}