import { Button, Row, Col, Spin } from "antd"
import { useEffect } from "react"
import { useSelector, useDispatch } from "react-redux"
import {loadProfile} from '../services/profileSlice'
import { useRouter } from "next/router"
import WSClientProvider from '../providers/wsclient'
import { findGame } from "../services/profileSlice"
const Profile = ()=>{
    const route = useRouter()
    const [profile, token] = useSelector(
        (state)=>[state.profile.profile, state.profile.token]
        )
    const dispatch = useDispatch()
    useEffect(()=>{
        if (!profile){
            dispatch(loadProfile(token))
        }
    }, [profile, token])

    const goToGame =()=>{
        route.push(`/game/${profile.CurrentGameId}`)
    }

    const onFindGame = () =>{
        dispatch(findGame({gameRequest:{
            PlayerRank: profile.Level,
            RankAbove:0,
            RankBelow:0,
            BoardSize:19,
            TimeoutInSeconds:300,
            OverTimeInSeconds:30,
            OverTimeCount:3
        }, token}))
    }

    if (profile){
    return (
        <WSClientProvider>
            <Row justify="center">
                <Col span={3} style={{marginTop:'300px'}}>
                    <Row>
                        <Col span={12}>
                            <b>User Id:</b>
                        </Col>
                        <Col span={12}>
                            {profile.UserId}
                        </Col>
                    </Row>
                    <Row>
                        <Col span={12}>
                            <b>Level:</b>
                        </Col>
                        <Col span={12}>
                            {profile.Level}
                        </Col>
                    </Row>
                    <Row>
                        <Col span={12}>
                            <b>Wins:</b>
                        </Col>
                        <Col span={12}>
                            {profile.Wins}
                        </Col>
                    </Row>
                    <Row>
                        <Col span={12}>
                            <b>Losses:</b>
                        </Col>
                        <Col span={12}>
                            {profile.Losses}
                        </Col>
                    </Row>

                    <Row>
                        <Col span={12}>
                            {profile.CurrentGameId ?
                            <Button onClick={goToGame}>Go to game</Button>
                            :
                            <Button onClick={onFindGame}>Find Game</Button>
                            }
                        </Col>
                    </Row>
                </Col>
            </Row>
        </WSClientProvider>)
        } else{
            return (<Row justify="center" align="center">
                    <Col span={2} style={{marginTop:'300px'}}>
                        <Spin size="large" tip="Loading user profile..."></Spin>
                    </Col>
                </Row>)
        }
}

export default Profile