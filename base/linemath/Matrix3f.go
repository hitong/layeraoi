package linemath

import (
	"errors"
	"math"
)

type Matrix3f struct {
	Data [3][3]float64
}

func (mat *Matrix3f) Clear() {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			mat.Data[i][j] = 0
		}
	}
}

func (mat Matrix3f) Det() float64 {
	res1 := mat.Data[0][0] * (mat.Data[1][1]*mat.Data[2][2] - mat.Data[1][2]*mat.Data[2][1])
	res2 := mat.Data[0][1] * (mat.Data[1][0]*mat.Data[2][2] - mat.Data[1][2]*mat.Data[2][0])
	res3 := mat.Data[0][2] * (mat.Data[1][0]*mat.Data[2][1] - mat.Data[1][1]*mat.Data[2][0])
	return res1 - res2 + res3
}

func (mat Matrix3f) IsRotationMatrix() bool {
	if math.Abs(float64(mat.Det()-1)) < 0.000001 {
		return true
	} else {
		return false
	}
}

func (mat Matrix3f) Mul(mulMat Matrix3f) Matrix3f {
	var res Matrix3f
	res.Clear()
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			for pos := 0; pos < 3; pos++ {
				res.Data[row][col] += mat.Data[row][pos] * mulMat.Data[pos][col]
			}
		}
	}
	return res
}

func (mat Matrix3f) ToEulerAngle() (Vector3, error) {
	if mat.IsRotationMatrix() == false {
		return Vector3_Invalid(), errors.New("This Matrix is not RotationMatrix")
	}
	var eulerAngle Vector3

	if mat.Data[0][2] < +1 {
		if mat.Data[0][2] > -1 {
			eulerAngle.Y = float32(math.Asin(mat.Data[0][2]))
			eulerAngle.X = -float32(math.Atan2(mat.Data[1][2], mat.Data[2][2]))
			eulerAngle.Z = float32(math.Atan2(mat.Data[0][1], mat.Data[0][0]))
		} else {
			eulerAngle.Y = -math.Pi / 2
			eulerAngle.X = float32(-math.Atan2(mat.Data[1][0], mat.Data[1][1]))
			eulerAngle.Z = 0
		}
	} else {
		eulerAngle.Y = math.Pi / 2
		eulerAngle.X = float32(math.Atan2(mat.Data[1][0], mat.Data[1][1]))
		eulerAngle.Z = 0
	}

	return eulerAngle, nil
}
