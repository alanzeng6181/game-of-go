import '../styles/globals.css'
import 'antd/dist/antd.css'
import { Provider } from 'react-redux';
import {store, persistor} from '../services/store'
import { PersistGate } from 'redux-persist/integration/react';

function MyApp({ Component, pageProps }) {
  
  return (
    <Provider store={store}>
      <PersistGate persistor={persistor}>
        <Component {...pageProps} />
      </PersistGate>
    </Provider>
  )
}

export default MyApp
