import { useEffect, useState } from "react"
import "./main.css"

type Friend = {
  userId: number,
  name: string,
  isOnline: boolean,
};

function reloadFriendList(
  sessionID:      number, 
  friendList:     Friend[],
  setFriendList:  React.Dispatch<React.SetStateAction<Friend[]>>
) {
  let newFriendList: Friend[] = [];
  friendList.forEach(friend=>{
    let newFriend: Friend = {
      userId: friend.userId,
      name: friend.name,
      isOnline: !friend.isOnline,
    };
    newFriendList.push(newFriend);
  })
  setFriendList(newFriendList);
}

function FriendList({sessionID}: {sessionID: number}) {
  const [friendList, setFriendList] = useState<Friend[]>([
    {userId: 0, name: 'Bob', isOnline: false},
    {userId: 1, name: 'Peter', isOnline: true},
  ]);
  const MAX_NAME_LENGTH: number = 30;
  const MAX_STAT_LENGTH: number = 15;

  useEffect(()=>{
    const id = setInterval(()=>{
      reloadFriendList(sessionID, friendList, setFriendList)
    }, 3000);
    return ()=>clearInterval(id);
  })

  return (
    <>
      <p>Friend List</p>
      <div>
        <span>{ 'Name' + ' '.repeat(MAX_NAME_LENGTH - 4) }</span>
        <span>{ 'Online status' + ' '.repeat(MAX_STAT_LENGTH - 13) }</span>
      </div>
      {friendList.map((friend) => { return (
          <div>
          <span>
            { friend.name + ' '.repeat(MAX_NAME_LENGTH - friend.name.length) }
          </span>
          <span>
            { friend.isOnline + ' '.repeat(MAX_STAT_LENGTH - friend.isOnline.toString().length ) }
          </span>
        </div>
      )})}
    </>
  )
}

function Dashboard({sessionID, codeRoomID, setLoggedIn, setSessionID, setCodeRoomID}: {
  sessionID:      number,
  codeRoomID:     number,
  setLoggedIn:    React.Dispatch<React.SetStateAction<boolean>>,
  setSessionID:   React.Dispatch<React.SetStateAction<number>>,
  setCodeRoomID:  React.Dispatch<React.SetStateAction<number>>,
}) {
  return (
    <>
      <div className="mainDiv" id="dashboardDiv">
        <div className="scrollableArea" id="friendList">
          <FriendList sessionID={sessionID}/>
        </div>
      </div>
    </>
  )
}

export default Dashboard;
