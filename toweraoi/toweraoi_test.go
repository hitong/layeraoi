package toweraoi

import (
	"aoi"
	"aoi/base/linemath"
	"fmt"
	"math/rand"
	"testing"
)

func TestTowerAOI_AddRemove(t *testing.T) {
	minPos := linemath.Vector2{X: -5, Y: -5}
	maxPos := linemath.Vector2{X: 5, Y: 5}
	ta, err := New(&Config{
		MinPos:    minPos,
		MaxPos:    maxPos,
		TowerSize: 1,
	})
	if err != nil {
		t.Fatal(err)
	}

	o1 := &aoi.TestWatcher{ID: "obj:1", Visual: 2, Pos: linemath.Vector2{X: -4.5, Y: -4.5}}
	t.Log("Obj1 Add")
	ta.AddToAOI(o1)

	o2 := &aoi.TestWatcher{ID: "obj:2", Visual: 2, Pos: linemath.Vector2{X: -3.5, Y: -3.5}}
	t.Log("Obj2 Add")
	ta.Add(o2)

	t.Log("Obj1 Remove")
	ta.Remove(o1)

	t.Log("Obj2 Remove")
	ta.Remove(o2)
}

func TestTowerAOI_Global(t *testing.T) {
	minPos := linemath.Vector2{X: -5, Y: -5}
	maxPos := linemath.Vector2{X: 5, Y: 5}
	ta, err := New(&Config{
		MinPos:    minPos,
		MaxPos:    maxPos,
		TowerSize: 1,
	})
	if err != nil {
		t.Fatal(err)
	}

	o1 := &aoi.TestWatcher{ID: "obj:1", Visual: 2}
	t.Log("Obj1 Add", o1.Pos)
	ta.Add(o1)

	global := &aoi.TestMarker{ID: "global"}
	t.Log("Marker Add normal", global.Pos)
	ta.Add(global)

	t.Log("Marker Add global")
	ta.AddGlobalMarker(global)

	global.Pos.X = -4
	global.Pos.Y = -4
	t.Log("Global Move", global.Pos)
	ta.Move(global)

	t.Log("Global Remove")
	ta.RemoveGlobalMarker(global)
}

func TestTowerAOI_Move(t *testing.T) {
	minPos := linemath.Vector2{X: -5, Y: -5}
	maxPos := linemath.Vector2{X: 5, Y: 5}
	ta, err := New(&Config{
		MinPos:    minPos,
		MaxPos:    maxPos,
		TowerSize: 1,
	})
	if err != nil {
		t.Fatal(err)
	}

	pos := linemath.Vector2{X: -5, Y: -5}
	o1 := &aoi.TestWatcher{ID: "obj:1", Visual: 2, Pos: pos}
	t.Log("Obj1 Add", o1.Pos)
	ta.Add(o1)

	pos.X = -3.5
	pos.Y = -3.5
	o2 := &aoi.TestWatcher{ID: "obj:2", Visual: 2, Pos: pos}
	t.Log("Obj2 Add", o2.Pos)
	ta.Add(o2)

	o1.Pos.X = -4.2
	o1.Pos.Y = -4.2
	t.Log("Obj1 Move", o1.Pos, "相同格子, 不触发")
	ta.Move(o1)

	o1.Pos.X = -3.5
	o1.Pos.Y = -3.5
	t.Log("Obj1 Move", o1.Pos, "跨格子, 但原来格子还在视野内")
	ta.Move(o1)

	o1.Pos.X = 4
	o1.Pos.Y = 4
	t.Log("Obj1 Move", o1.Pos, "跨格子, 原来的格子不在视野内")
	ta.Move(o1)

	t.Log("Obj1 Remove")
	ta.Remove(o1)

	t.Log("Obj2 Remove")
	ta.Remove(o2)
}

func TestTowerAOI_Group(t *testing.T) {
	minPos := linemath.Vector2{X: -5, Y: -5}
	maxPos := linemath.Vector2{X: 5, Y: 5}
	ta, err := New(&Config{
		MinPos:    minPos,
		MaxPos:    maxPos,
		TowerSize: 1,
	})
	if err != nil {
		t.Fatal(err)
	}

	o1 := &aoi.TestWatcher{ID: "obj:1", Visual: 2, Pos: linemath.Vector2{}}
	t.Log("Obj1 Add", o1.Pos)
	ta.Add(o1)

	o2 := &aoi.TestWatcher{ID: "obj:2", Visual: 2, Pos: linemath.Vector2{X: -3.5, Y: -3.5}}
	t.Log("Obj2 Add", o2.Pos)
	ta.Add(o2)

	t.Log("CreateGroup o1 o2")
	id, _ := ta.CreateGroup([]aoi.IObject{o1, o2})

	//o2.Pos.X = 0
	//o2.Pos.Y = 0
	//t.Log("Obj2 Move", o2.Pos)
	//ta.Move(o2)

	t.Log("DestroyGroup")
	ta.DestroyGroup(id)
}

func TestTowerAOI_Traversal(t *testing.T) {
	minPos := linemath.Vector2{X: -5, Y: -5}
	maxPos := linemath.Vector2{X: 5, Y: 5}
	ta, err := New(&Config{
		MinPos:    minPos,
		MaxPos:    maxPos,
		TowerSize: 1,
	})
	if err != nil {
		t.Fatal(err)
	}

	o1 := &aoi.TestWatcher{ID: "obj:1", Visual: 2}
	t.Log("Obj1 Add", o1.Pos)
	ta.Add(o1)

	o2 := &aoi.TestWatcher{ID: "obj:2", Visual: 2}
	t.Log("Obj2 Add", o2.Pos)
	ta.Add(o2)

	t.Log("Traversal o1, o2")
	ta.Traversal(o1, func(w aoi.IWatcher) bool {
		fmt.Println("Traversal", w.GetAOIID())
		return true
	})

	o2.Pos.X = -3.5
	o2.Pos.Y = -3.5
	t.Log("o2 move out")
	ta.Move(o2)

	t.Log("Traversal o1")
	ta.Traversal(o1, func(w aoi.IWatcher) bool {
		fmt.Println("Traversal", w.GetAOIID())
		return true
	})

	t.Log("CreateGroup o1, o2")
	groupID, _ := ta.CreateGroup([]aoi.IObject{o1, o2})

	t.Log("TraversalGroup o1, o2")
	ta.TraversalGroup(o1, groupID, func(w aoi.IWatcher) bool {
		fmt.Println("TraversalGroup", w.GetAOIID())
		return true
	})

	t.Log("Traversal o1, o2")
	ta.Traversal(o1, func(w aoi.IWatcher) bool {
		fmt.Println("Traversal", w.GetAOIID())
		return true
	})

	t.Log("Traversal o1 or o2")
	ta.Traversal(o1, func(w aoi.IWatcher) bool {
		fmt.Println("Traversal", w.GetAOIID())
		return false
	})
}

func BenchmarkTowerAOI_200000_10000_Traversal(b *testing.B) {
	b.ReportAllocs()

	minPos := linemath.Vector2{}
	maxPos := linemath.Vector2{X: 8000, Y: 8000}
	ta, err := New(&Config{
		MinPos:    minPos,
		MaxPos:    maxPos,
		TowerSize: 50,
	})
	if err != nil {
		b.Fatal(err)
	}

	objs := make(map[int]*aoi.TestMarker)
	for i := 0; i < 200000; i++ {
		o := &aoi.TestMarker{ID: fmt.Sprintf("obj:%d", i)}
		o.Pos.X = rand.Float32() * 8000
		o.Pos.Y = rand.Float32() * 8000
		ta.Add(o)

		objs[i] = o
	}

	for i := 0; i < 10000; i++ {
		w := &aoi.TestWatcher{ID: fmt.Sprintf("watcher:%d", i), Visual: 100}
		w.Pos.X = rand.Float32() * 8000
		w.Pos.Y = rand.Float32() * 8000
		ta.Add(w)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ta.Traversal(objs[rand.Intn(200000)], func(w aoi.IWatcher) bool {
			return true
		})
	}
}

func BenchmarkTowerAOI_200000_10000_Move(b *testing.B) {
	b.ReportAllocs()

	minPos := linemath.Vector2{}
	maxPos := linemath.Vector2{X: 8000, Y: 8000}
	ta, err := New(&Config{
		MinPos:    minPos,
		MaxPos:    maxPos,
		TowerSize: 50,
	})
	if err != nil {
		b.Fatal(err)
	}

	objs := make(map[int]*aoi.TestMarker)
	for i := 0; i < 200000; i++ {
		o := &aoi.TestMarker{ID: fmt.Sprintf("obj:%d", i)}
		o.Pos.X = rand.Float32() * 8000
		o.Pos.Y = rand.Float32() * 8000
		ta.Add(o)

		objs[i] = o
	}

	for i := 0; i < 10000; i++ {
		w := &aoi.TestWatcher{ID: fmt.Sprintf("watcher:%d", i), Visual: 100}
		w.Pos.X = rand.Float32() * 8000
		w.Pos.Y = rand.Float32() * 8000
		ta.Add(w)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := rand.Intn(200000)
		objs[index].Pos.X = rand.Float32() * 8000
		objs[index].Pos.Y = rand.Float32() * 8000
		ta.Move(objs[index])
	}
}

func BenchmarkTowerAOI_200000_10000_MoveWatcher(b *testing.B) {
	b.ReportAllocs()

	minPos := linemath.Vector2{}
	maxPos := linemath.Vector2{X: 8000, Y: 8000}
	ta, err := New(&Config{
		MinPos:    minPos,
		MaxPos:    maxPos,
		TowerSize: 50,
	})
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < 200000; i++ {
		o := &aoi.TestMarker{ID: fmt.Sprintf("obj:%d", i)}
		o.Pos.X = rand.Float32() * 8000
		o.Pos.Y = rand.Float32() * 8000
		ta.Add(o)
	}

	objs := make(map[int]*aoi.TestWatcher)
	for i := 0; i < 10000; i++ {
		w := &aoi.TestWatcher{ID: fmt.Sprintf("watcher:%d", i), Visual: 100}
		w.Pos.X = rand.Float32() * 8000
		w.Pos.Y = rand.Float32() * 8000
		ta.Add(w)

		objs[i] = w
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := rand.Intn(10000)
		objs[index].Pos.X = rand.Float32() * 8000
		objs[index].Pos.Y = rand.Float32() * 8000
		ta.Move(objs[index])
	}
}
