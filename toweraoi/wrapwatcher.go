package toweraoi

import (
	"aoi"
	log "github.com/cihub/seelog"
)

type _WrapWatcher struct {
	aoi.IWatcher

	notify func(watcher *_WrapWatcher)
	info   map[string]*_CacheInfo
}

type _CacheInfo struct {
	obj     aoi.IObject
	count   int
	entered bool
}

func newWrapWatcher(w aoi.IWatcher, notify func(watcher *_WrapWatcher)) *_WrapWatcher {
	wo := &_WrapWatcher{}
	wo.IWatcher = w
	wo.notify = notify
	wo.info = make(map[string]*_CacheInfo)

	return wo
}

func (wo *_WrapWatcher) OnObjectEnter(obj aoi.IObject) {
	if _, ok := wo.info[obj.GetAOIID()]; !ok {
		wo.info[obj.GetAOIID()] = &_CacheInfo{obj: obj}
	}

	wo.info[obj.GetAOIID()].count++
	wo.notify(wo)
}

func (wo *_WrapWatcher) OnObjectLeave(obj aoi.IObject) {
	if info, ok := wo.info[obj.GetAOIID()]; ok {
		info.count--
		wo.notify(wo)
	} else {
		log.Debug("非常奇怪的事情发生了, 检查代码和日志!")
	}
}

func (wo *_WrapWatcher) OnBatchEnter(objs []aoi.IObject) {
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

func (wo *_WrapWatcher) OnBatchLeave(objs []aoi.IObject) {
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

func (wo *_WrapWatcher) Flush() {
	enterList := make([]aoi.IObject, 0, 1)
	leaveList := make([]aoi.IObject, 0, 1)

	for _, info := range wo.info {
		if info.count > 0 && !info.entered {
			//wo.IWatcher.OnObjectEnter(info.obj)
			enterList = append(enterList, info.obj)
			info.entered = true
		} else if info.count <= 0 {
			if info.entered {
				//wo.IWatcher.OnObjectLeave(info.obj)
				leaveList = append(leaveList, info.obj)
			}
			delete(wo.info, info.obj.GetAOIID())
		}
	}

	if len(enterList) > 0 {
		wo.IWatcher.OnBatchEnter(enterList)
	}
	if len(leaveList) > 0 {
		wo.IWatcher.OnBatchLeave(leaveList)
	}
}
