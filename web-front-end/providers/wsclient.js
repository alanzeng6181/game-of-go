import { createContext } from "react";
import {useDispatch} from "react-redux"
import io from 'socket.io-client'
const WSClientContext = createContext(null)
export{WSClientContext}

export default ({children})=>{
    let wsclient
    const dispatch = useDispatch()
    const ws = io("ws://localhost:8080",{path:"/ws?gameId=1&playerId=1", transports: ["websocket"],})
    ws.on("event://get-message",e=>{
        if (e.ResponseType==='State'){
            dispatch({type:"update", payload:JSON.parse(e.Content)})
        }
    })
    const sendMessage = (message) =>{
        ws.emit("event://send-message",message)
    }
    wsclient = {
        ws : ws,
        sendMessage: sendMessage
    }

    return (<WSClientContext.Provider value={wsclient}>{children}</WSClientContext.Provider>)
}