package imageToDWT


func EmbedImageToHH(HH, secret [][]float64, alpha float64) [][]float64 {
	h := len(HH)
	w := len(HH[0])
	result := make([][]float64, h)
	for y := 0; y < h; y++ {
		result[y] = make([]float64, w)
		for x := 0; x < w; x++ {
			result[y][x] = HH[y][x] + alpha*secret[y][x]
		}
	}
	return result
}
