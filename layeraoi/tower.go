package layeraoi

import "aoi"

type Tower struct {
	objs     map[string]aoi.IObject
	watchers map[string]ILayerWatcher
}

func NewTower() *Tower {
	return &Tower{
		objs:     make(map[string]aoi.IObject),
		watchers: make(map[string]ILayerWatcher),
	}
}

func (t *Tower) Add(obj aoi.IObject,layer int) error {
	if _, ok := t.objs[obj.GetAOIID()]; ok {
		return aoi.ErrObjectExisted
	}

	t.objs[obj.GetAOIID()] = obj

	for _, w := range t.watchers {
		w.OnLayerObjectEnter(obj,layer)
	}

	return nil
}

func (t *Tower) Remove(obj aoi.IObject,layer int) error {
	if _, ok := t.objs[obj.GetAOIID()]; !ok {
		return aoi.ErrObjectNotExisted
	}

	delete(t.objs, obj.GetAOIID())

	for _, w := range t.watchers {
		w.OnLayerObjectLeave(obj,layer)
	}

	return nil
}

func (t *Tower) AddWatcher(w ILayerWatcher,layer int) error {
	if _, ok := t.watchers[w.GetAOIID()]; ok {
		return aoi.ErrObjectExisted
	}

	t.watchers[w.GetAOIID()] = w

	for _, o := range t.objs {
		w.OnLayerObjectEnter(o,layer)
	}

	return nil
}

func (t *Tower) RemoveWatcher(w ILayerWatcher,layer int) error {
	if _, ok := t.watchers[w.GetAOIID()]; !ok {
		return aoi.ErrObjectNotExisted
	}

	delete(t.watchers, w.GetAOIID())

	for _, o := range t.objs {
		w.OnLayerObjectLeave(o,layer)
	}

	return nil
}

func (t *Tower) Traversal(obj aoi.IObject, cb func(obj aoi.IWatcher) bool) {
	if !t.Existed(obj.GetAOIID()) {
		return
	}

	for _, w := range t.watchers {
		if !cb(w) {
			return
		}
	}
}

func (t *Tower) Clear(layer int) {
	for _, w := range t.watchers {
		for _, o := range t.objs {
			w.OnLayerObjectLeave(o,layer)
		}
	}

	t.objs = make(map[string]aoi.IObject)
	t.watchers = make(map[string]ILayerWatcher)
}

func (t *Tower) GetWatchers() map[string]ILayerWatcher {
	return t.watchers
}

func (t *Tower) GetWatchersLen() int {
	return len(t.watchers)
}

func (t *Tower) GetObjs() map[string]aoi.IObject {
	return t.objs
}

func (t *Tower) GetObjsLen() int {
	return len(t.objs)
}

func (t *Tower) Existed(id string) bool {
	_, ok := t.objs[id]
	return ok
}

func (t *Tower) ExistedWatcher(id string) bool {
	_, ok := t.watchers[id]
	return ok
}
