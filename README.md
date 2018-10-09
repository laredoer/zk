# zk
golang  的 zookeeper 工具
>安装zookeeper

####
利用docker可以方便的运行zookeeper,本项目已提供docker-compose.yml文件，只需要运行一下命令
```
//启动
    docker-compose up -d

//查看进程
    docker-compose ps

//停止
    docker-compose stop
```

>使用

```golang
main.go >>>

--------------------------------------------------

package main

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	zk "github.com/wule61/lib/zookeeper"
)

func main() {
	c, err := zk.New([]string{"127.0.0.1:2181", "127.0.0.1:2182", "127.0.0.1:2183"}, time.Second)
	if err != nil {
		logs.Error(err)
	}
	c.Connect()
    defer c.Close()
    //创建节点
    err = c.CreatePersistentNode("/node","123")
    if err != nil {
        logs.Error(err)
    }
    
    //监听节点
	snapshots, errs := c.WatchServerList("/node")
	go func() {
		for {
			select {
			case serverList := <-snapshots:
				fmt.Println(serverList)
			case erros := <-errs:
				fmt.Println(erros)
			}
		}
    }()
    
    //监听节点的值
	configs, errors := c.WatchGetData("/node")

	go func() {
		for {
			select {
			case configData := <-configs:
				fmt.Println(string(configData))
			case err = <-errors:
				fmt.Println(err)
			}
		}
	}()
	select {}
}


```
