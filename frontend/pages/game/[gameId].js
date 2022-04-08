import Board from "../board"
import Messenger from "../messegener"
import {Row, Col} from 'antd'
import WSClientProvider from '../../providers/wsclient'
import { useRouter } from "next/router"

const Game = ()=>{
    const route = useRouter()
    const {gameId}=route.query
    return (
        <WSClientProvider>
            <Row>
                <Col span={4}>
                </Col>
                <Col span={15}>
                    <Board gameId={gameId}></Board>
                </Col>
                <Col span={5}>
                    <Row style={{height:500}}>
                    </Row>
                    <Messenger gameId={gameId}></Messenger>
                </Col>
            </Row> 
        </WSClientProvider>)
}

export default Game