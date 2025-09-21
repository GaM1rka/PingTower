
import { useCheckers } from './api/checkers'
import AppBar from './components/appbar'
import Checker from './components/checker'
import Createchecker from './components/createChecker'

function App() {

  const { checkers } = useCheckers();

   checkers.map(checker =>
    <Checker id={checker.id} url={checker.site} status={checker.status} />
  )
  return (
    <>
      <AppBar/>
      <Createchecker></Createchecker>
      <Checker id={1} url="https://example.com" status={'ok'} />
      {checkers}
    </>
  )
}

export default App
