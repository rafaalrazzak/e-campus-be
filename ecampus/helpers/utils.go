package helpers

import "github.com/bwmarrin/snowflake"

func GenerateId() int64 {
	node, _ := snowflake.NewNode(1)
	return node.Generate().Int64()
}
