package watermark


import (
	"HW2/imageToDWT"

)


func ExtractWaterMarkBand(
	DWTImage imageToDWT.DWTResult, // Host image band
	original imageToDWT.DWTResult, // Original image band
	alpha,beta []float64, // Alpha and Beta values
)([][]float64){


	wmLL := make([][]float64,64)

	for y:=0;y<64;y++{
		wmLL[y] = make([]float64,64)

		for x:=0;x<64;x++{
			wmLL[y][x] = (DWTImage.LL3[y][x] - alpha[1]*original.LL3[y][x]) / beta[1]

			

		}
	}





	return wmLL
}


func Denormalization(data [][]float64){
	for y := range data {
		for x := range data[y] {
			v := data[y][x]
			data[y][x] = v*128 + 128
		}
	}
}


func AverageMatrix(a, b [][]float64) [][]float64 {
    h := len(a)
    w := len(a[0])
    result := make([][]float64, h)
    for y := 0; y < h; y++ {
        result[y] = make([]float64, w)
        for x := 0; x < w; x++ {
            result[y][x] = (a[y][x] + b[y][x]) / 2
        }
    }
    return result
}
