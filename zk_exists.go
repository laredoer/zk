package zk

import (
	"errors"
	"time"
)

type existsType struct {
	b   bool
	err error
}

//Exists 判断节点是否存在
func (client *ZookeeperClient) Exists(path string) (b bool, err error) {
	if !client.isConnect {
		err = errors.New("链接已经关闭，判断节点是否存在失败")
	}
	if client.done {
		err = errors.New("链接已经手动关闭")
	}
	//开启协程判断链接是否存在
	ch := make(chan interface{}, 1)
	go func(chan interface{}) {
		b, _, err = client.conn.Exists(path)
		ch <- existsType{b: b, err: err}
	}(ch)
	select {
	case <-time.After(time.Second * 2):
		err = errors.New("判断节点是否存在超时")
		return
	case data := <-ch:
		b = data.(existsType).b
		err = data.(existsType).err
		if err != nil {
			return false, err
		}
		return b, nil
	}
}
