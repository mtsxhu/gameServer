package core

import "fmt"
// AOI边界值
const (
	AOI_LEFT = 85
	AOI_RIGHT=410
	AOI_COUNTS_X=10
	AOI_UPPER=75
	AOI_LOWER=400
	AOI_COUNTS_Y=20
)

type AOIManager struct {
	// 地图左边界坐标
	Left int
	// 地图右边界
	Right int
	// x方向格子数量
	CountsX int
	// Y方向格子数量
	CountsY int
	// 地图上边界
	Upper int
	// 地图下边界
	Lower int
	// 当前区域格子对象
	Grids map[int]*Grid
}

func NewAOIManager(left,right,cntsX,cntsY,upper,lower int)*AOIManager{
	var aoiM *AOIManager=&AOIManager{
		Left:    left,
		Right:   right,
		CountsX: cntsX,
		CountsY: cntsY,
		Upper:   upper,
		Lower:   lower,
		Grids:   make(map[int]*Grid),
	}
	// 初始化格子
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			gid:=y*cntsX+x
			aoiM.Grids[gid]=NewGrid(gid,aoiM.Left+x*aoiM.gridWidth(),
				aoiM.Left+(x+1)*aoiM.gridWidth(),
				aoiM.Upper+y*aoiM.gridLength(),
				aoiM.Upper+(y+1)*aoiM.gridLength())
		}
	}
	return aoiM
}

// 得到每个格子在X轴方向的宽度
func (a *AOIManager)gridWidth()int{
	return (a.Right-a.Left)/a.CountsX
}

// 得到每个格子在Y轴方向的长度
func (a *AOIManager)gridLength()int{
	return (a.Lower-a.Upper)/a.CountsY
}

// 打印地图信息
func (a *AOIManager) getAoiManagerInfo() string{
	return ""
}

// GetSurroundGridsByGid 根据格子的gID得到当前周边的九宫格信息
func (m *AOIManager) GetSurroundGridsByGid(gID int) (grids []*Grid) {
	//判断gID是否存在
	if _, ok := m.Grids[gID]; !ok  {
		return
	}

	//将当前gid添加到九宫格中
	grids  = append(grids, m.Grids[gID])

	//根据gid得到当前格子所在的X轴编号
	idx := gID % m.CountsX
	//当前idx左边还有格子，就加入到九宫格中
	if idx > 0 {
		grids = append(grids, m.Grids[gID-1])
	}
	//当前的idx右边还有格子，就加入到九宫格中
	if idx < m.CountsX - 1 {
		grids = append(grids, m.Grids[gID+1])
	}

	//将x轴当前的格子都取出，进行遍历，再分别得到每个格子的上下是否有格子
	//得到当前x轴的格子id集合
	gidsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gidsX = append(gidsX, v.GID)
	}

	//遍历x轴格子
	for _, v := range gidsX {
		//计算该格子处于第几行
		idy := v / m.CountsX

		//判断当前的idy上边是否还有格子
		if idy > 0 {
			grids = append(grids, m.Grids[v-m.CountsX])
		}
		//判断当前的idy下边是否还有格子
		if idy < m.CountsY - 1 {
			grids = append(grids, m.Grids[v+m.CountsX])
		}
	}

	return
}

// GetGIDByPos 通过横纵坐标获取对应的格子ID
func (m *AOIManager) GetGIDByPos(x, y float32) int {
	gx := (int(x) - m.Left) / m.gridWidth()
	gy := (int(x) - m.Upper) / m.gridLength()

	return gy * m.CountsX + gx
}

// GetPIDsByPos 通过横纵坐标得到周边九宫格内的全部PlayerIDs
func (m *AOIManager) GetPIDsByPos(x, y float32) (playerIDs []int) {
	//根据横纵坐标得到当前坐标属于哪个格子ID
	gID := m.GetGIDByPos(x, y)

	//根据格子ID得到周边九宫格的信息
	grids := m.GetSurroundGridsByGid(gID)
	for _, v := range grids {
		playerIDs = append(playerIDs, v.GetAllPlayerIds()...)
		fmt.Printf("===> grid ID : %d, pids : %v  ====", v.GID, v.GetAllPlayerIds())
	}

	return
}

//通过GID获取当前格子的全部playerID
func (m *AOIManager) GetPidsByGid(gID int) (playerIDs []int) {
	playerIDs = m.Grids[gID].GetAllPlayerIds()
	return
}

//移除一个格子中的PlayerID
func (m *AOIManager) RemovePidFromGrid(pID, gID int) {
	m.Grids[gID].Remove(pID)
}

//添加一个PlayerID到一个格子中
func (m *AOIManager) AddPidToGrid(pID, gID int) {
	m.Grids[gID].Add(pID)
}

//通过横纵坐标添加一个Player到一个格子中
func (m *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := m.GetGIDByPos(x, y)
	grid := m.Grids[gID]
	grid.Add(pID)
}

//通过横纵坐标把一个Player从对应的格子中删除
func (m *AOIManager) RemoveFromGridByPos(pID int, x, y float32) {
	gID := m.GetGIDByPos(x, y)
	grid := m.Grids[gID]
	grid.Remove(pID)
}






