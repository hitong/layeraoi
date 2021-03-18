package toweraoi

import (
	"aoi"
	"aoi/base/linemath"
	"errors"
	"math"

	log "github.com/cihub/seelog"
)

// TowerAOI 基于灯塔的AOI实现
type TowerAOI struct {
	// 灯塔相关
	towers         [][]*Tower
	minPos, maxPos linemath.Vector2 // 地图边界
	towerSize      float32
	towerSizeX     int // 灯塔数量
	towerSizeY     int // 灯塔数量
	objs           map[string]*_CacheObject

	// 全局列表
	global *Tower

	// 群组
	groups      map[int]*Tower
	groupIDSeed int

	// 缓存需要清理的watcher
	dirtyWatchers map[string]*_WrapWatcher
}

type _CacheObject struct {
	X, Y        int
	wrapWatcher *_WrapWatcher
	groups      []int
}

var (
	ErrTowerConfigInvalid = errors.New("config invalid")
)

type Config struct {
	MinPos    linemath.Vector2
	MaxPos    linemath.Vector2
	TowerSize float32
}

func New(cfg *Config) (aoi.IAOI, error) {
	if cfg.MinPos.X > cfg.MaxPos.X || cfg.MinPos.Y > cfg.MaxPos.Y {
		return nil, ErrTowerConfigInvalid
	}

	ta := &TowerAOI{
		minPos:        cfg.MinPos,
		maxPos:        cfg.MaxPos,
		towerSize:     cfg.TowerSize,
		objs:          make(map[string]*_CacheObject),
		dirtyWatchers: make(map[string]*_WrapWatcher),
	}

	ta.towerSizeX = int(math.Ceil(float64((cfg.MaxPos.X - cfg.MinPos.X) / cfg.TowerSize)))
	ta.towerSizeY = int(math.Ceil(float64((cfg.MaxPos.Y - cfg.MinPos.Y) / cfg.TowerSize)))

	ta.towers = make([][]*Tower, ta.towerSizeX)
	for x := 0; x < ta.towerSizeX; x++ {
		ta.towers[x] = make([]*Tower, ta.towerSizeY)
		for y := 0; y < ta.towerSizeY; y++ {
			ta.towers[x][y] = NewTower()
		}
	}

	ta.global = NewTower()
	ta.groups = make(map[int]*Tower)

	return ta, nil
}

func (t *TowerAOI) AddToAOI(obj aoi.IObject) error {
	if obj == nil {
		return aoi.ErrObjectInvalid
	}

	if _, ok := t.objs[obj.GetAOIID()]; ok {
		return aoi.ErrObjectExisted
	}

	pos := obj.GetCoordPos()
	if !t.isInvalid(pos) {
		return aoi.ErrPosInvalid
	}

	x, y := t.transPos(pos)
	t.towers[x][y].Add(obj)
	t.objs[obj.GetAOIID()] = &_CacheObject{
		X: x,
		Y: y,
	}

	if w, ok := obj.(aoi.IWatcher); ok {
		ww := newWrapWatcher(w, t.notifyDirty)
		t.objs[obj.GetAOIID()].wrapWatcher = ww
		if visual := w.GetVisual(); visual > 0 {
			towerVisual := int(math.Ceil(float64(visual/t.towerSize))) - 1
			t.traversalTowerByVisual(x, y, towerVisual, func(towerX, towerY int, tower *Tower) {
				tower.AddWatcher(ww)
			})
		}

		t.global.AddWatcher(ww)
	}

	t.flushWatchers()

	return nil
}

func (t *TowerAOI) RemoveFromAOI(obj aoi.IObject) error {
	if obj == nil {
		return aoi.ErrObjectInvalid
	}

	cacheObj, ok := t.objs[obj.GetAOIID()]
	if !ok {
		return aoi.ErrObjectNotExisted
	}
	delete(t.objs, obj.GetAOIID())

	// 从所有group中移除
	for _, groupID := range cacheObj.groups {
		if group, existed := t.groups[groupID]; existed {
			if group.Existed(obj.GetAOIID()) {
				group.Remove(obj)
				if cacheObj.wrapWatcher != nil {
					group.RemoveWatcher(cacheObj.wrapWatcher)
				}

				if group.GetObjsLen() == 0 && group.GetWatchersLen() == 0 {
					delete(t.groups, groupID)
				}
			} else {
				log.Debug("奇怪的事情发生了, 检查代码")
			}
		}
	}

	if t.global.Existed(obj.GetAOIID()) {
		// 如果是全局object, 从全局系统中移除
		t.global.Remove(obj)
		if cacheObj.wrapWatcher != nil {
			t.global.RemoveWatcher(cacheObj.wrapWatcher)
		}
	} else {
		// 从灯塔系统中移除
		t.towers[cacheObj.X][cacheObj.Y].Remove(obj)
		if cacheObj.wrapWatcher != nil {
			if visual := cacheObj.wrapWatcher.GetVisual(); visual > 0 {
				towerVisual := int(math.Ceil(float64(visual/t.towerSize))) - 1
				t.traversalTowerByVisual(cacheObj.X, cacheObj.Y, towerVisual, func(towerX, towerY int, tower *Tower) {
					tower.RemoveWatcher(cacheObj.wrapWatcher)
				})
			}
		}
	}

	t.flushWatchers()

	return nil
}

func (t *TowerAOI) Move(obj aoi.IObject) error {
	if obj == nil {
		return aoi.ErrObjectInvalid
	}

	cacheObj, ok := t.objs[obj.GetAOIID()]
	if !ok {
		return aoi.ErrObjectNotExisted
	}
	oldX, oldY := cacheObj.X, cacheObj.Y

	pos := obj.GetCoordPos()
	if !t.isInvalid(pos) {
		return aoi.ErrPosInvalid
	}
	newX, newY := t.transPos(pos)
	if oldX == newX && oldY == newY {
		return nil
	}
	cacheObj.X = newX
	cacheObj.Y = newY

	if t.global.Existed(obj.GetAOIID()) {
		return nil
	}

	oldTower := t.towers[oldX][oldY]
	newTower := t.towers[newX][newY]
	oldTower.Remove(obj)
	newTower.Add(obj)

	if cacheObj.wrapWatcher != nil {
		if visual := cacheObj.wrapWatcher.GetVisual(); visual > 0 {
			towerVisual := int(math.Ceil(float64(visual/t.towerSize))) - 1
			t.traversalTowerByVisual(oldX, oldY, towerVisual, func(towerX, towerY int, tower *Tower) {
				tower.RemoveWatcher(cacheObj.wrapWatcher)
			})
			t.traversalTowerByVisual(newX, newY, towerVisual, func(towerX, towerY int, tower *Tower) {
				tower.AddWatcher(cacheObj.wrapWatcher)
			})
		}
	}

	t.flushWatchers()

	return nil
}

func (t *TowerAOI) AddGlobalMarker(obj aoi.IObject) error {
	if obj == nil {
		return aoi.ErrObjectInvalid
	}

	cacheObj, ok := t.objs[obj.GetAOIID()]
	if !ok {
		return aoi.ErrObjectNotExisted
	}

	if t.global.Existed(obj.GetAOIID()) {
		return aoi.ErrObjectExisted
	}

	t.towers[cacheObj.X][cacheObj.Y].Remove(obj)

	t.global.Add(obj)

	t.flushWatchers()

	return nil
}

func (t *TowerAOI) RemoveGlobalMarker(obj aoi.IObject) error {
	if obj == nil {
		return aoi.ErrObjectInvalid
	}

	if !t.global.Existed(obj.GetAOIID()) {
		return aoi.ErrObjectNotExisted
	}

	pos := obj.GetCoordPos()
	if !t.isInvalid(pos) {
		return aoi.ErrPosInvalid
	}

	t.global.Remove(obj)

	x, y := t.transPos(pos)
	t.towers[x][y].Add(obj)
	t.objs[obj.GetAOIID()].X = x
	t.objs[obj.GetAOIID()].Y = y

	t.flushWatchers()

	return nil
}

func (t *TowerAOI) Traversal(obj aoi.IObject, cb func(obj aoi.IWatcher) bool) {
	if obj == nil {
		return
	}

	cacheObject, ok := t.objs[obj.GetAOIID()]
	if !ok {
		return
	}

	// 全局对象, 直接遍历所有的watcher就可以
	if t.global.Existed(obj.GetAOIID()) {
		t.global.Traversal(obj, cb)
		return
	}

	calledWatchers := make(map[string]bool)
	for _, groupID := range cacheObject.groups {
		t.TraversalGroup(obj, groupID, func(watcher aoi.IWatcher) bool {
			if _, called := calledWatchers[watcher.GetAOIID()]; !called {
				calledWatchers[watcher.GetAOIID()] = true
				return cb(watcher)
			}
			return true
		})
	}

	tower := t.towers[cacheObject.X][cacheObject.Y]
	for _, w := range tower.GetWatchers() {
		if _, called := calledWatchers[w.GetAOIID()]; !called {
			if !cb(w) {
				return
			}
			//calledWatchers[w.GetID()] = true   noneed
		}
	}
}

func (t *TowerAOI) CreateGroup(objs []aoi.IObject) (int, error) {
	if len(objs) == 0 {
		return 0, aoi.ErrObjectInvalid
	}

	// 检查所有objs是否在灯塔系统中
	for _, o := range objs {
		if _, ok := t.objs[o.GetAOIID()]; !ok {
			return 0, aoi.ErrObjectNotExisted
		}
	}

	// 构建群组
	group := NewTower()
	t.groupIDSeed++
	t.groups[t.groupIDSeed] = group
	for _, o := range objs {
		group.Add(o)

		cacheObj := t.objs[o.GetAOIID()]
		if cacheObj.groups == nil {
			cacheObj.groups = make([]int, 0, 1)
		}
		cacheObj.groups = append(cacheObj.groups, t.groupIDSeed)

		if cacheObj.wrapWatcher != nil {
			group.AddWatcher(cacheObj.wrapWatcher)
		}
	}

	t.flushWatchers()

	return t.groupIDSeed, nil
}

func (t *TowerAOI) DestroyGroup(groupID int) error {
	group, ok := t.groups[groupID]
	if !ok {
		return aoi.ErrGroupNotExisted
	}

	for _, o := range group.GetObjs() {
		cacheObj := t.objs[o.GetAOIID()]
		for i, v := range cacheObj.groups {
			if v == groupID {
				cacheObj.groups = append(cacheObj.groups[:i], cacheObj.groups[i+1:]...)
			}
		}
	}

	group.Clear()
	delete(t.groups, groupID)

	t.flushWatchers()

	return nil
}

func (t *TowerAOI) AddToGroup(obj aoi.IObject, groupID int) error {
	if obj == nil {
		return aoi.ErrObjectInvalid
	}

	cacheObj, ok := t.objs[obj.GetAOIID()]
	if !ok {
		return aoi.ErrObjectNotExisted
	}

	group, ok := t.groups[groupID]
	if !ok {
		return aoi.ErrGroupNotExisted
	}

	if group.Existed(obj.GetAOIID()) {
		return nil
	}

	group.Add(obj)
	if cacheObj.wrapWatcher != nil {
		group.AddWatcher(cacheObj.wrapWatcher)
	}

	cacheObj.groups = append(cacheObj.groups, groupID)

	t.flushWatchers()

	return nil
}

func (t *TowerAOI) RemoveFromGroup(obj aoi.IObject, groupID int) error {
	if obj == nil {
		return aoi.ErrObjectInvalid
	}

	cacheObj, ok := t.objs[obj.GetAOIID()]
	if !ok {
		return aoi.ErrObjectNotExisted
	}

	group, ok := t.groups[groupID]
	if !ok {
		return aoi.ErrGroupNotExisted
	}

	if !group.Existed(obj.GetAOIID()) {
		return nil
	}

	group.Remove(obj)
	if cacheObj.wrapWatcher != nil {
		group.RemoveWatcher(cacheObj.wrapWatcher)
	}

	for i := range cacheObj.groups {
		if cacheObj.groups[i] == groupID {
			cacheObj.groups = append(cacheObj.groups[:i], cacheObj.groups[i+1:]...)
		}
	}

	if group.GetWatchersLen() == 0 && group.GetObjsLen() == 0 {
		delete(t.groups, groupID)
	}

	t.flushWatchers()

	return nil
}

func (t *TowerAOI) TraversalGroup(obj aoi.IObject, groupID int, cb func(watcher aoi.IWatcher) bool) {
	if obj == nil {
		return
	}

	if _, ok := t.objs[obj.GetAOIID()]; !ok {
		return
	}

	group, ok := t.groups[groupID]
	if !ok {
		return
	}

	group.Traversal(obj, cb)
}

func (t *TowerAOI) isInvalid(pos linemath.Vector2) bool {
	if pos.X < t.minPos.X || pos.X > t.maxPos.X || pos.Y < t.minPos.Y || pos.Y > t.maxPos.Y {
		return false
	}

	return true
}

func (t *TowerAOI) transPos(pos linemath.Vector2) (int, int) {
	x := math.Floor(float64((pos.X - t.minPos.X) / t.towerSize))
	y := math.Floor(float64((pos.Y - t.minPos.Y) / t.towerSize))
	return int(x), int(y)
}

func (t *TowerAOI) traversalTowerByVisual(x, y, visual int, cb func(towerX, towerY int, tower *Tower)) {
	startX, endX, startY, endY := t.getTowerRange(x, y, visual)
	for i := startX; i <= endX; i++ {
		for j := startY; j <= endY; j++ {
			cb(i, j, t.towers[i][j])
		}
	}
}

func (t *TowerAOI) getTowerRange(x, y, i int) (int, int, int, int) {
	startX := x - i
	if startX < 0 {
		startX = 0
	}
	endX := x + i
	if endX > t.towerSizeX-1 {
		endX = t.towerSizeX - 1
	}
	startY := y - i
	if startY < 0 {
		startY = 0
	}
	endY := y + i
	if endY > t.towerSizeY-1 {
		endY = t.towerSizeY - 1
	}

	return startX, endX, startY, endY
}

func (t *TowerAOI) notifyDirty(ww *_WrapWatcher) {
	t.dirtyWatchers[ww.GetAOIID()] = ww
}

func (t *TowerAOI) flushWatchers() {
	for id, w := range t.dirtyWatchers {
		w.Flush()
		delete(t.dirtyWatchers, id)
	}
}
