package imageToDWT


import (
	"math"
	"image"
	"image/color")

// 旋轉圖片（以中心旋轉指定角度，最近鄰補值）
func RotateImage(img *image.Gray, degree float64) *image.Gray {
    bounds := img.Bounds()
    w := bounds.Max.X
    h := bounds.Max.Y
    centerX := float64(w) / 2
    centerY := float64(h) / 2

    angle := degree * math.Pi / 180.0 // 角度轉弧度

    rotated := image.NewGray(image.Rect(0, 0, w, h))
    for y := 0; y < h; y++ {
        for x := 0; x < w; x++ {
            // 反向旋轉映射
            dx := float64(x) - centerX
            dy := float64(y) - centerY
            srcX := centerX + dx*math.Cos(-angle) - dy*math.Sin(-angle)
            srcY := centerY + dx*math.Sin(-angle) + dy*math.Cos(-angle)

            ix := int(math.Round(srcX))
            iy := int(math.Round(srcY))

            if ix >= 0 && ix < w && iy >= 0 && iy < h {
                rotated.SetGray(x, y, img.GrayAt(ix, iy))
            } else {
                rotated.SetGray(x, y, color.Gray{Y: 0}) // 超出範圍填黑
            }
        }
    }
    return rotated
}
