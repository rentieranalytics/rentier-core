package utils

func IndexOf(target int, list []int) int {
	for i, v := range list {
		if v == target {
			return i
		}
	}
	return -1
}
