package imageToDWT

import (
	
)




func AlignMatrix(mat [][]float64, offsetY, offsetX, targetH, targetW int) [][]float64 {
	res := make([][]float64, targetH)
	for y := 0; y < targetH; y++ {
		res[y] = make([]float64, targetW)
		for x := 0; x < targetW; x++ {
			if y+offsetY < len(mat) && x+offsetX < len(mat[0]) {
				res[y][x] = mat[y+offsetY][x+offsetX]
			} else {
				res[y][x] = 0 // padding
			}
		}
	}
	return res
}
