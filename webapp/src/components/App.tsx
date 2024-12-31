import { RouteWrapper } from './RouteWrapper';
import { AuthContext, AuthContextState } from '../contexts/AuthContext';
import { useLocalStorage } from '../hooks/UseLocalStorage';
import { IUser } from '../types';
import { ToastContainer } from 'react-toastify';
import { useState } from 'react';


function App() {

  const [authToken, setAuthToken] = useLocalStorage<string>('authToken', "")
  const [currentUser, setCurrentUser] = useLocalStorage<IUser>('currentUser', {
    id: -1,
    email: "",
    password: "",
    twoFactorSecret: {
      String: "", 
      Valid: false,
    },
    isTwoFactorEnabled: {
      Bool: false,
      Valid: true,
    },
    createdAt: "",
    updatedAt: "",
  })
  const [loginSessionId, setLoginSessionId] = useState<number>(-1);

  const value: AuthContextState = {
    authToken,
    setAuthToken,
    currentUser,
    setCurrentUser,
    loginSessionId,
    setLoginSessionId,
  }

  return (
    <div className="App">
      <AuthContext.Provider value={value}>
        <RouteWrapper />
        <ToastContainer />
      </AuthContext.Provider>
    </div>
  );
}

export default App;