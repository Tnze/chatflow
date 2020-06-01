package main

import (
	"fmt"
	"github.com/Tnze/chatflow"
	"math/rand"
	"strconv"
)

func main() {
	r := chatflow.New()
	r.Prefix("/猜数字", Game)

	/* cqp.GroupMsg = func() */
	{
		r.HandleMsg(groupSource{
			groupID:  100000,
			senderID: 100001,
		}, "/猜数字")
	}
}

type groupSource struct {
	groupID, senderID int64
}

func (g groupSource) Say(msg ...interface{}) {
	// cqp.SendGroupMsg(fmt.Sprint(msg...))
	fmt.Print(msg...)
}

func Game(c *chatflow.Context) {
	num := rand.Intn(100) + 1
	c.Say("我选好数字啦，你开始猜吧！")
	for {
		msg, _ := c.Next()
		if msg == "不玩了" {
			c.Say("那就不玩了吧")
			return
		}
		if n, err := strconv.Atoi(msg); err != nil {
			c.Sayf("%q不是一个数字呢，如果不想玩了记得跟我说\"不玩了\"哦", msg)
		} else if n < num {
			c.Say("太小了")
		} else if n > num {
			c.Say("太大了")
		} else {
			c.Sayf("猜对了！就是%d！", n)
			return
		}

	}
}
