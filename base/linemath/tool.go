package linemath

import (
	"math"
	"math/rand"
	"time"
)

func Clamp32(v float32, min float32, max float32) float32 {
	if v < min {
		return min
	} else if v > max {
		return max
	}

	return v
}

func Clamp64(v float64, min float64, max float64) float64 {
	if v < min {
		return min
	} else if v > max {
		return max
	}

	return v
}

// RandXZ 在XZ平面上半径为r的圆内选取一个随机点
func RandXZ(v Vector3, r float32) Vector3 {
	randSeed := rand.New(rand.NewSource(time.Now().UnixNano()))

	tarR := randSeed.Float64() * float64(r)
	angle := randSeed.Float64() * 2 * math.Pi

	pos := Vector3{}
	pos.Y = 0

	pos.X = float32(math.Cos(angle) * tarR)
	pos.Z = float32(math.Sin(angle) * tarR)

	return v.Add(pos)
}

const PI2 = math.Pi * 2
const PI_HALF = math.Pi * 0.5

// 四元数转欧拉角
func QuaternionToEuler(v Quaternion) Vector3 {
	if v.IsInValid() {
		return Vector3_Invalid()
	}

	angle := Vector3_Zero()
	x := math.Atan2(float64(2*(v.W*v.X+v.Y*v.Z)), 1-float64(2*(v.X*v.X+v.Y*v.Y)))

	angle.X = float32(x)

	var y float64
	sinp := Clamp64(float64(2*(v.W*v.Y-v.Z*v.X)), -1, 1)
	if math.Abs(sinp) >= 1 {
		y = math.Copysign(math.SqrtPi, sinp)
	} else {
		y = math.Asin(sinp)
	}

	angle.Y = float32(y)

	z := math.Atan2(float64(2*(v.W*v.Z+v.X*v.Y)), 1-float64(2*(v.Y*v.Y+v.Z*v.Z)))
	
	angle.Z = float32(z)

	return angle
}

// 欧拉角转四元数
func EulerToQuaternion(v Vector3) Quaternion {
	if v.IsInValid() {
		return Quaternion_Invalid()
	}

	var SP, SY, SR float64
	var CP, CY, CR float64

	SR, CR = math.Sincos(float64(v.X) / 2)
	SP, CP = math.Sincos(float64(v.Y) / 2)
	SY, CY = math.Sincos(float64(v.Z) / 2)

	quat := Quaternion_Zero()
	quat.W = float32(CY*CP*CR + SY*SP*SR)
	quat.X = float32(CY*CP*SR - SY*SP*CR)
	quat.Y = float32(SY*CP*SR + CY*SP*CR)
	quat.Z = float32(SY*CP*CR - CY*SP*SR)

	return quat
}

//归一化为(-180,180]角度
func ClampAngle(r float32) float32 {
	x := math.Mod(float64(r), 360)
	if x < 0 {
		x += 360
	}
	if x > 180 {
		x = x - 360
	}
	return float32(x)
}
