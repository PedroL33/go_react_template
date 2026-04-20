import axios from "axios";
import { 
  IBegin2faResponse, 
  IChangePasswordRequest, 
  IComplete2faRequest, 
  IComplete2faResponse, 
  IDisable2faRequest, 
  ILoginRequest, 
  ILoginResponse, 
  ILoginWithOtpCodeRequest, 
  ILoginWithRecoveryCodeRequest, 
  IRegenerateRecoveryCodesResponse 
} from "types";
import { SERVER_URL } from "utils/Environment";
import { parseApiError } from "utils/Error";

export const Login = async (loginRequest: ILoginRequest): Promise<ILoginResponse> => {
  try {
    const response = await axios.post(`${SERVER_URL}/login`, loginRequest, {
      headers: {
        "Content-Type": "application/json",
      },
    })
    return response.data.data;
  }catch(err: any) {
    throw parseApiError(err)
  }
};

export const LoginWithOptCode = async (verifyLoginRequest: ILoginWithOtpCodeRequest): Promise<ILoginResponse> => {
  try {
    const response = await axios.post(`${SERVER_URL}/login_otp`, verifyLoginRequest, {
      headers: {
        "Content-Type": "application/json",
      },
    })

    return response.data.data;
  }catch(err: any) {
    throw parseApiError(err)
  }
};

export const LoginWithRecoveryCode = async (verifyLoginRequest: ILoginWithRecoveryCodeRequest): Promise<ILoginResponse> => {
  try {
    const response = await axios.post(`${SERVER_URL}/login_recovery_code`, verifyLoginRequest, {
      headers: {
        "Content-Type": "application/json",
      },
    })

    return response.data.data;
  }catch(err: any) {
    throw parseApiError(err)
  }
};

export const EnableTwoFactorAuth = async (authToken: string): Promise<IBegin2faResponse> => {
  try {
    const response = await axios.post(`${SERVER_URL}/begin2fa`, undefined, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${authToken}`
      },
    })

    return response.data.data;
  }catch(err: any) {
    throw parseApiError(err)
  }
}

export const CompleteTwoFactorAuth = async (authToken: string, complete2faSetupRequest: IComplete2faRequest): Promise<IComplete2faResponse> => {
  try {
    const response = await axios.post(`${SERVER_URL}complete2fa`, complete2faSetupRequest, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${authToken}`
      },
    })
    
    return response.data.data
  }catch(err) {
    throw parseApiError(err);
  }
}

export const DisableTwoFactorAuth = async (authToken: string, disable2faRequest: IDisable2faRequest): Promise<string> => {
  try {
    const response = await axios.post(`${SERVER_URL}/disable2fa`, disable2faRequest, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${authToken}`
      },
    })

    return response.data.data
  }catch(err: any) {
    throw parseApiError(err);
  }
}

export const ChangePassword = async (authToken: string, changePasswordRequest: IChangePasswordRequest): Promise<string> => {
  try {
    const response = await axios.put(`${SERVER_URL}/change_password`, changePasswordRequest, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${authToken}`
      },
    })

    return response.data.data;
  }catch(err: any) {
    throw parseApiError(err);
  }
}

export const RegenerateRecoveryCodes = async (authToken: string): Promise<IRegenerateRecoveryCodesResponse> => {
  try {
    const response = await axios.post(`${SERVER_URL}/regenerate_recovery_codes`, undefined, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${authToken}`
      },
    })

    return response.data.data;
  }catch(err: any) {
    throw parseApiError(err);
  }
}