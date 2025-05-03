package imageToDWT

// 計算 NC 值，輸入為原始浮水印與提取浮水印
func CalcNC(original, extracted [][]float64) float64 {
	h := len(original)
	w := len(original[0])
	var numerator, denominator float64

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			numerator += original[y][x] * extracted[y][x]
			denominator += original[y][x] * original[y][x]
		}
	}

	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}
