package imageToDWT


import (

)


func DWT1D(input []float64) ([]float64, []float64) {
	n := len(input)
	if n%2 != 0 {
		panic("Input length must be even")
	}

	outputLow := make([]float64, n/2)
	outputHigh := make([]float64, n/2)

	for i := 0; i < n/2; i++ {
		outputLow[i] = (input[2*i] + input[2*i+1]) / 2
		outputHigh[i] = (input[2*i] - input[2*i+1]) / 2
	}

	return outputLow, outputHigh
}



// DWT2D performs a 2D Discrete Wavelet Transform on the input matrix.
// It returns four sub-bands: LL, LH, HL, and HH.
func DWT2D(matrix [][]float64) ([][]float64,[][]float64,[][]float64,[][]float64) {


	//取得原圖高與寬
	height := len(matrix)
	width := len(matrix[0])

	rowLow := make([][]float64, height)
	rowHigh := make([][]float64, height)

	
	//將每一列拆成低頻與高頻 "第一次水平分割"
	for y := 0 ;y<height;y++{
		approx, detail := DWT1D(matrix[y])
		rowLow[y] = approx
		rowHigh[y] = detail
	}

	//因為每做一次 DWT，大小都縮小一半，所以輸出是 height/2 x width/2
	halfW := width/2
	halfH := height/2


	LL := make([][]float64, halfH)
	HL := make([][]float64, halfH)
	LH := make([][]float64, halfH)
	HH := make([][]float64, halfH)


	for y := 0; y < halfH; y++ {
		LL[y] = make([]float64, halfW)
		HL[y] = make([]float64, halfW)
		LH[y] = make([]float64, halfW)
		HH[y] = make([]float64, halfW)
	}


	for x := 0;x<halfW;x++{
		//針對每一欄做垂直DWT
		colLow := make([]float64, height)
		colHigh := make([]float64, height)

		// 將L 區中的每一列組成一整欄
		for y := 0; y < height; y++ {
			colLow[y] = rowLow[y][x]
			colHigh[y] = rowHigh[y][x]
		}


		approxColLow, detailColLow := DWT1D(colLow) // LL & LH
		approxColHigh, detailColHeigh := DWT1D(colHigh) // HL & HH
		for y:=0 ;y<halfH;y++{
			LL[y][x] = approxColLow[y]
			LH[y][x] = detailColLow[y]
			HL[y][x] = approxColHigh[y]
			HH[y][x] = detailColHeigh[y]
		}

	}
	return LL, LH, HL, HH
}