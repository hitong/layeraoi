package linemath

import "math"

type Vector2 struct {
	X float32
	Y float32
}

// CreateVector2 创建一个新的Vector2
func CreateVector2(x, y float32) Vector2 {
	return Vector2{
		x,
		y,
	}
}

// NewVector2FromAngleX 角度转向量 angle=(-180~180]
func NewVector2FromAngleX(angle float32) Vector2 {
	e := float64(angle / 180 * math.Pi)

	return Vector2{
		float32(math.Cos(e)),
		float32(math.Sin(e)),
	}
}

// Vector2_Zero 返回零值
func Vector2_Zero() Vector2 {
	return Vector2{
		0,
		0,
	}
}

// Vector2_Invalid 返加一个无效的值 ，未赋值之前
func Vector2_Invalid() Vector2 {
	return Vector2{
		math.MaxFloat32,
		math.MaxFloat32,
	}
}

// IsInValid 是否有效
func (v Vector2) IsInValid() bool {
	return v.IsEqual(Vector2_Invalid())
}

// IsEqual 相等
func (v Vector2) IsEqual(r Vector2) bool {
	return v.X == r.X && v.Y == r.Y
}

// Add 加
func (v Vector2) Add(o Vector2) Vector2 {
	return Vector2{v.X + o.X, v.Y + o.Y}
}

// Reverse 反方向
func (v Vector2) Reverse() Vector2 {
	return Vector2{-v.X, -v.Y}
}

// AddS 加到自己身上
func (v *Vector2) AddS(o Vector2) {
	v.X += o.X
	v.Y += o.Y
}

// Sub 减
func (v Vector2) Sub(o Vector2) Vector2 {
	return Vector2{v.X - o.X, v.Y - o.Y}
}

// SubS 自已身上减
func (v *Vector2) SubS(o Vector2) {
	v.X -= o.X
	v.Y -= o.Y
}

// Mul 乘
func (v Vector2) Mul(m float32) Vector2 {
	return Vector2{
		v.X * m,
		v.Y * m,
	}
}

// Dot 点乘
func (v Vector2) Dot(o Vector2) float32 {
	return v.X*o.X + v.Y*o.Y
}

// Len 获取长度
func (v Vector2) Len() float32 {
	return float32(math.Sqrt(float64(v.Dot(v))))
}

// Cross 叉乘
func (v Vector2) Cross(o Vector2) float32 {
	return v.X*o.Y - v.Y*o.X
}

//Vector2 归一
func (v *Vector2) Normalize() {
	len := v.Len()
	if len < math.SmallestNonzeroFloat32 {
		return
	}

	v.X = v.X / len
	v.Y = v.Y / len
}

func (v Vector2) Normalized() Vector2 {
	len := v.Len()
	if len < math.SmallestNonzeroFloat32 {
		return v
	}

	return Vector2{v.X / len, v.Y / len}
}

//AngleX 与x的夹角 (-180~180]
func (v Vector2) AngleX() float32 {
	if v.X == 0 && v.Y == 0 {
		return 0
	}

	a := float32(180 / math.Pi * math.Acos(float64(v.X/v.Len())))
	if v.Y < 0 {
		a = -a
	}
	return a
}

//Rotation 向量顺时针围绕圆点旋转 angle 度
func (v Vector2) Rotation(angle float64) Vector2 {
	//角度转化为弧度
	angle = angle * math.Pi / 180
	targetV2 := Vector2{
		X: v.X*float32(math.Cos(angle)) + v.Y*float32(math.Sin(angle)),
		Y: -1*v.X*float32(math.Sin(angle)) + v.Y*float32(math.Cos(angle)),
	}
	return targetV2
}

//IncludedAngle 向量夹角(-180~180] v到tv的夹角
//func (v Vector2) BetweenAngle(tv Vector2) float32 {
//
//	cosVal := v.Dot(tv) / (v.Len() * tv.Len())
//	cross := v.Cross(tv)
//
//	var a float32
//	ac := math.Acos(float64(cosVal))
//	if ac == 0 {
//		a = 0
//	} else {
//		a = float32(180 / math.Pi * math.Acos(float64(cosVal)))
//	}
//	if cross < 0 {
//		a = -a
//	}
//	//fmt.Println("cs:", cross)
//	return a
//}

func (v Vector2) BetweenAngle(tv Vector2) float32 {
	val := math.Atan2(float64(tv.Y), float64(tv.X)) - math.Atan2(float64(v.Y), float64(v.X))
	a := float32(val) * 180 / math.Pi
	return ClampAngle(a)
}