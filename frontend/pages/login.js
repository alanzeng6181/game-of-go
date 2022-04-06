
const Login = () => {
    const [crendential, setCredential] = useState({})
    const onChange = (e) => {

    }
    const handleClick = (e)=>{

    }

    return (
        <Row>
            <Col span={8}>
                    <Input placeholder='enter username' value={inputMessage} onPressEnter={handleClick} onChange={onChange}/>
            </Col>
            <Col span={8}>
                <Input placeholder='enter username' value={inputMessage} onPressEnter={handleClick} onChange={onChange}/>
            </Col>
            <Col span={8}>
                <Button type="primary" onClick={handleClick}>Login</Button>
            </Col>
        </Row>
        )
}