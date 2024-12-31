import { useState, useContext } from 'react';
import { AuthContext } from '../contexts/AuthContext';
import { ILoginRequest, ILoginResponse } from '../types';
import { AuthenticatedPage } from '../components/AuthenticatedPage';

export const Dashboard = () => {

  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [error, setError] = useState<string>("");


  return (
    <AuthenticatedPage>
      <div className="p-10 m-auto flex flex-col items-center justify-center bg-cool-gray-700">
        <h1 className="text-9xl text-center mb-28">
          <div className="bg-gradient-to-r text-transparent bg-clip-text from-green-400 to-purple-500">
            Dashboard
          </div>
        </h1>
        <div className="color-red">{error}</div>
      </div>
    </AuthenticatedPage>
  )
}