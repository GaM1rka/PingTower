
import { useCheckers } from './hooks/checkers'
import AppBar from './components/appbar'
import Checker from './components/checker'
import Createchecker from './components/createChecker'
import { ToastContainer } from 'react-toastify'
import { useEffect, useState } from 'react'
function App() {

  const { checkers } = useCheckers();
  const [overlay, setOverlay] = useState(true);

  useEffect(() => {
      const value = localStorage.getItem("email");
      if (value) setOverlay(false);
    }, []);

  return (
    <>
      {overlay && (
        <div className='overlay' />
      )}
      <AppBar/>
      <Createchecker></Createchecker>
      <Checker id={1} url="https://example.com" status={'ok'} />
      {checkers.map(checker => (
        <Checker 
        key={checker.id}
        id={checker.id}
        url={checker.site}
        status={checker.status}
        />
      ))}
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
