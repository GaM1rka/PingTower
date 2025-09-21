
import { useCheckers } from './api/checkers'
import AppBar from './components/appbar'
import Checker from './components/checker'
import Createchecker from './components/createChecker'

function App() {

  const { checkers } = useCheckers();

  return (
    <>
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
    </>
  )
}

export default App
