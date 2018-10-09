package zk

import (
	"errors"
	"time"
)

//Update 更新数据
func (client *ZookeeperClient) Update(path string, data string, version int32) (err error) {
	if !client.isConnect {
		err = errors.New("链接已关闭，获取子节点失败")
		return
	}
	if client.done {
		err = errors.New("链接已手动关闭，获取子节点失败")
		return
	}
	b, err := client.Exists(path)
	if b == false || err != nil {
		err = errors.New("节点不存在")
		return
	}
	ch := make(chan error, 1)
	go func(chan error) {
		_, err = client.conn.Set(path, []byte(data), version)
		ch <- err
	}(ch)
	select {
	case <-time.After(time.Second * 2):
		err = errors.New("修改数据超时")
		return
	case err = <-ch:
		if err != nil {
			return err
		}
		return nil
	}
}
