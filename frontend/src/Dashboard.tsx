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

async function joinRoom(
  roomID: string,
  sessionID: string,
  setCodeRoomID: React.Dispatch<React.SetStateAction<string>>
): Promise<void> {
  try {
    const response = await fetch(`/api/rooms/${roomID}`);
    if (response.ok) {
      const data = await response.json();
      if (data.exists) {
        setCodeRoomID(roomID);
      } else {
        alert('Room not found');
      }
    } else {
      alert('Room not found');
    }
  } catch (error) {
    console.error('Join room error:', error);
    alert('Failed to join room');
  }
}

async function createRoom(
  sessionID: string,
  setCodeRoomID: React.Dispatch<React.SetStateAction<string>>
): Promise<void> {
  try {
    const response = await fetch('/api/rooms', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ sessionID }),
    });

    if (response.ok) {
      const data = await response.json();
      setCodeRoomID(data.roomID);
    } else {
      alert('Failed to create room');
    }
  } catch (error) {
    console.error('Create room error:', error);
    alert('Failed to create room');
  }
}

export default Dashboard;
