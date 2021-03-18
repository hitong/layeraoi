package main

import (
	"aoi/base/linemath"
	"aoi/layeraoi"
	"fmt"
	"github.com/hajimehoshi/ebiten/examples/resources/images/blocks"
	g "github.com/magicsea/gosprite"
	"github.com/magicsea/gosprite/ui"
	"image/color"
	"math"
	"math/rand"
	"strconv"
)

type ToolScane struct {
	g.Scene
	aoiScene2 layeraoi.LayerAOI

	objects map[string]*Object
	//watchers map[uint64]*Watcher
	//makers   map[uint64]*Marker

	idgen  uint64
	notice *ui.TextBox
}

func (s *ToolScane) initBg() error {
	//bg
	cell := 50
	var fcellsize = float64(cell)
	bg2, bgErr2 := g.NewSprite(blocks.Background_png)
	if bgErr2 != nil {
		fmt.Println("bg load error:", bgErr2)
		return bgErr2
	}
	bg2.SetSize(g.NewVector(screenW, screenH))
	bg2.SetSpriteType(g.SpriteTypeSlice)
	bg2.SetScale(g.NewVector(fcellsize/32.0, fcellsize/32.0))
	s.AddNode(bg2)

	//load line
	for i := 0; i < screenW/cell; i++ {
		from := g.NewVector(float64(i*cell), 0)
		to := g.NewVector(float64(i*cell), screenH)
		line := g.NewLine(from, to, 1, color.Black)
		line.SetDepth(2)
		s.AddNode(line)
	}
	for j := 0; j < screenH/cell; j++ {
		from := g.NewVector(0, float64(j*cell))
		to := g.NewVector(screenW, float64(j*cell))
		line := g.NewLine(from, to, 1, color.Black)
		line.SetDepth(2)
		s.AddNode(line)
	}
	return nil
}

var objLayerBits uint64 = (1 << 5) - 1
var watchLayerBits uint64 = (1 << 5) - 1

func (s *ToolScane) initUI() error {
	var fromx float64 = screenW - 130
	bg := ui.NewTextBox(g.NewVector(fromx, 0), g.NewVector(130, 260), "", ui.AliVertical_Mid, ui.AliHorizontal_Right, false)
	bg.SetColor(color.RGBA{R: 128, G: 128, A: 128})
	s.AddUINode(bg)

	fromx += 10
	tb := ui.NewTextBox(g.NewVector(fromx, 0), g.NewVector(70, 40), "Marker:", ui.AliVertical_Mid, ui.AliHorizontal_Right, false)
	tb.SetColor(color.RGBA{R: 128, G: 128, B: 128})
	s.AddUINode(tb)
	inf := ui.NewInputField(g.NewVector(screenW-50, 0), g.NewVector(50, 40), "100", ui.AliVertical_Mid, ui.AliHorizontal_Right, false)
	inf.SetColor(color.RGBA{R: 128, G: 128, B: 128, A: 128})
	s.AddUINode(inf)

	//添加marker
	btn := ui.NewButton(g.NewVector(fromx, 50), g.NewVector(60, 40), "Add")
	btn.SetColor(color.RGBA{R: 100, G: 100, B: 100, A: 128})
	btn.SetOnPressed(func(b *ui.Button) {
		num, err := strconv.Atoi(inf.ValueText)
		if num <= 0 || num > 200000 {
			num = 100
			fmt.Println("input num invalid,", err)
		}
		for i := 0; i < num; i++ {
			s.addObject(objLayerBits)
		}
	})
	s.AddUINode(btn)

	btnD := ui.NewButton(g.NewVector(fromx,100),g.NewVector(60,40),"randRemove")
	btnD.SetColor(color.RGBA{R: 100, G: 100, B: 100, A: 128})
	btnD.SetOnPressed(func(b *ui.Button) {
		num, err := strconv.Atoi(inf.ValueText)
		if num <= 0 || num > 200000 {
			num = 100
			fmt.Println("input num invalid,", err)
		}
		s.randRemove(num)
	})
	s.AddUINode(btnD)

	s.notice = ui.NewTextBox(g.NewVector(fromx, 200), g.NewVector(120, 50), "", ui.AliVertical_Mid, ui.AliHorizontal_Left, false)
	s.notice.SetColor(color.RGBA{R: 128, G: 128, B: 128, A: 128})
	s.AddUINode(s.notice)

	return nil
}

func (s *ToolScane)randRemove(num int){
out:
	for num > 0 && len(s.objects) > 0{
		num--
		for k,obj := range s.objects{
			if obj.GetAOIID() == "Obj:1"{
				continue
			}
			obj.Destroy()
			s.aoiScene2.RemoveFromAOI(obj)
			delete(s.objects, k)
			continue out
		}
	}
}

type LoadBalanceSystem struct {
}

func (l *LoadBalanceSystem)LayerBalance(v map[int]int)(int,int){
	var minIdx = 0
	var minNum = math.MaxInt64
	for idx, v := range v {
		if minNum > v {
			minIdx = idx
			minNum = v
		}
	}

	return minIdx,minNum
}


func (s *ToolScane) Init() error {
	s.objects = make(map[string]*Object)
	s.aoiScene2 = layeraoi.New()
	layer1,_ := layeraoi.NewTowerAoi(&layeraoi.Config{MinPos: linemath.Vector2{
		X: 0,
		Y: 0,
	}, MaxPos: linemath.Vector2{
		X: screenW,
		Y: screenH,
	}, TowerSize: 50,LayerLimit: 10,LoadBalanceCfg:&layeraoi.LoadBalanceConfig{MethodObj:LoadBalanceSystem{} }})
	layer2,_ := layeraoi.NewTowerAoi(&layeraoi.Config{MinPos: linemath.Vector2{
		X: 0,
		Y: 0,
	}, MaxPos: linemath.Vector2{
		X: screenW,
		Y: screenH,
	}, TowerSize: 50,LayerLimit: 10,LoadBalanceCfg:&layeraoi.LoadBalanceConfig{MethodObj:LoadBalanceSystem{} }})
	layer3,_ := layeraoi.NewTowerAoi(&layeraoi.Config{MinPos: linemath.Vector2{
		X: 0,
		Y: 0,
	}, MaxPos: linemath.Vector2{
		X: screenW,
		Y: screenH,
	}, TowerSize: 50,LayerLimit: 1000,LoadBalanceCfg:&layeraoi.LoadBalanceConfig{MethodObj:LoadBalanceSystem{} }})
	layer4,_ := layeraoi.NewTowerAoi(&layeraoi.Config{MinPos: linemath.Vector2{
		X: 0,
		Y: 0,
	}, MaxPos: linemath.Vector2{
		X: screenW,
		Y: screenH,
	}, TowerSize: 50,LayerLimit: 10000,LoadBalanceCfg:&layeraoi.LoadBalanceConfig{MethodObj:LoadBalanceSystem{} }})
	layer5,_ := layeraoi.NewTowerAoi(&layeraoi.Config{MinPos: linemath.Vector2{
		X: 0,
		Y: 0,
	}, MaxPos: linemath.Vector2{
		X: screenW,
		Y: screenH,
	}, TowerSize: 50,LayerLimit: 100,LoadBalanceCfg:&layeraoi.LoadBalanceConfig{MethodObj:LoadBalanceSystem{} }})
	s.aoiScene2.AddLayer((1 << 5) - 1,layer1,layer2,layer3,layer4,layer5)
	//s.aoiScene2.AddLayer(1,layer1)
	//watcher := NewObject(s, s.GenMID(), g.NewVector(400, 300), 100)
	//watcher.SetControl(true)
	//s.objects[watcher.ID] = watcher
	//s.aoiScene.AddToAOI(watcher)
	s.addWatcher(true,watchLayerBits)
	s.addWatcher(false,(1 << 5) - 1 - 3)
	for i := 0; i < 100; i++ {
		s.addObject(objLayerBits)
	}

	s.initBg()
	s.initUI()
	return nil
}

func (s *ToolScane) addObject(layerBits uint64) {
	fx := rand.Float64() * screenW
	fy := rand.Float64() *screenH
	o := NewObject(s, s.GenMID(), g.Vector{X: fx, Y: fy}, 50)
	o.LayerBits = layerBits
	fmt.Println(o.GetAOIID())
	s.objects[o.ID] = o
	s.aoiScene2.AddToAOI(o)
}

func (s *ToolScane) addWatcher(ctl bool,layerBits uint64) {
	o := NewObject(s, s.GenMID(), g.NewVector(400, 300), 4)
	w := Watcher{
		Visual: make(map[int]float32),
		Mark:  o,
	}
	o.SetControl(ctl)
	o.LayerBits = layerBits
	s.objects[o.ID] = o
	w.Visual[0] = 250
	w.Visual[1] = 200
	w.Visual[2] = 150
	w.Visual[3] = 100
	w.Visual[4] = 50
	s.aoiScene2.AddToAOI(&w)
}

func (s *ToolScane) GenMID() string {
	s.idgen++
	return fmt.Sprintf("Obj:%d", s.idgen)
}

func (s *ToolScane) Update(detaTime float64) {
	for _, o := range s.objects {
		o.Update(detaTime)
	}
}
