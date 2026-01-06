import { useState } from "react"
import "./Login.css"
import "./main.css"

function Login(
  {setLoggedIn, setSessionID} : {
    setLoggedIn:  React.Dispatch<React.SetStateAction<boolean>>,
    setSessionID: React.Dispatch<React.SetStateAction<string>>,
  },
) {
  const [usrname, setUsrname] = useState("");
  const [password, setPassword] = useState("");

  return (
    <>
      <div className="mainDiv" id="loginDiv">
        <label htmlFor="usernameInput">
          USERNAME:
          <input name="usernameInput"
            value={usrname}
            onChange={e => setUsrname(e.target.value)}
          />
        </label>
        <label htmlFor="passwordInput">
          PASSWORD:
          <input name="passwordInput"
            value={password}
            onChange={e => setPassword(e.target.value)}
          />
        </label>
        <button onClick={()=>loginAttempt(usrname, password, { setLoggedIn, setSessionID })}>
          Login
        </button>
      </div>
    </>
  )
} 

async function verifyLogin(usrname: string, password: string): Promise<string> {
  if (usrname === "usrname" && password === "password") {
    return "000001";
  }
  return "";
}

async function loginAttempt(
  usrname: string,
  password: string,
  { setLoggedIn, setSessionID } : {
    setLoggedIn:  React.Dispatch<React.SetStateAction<boolean>>,
    setSessionID: React.Dispatch<React.SetStateAction<string>>
  }
): Promise<void> {

  let sessionID: string = await verifyLogin(usrname, password);
  if (sessionID.length === 0) {
    return
  } else {
    setLoggedIn(true);
    setSessionID(sessionID);
  }

}

export default Login;
