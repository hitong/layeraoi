package layeraoi

import "aoi"

type ILayerAOIBase interface {
	GetLayer()int
	SetLayer(int)
	aoi.IAOIBase
}

type ILayerObject interface {
	GetLayerBits() uint64
	aoi.IObject
}

type ILayerWatcher interface {
	ILayerObject
	aoi.IWatcher
	GetLayerVisual(layer int) float32
	OnLayerObjectEnter(obj aoi.IObject,layer int)
	OnLayerObjectLeave(obj aoi.IObject,layer int)
	OnLayerBatchEnter(objs []aoi.IObject,layer int)
	OnLayerBatchLeave(objs []aoi.IObject,layer int)
}

type LayerType uint32

const (
	StaticLayer LayerType = 0
	DynamicLayer LayerType = 1

	ArrayLayer LayerType = 2
	MapLayer LayerType = 3
)

