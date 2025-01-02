import { AuthenticatedPage } from '../components/AuthenticatedPage';

export const Dashboard = () => {

  return (
    <AuthenticatedPage>
      <div className="p-10 m-auto flex flex-col items-center justify-center bg-cool-gray-700">
        <h1 className="text-9xl text-center mb-28">
          <div className="bg-gradient-to-r text-transparent bg-clip-text from-green-400 to-purple-500">
            Dashboard
          </div>
        </h1>
      </div>
    </AuthenticatedPage>
  )
}