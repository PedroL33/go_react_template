import { OtpCodeLogin } from "components/verify_login/OtpCodeLogin";
import { RecoveryCodeLogin } from "components/verify_login/RecoveryCodeLogin";
import { useState } from "react";

export const VerifyLoginPage: React.FC = ( ) => {

  const [useRecoveryCode, setUserRecoveryCode] = useState<boolean>(false);

  return (
    <div className="h-screen flex items-center justify-center">
      {
        useRecoveryCode ? <RecoveryCodeLogin setUseRecoveryCode={setUserRecoveryCode}/>: <OtpCodeLogin setUseRecoveryCode={setUserRecoveryCode}/>
      }
    </div>
  )
}