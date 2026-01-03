import { useState } from 'react'
import './App.css'
import Login from './Login.tsx'

function App() {
  const [loggedIn, setLoggedIn] = useState(false);
  const [sessionID, setSessionID] = useState(0);

  if (!loggedIn) {
    return (
        <>
            <Login setLoggedIn={setLoggedIn} setSessionID={ setSessionID } />
        </>
    )
  } else if (loggedIn && sessionID == 0) {
    return (
        <>
            
        </>
    )
  }
}

export default App
