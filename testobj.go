package aoi

import (
	"aoi/base/linemath"
	"fmt"
)

type TestMarker struct {
	ID  string
	Pos linemath.Vector2
}

func (m *TestMarker) GetAOIID() string {
	return m.ID
}

func (m *TestMarker) GetCoordPos() linemath.Vector2 {
	return m.Pos
}

type TestWatcher struct {
	ID     string
	Pos    linemath.Vector2
	Visual float32
}
//
//type TestWatcherWrap struct {
//	layeraoi.IWatcher
//	TestWatcher
//	Type uint64
//}
//func (tw *TestWatcherWrap)	GetCoordPos() linemath.Vector2{
//	return tw.TestWatcher.GetCoordPos()
//}
//func (tw *TestWatcherWrap)	GetVisual(layer int) float32{
//	if layer == 1{
//		return 100
//	}
//	return 200
//}
//func (tw *TestWatcherWrap)OnObjectEnter(obj layeraoi.IObject,layer int){
//	fmt.Println("OnObjectEnter ",obj.GetAOIID(), " layer ", layer)
//}
//func (tw *TestWatcherWrap)OnObjectLeave(obj layeraoi.IObject,layer int){
//	fmt.Println("OnObjectLeave ",obj.GetAOIID(), " layer ", layer)
//}
//func (tw *TestWatcherWrap)OnBatchEnter(objs []layeraoi.IObject,layer int){
//	fmt.Println("OnBatchEnter ",tw.GetAOIID(), " layer ", layer)
//
//}
//func (tw *TestWatcherWrap)OnBatchLeave(objs []layeraoi.IObject,layer int){
//	fmt.Println("OnBatchLeave ",tw.GetAOIID(), " layer ", layer)
//}
//func (tw *TestWatcherWrap)GetAOIID() string{
//	return tw.TestWatcher.GetAOIID()
//}
//func (tw *TestWatcherWrap)GetType()uint64{
//	return tw.Type
//}

func (obj *TestWatcher) GetAOIID() string {
	return obj.ID
}

func (obj *TestWatcher) GetCoordPos() linemath.Vector2 {
	return obj.Pos
}

func (obj *TestWatcher) GetVisual() float32 {
	return obj.Visual
}

func (obj *TestWatcher) OnObjectEnter(o IObject) {
	fmt.Println(o.GetAOIID(), "Enter", obj.GetAOIID())
}

func (obj *TestWatcher) OnObjectLeave(o IObject) {
	fmt.Println(o.GetAOIID(), "Leave", obj.GetAOIID())
}

func (obj *TestWatcher) OnBatchEnter(objs []IObject) {
	if len(objs) == 0 {
		return
	}

	ids := ""
	for _, o := range objs {
		ids += o.GetAOIID()
		ids += " "
	}

	fmt.Println(ids, "BatchEnter", obj.GetAOIID())
}

func (obj *TestWatcher) OnBatchLeave(objs []IObject) {
	if len(objs) == 0 {
		return
	}

	ids := ""
	for _, o := range objs {
		ids += o.GetAOIID()
		ids += " "
	}

	fmt.Println(ids, "BatchLeave", obj.GetAOIID())
}
