package zk

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConnetc(t *testing.T) {

	Convey("zk测试", t, func() {
		c, err := New([]string{"localhost:2181"}, time.Second)
		if err != nil {
			fmt.Println(err)
		}
		err = c.Connect()
		fmt.Println(c.useCount)
		if err != nil {
			fmt.Println(err)
		}
		Convey("测试是否链接成功", func() {
			So(c.IsConnected(), ShouldEqual, true)
		})
		Convey("链接数量测试", func() {

			So(c.useCount, ShouldEqual, 1)
		})
		Convey("关闭链接", func() {
			c.Close()
			So(c.done, ShouldEqual, true)
			So(c.isConnect, ShouldEqual, false)
		})
	})
}
