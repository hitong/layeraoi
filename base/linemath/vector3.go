package linemath

import (
	"errors"
	"fmt"
	"math"
)

type Vector3 struct {
	X float32
	Y float32
	Z float32
}

// CreateVector3 创建一个新的矢量
func CreateVector3(x, y, z float32) Vector3 {
	return Vector3{
		x,
		y,
		z,
	}
}

// Vector3_Zero 返回零值
func Vector3_Zero() Vector3 {
	return Vector3{
		0,
		0,
		0,
	}
}

// Vector3_Invalid 返加一个无效的值 ，未赋值之前
func Vector3_Invalid() Vector3 {
	return Vector3{
		math.MaxFloat32,
		math.MaxFloat32,
		math.MaxFloat32,
	}
}

func (v Vector3) String() string {
	return fmt.Sprintf("{X:%f Y:%f Z:%f}", v.X, v.Y, v.Z)
}

func (v Vector3) ToV2() Vector2 {
	return Vector2{X: v.X, Y: v.Y}
}

// IsInValid 是否有效
func (v Vector3) IsInValid() bool {
	return v.IsEqual(Vector3_Invalid())
}

// IsZero 是否默认
func (v Vector3) IsZero() bool {
	return v.IsEqual(Vector3_Zero())
}

// IsEqual 相等
func (v Vector3) IsEqual(r Vector3) bool {
	if v.X-r.X > math.SmallestNonzeroFloat32 ||
		v.X-r.X < -math.SmallestNonzeroFloat32 ||
		v.Y-r.Y > math.SmallestNonzeroFloat32 ||
		v.Y-r.Y < -math.SmallestNonzeroFloat32 ||
		v.Z-r.Z > math.SmallestNonzeroFloat32 ||
		v.Z-r.Z < -math.SmallestNonzeroFloat32 {
		return false
	}

	return true
}

// Add 加
func (v Vector3) Add(o Vector3) Vector3 {
	return Vector3{v.X + o.X, v.Y + o.Y, v.Z + o.Z}
}

// AddS 加到自己身上
func (v *Vector3) AddS(o Vector3) {
	v.X += o.X
	v.Y += o.Y
	v.Z += o.Z
}

// Sub 减
func (v Vector3) Sub(o Vector3) Vector3 {
	return Vector3{v.X - o.X, v.Y - o.Y, v.Z - o.Z}
}

// Distance 距离
func (v Vector3) Distance(o Vector3) float32 {
	return v.Sub(o).Len()
}

// DistanceV2 2d距离
func (v Vector3) DistanceV2(o Vector3) float32 {
	return v.ToV2().Sub(o.ToV2()).Len()
}

// SubS 自已身上减
func (v *Vector3) SubS(o Vector3) {
	v.X -= o.X
	v.Y -= o.Y
	v.Z -= o.Z
}

// Mul 乘
func (v Vector3) Mul(o float32) Vector3 {
	return Vector3{v.X * o, v.Y * o, v.Z * o}
}

// MulS 自己乘
func (v *Vector3) MulS(o float32) {
	v.X *= o
	v.Y *= o
	v.Z *= o
}

// Cross 叉乘
func (v Vector3) Cross(o Vector3) Vector3 {
	return Vector3{v.Y*o.Z - v.Z*o.Y, v.Z*o.X - v.X*o.Z, v.X*o.Y - v.Y*o.X}
}

// Dot 点乘
func (v Vector3) Dot(o Vector3) float32 {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}

// Len 获取长度
func (v Vector3) Len() float32 {
	return float32(math.Sqrt(float64(v.Dot(v))))
}

// AngleFromXYFloor XY平面上仰角
func (v Vector3) AngleFromXYFloor() (float32, error) {
	if v.Len() == 0 {
		return 0, errors.New("zero len")
	}
	a := math.Asin(Clamp64(float64(v.Z/v.Len()), -1, 1))
	return float32(a) * 180 / math.Pi, nil
}

// AngleFromXYFloorEx 返回向量与XY平面夹角的弧度
func (v Vector3) AngleFromXYFloorEx() (float32, error) {
	if v.Len() == 0 {
		return 0, errors.New("zero len")
	}
	a := math.Asin(Clamp64(float64(v.Z/v.Len()), -1, 1))
	return float32(a), nil
}

func (v *Vector3) Normalize() {
	len := v.Len()

	if len < math.SmallestNonzeroFloat32 {
		return
	}

	v.X = v.X / len
	v.Y = v.Y / len
	v.Z = v.Z / len
}

func (v Vector3) Normalized() Vector3 {
	len := v.Len()

	if len < math.SmallestNonzeroFloat32 {
		return Vector3_Zero()
	}
	newv := Vector3{
		X: v.X / len,
		Y: v.Y / len,
		Z: v.Z / len,
	}

	return newv
}

func (v Vector3) ToRotationMatrix() Matrix3f {
	var (
		SP float64 = math.Sin(float64(v.Y))
		SY float64 = math.Sin(float64(v.Z))
		SR float64 = math.Sin(float64(v.X))
		CP float64 = math.Cos(float64(v.Y))
		CY float64 = math.Cos(float64(v.Z))
		CR float64 = math.Cos(float64(v.X))
	)
	result := Matrix3f{
		[3][3]float64{
			{CP * CY, CP * SY, SP},
			{SR*SP*CY - CR*SY, SR*SP*SY + CR*CY, -SR * CP},
			{-(CR*SP*CY + SR*SY), CY*SR - CR*SP*SY, CR * CP},
		},
	}
	// P:Y Y:Z R:X

	return result
}