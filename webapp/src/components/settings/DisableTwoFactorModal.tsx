import { DisableTwoFactorAuth } from "api_calls/users";
import { AuthContext } from "contexts/AuthContext";
import { useContext, useState } from "react";
import { toast } from "react-toastify";
import { IDisable2faRequest, IUser } from "types";
import { isLocalStorageAvailable, setLocalStorageCurrentUser } from "utils/Storage";

export const DisableTwoFactorModal: React.FC = () => {

  const { authToken, setCurrentUser, currentUser } = useContext(AuthContext);
  const [showModal, setShowModal] = useState<boolean>(false);
  const [password, setPassword] = useState<string>("");

  const disable2FA = async (e: any) => {
    e.preventDefault()
    try {
      const disable2faRequest: IDisable2faRequest = {
        password: password
      }
      const response = await DisableTwoFactorAuth(authToken, disable2faRequest)
      toast(response)
      setShowModal(false);
      const updatedUser: IUser = {...currentUser, isTwoFactorEnabled: {Bool: false, Valid: true}};
      isLocalStorageAvailable() ? setLocalStorageCurrentUser(updatedUser, setCurrentUser): setCurrentUser(updatedUser);
    }catch(err: any) {
      console.log(err)
      toast(err.message)
    }
  }

  const handleClose = () => {
    setShowModal(false);
    setPassword("");
  }

  return (
    <div>
      <button onClick={(e) => setShowModal(true)} className="block text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800" type="button">
        Disable
      </button>
      {
        showModal && <div tabIndex={-1} className="overflow-y-auto overflow-x-hidden fixed top-0 right-0 left-0 z-50 justify-center items-center w-full md:inset-0 h-[calc(100%-1rem)] max-h-full">
          <div className="h-full flex items-center justify-center">
            <div className="max-w-2xl relative bg-white rounded-lg shadow dark:bg-gray-700">
              <div className="flex items-center justify-between p-4 md:p-5 border-b rounded-t dark:border-gray-600">
                <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
                  Disable Two Factor Authentication
                </h3>
                <button type="button" onClick={handleClose} className="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm w-8 h-8 ms-auto inline-flex justify-center items-center dark:hover:bg-gray-600 dark:hover:text-white">
                  <svg className="w-3 h-3" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 14 14">
                    <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"/>
                  </svg>
                  <span className="sr-only">Close modal</span>
                </button>
              </div>
              <div className="p-4 md:p-5 space-y-4">
                <p className="text-base leading-relaxed text-gray-500 dark:text-gray-400">
                  Are you sure you want to disable two factor authentication?
                </p>
                <p className="text-base leading-relaxed text-gray-500 dark:text-gray-400">
                  Retype your password to confirm.
                </p>
              </div>
              <form onSubmit={disable2FA} className="max-w-sm mx-auto pb-10">
                <div className="mb-5">
                  <label htmlFor="email" className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Password</label>
                  <input type="password" onChange={(e) => setPassword(e.target.value)} className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" placeholder="****" required />
                </div>
                <div className="flex justify-center">
                  <button disabled={password.length < 4} type="submit" className="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">Confirm</button>
                </div>
              </form>

            </div>
          </div>
        </div>
      }
    </div>
  )
}