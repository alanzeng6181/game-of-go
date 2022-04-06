import React, { createContext } from "react";
import {useDispatch} from "react-redux"
import {update} from '../services/gameStateSlice'

const WSClientContext = createContext(null)
export{WSClientContext}


export default ({children})=>{
    let wsclient
    let ws
    const dispatch = useDispatch()
    if (typeof window !== "undefined"){
        ws = new WebSocket("ws://localhost:8080/api/ws?gameId=1&playerId=1")
        ws.onmessage=e=>{
            if (!e.data.includes("invalid"))
                dispatch(update(JSON.parse(e.data)))
            /*
            if (e.ResponseType==='State'){
                dispatch({type:"update", payload:JSON.parse(e.Content)})
            }*/
        };
    }
    const sendMessage = (message) =>{
        ws.send(message)
    }
    
    wsclient = {
        ws : ws,
        sendMessage: sendMessage
    }

    return (<WSClientContext.Provider value={wsclient}>{children}</WSClientContext.Provider>)
}