package imageToDWT


func ExtractImageFromHH(HHWithSecret, originalHH [][]float64, alpha float64) [][]float64 {
	h := len(HHWithSecret)
	w := len(HHWithSecret[0])
	secret := make([][]float64, h)
	for y := 0; y < h; y++ {
		secret[y] = make([]float64, w)
		for x := 0; x < w; x++ {
			secret[y][x] = (HHWithSecret[y][x] - originalHH[y][x]) / alpha
		}
	}
	return secret
}
