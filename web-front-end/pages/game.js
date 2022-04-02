import Board from "./board"
import Messenger from "./messegener"
import {Row, Col} from 'antd'

const Game = ()=>{
    return (
            <Row>
            <Col span={4}>
            </Col>
            <Col span={15}>
                <Board></Board>
            </Col>
            <Col span={5}>
                <Row style={{height:500}}>

                </Row>
                <Messenger></Messenger>
            </Col>
        </Row> )
}

export default Game