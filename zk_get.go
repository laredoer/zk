package zk

import (
	"errors"
	"fmt"
	"time"
)

type getType struct {
	data    []byte
	version int32
	err     error
}

//Get 获取节点的值
func (client *ZookeeperClient) Get(path string) (value string, version int32, err error) {
	if !client.isConnect {
		err = errors.New("链接已断开，获取节点失败")
		return
	}
	if client.done {
		err = errors.New("链接已经手动关闭，获取节点失败")
	}
	//用协程获取节点的值，超时则退出
	ch := make(chan interface{}, 1)
	go func(ch chan interface{}) {
		data, stat, err := client.conn.Get(path)
		ch <- getType{data: data, version: stat.Version, err: err}
	}(ch)
	select {
	//2秒超时
	case <-time.After(time.Second * 2):
		err = errors.New("获取节点的数据超时")
		return
	case data := <-ch:
		if client.done {
			err = errors.New("链接已经手动关闭")
			return
		}
		err = data.(getType).err
		if err != nil {
			err = fmt.Errorf("get node:%s error(err:%v)", path, err)
			return
		}
		value = string(data.(getType).data)
		version = data.(getType).version
		return
	}
}

type GetChildrenType struct {
	data    []string
	version int32
	err     error
}

type AllChildren struct {
	Path []string
}

//GetChildren 获取节点下面的子节点
func (client *ZookeeperClient) GetChildren(path string) (paths []string, version int32, err error) {
	if !client.isConnect {
		err = errors.New("链接已关闭，获取子节点失败")
		return
	}
	if client.done {
		err = errors.New("链接已手动关闭，获取子节点失败")
		return
	}
	//开启协程获取子节点
	ch := make(chan interface{}, 1)
	go func(ch chan interface{}) {
		data, stat, err := client.conn.Children(path)
		ch <- GetChildrenType{data: data, version: stat.Version, err: err}
	}(ch)

	select {
	case <-time.After(time.Second * 2):
		err = fmt.Errorf("获取子节超时")
		return
	case data := <-ch:
		paths = data.(GetChildrenType).data
		version = data.(GetChildrenType).version
		err = data.(GetChildrenType).err
		if err != nil {
			err = fmt.Errorf("get node(%s) children error(err:%v)", path, err)
		}
		return
	}
}

//GetAllChildren 获取所有子节点
func (client *ZookeeperClient) GetAllChildren(path string, c *AllChildren) {

	data, _, _ := client.conn.Children(path)
	//fmt.Println(data)
	//fmt.Println(path)
	c.Path = append(c.Path, path)
	for _, v := range data {
		client.GetAllChildren(path+"/"+v, c)
	}
}
