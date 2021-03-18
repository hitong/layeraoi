package layeraoi

import (
	"aoi"
	"aoi/base/linemath"
	"errors"
	"math"
	"time"

	log "github.com/cihub/seelog"
)

type TowerAOILoadBalanceLayer interface {
	LayerBalance(v map[int]int)(int,int)
}

type TowerAOILoadBalanceTower interface {
	TowerBalance(v map[int]int)(int,int)
}

type DefaultLoadBalance struct {
	TowerAOILoadBalanceLayer
	TowerAOILoadBalanceTower
}

func (d *DefaultLoadBalance)LayerBalance(v map[int]int)(int,int){
	var minIdx = 0
	var minNum = math.MaxInt64
	for idx, v := range v {
		if minNum > v {
			minIdx = idx
			minNum = v
		}
	}

	return minIdx,minNum
}

type LoadBalanceConfig struct {
	MethodObj interface{}
}


type mapTowerLayer map[int]map[int]*Tower
type arrayTowerLayer [][]*Tower

type towerLayer interface {
	getLayerType() LayerType
	getTower(x, y int) *Tower
	traversal(layer towerLayer, f func(x, y int, layer towerLayer, tower *Tower))
}

func (mt *mapTowerLayer) getTower(x, y int) *Tower {
	if _, ok := (*mt)[x]; !ok {
		(*mt)[x] = make(map[int]*Tower)

	}
	if _,ok := (*mt)[x][y];!ok{
		(*mt)[x][y] = NewTower()
	}

	return (*mt)[x][y]
}

func (mt *mapTowerLayer) traversal(layer towerLayer, f func(x, y int, layer towerLayer, tower *Tower)) {
	for x, m := range *mt {
		for y, v := range m {
			f(x, y, layer, v)
		}
	}
}

func (t *TowerAOILayer)cpTower(src *Tower, target *Tower) {
	if target == nil {
		target = src
	} else {
		for _, watcher := range src.watchers {
			src.RemoveWatcher(watcher,t.GetLayer())
			target.AddWatcher(watcher,t.GetLayer())
		}

		for _, obj := range src.objs {
			src.Remove(obj,t.GetLayer())
			target.Add(obj,t.GetLayer())
		}
	}
}

func (mt *mapTowerLayer) getLayerType() LayerType {
	return MapLayer
}

func (at *arrayTowerLayer) getTower(x, y int) *Tower {
	return (*at)[x][y]
}

func (at *arrayTowerLayer) getLayerType() LayerType {
	return ArrayLayer
}

func (at *arrayTowerLayer) traversal(layer towerLayer, f func(x, y int, layer towerLayer, tower *Tower)) {
	for x, m := range *at {
		for y, v := range m {
			f(x, y, layer, v)
		}
	}
}

// TowerAOI 基于灯塔的AOI实现
type TowerAOILayer struct {
	// 灯塔相关
	//towers         arrayTowerLayer
	minPos, maxPos linemath.Vector2 // 地图边界
	towerSize      float32
	towerSizeX     int // 灯塔数量
	towerSizeY     int // 灯塔数量
	objs           map[string]*_CacheObject
	towerLayers    map[int]towerLayer
	layersNums     map[int]int //方便快速计算层负载

	// 全局列表
	global *Tower

	// 群组
	groups      map[int]*Tower
	groupIDSeed int

	// 缓存需要清理的watcher
	dirtyWatchers map[string]*_WrapWatcher

	towerLimit int
	layer      int //layerAOI 使用
	layerLimit int
	layerID int //增加layer时+1，layer合并时-1
	lastAdjustment time.Time //上一次层调整时间，每次调整单层，从上至下
	loadBalancing interface{}
}

func (t *TowerAOILayer)nextLayerID()int{
	t.layerID += 1
	return t.layerID
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
	LoadBalanceCfg *LoadBalanceConfig
	LayerLimit int
}

func NewTowerAoi(cfg *Config) (ILayerAOIBase, error) {
	if cfg.MinPos.X > cfg.MaxPos.X || cfg.MinPos.Y > cfg.MaxPos.Y {
		return nil, ErrTowerConfigInvalid
	}

	ta := &TowerAOILayer{
		minPos:        cfg.MinPos,
		maxPos:        cfg.MaxPos,
		towerSize:     cfg.TowerSize,
		objs:          make(map[string]*_CacheObject),
		dirtyWatchers: make(map[string]*_WrapWatcher),
	}

	ta.towerSizeX = int(math.Ceil(float64((cfg.MaxPos.X - cfg.MinPos.X) / cfg.TowerSize)))
	ta.towerSizeY = int(math.Ceil(float64((cfg.MaxPos.Y - cfg.MinPos.Y) / cfg.TowerSize)))

	ta.towerLayers = make(map[int]towerLayer)
	ta.layersNums = make(map[int]int)
	ta.AddLayer(ta.nextLayerID(),ArrayLayer)

	ta.global = NewTower()
	ta.groups = make(map[int]*Tower)
	ta.layerLimit = cfg.LayerLimit
	if cfg.LoadBalanceCfg != nil{
		ta.loadBalancing = cfg.LoadBalanceCfg.MethodObj
	}
	return ta, nil
}

func newArrLayer(sizeX int, sizeY int) *arrayTowerLayer{
	var at arrayTowerLayer = make([][]*Tower, sizeX)
	for i := 0; i < sizeX; i++ {
		at[i] = make([]*Tower, sizeY)
		for j := 0; j < sizeY; j++ {
			at[i][j] = NewTower()
		}
	}
	return &at
}

func newMapLayer() *mapTowerLayer{
	var mt mapTowerLayer = make(map[int]map[int]*Tower)
	return &mt
}

func (t *TowerAOILayer) AddLayer(index int, layerType LayerType) bool {
	//if t.loadBalancing == nil && len(t.towerLayers) > 0{
	//	return false
	//}
	if _, ok := t.layersNums[index]; ok {
		return false
	}

	if layerType == ArrayLayer {
		t.towerLayers[index] = newArrLayer(t.towerSizeX, t.towerSizeY)
	} else if layerType == MapLayer {
		t.towerLayers[index] = newMapLayer()
	}

	t.layersNums[index] = 0
	return true
}

func (t *TowerAOILayer) DelTowerLayer(index int) {
	if index == 0 {
		panic("Can not del zero layer")
	}
	delete(t.towerLayers, index)
	delete(t.layersNums, index)
	t.layerID -= 1
}

func (t *TowerAOILayer) MergeLayer(srcIdx, targetIdx int) {
	t.towerLayers[srcIdx].traversal(t.towerLayers[targetIdx], func(x, y int, layer towerLayer, tower *Tower) {
		t.cpTower(tower, layer.getTower(x, y))
	})

	t.layersNums[targetIdx] += t.layersNums[srcIdx]
	t.DelTowerLayer(srcIdx)
}

func (t *TowerAOILayer) GetLayer() int {
	return t.layer
}

func (t *TowerAOILayer) SetLayer(layer int) {
	t.layer = layer
}

func (t *TowerAOILayer) layerLoadBalancing() (int,int) {
	if s,ok := t.loadBalancing.(TowerAOILoadBalanceLayer);ok{
		return s.LayerBalance(t.layersNums)
	}

	if _,ok := t.loadBalancing.(TowerAOILoadBalanceTower);ok{
		//todo:tower均衡负载
		return -1,-1
	}
	var de DefaultLoadBalance
	return de.LayerBalance(t.layersNums)
}

func (t *TowerAOILayer) AddToAOI(obj aoi.IObject) error {
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
	layerIdx,layerLoad := t.layerLoadBalancing()

	//todo：优化策略
	if layerLoad > t.layerLimit && t.loadBalancing != nil{
		newLayerID := t.nextLayerID()
		if t.AddLayer(newLayerID,MapLayer){
			layerIdx = newLayerID
		} else {
			t.layerID -= 1
		}
	}

	towerLayer := t.towerLayers[layerIdx]
	towerLayer.getTower(x, y).Add(obj, t.GetLayer())
	t.layersNums[layerIdx] += 1
	t.objs[obj.GetAOIID()] = &_CacheObject{
		X: x,
		Y: y,
	}

	if w, ok := obj.(ILayerWatcher); ok {
		ww := newWrapWatcher(w, t.notifyDirty)
		t.objs[obj.GetAOIID()].wrapWatcher = ww
		if visual := w.GetLayerVisual(t.GetLayer()); visual > 0 {
			towerVisual := int(math.Ceil(float64(visual/t.towerSize))) - 1
			t.traversalTowerByVisual(x, y, towerVisual, towerLayer, func(towerX, towerY int, tower *Tower) {
				tower.AddWatcher(ww, t.GetLayer())
			})
		}

		t.global.AddWatcher(ww, t.GetLayer())
	}

	t.flushWatchers()

	return nil

}

func (t *TowerAOILayer) findObjLayer(obj aoi.IObject) int {
	x, y := t.transPos(obj.GetCoordPos())
	return t.findObjLayerByXY(x, y, obj.GetAOIID())
}

func (t *TowerAOILayer) findObjLayerByXY(x, y int, aoiID string) int {
	for idx, layer := range t.towerLayers {
		if layer.getTower(x, y).objs[aoiID] != nil || layer.getTower(x, y).watchers[aoiID] != nil{
			return idx
		}
	}

	return math.MinInt32
}

func (t *TowerAOILayer) RemoveFromAOI(obj aoi.IObject) error {
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
				group.Remove(obj, t.GetLayer())
				if cacheObj.wrapWatcher != nil {
					group.RemoveWatcher(cacheObj.wrapWatcher, t.GetLayer())
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
		t.global.Remove(obj, t.GetLayer())
		if cacheObj.wrapWatcher != nil {
			t.global.RemoveWatcher(cacheObj.wrapWatcher, t.GetLayer())
		}
	} else {
		// 从灯塔系统中移除
		layerIdx := t.findObjLayer(obj)
		if layerIdx == math.MinInt32 {
			return errors.New("Not Found Obj " + obj.GetAOIID())
		}
		x, y := t.transPos(obj.GetCoordPos())
		t.layersNums[layerIdx] -= 1
		t.towerLayers[layerIdx].getTower(x, y).Remove(obj, t.GetLayer())
		if cacheObj.wrapWatcher != nil {
			if visual := cacheObj.wrapWatcher.GetLayerVisual(t.GetLayer()); visual > 0 {
				towerVisual := int(math.Ceil(float64(visual/t.towerSize))) - 1
				t.traversalTowerByVisual(cacheObj.X, cacheObj.Y, towerVisual, t.towerLayers[layerIdx], func(towerX, towerY int, tower *Tower) {
					tower.RemoveWatcher(cacheObj.wrapWatcher, t.GetLayer())
				})
			}
		}
	}

	t.flushWatchers()

	return nil
}

func (t *TowerAOILayer)checkTowerLoad(x,y,layer int,tower *Tower){
	if len(tower.watchers) + len(tower.objs) > t.towerLimit{

	}
}

func (t *TowerAOILayer) Move(obj aoi.IObject) error {
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

	if t.loadBalancing != nil && time.Now().After(t.lastAdjustment.Add(time.Second * 3)){
		for t.layerID > 1 && (t.layersNums[t.layerID] + t.layersNums[t.layerID - 1]) <= (t.layerLimit >> 1){
			t.MergeLayer(t.layerID,t.layerID - 1)
			println("adjustment", t.GetLayer(),"  ",t.layerID)
		}
		t.lastAdjustment = time.Now()
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
	layerIdx := t.findObjLayerByXY(oldX, oldY, obj.GetAOIID())
	minLayerIdx,minLoad := t.layerLoadBalancing()
	if t.layersNums[layerIdx] - minLoad > 1 { //int(float32(t.layerLimit) * 0.1)
		t.layersNums[minLayerIdx]++
		t.layersNums[layerIdx]--
	} else {
		minLayerIdx = layerIdx
	}

	oldTower := t.towerLayers[layerIdx].getTower(oldX, oldY)
	newTower := t.towerLayers[minLayerIdx].getTower(newX, newY)
	oldTower.Remove(obj, t.GetLayer())
	newTower.Add(obj, t.GetLayer())

	if cacheObj.wrapWatcher != nil {
		if visual := cacheObj.wrapWatcher.GetLayerVisual(t.GetLayer()); visual > 0 {
			towerVisual := int(math.Ceil(float64(visual/t.towerSize))) - 1
			t.traversalTowerByVisual(oldX, oldY, towerVisual, t.towerLayers[layerIdx], func(towerX, towerY int, tower *Tower) {
				tower.RemoveWatcher(cacheObj.wrapWatcher, t.GetLayer())
			})
			t.traversalTowerByVisual(newX, newY, towerVisual, t.towerLayers[minLayerIdx], func(towerX, towerY int, tower *Tower) {
				tower.AddWatcher(cacheObj.wrapWatcher, t.GetLayer())
			})
		}
	}

	t.flushWatchers()

	return nil
}

func (t *TowerAOILayer) AddGlobalMarker(obj aoi.IObject) error {
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
	layerIdx := t.findObjLayer(obj)
	if layerIdx == math.MinInt32 {
		return errors.New("Not Found obj " + obj.GetAOIID())
	}
	t.towerLayers[layerIdx].getTower(cacheObj.X, cacheObj.Y).Remove(obj, t.GetLayer())

	t.global.Add(obj, t.GetLayer())

	t.flushWatchers()

	return nil
}

func (t *TowerAOILayer) RemoveGlobalMarker(obj aoi.IObject) error {
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

	t.global.Remove(obj, t.GetLayer())

	x, y := t.transPos(pos)
	layerIdx := t.findObjLayer(obj)
	layer := t.towerLayers[layerIdx]
	layer.getTower(x, y).Add(obj, t.GetLayer())
	t.objs[obj.GetAOIID()].X = x
	t.objs[obj.GetAOIID()].Y = y

	t.flushWatchers()

	return nil
}

func (t *TowerAOILayer) Traversal(obj aoi.IObject, cb func(obj aoi.IWatcher) bool) {
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

	layerIdx := t.findObjLayer(obj)
	tower := t.towerLayers[layerIdx].getTower(cacheObject.X, cacheObject.Y)
	for _, w := range tower.GetWatchers() {
		if _, called := calledWatchers[w.GetAOIID()]; !called {
			if !cb(w) {
				return
			}
			//calledWatchers[w.GetID()] = true   noneed
		}
	}
}

func (t *TowerAOILayer) CreateGroup(objs []aoi.IObject) (int, error) {
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
		group.Add(o, t.GetLayer())

		cacheObj := t.objs[o.GetAOIID()]
		if cacheObj.groups == nil {
			cacheObj.groups = make([]int, 0, 1)
		}
		cacheObj.groups = append(cacheObj.groups, t.groupIDSeed)

		if cacheObj.wrapWatcher != nil {
			group.AddWatcher(cacheObj.wrapWatcher, t.GetLayer())
		}
	}

	t.flushWatchers()

	return t.groupIDSeed, nil
}

func (t *TowerAOILayer) DestroyGroup(groupID int) error {
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

	group.Clear(t.GetLayer())
	delete(t.groups, groupID)

	t.flushWatchers()

	return nil
}

func (t *TowerAOILayer) AddToGroup(obj aoi.IObject, groupID int) error {
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

	group.Add(obj, t.GetLayer())
	if cacheObj.wrapWatcher != nil {
		group.AddWatcher(cacheObj.wrapWatcher, t.GetLayer())
	}

	cacheObj.groups = append(cacheObj.groups, groupID)

	t.flushWatchers()

	return nil
}

func (t *TowerAOILayer) RemoveFromGroup(obj aoi.IObject, groupID int) error {
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

	group.Remove(obj, t.GetLayer())
	if cacheObj.wrapWatcher != nil {
		group.RemoveWatcher(cacheObj.wrapWatcher, t.GetLayer())
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

func (t *TowerAOILayer) TraversalGroup(obj aoi.IObject, groupID int, cb func(watcher aoi.IWatcher) bool) {
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

func (t *TowerAOILayer) isInvalid(pos linemath.Vector2) bool {
	if pos.X < t.minPos.X || pos.X > t.maxPos.X || pos.Y < t.minPos.Y || pos.Y > t.maxPos.Y {
		return false
	}

	return true
}

func (t *TowerAOILayer) transPos(pos linemath.Vector2) (int, int) {
	x := math.Floor(float64((pos.X - t.minPos.X) / t.towerSize))
	y := math.Floor(float64((pos.Y - t.minPos.Y) / t.towerSize))
	return int(x), int(y)
}

func (t *TowerAOILayer) traversalTowerByVisual(x, y, visual int, towerLayer towerLayer, cb func(towerX, towerY int, tower *Tower)) {
	startX, endX, startY, endY := t.getTowerRange(x, y, visual)
	for i := startX; i <= endX; i++ {
		for j := startY; j <= endY; j++ {
			cb(i, j, towerLayer.getTower(i, j))
		}
	}
}

func (t *TowerAOILayer) getTowerRange(x, y, i int) (int, int, int, int) {
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

func (t *TowerAOILayer) notifyDirty(ww *_WrapWatcher) {
	t.dirtyWatchers[ww.GetAOIID()] = ww
}

func (t *TowerAOILayer) flushWatchers() {
	for id, w := range t.dirtyWatchers {
		w.Flush(t.GetLayer())
		delete(t.dirtyWatchers, id)
	}
}
