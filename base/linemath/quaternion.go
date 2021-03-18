package linemath

import (
	"fmt"
	"math"
)

type Quaternion struct {
	X float32
	Y float32
	Z float32
	W float32
}

// 创建一个新的四元数
func CreateQuaternion(x, y, z, w float32) Quaternion {
	return Quaternion{
		x,
		y,
		z,
		w,
	}
}

// 返回零值
func Quaternion_Zero() Quaternion {
	return Quaternion{
		0,
		0,
		0,
		0,
	}
}

// 返加一个无效的值 ，未赋值之前
func Quaternion_Invalid() Quaternion {
	return Quaternion{
		math.MaxFloat32,
		math.MaxFloat32,
		math.MaxFloat32,
		math.MaxFloat32,
	}
}

func (v Quaternion) String() string {
	return fmt.Sprintf("X:%f Y:%f Z:%f W:%f", v.X, v.Y, v.Z, v.W)
}

// 是否有效
func (v Quaternion) IsInValid() bool {
	return v.IsEqual(Quaternion_Invalid())
}

// 相等
func (v Quaternion) IsEqual(r Quaternion) bool {
	if v.X-r.X > math.SmallestNonzeroFloat32 ||
		v.X-r.X < -math.SmallestNonzeroFloat32 ||
		v.Y-r.Y > math.SmallestNonzeroFloat32 ||
		v.Y-r.Y < -math.SmallestNonzeroFloat32 ||
		v.Z-r.Z > math.SmallestNonzeroFloat32 ||
		v.Z-r.Z < -math.SmallestNonzeroFloat32 ||
		v.W-r.W > math.SmallestNonzeroFloat32 ||
		v.W-r.W < -math.SmallestNonzeroFloat32 {
		return false
	}

	return true
}