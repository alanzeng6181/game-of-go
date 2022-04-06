import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";

const Login = createAsyncThunk('/login', async (credential)=>{
    return fetch('api/login', {
        method:'POST',
        headers:{
            'Content-Type':'application/json'
        },
        body: credential
    })
})

const ProfileSlice=createSlice({
    name:'Profile',
    extraReducers(builder){
        builder.addCase(Login.pending, (state,action) =>{})
        .addCase(Login.fulfilled, (state, action)=>{
            state.profile = action.payload
        })
        .addCase(Login.rejected, (state,action)=>{} )
    }
})

export {Login}
export default ProfileSlice.reducer