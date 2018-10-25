package zk

import (
	"errors"
	"fmt"

	"github.com/samuel/go-zookeeper/zk"
)

//eventWatch zk事件
func (client *ZookeeperClient) eventWatch() (chan string, chan error) {
	snapshots := make(chan string)
	errs := make(chan error)
	go func() {
		for {
			select {
			case <-client.CloseCh:
				errs <- errors.New("已经关闭")
			case v, ok := <-client.eventChan:
				if ok {
					switch v.State {
					case zk.StateAuthFailed:

						snapshots <- "链接失败"
						client.isConnect = false
					// 已经连接成功
					case zk.StateConnected:

						snapshots <- "链接成功"
						client.isConnect = true
					// 连接Session失效
					case zk.StateExpired:
						snapshots <- "连接Session失效"
						client.isConnect = false
					// 网络连接不成功
					case zk.StateDisconnected:
						snapshots <- fmt.Sprintf("zk已断开连接:%v", client.servers)
						client.isConnect = false
					// 网络断开，正在连接
					case zk.StateConnecting:
						snapshots <- "网络断开，正在连接"
						client.isConnect = false
					case zk.StateHasSession:
						client.isConnect = true
					}
				} else {
					snapshots <- "网络断开"
					client.isConnect = false
				}
			}
		}
	}()
	return snapshots, errs
}

//WatchServerList is
func (client *ZookeeperClient) WatchServerList(path string) (chan []string, chan error) {
	snapshots := make(chan []string)
	errors := make(chan error)

	go func() {
		for {
			snapshot, _, events, err := client.conn.ChildrenW(path)
			if err != nil {
				errors <- err
				return
			}
			snapshots <- snapshot
			evt := <-events
			if evt.Err != nil {
				errors <- evt.Err
				return
			}
		}
	}()

	return snapshots, errors
}

//WatchGetData is
func (client *ZookeeperClient) WatchGetData(path string) (chan []byte, chan error) {
	snapshots := make(chan []byte)
	errors := make(chan error)

	go func() {
		for {
			dataBuf, _, event, err := client.conn.GetW(path)
			if err != nil {
				errors <- err
				return
			}
			snapshots <- dataBuf
			evt := <-event
			if evt.Err != nil {
				errors <- evt.Err
				return
			}
		}

	}()
	return snapshots, errors
}
