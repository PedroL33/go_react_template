import { ChangePassword } from "api_calls/users";
import { AuthContext } from "contexts/AuthContext";
import { useContext, useState } from "react";
import { toast } from "react-toastify";
import { IChangePasswordRequest } from "types";

export const ChangePasswordModal: React.FC = () => {

  const { authToken } = useContext(AuthContext);
  const [showModal, setShowModal] = useState<boolean>(false);
  const [currentPassword, setCurrentPassword] = useState<string>("");
  const [newPassword, setNewPassword] = useState<string>("");
  const [retypeNewPassword, setRetypeNewPassword] = useState<string>("");
  const [errors, setErrors] = useState<Partial<IChangePasswordRequest>>();

  const updatePassword = async (e: any): Promise<void> => {
    e.preventDefault()
    try {
      setErrors(undefined)
      if(retypeNewPassword !== newPassword) {
        setErrors({newPassword: "passwords do not match"})
        return
      }
      const changePasswordRequest: IChangePasswordRequest = {
        currentPassword,
        newPassword
      }
      const response = await ChangePassword(authToken, changePasswordRequest)
      toast(response)
      setShowModal(false);
    }catch(err: any) {
      if(err.status == 422) {
        setErrors(err.error)
      }else {
        toast(err.message)
      }
    }
  }

  const handleClose = (): void => {
    setShowModal(false)
    setCurrentPassword("")
    setNewPassword("")
    setRetypeNewPassword("")
    setErrors(undefined);
  }

  return (
    <div>
      <button onClick={() => setShowModal(true)} className="block text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800" type="button">
        Change
      </button>
      {
        showModal && <div tabIndex={-1} className="overflow-y-auto overflow-x-hidden fixed top-0 right-0 left-0 z-50 justify-center items-center w-full md:inset-0 h-[calc(100%-1rem)] max-h-full">
          <div className="h-full flex items-center justify-center">
            <div className="max-w-2xl w-96 relative bg-white rounded-lg shadow dark:bg-gray-700 p-4">
              <div className="flex items-center justify-between p-4 mb-4 md:p-5 border-b rounded-t dark:border-gray-600">
                <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
                  Change Password
                </h3>
                <button type="button" onClick={handleClose} className="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm w-8 h-8 ms-auto inline-flex justify-center items-center dark:hover:bg-gray-600 dark:hover:text-white" data-modal-hide="default-modal">
                  <svg className="w-3 h-3" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 14 14">
                    <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"/>
                  </svg>
                  <span className="sr-only">Close modal</span>
                </button>
              </div>
              <form onSubmit={updatePassword} className="max-w-sm mx-auto">
                <div className="mb-5">
                  <label htmlFor="currentPassword" className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Current Password</label>
                  <input id="currentPassword" onChange={(e) => setCurrentPassword(e.target.value)} type="password" className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" placeholder="****" required />
                </div>
                <div className="mb-5">
                  <label htmlFor="newPassword" className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">New Password</label>
                  <input id="newPassword" onChange={(e) => setNewPassword(e.target.value)} type="password" className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" placeholder="****" required />
                  {errors?.newPassword && <p className="mt-2 text-sm text-red-600 dark:text-red-500"><span className="font-medium">Oops! </span>{errors.newPassword}</p>}
                </div>
                <div className="mb-5">
                  <label htmlFor="retypeNewPassword" className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Re-type New Password</label>
                  <input id="retypeNewPassword" onChange={(e) => setRetypeNewPassword(e.target.value)} type="password" className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" placeholder="****" required />
                </div>
                <div className="flex justify-center">
                  <button type="submit" className="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">Confirm</button>
                </div>
              </form>
            </div>
          </div>
        </div>
      }
    </div>
  )
}