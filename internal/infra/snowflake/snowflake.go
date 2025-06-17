package snowflake

import "github.com/bwmarrin/snowflake"

type IdGenerator struct {
	node *snowflake.Node
}

func NewIdGenerator(nodeId int64) (*IdGenerator, error) {
	node, err := snowflake.NewNode(nodeId)
	if err != nil {
		return nil, err
	}
	return &IdGenerator{node: node}, nil
}

func (ig *IdGenerator) GenerateId() string {
	return ig.node.Generate().Base64()
}
