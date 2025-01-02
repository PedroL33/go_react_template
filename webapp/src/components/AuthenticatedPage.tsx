import { useContext } from "react";
import { Navbar } from "./Navbar";

interface AuthenticatedPageProps {
  children: React.ReactNode
}

export const AuthenticatedPage: React.FC<AuthenticatedPageProps> = ({children}) => {

  return (
    <div className="min-h-screen min-w-screen flex justify-items-stretch flex-col">
      <Navbar />
      {children}
    </div>
  )
}