import WSClientProvider from '../providers/wsclient'
import Game from './game'
import { useSelector } from 'react-redux'

export default function Home() {
  const gameId = useSelector((state)=>state.session.gameId) ?? new Array(361)
  
  return ( 
    <WSClientProvider>
      <Game></Game>
    </WSClientProvider>
  )
}