import { useState } from "react"  

async function verifyLogin(usrname: string, password: string): Promise<number> {
  if (usrname === "usrname" && password === "password") {
    return 1;
  }
  return -1;
}

async function loginAttempt(
  usrname: string,
  password: string,
  { setLoggedIn, setSessionID } : {
    setLoggedIn:  React.Dispatch<React.SetStateAction<boolean>>,
    setSessionID: React.Dispatch<React.SetStateAction<number>>
  }
): Promise<void> {

  let sessionID: number = await verifyLogin(usrname, password);
  if (sessionID == -1) {
    return
  } else {
    setLoggedIn(true);
    setSessionID(sessionID);
  }

}

function Login({ setLoggedIn, setSessionID } : {
    setLoggedIn:  React.Dispatch<React.SetStateAction<boolean>>,
    setSessionID: React.Dispatch<React.SetStateAction<number>>
}) {
  const [usrname, setUsrname] = useState("");
  const [password, setPassword] = useState("");

  return (
    <>
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
      <button onClick={()=>loginAttempt(usrname, password, { setLoggedIn, setSessionID })}/>
    </>
  )
} 

export default Login;
