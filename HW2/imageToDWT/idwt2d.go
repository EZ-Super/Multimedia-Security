package imageToDWT

func IDWT1D(approx,detail []float64) []float64 {
	n := len(approx)

	result := make ([]float64,n*2)

	for i:=0 ;i<n;i++{
		result[2*i] = approx[i] + detail[i]
		result[2*i+1] = approx[i] - detail[i]
	}
	return result
}

func IDWT2D(LL, LH, HL, HH [][]float64) [][]float64 {
	halfH := len(LL)
	halfW := len(LL[0])
	height := halfH * 2
	width := halfW * 2

	// Step 1: 垂直方向重建 → 取得中間的 rowLow, rowHigh
	rowLow := make([][]float64, height)
	rowHigh := make([][]float64, height)

	for y := 0; y < halfH; y++ {
		rowLow[2*y] = make([]float64, halfW)
		rowLow[2*y+1] = make([]float64, halfW)
		rowHigh[2*y] = make([]float64, halfW)
		rowHigh[2*y+1] = make([]float64, halfW)

		for x := 0; x < halfW; x++ {
			// 垂直方向的 IDWT：每一欄中兩個值合成一列
			rowLow[2*y][x]   = LL[y][x] + LH[y][x]
			rowLow[2*y+1][x] = LL[y][x] - LH[y][x]
			rowHigh[2*y][x]   = HL[y][x] + HH[y][x]
			rowHigh[2*y+1][x] = HL[y][x] - HH[y][x]
		}
	}

	// Step 2: 水平方向重建 → 將 rowLow + rowHigh 合併還原成原始矩陣
	result := make([][]float64, height)
	for y := 0; y < height; y++ {
		result[y] = make([]float64, width)
		for x := 0; x < halfW; x++ {
			result[y][2*x]   = rowLow[y][x] + rowHigh[y][x]   // even
			result[y][2*x+1] = rowLow[y][x] - rowHigh[y][x]   // odd
		}
	}

	return result
}