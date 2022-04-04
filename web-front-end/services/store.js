import { configureStore} from '@reduxjs/toolkit'
import gameStateReduer from './gameStateSlice'
export default configureStore({
    reducer:{
        gameState:gameStateReduer
    }
})