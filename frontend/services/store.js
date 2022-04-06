import { configureStore} from '@reduxjs/toolkit'
import gameStateReduer from './gameStateSlice'
import profileReducer from './profileSlice'
import { composeWithDevTools } from 'redux-devtools-extension'
import { applyMiddleware } from '@reduxjs/toolkit'
import logger from 'redux-logger'
export default configureStore({
    reducer:{
        gameState:gameStateReduer,
        profile:profileReducer
    },
    enhancers: ()=>composeWithDevTools(applyMiddleware(logger))
})