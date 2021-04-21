package core

import (
	"fmt"
	"sync"
)

type Grid struct {
	// 格子ID
	GID int
	// 格子的左边界
	Left int
	// 格子右边界
	Right int
	// 格子上边界
	Upper int
	// 格子下边界
	Lower int
	// 玩家或物品ID集合
	PlayerIDs sync.Map
}
// 初始化格子
func NewGrid(gid,left,right,upper,lower int)*Grid{
	return &Grid{
		GID:       gid,
		Left:      left,
		Right:     right,
		Upper:     upper,
		Lower:     lower,
		PlayerIDs: sync.Map{},
	}
}
// 格子内添加一个玩家
func (g *Grid)Add(playerId int){
	g.PlayerIDs.Store(playerId,true)
}
// 格子内删除一个玩家
func (g* Grid) Remove(playerId int)  {
	g.PlayerIDs.Delete(playerId)
}
// 获取格子内所有玩家
func (g *Grid)GetAllPlayerIds()[]int{
	var playerIds []int
	g.PlayerIDs.Range(func(key, value interface{}) bool {
		playerIds=append(playerIds,key.(int))
		return true
	})
	return playerIds
}
// 获得格子基本信息
func (g *Grid)getGridInfo()string{
	return fmt.Sprintf("grid id :%d")
}