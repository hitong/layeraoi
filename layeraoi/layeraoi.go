package layeraoi

import (
	"aoi"
	"errors"
)

type LayerAOI struct {
	AllAoi map[int]ILayerAOIBase
	aoi.IAOIBase
}

func New() LayerAOI {
	return LayerAOI{AllAoi: make(map[int]ILayerAOIBase)}
}

func (l *LayerAOI) AddLayer(layers uint64, layer ...ILayerAOIBase) error {
	numZero := Ctz64(layers)
	var times = 0
	for numZero < 64 {
		if _, ok := l.AllAoi[numZero]; ok {
			return errors.New("Repeat layer key ")
		}
		layer[times].SetLayer(numZero)
		l.AllAoi[numZero] = layer[times]
		layers = SetNZero(layers, numZero)
		numZero = Ctz64(layers)
		times++
	}

	return nil
}

//返回分割至的位置，-1代表分割失败
func (l *LayerAOI) SplitLayer(layer int)int{
	return -1
}

//src + tar -> tar
func (l *LayerAOI) MergeLayer(src int,tar int) error{
	return nil
}

func (l *LayerAOI) GetLayersBits() (num uint64) {
	for _, layer := range l.AllAoi {
		num += 1 << layer.GetLayer()
	}
	return
}

func (l *LayerAOI) AddToAOI(obj aoi.IObject) error {
	typeNum := getLayerBits(obj)
	numZero := Ctz64(typeNum)
	for numZero != 64 {
		l.AllAoi[numZero].AddToAOI(obj)
		typeNum = SetNZero(typeNum, numZero)
		numZero = Ctz64(typeNum)
	}
	return nil
}

func (l *LayerAOI) RemoveFromAOI(obj aoi.IObject) error {
	typeNum := getLayerBits(obj)
	numZero := Ctz64(typeNum)
	for numZero < 64 {
		l.AllAoi[numZero].RemoveFromAOI(obj)
		typeNum = SetNZero(typeNum, numZero)
		numZero = Ctz64(typeNum)
	}
	return nil
}

func (l *LayerAOI) Move(obj aoi.IObject) error {
	typeNum := getLayerBits(obj)
	numZero := Ctz64(typeNum)
	for numZero < 64 {
		l.AllAoi[numZero].Move(obj)
		typeNum = SetNZero(typeNum, numZero)
		numZero = Ctz64(typeNum)
	}
	return nil
}

func (l *LayerAOI) Traversal(obj aoi.IObject, cb func(watcher aoi.IWatcher) bool) {
	typeNum := getLayerBits(obj)
	numZero := Ctz64(typeNum)
	for numZero < 64 {
		l.AllAoi[numZero].Traversal(obj, cb)
		typeNum = SetNZero(typeNum, numZero)
		numZero = Ctz64(typeNum)
	}
}
