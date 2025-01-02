import { createContext } from "react";
import { IUser } from "../types";

export interface AuthContextState {
  authToken: string,
  setAuthToken: (authToken: string) => void,
  currentUser: IUser,
  setCurrentUser: (user: IUser) => void,
  loginSessionId: number,
  setLoginSessionId: (loginSessionId: number) => void
}

export function BuildAuthContextInitialState(): AuthContextState {
  return {
    authToken: "",
    setAuthToken: () => {},
    currentUser: {
      id: -1,
      username: "",
      password: "",
      twoFactorSecret: {
        String: "", 
        Valid: false,
      },
      isTwoFactorEnabled: {
        Bool: false,
        Valid: true,
      },
      createdAt: "",
      updatedAt: "",
    }, 
    setCurrentUser: () => {},
    loginSessionId: -1,
    setLoginSessionId: () => {},
  }
}

export const AuthContext = createContext<AuthContextState>(BuildAuthContextInitialState());