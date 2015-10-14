package taskbytime

func Min(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
}

func Min64(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int32) int32 {
	if x > y {
		return x
	}
	return y
}

func Max64(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func Abs(x int32) int32 {
	if x < 0 {
		return -x
	}

	return x
}

func Abs64(x int64) int64 {
	if x < 0 {
		return -x
	}

	return x
}
