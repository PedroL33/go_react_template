import { LoginWithOptCode } from "api_calls/users";
import { AuthContext } from "contexts/AuthContext";
import { useContext, useState } from "react";
import { useNavigate } from "react-router-dom";
import { toast } from "react-toastify";
import { ILoginResponse, ILoginWithOtpCodeRequest } from "types";
import { isLocalStorageAvailable, setLocalStorageAuthToken, setLocalStorageCurrentUser } from "utils/Storage";

interface OtpCodeLoginProps {
  setUseRecoveryCode: (setUserRecoveryCode: boolean) => void
}

export const OtpCodeLogin: React.FC<OtpCodeLoginProps> = ( {setUseRecoveryCode} ) => {

  const { loginSessionId, setAuthToken, setCurrentUser } = useContext(AuthContext);
  const [otpCode, setOptCode] = useState<string>("");
  const navigate = useNavigate();

  const handleSubmit = async (e: any) => {
    e.preventDefault()
    try {
      const verifyLoginRequest: ILoginWithOtpCodeRequest = {
        loginSessionId,
        otpCode,
      }
      const response: ILoginResponse = await LoginWithOptCode(verifyLoginRequest) 
      if(response.token != null && response.user != null) {
        setAuthToken(response.token)
        setCurrentUser(response.user)
        if(isLocalStorageAvailable()) {
          setLocalStorageAuthToken(response.token, setAuthToken);
          setLocalStorageCurrentUser(response.user, setCurrentUser);
        }
      }
    }catch(err: any) {
      if(err.status === 410) {
        toast("Session expired.")
        navigate("/login")
      }
      toast(err.message)
    }
  }

  return (
    <div className="w-full m-auto max-w-sm p-4 bg-white border border-gray-200 rounded-lg shadow sm:p-6 md:p-8 dark:bg-gray-800 dark:border-gray-700">
      <form onSubmit={handleSubmit} className="space-y-6" action="#">
        <h5 className="text-xl font-medium text-gray-900 dark:text-white">Two Factor Authentication</h5>
        <div>
          <label htmlFor="username" className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">One Time Passcode</label>
          <input type="text" name="text" onChange={(e) => setOptCode(e.target.value)} className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-600 dark:border-gray-500 dark:placeholder-gray-400 dark:text-white" placeholder="123456..." required />
        </div>
        <button type="submit" className="w-full text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">Login</button>
        <div className="mt-2 text-sm text-gray-500 dark:text-gray-400 flex justify-between">
          <button onClick={() => setUseRecoveryCode(true)} className="font-medium text-blue-600 hover:underline dark:text-blue-500">Recovery Code?</button>
          <a href="/login" className="font-medium text-blue-600 hover:underline dark:text-blue-500">Back</a>
        </div>
      </form>
    </div>
  )
}