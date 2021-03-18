package main

import (
	"aoi"
	"aoi/base/linemath"
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/resources/images"
	g "github.com/magicsea/gosprite"
	"image"
	"image/color"
	"math/rand"
)

type Object struct {
	ID string

	Pos    linemath.Vector2
	Visual float32

	root  *g.EmptyNode
	scene *ToolScane
	body  *g.Circle

	moveVel g.Vector
	aiTimer int

	inView    int
	isControl bool
	LayerBits uint64
}

func (o *Object)GetLayerBits()uint64{
	return o.LayerBits
}

func NewObject(scene *ToolScane, id string, pos g.Vector, visual float32) *Object {
	o := &Object{
		scene: scene,
	}
	o.ID = id
	o.Visual = visual
	o.LayerBits = 1
	o.Init(pos)
	return o
}

func (o *Object) Init(pos g.Vector) {
	o.root = g.NewEmptyNode()
	o.root.SetPosition(pos)
	o.root.SetDepth(100)
	o.scene.AddNode(o.root)
	//load ani
	frameAniInfo := map[string]*g.FrameAniData{
		"idle": &g.FrameAniData{
			FrameInterval: 10,
			AniName:       "idle",
			FrameRects: []image.Rectangle{
				image.Rect(0, 0, 32, 32),
				image.Rect(1*32, 0, 2*32, 32),
				image.Rect(2*32, 0, 3*32, 32),
				image.Rect(3*32, 0, 4*32, 32),
			},
		},
		"run": &g.FrameAniData{
			FrameInterval: 10,
			AniName:       "run",
			FrameRects: []image.Rectangle{
				image.Rect(0, 32, 32, 64),
				image.Rect(1*32, 32, 2*32, 64),
				image.Rect(2*32, 32, 3*32, 64),
				image.Rect(3*32, 32, 4*32, 64),
				image.Rect(4*32, 32, 5*32, 64),
				image.Rect(5*32, 32, 6*32, 64),
				image.Rect(6*32, 32, 7*32, 64),
				image.Rect(7*32, 32, 8*32, 64),
			},
		},
		"jump": &g.FrameAniData{
			FrameInterval: 10,
			AniName:       "jump",
			FrameRects: []image.Rectangle{
				image.Rect(0, 64, 32, 96),
				image.Rect(1*32, 64, 2*32, 96),
				image.Rect(2*32, 64, 3*32, 96),
				image.Rect(3*32, 64, 4*32, 96),
			},
		},
	}

	as, _ := g.NewAniSprite(images.Runner_png, frameAniInfo)
	//as.SetScale(g.NewVector(3,3))
	as.SetDepth(1)
	as.SetParent(o.root)
	as.SetLocalPosition(g.VectorZero())
	as.Play("run")

	c := color.RGBA{G: 128, A: 128}
	o.body = g.NewCircle(10, c)
	o.body.SetDepth(1)
	o.body.SetParent(o.root)
	o.body.SetLocalPosition(g.VectorZero())
	//txt
	txt := g.NewText(o.ID, 8, color.RGBA{R: 255, A: 255})
	txt.SetDepth(2)
	txt.SetParent(o.root)
	txt.SetLocalPosition(g.NewVector(-4, 4))

	o.TransPos()
}

func (o *Object) TransPos() {
	o.Pos.X = float32(o.root.GetPosition().X)
	o.Pos.Y = float32(o.root.GetPosition().Y)
}

func (o *Object) Update(detaTime float64) {

	o.TransPos()

	var speed = 60 * detaTime

	if !o.isControl {
		c := color.RGBA{B: 255, A: 0}
		if o.inView & 1 != 0 {
			c = color.RGBA{R: 255, A: 255}
		}
		if o.inView & 2 != 0{
			c = color.RGBA{R: 100,B:100, A: 255}
		}
		if o.inView & 4 != 0 {
			c = color.RGBA{R: 100, B: 100, G: 100, A: 255}
		}
		if o.inView & 8 != 0{
			c = color.RGBA{R: 50,B:20,G:255, A: 255}
		}
		if o.inView & 16 != 0{
			c = color.RGBA{R: 160,B:250,G:100, A: 255}
		}
		o.body.SetColor(c)

		o.aiTimer--
		if o.aiTimer < 0 {
			o.RecountAI()
			o.aiTimer = rand.Int()%100 + 200
		}

		pos := o.root.GetPosition().Add(o.moveVel.Mul(detaTime))
		if pos.X > screenW || pos.X < 0 {
			o.RecountAI()
			return
		}

		if pos.Y > screenH || pos.Y < 0 {
			o.RecountAI()
			return
		}

		o.root.SetPosition(pos)

		o.TransPos()

		o.scene.aoiScene2.Move(o)
		return
	}

	var charX = o.root.GetPosition().X
	var charY = o.root.GetPosition().Y
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		charX -= speed
	} else if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		charX += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		charY -= speed
	} else if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		charY += speed
	}
	o.root.SetPosition(g.NewVector(charX, charY))

	o.TransPos()
	o.scene.aoiScene2.Move(o)
}

func (o *Object) SetControl(b bool) {
	o.isControl = b
}

func (o *Object) GetAOIID() string {
	return o.ID
}

func (o *Object) RecountAI() {
	var speed float64 = 60
	o.moveVel = g.NewVector(rand.Float64()*2-1, rand.Float64()*2-1).Normal().Mul(speed)
}

func (o *Object) SetInView(b int) {
	o.inView += b
}

func (o *Object) Destroy() {
	o.root.Destory()
	fmt.Println("del" + o.GetAOIID())
}

func (o *Object) GetVisual() float32 {
	return o.Visual
}

func (o *Object) OnLayerObjectEnter(obj aoi.IObject,layer int) {
	switch obj.(type) {
	case *Watcher:
		obj.(*Watcher).Mark.inView |= 1 << layer
	case *Object:
		obj.(*Object).inView |= 1 << layer
	}
}

func (o *Object) OnLayerObjectLeave(obj aoi.IObject,layer int) {
	switch obj.(type) {
	case *Watcher:
		obj.(*Watcher).Mark.inView = obj.(*Watcher).Mark.inView &^ (1 << layer)
	case *Object:
		if (layer == 2 && obj.(*Object).inView  == 5) || (layer == 3 && obj.(*Object).inView  == 9) {
			print("")
		}
		obj.(*Object).inView = obj.(*Object).inView &^ (1 << layer)
	}
}

func (o *Object) GetCoordPos() linemath.Vector2 {
	return o.Pos
}
