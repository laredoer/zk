package zk

import (
	"errors"
	"time"
)

//Delete 删除节点
func (client *ZookeeperClient) Delete(path string) (err error) {
	if !client.isConnect {
		err = errors.New("链接已经关闭，判断节点是否存在失败")
	}
	if client.done {
		err = errors.New("链接已经手动关闭")
	}
	//获取此节点下所有子节点
	var d AllChildren
	client.GetAllChildren(path, &d)
	//反转数组
	reverse(d.Path)
	//开启协程删除节点
	ch := make(chan error, 1)
	go func(chan error) {
		for _, v := range d.Path {
			err = client.conn.Delete(v, -1)
		}
		ch <- err
	}(ch)
	select {
	case <-time.After(time.Second * 2):
		err = errors.New("删除节点超时")
		return
	case err = <-ch:
		return
	}
}

//reverse 反转数组
func reverse(x []string) {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}
}
