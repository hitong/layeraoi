package defaultaoi

import (
	"aoi"
	"testing"
)

func TestDefaultAOI_AddRemove(t *testing.T) {
	da := New()

	o1 := &aoi.TestWatcher{ID: "obj:1", Visual: 1}
	o2 := &aoi.TestWatcher{ID: "obj:2", Visual: 1}

	t.Log("AddToAOI o1")
	da.AddToAOI(o1)

	t.Log("AddToAOI o2")
	da.AddToAOI(o2)

	t.Log("RemoveFromAOI o1")
	da.RemoveFromAOI(o1)

	t.Log("RemoveFromAOI o2")
	da.RemoveFromAOI(o2)
}
