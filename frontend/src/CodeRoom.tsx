import React, { useEffect, useRef, useState } from "react"
import Editor from "@monaco-editor/react"
import * as Y from "yjs"
import { MonacoBinding } from "y-monaco"
import type { editor } from "monaco-editor"

import "./main.css"

function CodeRoom(
  {sessionID, codeRoomID, setCodeRoomID}: {
    sessionID: string,
    codeRoomID: string,
    setCodeRoomID: React.Dispatch<React.SetStateAction<string>>,
  },
) {
  const editorRef = useRef<editor.IStandaloneCodeEditor | null>(null)
  const [yDoc] = useState(() => new Y.Doc())
  const [isConnected, setIsConnected] = useState(false)
  const wsRef = useRef<WebSocket | null>(null)
  const bindingRef = useRef<MonacoBinding | null>(null)

  useEffect(() => {
    // Only setup WebSocket if we have a room ID
    if (!codeRoomID || !editorRef.current) return

    const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = `${wsProtocol}//${window.location.host}/api/ws?roomID=${codeRoomID}`

    const ws = new WebSocket(wsUrl)
    wsRef.current = ws

    ws.onopen = () => {
      console.log('WebSocket connected')
      setIsConnected(true)

      // Send join message
      const joinMsg = {
        SessionID: parseInt(sessionID, 16) || 0,
        Op: "create",
        RoomID: codeRoomID,
      }
      ws.send(JSON.stringify(joinMsg))

      // Setup Yjs update handling
      const yText = yDoc.getText('monaco')

      // Send local updates to server
      yDoc.on('update', (update: Uint8Array) => {
        if (ws.readyState === WebSocket.OPEN) {
          const updateMsg = {
            SessionID: parseInt(sessionID, 16) || 0,
            Op: "update",
            YjsBytes: Array.from(update),
          }
          ws.send(JSON.stringify(updateMsg))
        }
      })

      // Bind Monaco editor to Yjs
      if (editorRef.current) {
        bindingRef.current = new MonacoBinding(
          yText,
          editorRef.current.getModel()!,
          new Set([editorRef.current]),
        )
      }
    }

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data)
        if (msg.YjsBytes && msg.Op === "update") {
          // Apply remote updates
          const update = new Uint8Array(msg.YjsBytes)
          Y.applyUpdate(yDoc, update)
        }
      } catch (error) {
        console.error('Failed to parse message:', error)
      }
    }

    ws.onerror = (error) => {
      console.error('WebSocket error:', error)
      setIsConnected(false)
    }

    ws.onclose = () => {
      console.log('WebSocket disconnected')
      setIsConnected(false)
    }

    // Cleanup on unmount
    return () => {
      if (bindingRef.current) {
        bindingRef.current.destroy()
      }
      if (ws.readyState === WebSocket.OPEN) {
        const exitMsg = {
          SessionID: parseInt(sessionID, 16) || 0,
          Op: "exit",
          YjsBytes: [],
        }
        ws.send(JSON.stringify(exitMsg))
        ws.close()
      }
    }
  }, [codeRoomID, sessionID, yDoc])

  function handleEditorDidMount(editor: editor.IStandaloneCodeEditor) {
    editorRef.current = editor
  }

  return (
    <>
      <div style={{ padding: '10px' }}>
        <div style={{ marginBottom: '10px' }}>
          <span>Room ID: {codeRoomID} | </span>
          <span style={{ color: isConnected ? 'green' : 'red' }}>
            {isConnected ? 'Connected' : 'Disconnected'}
          </span>
          <button
            style={{ marginLeft: '10px' }}
            onClick={() => setCodeRoomID('')}
          >
            Leave Room
          </button>
        </div>
        <Editor
          height="60vh"
          defaultLanguage="python"
          defaultValue="# Welcome to Codelaborate!\n# Start typing to see real-time collaboration\n"
          onMount={handleEditorDidMount}
          options={{
            automaticLayout: true,
            minimap: { enabled: false },
          }}
        />
      </div>
    </>
  )
}

export default CodeRoom;
