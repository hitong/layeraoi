package main

import (
	"aoi"
	"aoi/base/linemath"
	"aoi/layeraoi"
	"fmt"
)

type Watcher struct {
	//aoi.IObject
	//aoi.IWatcher
	Mark   *Object
	Visual map[int]float32
	layeraoi.ILayerWatcher
}

func (tw *Watcher) GetCoordPos() linemath.Vector2 {
	return tw.Mark.GetCoordPos()
}
func (tw *Watcher) GetLayerVisual(layer int) float32 {
	return tw.Visual[layer]
}
func (tw *Watcher) OnLayerObjectEnter(obj aoi.IObject, layer int) {
	fmt.Println("OnLayerObjectEnter ", obj.GetAOIID(), " layer ", layer)
}
func (tw *Watcher) OnLayerObjectLeave(obj aoi.IObject, layer int) {
	fmt.Println("OnLayerObjectLeave ", obj.GetAOIID(), " layer ", layer)
}
func (tw *Watcher) OnLayerBatchEnter(objs []aoi.IObject, layer int) {
	//	fmt.Println("OnBatchEnter ",tw.GetAOIID(), " layer ", layer)
	for _, obj := range objs {
		tw.Mark.OnLayerObjectEnter(obj, layer)
	}

}
func (tw *Watcher) OnLayerBatchLeave(objs []aoi.IObject, layer int) {
	//	fmt.Println("OnBatchLeave ",tw.GetAOIID(), " layer ", layer)
	for _, obj := range objs {
		tw.Mark.OnLayerObjectLeave(obj, layer)
	}
}
func (tw *Watcher) GetAOIID() string {
	return tw.Mark.GetAOIID()
}
func (tw *Watcher) GetLayerBits() uint64 {
	return tw.Mark.GetLayerBits()
}
