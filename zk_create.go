package zk

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

type createType struct {
	rpath string
	err   error
}

//CreatePersistentNode 创建永久节点
func (client *ZookeeperClient) CreatePersistentNode(path string, data string) (err error) {
	if !client.isConnect {
		err = errors.New("链接已关闭")
		return
	}
	if client.done {
		err = errors.New("链接已手动关闭")
		return
	}
	//检查path是否存在
	b, err := client.Exists(path)
	if b == true || err != nil {
		err = fmt.Errorf("创建节点失败，节点已经存在")
		return
	}
	if path == "/" {
		return nil
	}
	//获取每级目录并检查是否存在，不存在则创建
	paths := client.getPaths(path)
	for i := 0; i < len(paths)-1; i++ {
		b, err := client.Exists(paths[i])
		if err != nil {
			return err
		}
		if b {
			continue
		}
		_, err = client.Create(paths[i], "", int32(0), zk.WorldACL(zk.PermAll))
		if err != nil {
			return err
		}
	}
	//创建最后一级目录
	_, err = client.Create(path, data, int32(0), zk.WorldACL(zk.PermAll))
	if err != nil {
		return
	}
	return nil
}

//CreateTempNode 创建零时节点
func (client *ZookeeperClient) CreateTempNode(path string, data string) (err error) {
	if !client.isConnect {
		err = errors.New("链接已关闭")
		return
	}
	if client.done {
		err = errors.New("链接已手动关闭")
		return
	}
	//检查path是否存在
	b, err := client.Exists(path)
	if b == true || err != nil {
		err = fmt.Errorf("创建节点失败，节点已经存在")
		return
	}
	if path == "/" {
		return nil
	}
	_, err = client.Create(path, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		return
	}
	return nil
}

//Create 创建节点
func (client *ZookeeperClient) Create(path string, data string, flags int32, acl []zk.ACL) (rpath string, err error) {
	if !client.isConnect {
		err = errors.New("链接已关闭")
		return
	}
	if client.done {
		err = errors.New("链接已手动关闭")
		return
	}
	//开启协程创建节点
	ch := make(chan interface{}, 1)
	go func(chan interface{}) {
		data, err := client.conn.Create(path, []byte(data), flags, acl)
		ch <- createType{rpath: data, err: err}
	}(ch)
	select {
	case <-time.After(time.Second * 2):
		err = errors.New("创建节点超时")
		return
	case data := <-ch:
		err = data.(createType).err
		if err != nil {
			return "", err
		}
		rpath = data.(createType).rpath
		return rpath, nil
	}
}

//getPaths 获取当前路径的所有子路径
func (client *ZookeeperClient) getPaths(path string) []string {
	nodes := strings.Split(path, "/")
	len := len(nodes)
	paths := make([]string, 0, len-1)
	for i := 1; i < len; i++ {
		npath := "/" + strings.Join(nodes[1:i+1], "/")
		paths = append(paths, npath)
	}
	return paths
}
