package defaultaoi

import "aoi"

// _DefaultAOI 最基础的AOI, 全局AOI
type _DefaultAOI struct {
	objs     map[string]aoi.IObject
	watchers map[string]aoi.IWatcher
}

func New() aoi.IAOI {
	return &_DefaultAOI{
		objs:     make(map[string]aoi.IObject),
		watchers: make(map[string]aoi.IWatcher),
	}
}

func (da *_DefaultAOI) AddToAOI(obj aoi.IObject) error {
	if obj == nil {
		return aoi.ErrObjectInvalid
	}

	if _, ok := da.objs[obj.GetAOIID()]; ok {
		return aoi.ErrObjectExisted
	}

	for _, w := range da.watchers {
		w.OnObjectEnter(obj)
	}

	if w, ok := obj.(aoi.IWatcher); ok {
		da.watchers[w.GetAOIID()] = w

		var objList []aoi.IObject
		for _, o := range da.objs {
			objList = append(objList, o)
		}

		if len(objList) > 0 {
			w.OnBatchEnter(objList)
		}
	}

	da.objs[obj.GetAOIID()] = obj

	return nil
}

func (da *_DefaultAOI) RemoveFromAOI(obj aoi.IObject) error {
	if obj == nil {
		return aoi.ErrObjectInvalid
	}

	if _, ok := da.objs[obj.GetAOIID()]; !ok {
		return aoi.ErrObjectNotExisted
	}

	delete(da.objs, obj.GetAOIID())

	for _, w := range da.watchers {
		if w.GetAOIID() == obj.GetAOIID() {
			var objList []aoi.IObject
			for _, o := range da.objs {
				objList = append(objList, o)
			}

			if len(objList) > 0 {
				w.OnBatchLeave(objList)
			}
			delete(da.watchers, obj.GetAOIID())
		} else {
			w.OnObjectLeave(obj)
		}
	}

	return nil
}

func (da *_DefaultAOI) Traversal(obj aoi.IObject, cb func(watcher aoi.IWatcher) bool) {
	if obj == nil {
		return
	}

	if _, ok := da.objs[obj.GetAOIID()]; !ok {
		return
	}

	for _, w := range da.watchers {
		cb(w)
	}
}

func (da *_DefaultAOI) Move(obj aoi.IObject) error                         { return nil }
func (da *_DefaultAOI) AddGlobalMarker(obj aoi.IObject) error              { return nil }
func (da *_DefaultAOI) RemoveGlobalMarker(obj aoi.IObject) error           { return nil }
func (da *_DefaultAOI) CreateGroup([]aoi.IObject) (int, error)             { return 0, nil }
func (da *_DefaultAOI) DestroyGroup(groupID int) error                     { return nil }
func (da *_DefaultAOI) AddToGroup(obj aoi.IObject, groupID int) error      { return nil }
func (da *_DefaultAOI) RemoveFromGroup(obj aoi.IObject, groupID int) error { return nil }
func (da *_DefaultAOI) TraversalGroup(obj aoi.IObject, groupID int, cb func(watcher aoi.IWatcher) bool) {
}
