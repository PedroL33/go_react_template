import { RegenerateRecoveryCodes } from "api_calls/users";
import { AuthContext } from "contexts/AuthContext";
import React, { useContext, useState } from "react";
import { toast } from "react-toastify";

export const RegenerateRecoveryCodeModal: React.FC = () => {

  const [showModal, setShowModal] = useState<boolean>(false)
  const [confirmMessage, setConfirmMessage] = useState<string>("");
  const [error, setError] = useState<string>("");
  const [recoveryCodes, setRecoveryCodes] = useState<string[]>([]);

  const { authToken } = useContext(AuthContext);

  const regenerateCodes = async (e: any) => {
    e.preventDefault();

    try {
      if(confirmMessage !== "regenerate") {
        setError("Please type the confirm message.")
        return
      }
      const response = await RegenerateRecoveryCodes(authToken);
      setRecoveryCodes(response.recoveryCodes);
      toast("Successfully regenerated recovery codes.")
    }catch(err: any) {
      toast(err.message)
    }
  }

  const handleClose = () => {
    setShowModal(false);
    setConfirmMessage("");
    setRecoveryCodes([]);
    setError("");
  }

  return (
    <div>
      <button onClick={() => setShowModal(true)} className="block text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800" type="button">
        Regenerate
      </button>
      {
        showModal && <div tabIndex={-1} className="overflow-y-auto overflow-x-hidden fixed top-0 right-0 left-0 z-50 justify-center items-center w-full md:inset-0 h-[calc(100%-1rem)] max-h-full">
          <div className="h-full flex items-center justify-center">
            <div className="max-w-2xl relative bg-white rounded-lg shadow dark:bg-gray-700">
              <div className="flex items-center justify-between p-4 md:p-5 border-b rounded-t dark:border-gray-600">
                <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
                  Regenerate Recovery Codes
                </h3>
                <button type="button" onClick={handleClose} className="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm w-8 h-8 ms-auto inline-flex justify-center items-center dark:hover:bg-gray-600 dark:hover:text-white">
                  <svg className="w-3 h-3" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 14 14">
                    <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"/>
                  </svg>
                  <span className="sr-only">Close modal</span>
                </button>
              </div>
              {
                recoveryCodes.length ? <div>
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
                </div>:
                <div>
                  <div className="p-4 md:p-5 mb-5 mt-10 space-y-4">
                    <p className="text-base leading-relaxed text-gray-500 dark:text-gray-400">
                      Are you sure you want to regenerate your recovery codes?
                    </p>
                    <p className="text-base leading-relaxed text-gray-500 dark:text-gray-400 text-center">
                      Type "regenerate" to confirm.
                    </p>
                  </div>
                  <form onSubmit={regenerateCodes} className="max-w-sm mx-auto pb-10">
                    <div className="mb-5">
                      {error.length > 0 && <p className="mt-2 text-sm text-red-600 dark:text-red-500"><span className="font-medium">Oops! </span>{error}</p>}
                      <input type="text" onChange={(e) => setConfirmMessage(e.target.value)} className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" placeholder='Type "regenerate" to confirm' required />
                    </div>
                    <div className="flex justify-center">
                      <button type="submit" className="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">Confirm</button>
                    </div>
                  </form>
                </div>
              }
            </div>
          </div>
        </div>
      }
    </div>
  )
}