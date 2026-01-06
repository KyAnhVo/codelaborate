import React from "react"
import Editor from "@monaco-editor/react"

import "./main.css"

function CodeRoom(
  {sessionID, codeRoomID, setCodeRoomID}: {
    sessionID: string,
    codeRoomID: string,
    setCodeRoomID: React.Dispatch<React.SetStateAction<string>>,
  },
) {
  return (
    <>
      <Editor 
        height="60vh"
        defaultLanguage="python"
        defaultValue={"// hello\n"}
      />
    </>
  )
}

export default CodeRoom;
