import WSClientProvider from '../providers/wsclient'
import Game from './game'

export default function Home() {
  return ( 
    <WSClientProvider>
      <Game></Game>
    </WSClientProvider>
  )
}