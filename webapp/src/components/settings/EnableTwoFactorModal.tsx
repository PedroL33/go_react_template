import { CompleteTwoFactorAuth, EnableTwoFactorAuth } from "api_calls/users";
import { AuthContext } from "contexts/AuthContext";
import { useContext, useState } from "react";
import { toast } from "react-toastify";
import { IComplete2faRequest, IUser } from "types";
import { isLocalStorageAvailable, setLocalStorageCurrentUser } from "utils/Storage";

export const EnableTwoFactorModal: React.FC = () => {

  const { authToken, currentUser, setCurrentUser } = useContext(AuthContext)
  const [showModal, setShowModal] = useState(false);
  const [image, setImage] = useState<string>("");
  const [code, setCode] = useState<string>("");
  const [recoveryCodes, setRecoveryCodes] = useState<string[]>([]);
  const [updatedUser, setUpdatedUser] = useState<IUser>(currentUser);

  const begin2faSetup = async (e: any) => {
    e.preventDefault();
    try {
      setShowModal(true);
      const response = await EnableTwoFactorAuth(authToken);
      setImage(response.base64Qrcode);
    }catch(err: any) {
      toast(err.message)
      console.log(err)
    }
  }

  const handleClose = () => {
    setImage("");
    setRecoveryCodes([]);
    setShowModal(false);
    setCode("")
    setUpdatedUser(currentUser)
    if(isLocalStorageAvailable()) {
      setLocalStorageCurrentUser(updatedUser, setCurrentUser);
    }
  }

  const complete2faSetup = async (e: any) => {
    e.preventDefault();
    try {
      const body: IComplete2faRequest = {
        otpCode: code
      }
      const response = await CompleteTwoFactorAuth(authToken, body);
      const updatedUser: IUser = {...currentUser, isTwoFactorEnabled: {Bool: true, Valid: true}};
      setUpdatedUser({...updatedUser, isTwoFactorEnabled: {Bool: true, Valid: true}});
      toast("Successfully enrolled in two factor authentication.")
      setRecoveryCodes(response.recoveryCodes)
    }catch(err: any) {
      toast(err.message)
      console.log(err)
    }
  }

  return (
    <div>
      <button onClick={begin2faSetup} className="block text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800" type="button">
        Enable
      </button>
      {
        showModal && <div tabIndex={-1} className="overflow-y-auto overflow-x-hidden fixed top-0 right-0 left-0 z-50 justify-center items-center w-full md:inset-0 h-[calc(100%-1rem)] max-h-full">
          <div className="h-full flex items-center justify-center">
            <div className="max-w-2xl relative bg-white rounded-lg shadow dark:bg-gray-700">
              <div className="flex items-center justify-between p-4 md:p-5 border-b rounded-t dark:border-gray-600">
                <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
                  Enable Two Factor Authentication
                </h3>
                <button type="button" onClick={handleClose} className="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm w-8 h-8 ms-auto inline-flex justify-center items-center dark:hover:bg-gray-600 dark:hover:text-white" data-modal-hide="default-modal">
                  <svg className="w-3 h-3" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 14 14">
                    <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"/>
                  </svg>
                  <span className="sr-only">Close modal</span>
                </button>
              </div>
              {
              !recoveryCodes.length ? 
              <div>
                <div className="p-4 md:p-5 space-y-4">
                  <img src={`data:image/png;base64,${image}`} className="m-auto p-10"/>
                  <p className="text-base leading-relaxed text-gray-500 dark:text-gray-400">
                    Scan this Qr code with your favorite authenticator and submit the associated one time passcode.
                  </p>
                </div>
                <div className="flex flex-col px-24 pb-8">
                  <label className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">One Time Passcode</label>
                  <input type="text" onChange={(e) => setCode(e.target.value)} name="code" className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-2.5 dark:bg-gray-600 dark:border-gray-500 dark:placeholder-gray-400 dark:text-white" placeholder="..." required />
                </div>
                <div className="flex items-center p-4 md:p-5 border-t border-gray-200 rounded-b dark:border-gray-600">
                  <form className="space-y-6" action="#">
                    <button disabled={code.length < 6} onClick={complete2faSetup} type="button" className="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">Submit</button>
                  </form>
                </div>
              </div>:
              <div className="w-full p-4 bg-white shadow sm:p-6 dark:bg-gray-800 dark:border-gray-700">
                <h5 className="mb-3 text-base font-semibold text-gray-900 md:text-xl dark:text-white">
                  Recovery Codes
                </h5>
                <p className="text-sm font-normal text-gray-500 dark:text-gray-400">Save these recovery codes somewhere safe.</p>
                <ul className="my-4 space-y-3">
                  {
                    recoveryCodes.map((code, idx) => <li key={idx}>
                      <div className="flex items-center text-base font-bold text-gray-900 rounded-lg bg-gray-50 hover:bg-gray-100 group hover:shadow dark:bg-gray-600 dark:hover:bg-gray-500 dark:text-white">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth="1.5" stroke="currentColor" className="size-6">
                          <path strokeLinecap="round" strokeLinejoin="round" d="M9 12.75 11.25 15 15 9.75m-3-7.036A11.959 11.959 0 0 1 3.598 6 11.99 11.99 0 0 0 3 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285Z" />
                        </svg>
                        <span className="inline-flex items-center justify-center px-2 py-0.5 ms-3 text-m font-large text-gray-500 rounded dark:bg-gray-700 dark:text-gray-400">{code}</span>
                      </div>
                    </li>)
                  }
                </ul>
              </div>
              }   
            </div>
          </div>
        </div>
      }
    </div>
  )
}