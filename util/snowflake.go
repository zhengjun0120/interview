package util

import (
	"ai_interview/conf"
	"fmt"
	"github.com/bwmarrin/snowflake"
)

var note *snowflake.Node //全局节点

func InitSnowflake() {
	noteID := conf.GetConfig().Snowflake.NodeId
	if noteID < 1 || noteID > 1023 {
		panic("雪花节点id不在1-1023之间")
	}

	n, err := snowflake.NewNode(noteID)
	if err != nil {
		panic("创建雪花节点失败" + err.Error())
	}
	note = n

	fmt.Println("雪花节点初始化完成")
}

func GenerateStringID() string {
	return note.Generate().String()
}
