package toweraoi

import (
	"aoi"
)

type Tower struct {
	objs     map[string]aoi.IObject
	watchers map[string]aoi.IWatcher
}

func NewTower() *Tower {
	return &Tower{
		objs:     make(map[string]aoi.IObject),
		watchers: make(map[string]aoi.IWatcher),
	}
}

func (t *Tower) Add(obj aoi.IObject) error {
	if _, ok := t.objs[obj.GetAOIID()]; ok {
		return aoi.ErrObjectExisted
	}

	t.objs[obj.GetAOIID()] = obj

	for _, w := range t.watchers {
		w.OnObjectEnter(obj)
	}

	return nil
}

func (t *Tower) Remove(obj aoi.IObject) error {
	if _, ok := t.objs[obj.GetAOIID()]; !ok {
		return aoi.ErrObjectNotExisted
	}

	delete(t.objs, obj.GetAOIID())

	for _, w := range t.watchers {
		w.OnObjectLeave(obj)
	}

	return nil
}

func (t *Tower) AddWatcher(w aoi.IWatcher) error {
	if _, ok := t.watchers[w.GetAOIID()]; ok {
		return aoi.ErrObjectExisted
	}

	t.watchers[w.GetAOIID()] = w

	for _, o := range t.objs {
		w.OnObjectEnter(o)
	}

	return nil
}

func (t *Tower) RemoveWatcher(w aoi.IWatcher) error {
	if _, ok := t.watchers[w.GetAOIID()]; !ok {
		return aoi.ErrObjectNotExisted
	}

	delete(t.watchers, w.GetAOIID())

	for _, o := range t.objs {
		w.OnObjectLeave(o)
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

func (t *Tower) Clear() {
	for _, w := range t.watchers {
		for _, o := range t.objs {
			w.OnObjectLeave(o)
		}
	}

	t.objs = make(map[string]aoi.IObject)
	t.watchers = make(map[string]aoi.IWatcher)
}

func (t *Tower) GetWatchers() map[string]aoi.IWatcher {
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
