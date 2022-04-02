import Head from 'next/head'
import Image from 'next/image'
import styles from '../styles/Home.module.css'
import {Button, Space, DatePicker, Card} from 'antd';
import {CiCircleFilled} from '@ant-design/icons';
import Game from './game'
export default function Home() {
  const onChange =()=>{};
  return (
    <Game></Game>
  )
}