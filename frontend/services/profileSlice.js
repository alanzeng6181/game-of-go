import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import thunk from "redux-thunk";

const login = createAsyncThunk('/login', async (credential, thunkAPI)=>{
    const response = await fetch('api/login', {
        method:'POST',
        headers:{
            'Content-Type':'application/json'
        },
        body: JSON.stringify(credential)
    });
    if (response.ok){
        return response.text();
    } else{
        return thunkAPI.rejectWithValue(response.status==500 ? response.statusText : response.json())
    }
})

const loadProfile = createAsyncThunk('/getProfile', async ( token,thunkAPI)=>{
    //TODO: state is undefined
    const state = thunkAPI.getState();
    const response = await fetch('api/profile', {
        method:'GET',
        headers:{
            'Content-Type':'application/json',
            'Token': token ? `${token}` : null
        }
    });
    if (response.ok){
        return response.json()
    } else{
        return thunkAPI.rejectWithValue(response.status==500 ? response.statusText : response.json())
    }
})

const findGame = createAsyncThunk('/findgame', async(args,thunkAPI)=>{
    const state = thunkAPI.getState();
    const response = await fetch('/api/findgame', {
        method:'POST',
        headers:{
            'Token': args.token??''
        },
        body: JSON.stringify(args.gameRequest)
    });
    if (response.ok){
        return response.json()
    } else{
        return thunkAPI.rejectWithValue(response.status==500 ? response.statusText : response.json())
    }
})

const ProfileSlice = createSlice({
    name:'profile',
    initialState:{},
    extraReducers(builder){
        builder.addCase(login.pending, (state,action) =>{

        })
        .addCase(login.fulfilled, (state, action)=>{
            state.token = action.payload;
        })
        .addCase(login.rejected, (state,action)=>{} )
        .addCase(loadProfile.fulfilled, (state, action)=>{
            state.profile = action.payload;
        })
        .addCase(findGame.fulfilled, (state, action)=>{
            state.profile.CurrentGameId = action.payload.gameId;
        })
    }
})

export {login, loadProfile, findGame}
export default ProfileSlice.reducer