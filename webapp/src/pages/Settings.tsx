import { AuthenticatedPage } from "components/AuthenticatedPage"
import { ChangePasswordModal } from "components/settings/ChangePasswordModal"
import { DisableTwoFactorModal } from "components/settings/DisableTwoFactorModal"
import { EnableTwoFactorModal } from "components/settings/EnableTwoFactorModal"
import { RegenerateRecoveryCodeModal } from "components/settings/RegenerateRecoveryCodes"
import { AuthContext } from "contexts/AuthContext"
import { useContext } from "react"

export const Settings = () => {

  const { currentUser } = useContext(AuthContext)

  return (
    <AuthenticatedPage>
      <div className="m-auto w-full flex items-center justify-center">
        <div className="w-full max-w-xl p-4 bg-white border border-gray-200 rounded-lg shadow sm:p-8 dark:bg-gray-800 dark:border-gray-700">
          <div className="flex items-center justify-between mb-4">
            <h5 className="text-xl font-bold leading-none text-gray-900 dark:text-white">Settings</h5>
          </div>
          <div className="flow-root">
            <ul role="list" className="divide-y divide-gray-200 dark:divide-gray-700">
              <li className="py-3 sm:py-4">
                <div className="flex items-center">
                  <div className="flex-1 min-w-0 ms-4">
                    <p className="text-sm font-medium text-gray-900 truncate dark:text-white">
                      Change Password
                    </p>
                  </div>
                  <ChangePasswordModal />
                </div>
              </li>
              <li className="pt-3 sm:py-4">
                <div className="flex items-center ">
                  <div className="flex-1 min-w-0 ms-4">
                    <p className="text-sm font-medium text-gray-900 truncate dark:text-white">
                      {
                        currentUser.isTwoFactorEnabled.Bool ? "Disable Two Factor Authentication": "Enable Two Factor Authentication"
                      }
                    </p>
                  </div>
                  <div className="inline-flex items-center text-base font-semibold text-gray-900 dark:text-white">
                    {
                      currentUser.isTwoFactorEnabled.Bool ? <DisableTwoFactorModal />: <EnableTwoFactorModal />
                    }
                  </div>
                </div>
              </li>
              {
                currentUser.isTwoFactorEnabled.Bool && <li className="pt-3 pb-0 sm:pt-4">
                  <div className="flex items-center ">
                    <div className="flex-1 min-w-0 ms-4">
                      <p className="text-sm font-medium text-gray-900 truncate dark:text-white">
                        Regenerate Recovery Codes
                      </p>
                    </div>
                    <div className="inline-flex items-center text-base font-semibold text-gray-900 dark:text-white">
                      <RegenerateRecoveryCodeModal />
                    </div>
                  </div>
                </li>
              }
            </ul>
          </div>
        </div>
      </div>
    </AuthenticatedPage>
  )
}