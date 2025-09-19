
import AppBar from './components/appbar'
import Checker from './components/checker'
import Createchecker from './components/createChecker'

function App() {

  return (
    <>
      <AppBar/>
      <Createchecker></Createchecker>
      <Checker id={1} url="https://example.com" status={true} />
    </>
  )
}

export default App
