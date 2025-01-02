import { useState, useContext } from 'react';
import { AuthContext } from '../contexts/AuthContext';
import { ILoginRequest } from '../types';
import { isLocalStorageAvailable, setLocalStorageAuthToken, setLocalStorageCurrentUser } from '../utils/Storage';
import { Login } from 'api_calls/users';
import { toast } from 'react-toastify';
import { useNavigate } from 'react-router-dom';

export const LoginPage = () => {

  const [username, setUsername] = useState<string>("");
  const [password, setPassword] = useState<string>("");

  const { setAuthToken, setCurrentUser, setLoginSessionId } = useContext(AuthContext); 

  const navigate = useNavigate();

  const handleSubmit = async (e: any): Promise<void> => {
    e.preventDefault();
    const loginRequest: ILoginRequest = {
      username,
      password
    }
    
    try {
      const result = await Login(loginRequest)
      if(result.loginSessionId != null) {
        setLoginSessionId(result.loginSessionId)
        navigate("/two_factor_authentication")
      }else {
        setAuthToken(result.token!)
        setCurrentUser(result.user!)
        if(isLocalStorageAvailable()) {
          setLocalStorageAuthToken(result.token!, setAuthToken);
          setLocalStorageCurrentUser(result.user!, setCurrentUser);
        }
      }
    }catch(err: any) {
      toast(err.message)
      console.log(err)
    }
  }

  return (
    <div className="h-screen flex items-center justify-center">
      <div className="w-full m-auto max-w-sm p-4 bg-white border border-gray-200 rounded-lg shadow sm:p-6 md:p-8 dark:bg-gray-800 dark:border-gray-700">
        <form onSubmit={handleSubmit} className="space-y-6" action="#">
          <h5 className="text-xl font-medium text-gray-900 dark:text-white">Sign in</h5>
          <div>
            <label htmlFor="username" className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Username</label>
            <input type="text" name="username" id="username" onChange={(e) => setUsername(e.target.value)} className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-600 dark:border-gray-500 dark:placeholder-gray-400 dark:text-white" placeholder="username..." required />
          </div>
          <div>
            <label htmlFor="password" className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Password</label>
            <input type="password" name="password" id="password" placeholder="••••••••" onChange={(e) => setPassword(e.target.value)} className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-600 dark:border-gray-500 dark:placeholder-gray-400 dark:text-white" required />
          </div>
          <button type="submit" className="w-full text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">Login</button>
        </form>
      </div>
    </div>
  )
}