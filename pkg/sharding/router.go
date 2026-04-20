package sharding

import (
	"hash/fnv"

	"github.com/google/uuid"
)

const SHARDS_UNITS int = 2

type ShardRouter struct {
	ShardUnits int
}

type IShardRouter interface {
	RouteForClientID(clientID uuid.UUID) int
}

func NewShardingRouter(shardUnit int) *ShardRouter {
	return &ShardRouter{
		ShardUnits: shardUnit,
	}
}

func (s *ShardRouter) RouteForClientID(clientID uuid.UUID) int {
	h := fnv.New32a()
	h.Write([]byte(clientID.String()))
	return int(h.Sum32()) % s.ShardUnits
}
