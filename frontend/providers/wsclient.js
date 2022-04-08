import React, { createContext } from "react";
import {useDispatch, useSelector} from "react-redux"
import {update} from '../services/gameStateSlice'

const WSClientContext = createContext(null)
export{WSClientContext}


export default ({children})=>{
    let wsclient
    const dispatch = useDispatch()
    if (typeof window !== "undefined" && !window.app_websocket){
        const token = useSelector(state=>state.profile.token)
        window.app_websocket = new WebSocket(`ws://localhost:8080/api/ws?token=${token??''}`)
        window.app_websocket.onmessage=e=>{
            if (!e.data.includes("invalid"))
                dispatch(update(JSON.parse(e.data)))
            /*
            if (e.ResponseType==='State'){
                dispatch({type:"update", payload:JSON.parse(e.Content)})
            }*/
        };
    }
    const sendMessage = (message) =>{
        window.app_websocket.send(message)
    }
    
    wsclient = {
        sendMessage: sendMessage
    }

    return (<WSClientContext.Provider value={wsclient}>{children}</WSClientContext.Provider>)
}