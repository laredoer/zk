package zk

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/samuel/go-zookeeper/zk"
)

//ZookeeperClient zk对象
type ZookeeperClient struct {
	servers   []string        //zk地址
	timeout   time.Duration   //超时时间
	conn      *zk.Conn        //zk连接实例
	eventChan <-chan zk.Event //zk事件
	useCount  int32           //当前使用人数
	isConnect bool            //是否连接
	once      sync.Once       //一个锁
	CloseCh   chan struct{}   //关闭通道
	done      bool            //是否手动关闭
}

//New 创建Zookeeper实例
func New(servers []string, timeout time.Duration) (*ZookeeperClient, error) {
	client := &ZookeeperClient{servers: servers, timeout: timeout, useCount: 0}
	client.CloseCh = make(chan struct{})
	return client, nil
}

//Connect 链接zk服务器
func (client *ZookeeperClient) Connect() (err error) {
	if client.conn == nil {
		conn, eventChan, err := zk.Connect(client.servers, client.timeout)
		if err != nil {
			return err
		}
		client.conn = conn
		client.eventChan = eventChan
		go client.eventWatch()
	}
	atomic.AddInt32(&client.useCount, 1) //原子操作加一
	time.Sleep(time.Second)
	client.isConnect = true
	return
}

//IsConnected 是否已连接到服务器
func (client *ZookeeperClient) IsConnected() bool {
	return client.isConnect
}

//Close 关闭链接
func (client *ZookeeperClient) Close() (err error) {
	atomic.AddInt32(&client.useCount, -1)
	if client.useCount > 0 {
		//别人在使用，不关闭真正的链接
		return nil
	}
	//关闭真正的链接,只执行一次
	if client.conn != nil {
		client.once.Do(client.conn.Close)
	}
	//修改结构体数据
	client.isConnect = false
	client.done = true
	client.once.Do(func() {
		close(client.CloseCh)
	})
	logs.Alert("链接已关闭")
	return nil
}

//Reconnect 重新链接
func (client *ZookeeperClient) Reconnect() (err error) {
	client.isConnect = false
	if client.conn != nil {
		client.Close()
	}
	client.done = false
	return client.Connect()
}
