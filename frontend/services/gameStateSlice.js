import {createSlice} from "@reduxjs/toolkit"

export const GameStateSlice = createSlice({
    name : "gameState",
    initialState:{},
    reducers: {
        update: (state, action) =>{
            return action.payload
        }
    }
})

export const {update} =GameStateSlice.actions

export default GameStateSlice.reducer