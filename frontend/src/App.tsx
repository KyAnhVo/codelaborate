import React, { useState } from 'react'
import Login from './Login.tsx'
import Dashboard from './Dashboard.tsx';

function App() {
  const [loggedIn, setLoggedIn] = useState(false);
  const [sessionID, setSessionID] = useState(0);
  const [codeRoomID, setCodeRoomID] = useState(0);

  if (!loggedIn) {
    return (
      <>
        <button onClick={()=>reset(setLoggedIn, setSessionID, setCodeRoomID)}>Reset</button>
        <Login setLoggedIn={setLoggedIn} setSessionID={ setSessionID } />
      </>
    )
  } else {
    return (
      <>
        <button onClick={()=>reset(setLoggedIn, setSessionID, setCodeRoomID)}>Reset</button>
        <Dashboard 
          sessionID={sessionID}
          codeRoomID={codeRoomID}
          setLoggedIn={setLoggedIn} 
          setSessionID={setSessionID}
          setCodeRoomID={setCodeRoomID}
        />
      </>
    )
  }
}

function reset(
  setLoggedIn: React.Dispatch<React.SetStateAction<boolean>>,
  setSessionID: React.Dispatch<React.SetStateAction<number>>,
  setCodeRoomID: React.Dispatch<React.SetStateAction<number>>,
) {
  setLoggedIn(false);
  setSessionID(0);
  setCodeRoomID(0);
}

export default App
