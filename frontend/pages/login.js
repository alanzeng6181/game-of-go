import {useDispatch} from "react-redux"
import { login} from "../services/profileSlice"
import React, { useEffect, useState } from "react"
import {Row,Col, Input, Button} from 'antd'
import { useSelector } from "react-redux"
import { useRouter } from "next/router"

const Login = () => {
    const token = useSelector(state=> state.profile.token)
    const dispatch = useDispatch()
    const router = useRouter()
    const [username, setUsername] = useState()
    const [password, setPassword] = useState()
    const onUsernameChange = (e) => {
        setUsername(e.target.value)
    }
    const onPasswordChange = (e) => {
        setPassword(e.target.value)
    }
    const handleClick = (e)=>{
        var credential = {username:username, password:password}
        const loginThunk = login(credential)
        dispatch(loginThunk)
    }

    useEffect(()=>{
        if (token){
            router.push('/profile')
        }
    }, [token])

    return (
        <div>
            <Row justify="space-around" gutter={[30, 30]}>
                <Col span={4}  style={{marginTop:'200px', marginBottom:'30px'}}>
                    <img src="gameofgo.png"></img>
                </Col>
            </Row>
            <Row justify="space-around"  gutter={[30, 30]}>
                <Col span={2} >
                    <Row justify="end" align="middle" gutter={[15, 15]}>
                        <Col className="gutter-row" span={24}>
                            <Input placeholder='enter username' value={username} onPressEnter={handleClick} onChange={onUsernameChange}/>
                        </Col>
                        <Col className="gutter-row" span={24}>
                            <Input type="password" placeholder='enter password' value={password} onPressEnter={handleClick} onChange={onPasswordChange}/>
                        </Col>
                        <Col className="gutter-row" span={24}>
                            <Button type="primary" onClick={handleClick}>Login</Button>
                        </Col>
                    </Row>
                </Col>
            </Row>
        </div>
    )
}

export default Login