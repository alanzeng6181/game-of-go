import { useSelector } from 'react-redux'
import {useRouter} from 'next/router'
import { useEffect } from 'react'

export default function Home() {
  const currentGameId = useSelector((state)=>state.profile.currentGameId)
  const token = useSelector((state)=>state.token)
  const router = useRouter()
  useEffect(()=>{
    if (currentGameId){
      router.push('/game')
    } else if (token){
      router.push('/profile')
    } else{
      router.push('/login')
    }
  }, [currentGameId, token])
  return <></>
}