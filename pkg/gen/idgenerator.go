package gen

import "github.com/bwmarrin/snowflake"

type IDGenerator interface {
	Generate() string
}

type SnowFlakeGenerator struct {
	node *snowflake.Node
}

func (generator SnowFlakeGenerator) Generate() string {
	return generator.node.Generate().String()
}

func NewSnowFlakeGenerator() IDGenerator {
	node, _ := snowflake.NewNode(1)
	return SnowFlakeGenerator{node: node}
}
