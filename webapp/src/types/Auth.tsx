export interface IUser {
  id: number,
  email: string,
  password: string,
  twoFactorSecret: {
    String: string,
    Valid: boolean,
  },
  isTwoFactorEnabled: {
    Bool: boolean,
    Valid: boolean,
  },
  createdAt: string,
  updatedAt: string,
}

export interface ILoginResponse {
  user?: IUser,
  token?: string,
  loginSessionId?: number
}

export interface ILoginRequest {
  email: string,
  password: string,
}

export interface ILoginWithOtpCodeRequest {
  otpCode: string,
  loginSessionId: number,
}

export interface ILoginWithRecoveryCodeRequest {
  recoveryCode: string,
  loginSessionId: number,
}

export interface IBegin2faResponse {
  base64Qrcode: string
}

export interface IComplete2faRequest {
  otpCode: string
}

export interface IComplete2faResponse {
  recoveryCodes: string[]
}

export interface IDisable2faRequest {
  password: string
}

export interface IChangePasswordRequest {
  currentPassword: string,
  newPassword: string
}

export interface IRegenerateRecoveryCodesResponse {
  recoveryCodes: string[]
}