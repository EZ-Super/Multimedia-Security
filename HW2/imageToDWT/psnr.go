package imageToDWT


import (
	"math"
)

// 計算 MSE（均方誤差）
func MSE(img1, img2 [][]float64) float64 {
	h := len(img1)
	w := len(img1[0])
	var sum float64

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			diff := img1[y][x] - img2[y][x]
			sum += diff * diff
		}
	}

	return sum / float64(h*w)
}

// 計算 PSNR
func PSNR(img1, img2 [][]float64) float64 {
	mse := MSE(img1, img2)
	if mse == 0 {
		return math.Inf(1) // 完全相同
	}
	maxPixel := 255.0
	return 10 * math.Log10((maxPixel * maxPixel) / mse)
}
