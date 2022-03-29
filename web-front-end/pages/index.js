import Head from 'next/head'
import Image from 'next/image'
import styles from '../styles/Home.module.css'
import {Button, Space, DatePicker, Card} from 'antd';
import {CiCircleFilled} from '@ant-design/icons';

export default function Home() {
  const onChange =()=>{};
  return (
    //TODO: delete. just for checking ant design works.
    <div style={{padding:100}}>
      <Space direction='vertical'>
        <Button type="primary">Primary Button</Button>
        <Button type="ghost">Ghost Button</Button>
        <DatePicker onChange={onChange}></DatePicker>
        <CiCircleFilled />
      </Space>
    </div>
  )
}
