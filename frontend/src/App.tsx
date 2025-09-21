
import { useCheckers } from './hooks/checkers'
import AppBar from './components/appbar'
import Checker from './components/checker'
import Createchecker from './components/createChecker'
import { ToastContainer } from 'react-toastify'
import { useEffect } from 'react'

function App() {

  const { checkers, refresh } = useCheckers();

  useEffect(() => {
    const interval = setInterval(() => {
      refresh();
    }, 10000);
    return () => clearInterval(interval);
  }, [refresh]);

  return (
    <>
      <AppBar/>
      <Createchecker></Createchecker>
      <Checker id={1} url="https://example.com" status={'ok'} />
      <div className='centdiv'>
      {Array.isArray(checkers) && checkers.length > 0 ? (
        checkers.map(checker => (
          <Checker
          key={checker.id}
          id={checker.id}
          url={checker.site}
          status={checker.status}
          />
      ))) : (
      <p>Your checkers will appear here!</p>
      )}
      </div>
      <ToastContainer
        position="bottom-right"
        autoClose={2000}
        hideProgressBar={true}
        newestOnTop={true}
        closeOnClick
        theme="colored"
      />
    </>
  )
}

export default App
