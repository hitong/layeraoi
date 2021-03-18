package layeraoi

import (
	"aoi"
	log "github.com/cihub/seelog"
)

type _WrapWatcher struct {
	ILayerWatcher

	notify func(watcher *_WrapWatcher)
	info   map[string]*_CacheInfo
}

type _CacheInfo struct {
	obj     aoi.IObject
	count   int
	entered bool
}

func newWrapWatcher(w ILayerWatcher, notify func(watcher *_WrapWatcher)) *_WrapWatcher {
	wo := &_WrapWatcher{}
	wo.ILayerWatcher = w
	wo.notify = notify
	wo.info = make(map[string]*_CacheInfo)
	return wo
}

func (wo *_WrapWatcher) OnLayerObjectEnter(obj aoi.IObject,layer int) {
	if _, ok := wo.info[obj.GetAOIID()]; !ok {
		wo.info[obj.GetAOIID()] = &_CacheInfo{obj: obj}
	}

	wo.info[obj.GetAOIID()].count++
	wo.notify(wo)
}

func (wo *_WrapWatcher) OnLayerObjectLeave(obj aoi.IObject,layer int) {
	if info, ok := wo.info[obj.GetAOIID()]; ok {
		info.count--
		wo.notify(wo)
	} else {
		log.Debug("非常奇怪的事情发生了, 检查代码和日志!")
	}
}

func (wo *_WrapWatcher) OnLayerBatchEnter(objs []aoi.IObject,layer int) {
	if len(objs) == 0 {
		return
	}

	for _, o := range objs {
		if _, ok := wo.info[o.GetAOIID()]; !ok {
			wo.info[o.GetAOIID()] = &_CacheInfo{obj: o}
		}
		wo.info[o.GetAOIID()].count++
	}

	wo.notify(wo)
}

func (wo *_WrapWatcher) OnLayerBatchLeave(objs []aoi.IObject,layer int) {
	if len(objs) == 0 {
		return
	}

	for _, o := range objs {
		if info, ok := wo.info[o.GetAOIID()]; ok {
			info.count--
		} else {
			log.Debug("非常奇怪的事情发生了, 检查代码和日志!")
		}
	}

	wo.notify(wo)
}

func (wo *_WrapWatcher) Flush(layer int) {
	enterList := make([]aoi.IObject, 0, 1)
	leaveList := make([]aoi.IObject, 0, 1)

	for _, info := range wo.info {
		if info.count > 0 && !info.entered {
			//wo.IWatcher.OnObjectEnter(info.obj)
			enterList = append(enterList, info.obj)
			info.entered = true
		} else if info.count <= 0 {
			if info.entered {
				//wo.IWatcher.OnObjectLeave(info.obj,layer)
				leaveList = append(leaveList, info.obj)
			}
			delete(wo.info, info.obj.GetAOIID())
		}
	}

	if len(enterList) > 0 {
		wo.ILayerWatcher.OnLayerBatchEnter(enterList,layer)
	}
	if len(leaveList) > 0 {
		wo.ILayerWatcher.OnLayerBatchLeave(leaveList,layer)
	}
}
