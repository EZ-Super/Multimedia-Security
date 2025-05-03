package imageToDWT

import (
	"image"
	"image/color"
)

// CropPatch 將輸入圖片中的指定區塊填為黑色（模擬裁切）
func CropPatch(img *image.Gray, x, y, cropW, cropH int) {
	bounds := img.Bounds()

	for j := 0; j < cropH; j++ {
		for i := 0; i < cropW; i++ {
			cx := x + i
			cy := y + j
			// 確保不會超出圖片邊界
			if cx >= bounds.Min.X && cx < bounds.Max.X && cy >= bounds.Min.Y && cy < bounds.Max.Y {
				img.SetGray(cx, cy, color.Gray{Y: 0}) // 裁切區設為黑色（值為 0）
			}
		}
	}
}
