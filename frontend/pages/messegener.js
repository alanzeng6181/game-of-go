import {Row, Col, List, Skeleton, Avatar, Input, Button} from 'antd'
import InfiniteScroll from 'react-infinite-scroll-component'
import React, {useState, useRef} from 'react'
import { useSelector } from 'react-redux'
import {useContext} from 'react'
import { WSClientContext } from '../providers/wsclient'

const Messenger = (props)=>{
    const sendMessage = (useContext(WSClientContext)).sendMessage;
    const [inputMessage, setInputMessage] = useState()
    const messages = useSelector(
        state=>state.gameState.Comments
        ) ?? [];
    const handleClick=()=>{
        sendMessage(JSON.stringify({
            "Command":"Comment",
            "Arguments":[`${props.gameId}`,`${inputMessage}`]
        }));
    };

    const onChange = (e)=>{
        setInputMessage(e.target.value)
    }
    const loadMoreData=()=>{
    }

    return(
        <div>
            <Row>
                <Col span={24}>
                <div>
                    <InfiniteScroll dataLength={messages.length}
                        next={loadMoreData}
                        height={600}
                        hasMore={messages.length<500}>
                        <List
                            itemLayout='horizontal'
                            dataSource={messages}
                            split = {false}
                            renderItem={(item, index) => {
                            return index % 2 === 0 ? (
                                <List.Item>
                                    <Row style={{width:'100%'}}>
                                        <Col span={2}>
                                            <Avatar src="https://joeschmoe.io/api/v1/random" />
                                        </Col>
                                        <Col span={15}>
                                            <span>{item.content}</span>
                                        </Col>
                                    </Row>
                                </List.Item>
                            ) : (
                                <List.Item>
                                    <Row style={{width:'100%'}}>
                                        <Col span={15} offset={7} style={{textAlign:'right'}}>
                                            <span>{item.content}</span>
                                        </Col>
                                        <Col span={2}>
                                            <Avatar src="https://joeschmoe.io/api/v1/random" />
                                        </Col>
                                    </Row>
                                </List.Item>
                            );
                            }}
                        />
                    </InfiniteScroll>
                    </div>
                </Col>
            </Row>
            <Row>
                <Col span={18}>
                    <Input placeholder='type in message' value={inputMessage} onPressEnter={handleClick} onChange={onChange}/>
                </Col>
                <Col span={6}>
                    <Button type="primary" onClick={handleClick}>Send</Button>
                </Col>
            </Row>
        </div>
    )
}

export default Messenger