package main

import (
	"aoi"
	"aoi/base/linemath"
	"aoi/toweraoi"
	"fmt"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
)

type Marker struct {
	ID  string
	Pos linemath.Vector2
}

func (m *Marker) GetAOIID() string              { return m.ID }
func (m *Marker) GetCoordPos() linemath.Vector2 { return m.Pos }

type Watcher struct {
	ID  string
	Pos linemath.Vector2
}

func (w *Watcher) GetAOIID() string              { return w.ID }
func (w *Watcher) GetCoordPos() linemath.Vector2 { return w.Pos }
func (w *Watcher) GetVisual() float32            { return 100 }
func (w *Watcher) OnObjectEnter(aoi.IObject)     {}
func (w *Watcher) OnObjectLeave(aoi.IObject)     {}

func main() {
	minPos := linemath.Vector2{}
	maxPos := linemath.Vector2{X: 8000, Y: 8000}
	towerSize := float32(50)

	t, err := toweraoi.New(&toweraoi.Config{MinPos: minPos, MaxPos: maxPos, TowerSize: towerSize})
	if err != nil {
		panic(err)
	}

	// 加入对象, 没有视野,
	objs := make([]*Marker, 200000)
	for i := 0; i < 200000; i++ {
		o := &Marker{}
		o.ID = fmt.Sprintf("obj:%d", i)
		o.Pos.X = rand.Float32() * 8000
		o.Pos.Y = rand.Float32() * 8000
		objs[i] = o
		//t.Add(o)
		t.AddToAOI(o)
	}

	// 加入观察者, 视野100, 周边1个格子
	watchers := make([]*Watcher, 10000)
	for i := 0; i < 10000; i++ {
		w := &Watcher{}
		w.ID = fmt.Sprintf("watcher:%d", i)
		w.Pos.X = rand.Float32() * 8000
		w.Pos.Y = rand.Float32() * 8000
		watchers[i] = w
		t.AddGlobalMarker(w)
	}

	go func() {
		_ = http.ListenAndServe("0.0.0.0:30881", nil)
	}()

	fmt.Println("动起来, 动次打次")

	for {
		w := watchers[rand.Intn(10000)]
		w.Pos.X = rand.Float32() * 8000
		w.Pos.Y = rand.Float32() * 8000
		t.Move(w)
	}
}
