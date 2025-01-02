import { IUser } from "../types";

export const setLocalStorageAuthToken = (authToken: string, setAuthToken: (setAuthToken: string) => void): void => {
  localStorage.setItem('authToken', JSON.stringify(authToken));
  setAuthToken(authToken);
}

export const deleteLocalStorageAuthToken = (setAuthToken: (authToken: string) => void): void => {
  localStorage.removeItem('authToken');
  setAuthToken("");
}

export const setLocalStorageCurrentUser = (currentUser: IUser, setCurrentUser: (currentUser: IUser) => void): void => {
  localStorage.setItem('currentUser', JSON.stringify(currentUser));
  setCurrentUser(currentUser);
}

export const deleteLocalStorageCurrentUser = (setCurrentUser: (currentUser: IUser) => void): void => {
  localStorage.removeItem('currentUser');
  setCurrentUser({
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
  });
}

export const isLocalStorageAvailable = (): boolean => {
  return typeof window !== 'undefined';
}