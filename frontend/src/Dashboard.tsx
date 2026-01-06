import { useState } from "react"
import "./main.css"

function Dashboard(
  {sessionID, codeRoomID, setLoggedIn, setSessionID, setCodeRoomID}: {
    sessionID:      string,
    codeRoomID:     string,
    setLoggedIn:    React.Dispatch<React.SetStateAction<boolean>>,
    setSessionID:   React.Dispatch<React.SetStateAction<string>>,
    setCodeRoomID:  React.Dispatch<React.SetStateAction<string>>,
  }
) {
  const [inputRoomID, setInputRoomID] = useState("");
  return (
    <>
      <div className="mainDiv" id="dashboardDiv">
        <div id="roomJoinDiv">
          <label htmlFor="roomIdInput">
            Room ID:
            <input name="roomIdInput"
              value={inputRoomID}
              onChange={e=>setInputRoomID(e.target.value.replace(/\D/g, ""))}
            />
          </label>
          <button onClick={()=>joinRoom(inputRoomID, sessionID, setCodeRoomID)}>
            Find room
          </button>
          <button onClick={()=>createRoom(sessionID, setCodeRoomID)}>
            Create room
          </button>
        </div>
      </div>
    </>
  )
}

function joinRoom(
  roomID: string,
  sessionID: string,
  setCodeRoomID: React.Dispatch<React.SetStateAction<string>>
): void {
  return
}

function createRoom(
  sessionID: string,
  setCodeRoomID: React.Dispatch<React.SetStateAction<string>>
): void {
  setCodeRoomID('000001')
}

export default Dashboard;
