package collider

type Collider interface {
	Rect() (float64, float64, float64, float64)
}

func CheckPointInsideRect(sx, sy float64, targetObj Collider) bool {
	tx, ty, tw, th := targetObj.Rect()
	inSide := tx < sx && sx < tx+tw && ty < sy && sy < ty+th
	return inSide
}
