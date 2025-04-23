package imageToDWT

import (
	"image"
	"image/color"
)

func CombinDWTBandsToImage(LL, LH , HL ,HH [][]float64) image.Image{
	h := len(LL)
	w := len(LL[0])
	out := image.NewGray(image.Rect(0,0,w*2,h*2))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// 將 LL 區域的值轉換為灰階顏色
			grayValue := uint8(LL[y][x])
			out.SetGray(x, y, color.Gray{Y: grayValue})
			// 將 LH 區域的值轉換為灰階顏色
			grayValue = uint8(LH[y][x])
			out.SetGray(x+w, y, color.Gray{Y: grayValue})
			// 將 HL 區域的值轉換為灰階顏色
			grayValue = uint8(HL[y][x])
			out.SetGray(x, y+h, color.Gray{Y: grayValue})
			// 將 HH 區域的值轉換為灰階顏色
			grayValue = uint8(HH[y][x])
			out.SetGray(x+w, y+h, color.Gray{Y: grayValue})
		}
	}
	return out
}

