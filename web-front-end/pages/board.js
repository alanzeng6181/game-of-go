import {useSelector} from 'react-redux'
import {useRef, useEffect, useContext} from 'react'
import styles from '../styles/Board.module.css'
import {WSClientContext} from '../providers/wsclient'
const Board = props =>{
    const sendMessage = (useContext(WSClientContext)).sendMessage
    const arr = useSelector((state)=>state.gameState.Positions) ?? new Array(361)
    const size =props.boardRows ?? 1000;
    const boardRows = props.boardRows ?? 19;
    const margin = Math.ceil(size * 0.05);
    const space = Math.floor((size - 2*margin)/(boardRows-1))

    const handleClick = (e) =>{
        const bounds = e.target.getBoundingClientRect()
        const X = e.clientX -bounds.left
        const Y = e.clientY - bounds.top
        const row = Math.round((X - margin)/space)
        const col = Math.round((Y -margin)/space)

        const intersectX=margin+row*space
        const intersectY=margin+col*space

        const sensitivity = 0.38*space

        const distance = Math.sqrt(Math.pow(X-intersectX,2)+Math.pow(Y-intersectY,2))

        if (distance<=sensitivity){      
            sendMessage(JSON.stringify({
                "Command":"Move",
                "Arguments":[row, col]
            }))
        } else{
            console.log(`clicked, but position unambigous`);
        }
    }

    const boardCanvasRef=useRef(null)
    const gridCanvasRef=useRef(null)
    const stoneCanvasRef=useRef(null)
    useEffect(()=>{
        
        const boardCanvas = boardCanvasRef.current
        const boardCanvasContext = boardCanvas.getContext('2d')
        const background = new Image()
        background.src = "woodveneer.png"
        background.onload=()=>{
            boardCanvasContext.drawImage(background,0,0, boardCanvas.width, boardCanvas.height)
        }
        
        const gridCanvas = gridCanvasRef.current;
        const gridContext = gridCanvas.getContext('2d');

        gridContext.lineWidth=1;
        gridContext.beginPath();
        for(let i=0; i<boardRows; i++){
            gridContext.moveTo(margin+i*space, margin)
            gridContext.lineTo(margin+i*space, margin+(boardRows-1)*space)

            gridContext.moveTo(margin, margin+i*space)
            gridContext.lineTo(margin+(boardRows-1)*space, margin+i*space)
        }
        gridContext.stroke();

        const stoneCanvas = stoneCanvasRef.current
        const stoneCanvasContext = stoneCanvas.getContext('2d')
        //stoneCanvasContext.globalAlpha=0.3
        //TODO: null handling
        
        const stoneSize=Math.floor(space*0.92);

        const placeStones = (image, stoneColor, positionArray)=>{
            positionArray.forEach((stone,index)=>{
                if (stone===stoneColor){
                    const row = Math.floor(index/boardRows);
                    const col = Math.floor(index%boardRows);
                    stoneCanvasContext.drawImage(image, margin+row*space - Math.floor(stoneSize/2), margin+col*space- Math.floor(stoneSize/2),  stoneSize, stoneSize)
                }
            })
        }

        const blackStoneImage = new Image(stoneSize, stoneSize);
        blackStoneImage.src = "blackstone.png"
        blackStoneImage.onload=()=> placeStones(blackStoneImage, 1, arr);

        const whiteStoneImage = new Image(stoneSize, stoneSize);
        whiteStoneImage.src = "whitestone.png"
        whiteStoneImage.onload=()=>placeStones(whiteStoneImage, 2, arr)
    })
    return  (<div>
                <canvas className={styles.bottomLayer} width={1000} height={1000} ref={boardCanvasRef} ></canvas>
                <canvas className={styles.middleLayer} width={1000} height={1000} ref={gridCanvasRef} ></canvas>
                <canvas onClick={handleClick} className={styles.topLayer} width={1000} height={1000} ref={stoneCanvasRef} ></canvas>
           </div>)
}

export default Board