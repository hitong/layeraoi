package linemath

func NewVector2() *Vector2 {
	return &Vector2{
		X: 0,
		Y: 0,
	}
}

func NewVector3() *Vector3 {
	return &Vector3{
		X: 0,
		Y: 0,
		Z: 0,
	}
}

func NewQuaternion() *Quaternion {
	return &Quaternion{
		X: 0,
		Y: 0,
		Z: 0,
		W: 0,
	}
}
