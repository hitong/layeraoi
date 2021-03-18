package aoi

import (
	"aoi/base/linemath"
	"errors"
)

// IAOIBase 最基本的AOI接口
type IAOIBase interface {
	AddToAOI(obj IObject) error
	RemoveFromAOI(obj IObject) error
	Move(obj IObject) error

	Traversal(obj IObject, cb func(watcher IWatcher) bool)
}

// IAOI 完整的AOI接口
type IAOI interface {
	IAOIBase

	AddGlobalMarker(obj IObject) error
	RemoveGlobalMarker(obj IObject) error

	CreateGroup([]IObject) (int, error)
	DestroyGroup(groupID int) error
	AddToGroup(obj IObject, groupID int) error
	RemoveFromGroup(obj IObject, groupID int) error
	TraversalGroup(obj IObject, groupID int, cb func(watcher IWatcher) bool)
}

type IObject interface {
	GetAOIID() string
	GetCoordPos() linemath.Vector2
}

type IWatcher interface {
	IObject

	GetVisual() float32
	OnObjectEnter(obj IObject)
	OnObjectLeave(obj IObject)
	OnBatchEnter(objs []IObject)
	OnBatchLeave(objs []IObject)
}

var (
	ErrObjectInvalid    = errors.New("object invalid")
	ErrObjectExisted    = errors.New("object existed")
	ErrObjectNotExisted = errors.New("object not existed")

	ErrPosInvalid = errors.New("pos invalid")

	ErrGroupNotExisted = errors.New("group not existed")
)
