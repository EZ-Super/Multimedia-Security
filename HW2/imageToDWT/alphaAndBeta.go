package imageToDWT

import (
	"math"
)


var r [4]float32 = [4]float32{3.78,7,5.30,5.30}


func AlphaAndBeta(loacl int,D float64) (float64, float64) {
	// 計算 alpha 和 beta 的值
    alpha := 1 - (math.Pow(7.2-float64(r[loacl]), 2) / math.Pow(7.2, 2)) + float64(D)
    beta := 0.01 + (math.Pow(7.2-float64(r[loacl]), 2) / math.Pow(7.2, 2)) + (2 * float64(D))
	return alpha, beta
}