package imageToDWT

import (
	"image"
	"image/color"
	"math"
)

// 輔助函數：將float64矩陣轉為灰階圖像
func MatrixToImage(matrix [][]float64) *image.Gray {
    h := len(matrix)
    w := len(matrix[0])
    img := image.NewGray(image.Rect(0, 0, w, h))
    for y := 0; y < h; y++ {
        for x := 0; x < w; x++ {
            // 取整
            val := int(math.Round(matrix[y][x]))
            if val < 0 { val = 0 }
            if val > 255 { val = 255 }

            img.SetGray(x, y, color.Gray{Y: uint8(val)})
        }
    }
    return img
}


func DWTtoImage(LL,LH,HL,HH [][]float64) *image.Gray{
    
    h := len(LL)
    w := len(LL[0])
    img := image.NewGray(image.Rect(0, 0, w*2, h*2))

    for y := 0; y < h; y++ {
        for x := 0; x < w; x++ {
            val := int(math.Round(LL[y][x]))
            if val < 0 { val = 0 }
            if val > 255 { val = 255 }
            img.SetGray(x, y, color.Gray{Y: uint8(val)})
            img.SetGray(x+w, y, color.Gray{Y: uint8(math.Round(LH[y][x]))})
            img.SetGray(x, y+h, color.Gray{Y: uint8(math.Round(HL[y][x]))})
            img.SetGray(x+w, y+h, color.Gray{Y: uint8(math.Round(HH[y][x]))})
        }
    }
    return img
}

func VisualizeFullDWT(dwt DWTResult) *image.Gray {
	canvasSize := 512
	out := image.NewGray(image.Rect(0, 0, canvasSize, canvasSize))
	// Level 3
	putMatrixToImage(dwt.LL3, out, 0, 0)                   // 左上
	putMatrixToImage(dwt.HL3, out, 64, 0)                 // 右上
	putMatrixToImage(dwt.LH3, out, 0, 64)                 // 左下
	putMatrixToImage(dwt.HH3, out, 64, 64)               // 右下
	// Level 2 (排在右半)
	putMatrixToImage(dwt.HL2, out, 128, 0)                 // 下方
	putMatrixToImage(dwt.LH2, out, 0, 128)
    putMatrixToImage(dwt.HH2, out, 128, 128) // 右下
	// Level 1 (排在更下方)
	putMatrixToImage(dwt.HL1, out,256, 0)
	putMatrixToImage(dwt.LH1, out, 0, 256)
    putMatrixToImage(dwt.HH1, out, 256, 256) // 右下

	return out
}
func putMatrixToImage(mat [][]float64, dst *image.Gray, offsetX, offsetY int) {
	h := len(mat)
	w := len(mat[0])
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			val := math.Round(mat[y][x])
			if val < 0 {
				val = 0
			} else if val > 255 {
				val = 255
			}
			dst.SetGray(offsetX+x, offsetY+y, color.Gray{Y: uint8(val)})
		}
	}
}


func ClampRound(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v + 0.5) // 正確做四捨五入
}