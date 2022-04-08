import { combineReducers, configureStore, createStore} from '@reduxjs/toolkit'
import gameStateReduer from './gameStateSlice'
import profileReducer from './profileSlice'
import { composeWithDevTools } from 'redux-devtools-extension'
import { applyMiddleware } from '@reduxjs/toolkit'
import thunkMiddleware from 'redux-thunk'
import logger from 'redux-logger'
import {
    persistStore,
    persistReducer,
    FLUSH,
    REHYDRATE,
    PAUSE,
    PERSIST,
    PURGE,
    REGISTER,} from 'redux-persist'
import storage from 'redux-persist/lib/storage'

const persistConfig = {
    key:'root',
    storage,
}

const persistedReducer = persistReducer(persistConfig, combineReducers({
    gameState:gameStateReduer,
    profile:profileReducer
}));
const store = configureStore({
    reducer:persistedReducer,
    enhancer:composeWithDevTools(applyMiddleware(thunkMiddleware,logger)),
    middleware: (getDefaultMiddleware) =>
        getDefaultMiddleware({
            serializableCheck: {
                ignoredActions: [FLUSH, REHYDRATE, PAUSE, PERSIST, PURGE, REGISTER],
            },
        }),
});

const persistor = persistStore(store)

export {store, persistor}