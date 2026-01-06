import React, { useState } from 'react'

import Login      from './Login.tsx'
import Dashboard  from './Dashboard.tsx'
import CodeRoom   from './CodeRoom.tsx'

function App() {
  const [loggedIn, setLoggedIn] = useState<boolean>(false);
  const [sessionID, setSessionID] = useState<string>("");
  const [codeRoomID, setCodeRoomID] = useState<string>("");

  if (!loggedIn) {

    return (
      <>
        <a href="#" onClick={(e)=>{
          e.preventDefault();
          reset(setLoggedIn, setSessionID, setCodeRoomID);
        }}>Reset</a>
        <Login setLoggedIn={setLoggedIn} setSessionID={ setSessionID } />
      </>
    )

  } else if (codeRoomID.length === 0) {

    return (
      <>
        <a href="#" onClick={(e)=>{
          e.preventDefault();
          reset(setLoggedIn, setSessionID, setCodeRoomID);
        }}>Reset</a>
        <Dashboard 
          sessionID={sessionID}
          codeRoomID={codeRoomID}
          setLoggedIn={setLoggedIn} 
          setSessionID={setSessionID}
          setCodeRoomID={setCodeRoomID}
        />
      </>
    )

  } else {
    
    return (
      <>
        <a href="#" onClick={(e)=>{
          e.preventDefault();
          reset(setLoggedIn, setSessionID, setCodeRoomID);
        }}>Reset</a>
        <CodeRoom
          sessionID={sessionID}
          codeRoomID={codeRoomID}
          setCodeRoomID={setCodeRoomID}
        />
      </>
    )

  }
}

function reset(
  setLoggedIn: React.Dispatch<React.SetStateAction<boolean>>,
  setSessionID: React.Dispatch<React.SetStateAction<string>>,
  setCodeRoomID: React.Dispatch<React.SetStateAction<string>>,
) {
  setLoggedIn(false);
  setSessionID("");
  setCodeRoomID("");
}

export default App
