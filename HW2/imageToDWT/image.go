package imageToDWT

import (
	"fmt"
    "image/draw"
	"image/png"
	"image"
	"os"
	//"math"
)


func ImageToMatrix(path string) ([][]float64,error) {
	f,err := os.Open(path)

	if err != nil {
		fmt.Println("Error opening image:", err)
		return nil ,err
	}


	defer f.Close()
	img,err :=  png.Decode(f)

	if err != nil {
		fmt.Println("Error decoding image:", err)
		return nil ,err
	}



	bound := img.Bounds()
	w,h := bound.Max.X , bound.Max.Y
	mat := make([][]float64,h)

	for y:=0 ;y<h;y++{
		mat[y] = make([]float64, w)
		for x:=0 ;x<w;x++{
			r,g,b,_ := img.At(x,y).RGBA()
			//fmt.Printf("r: %d, g: %d, b: %d\n", r>>8, g>>8, b>>8)
			gray := float64((r+g+b)/3.0 >>8)
			mat[y][x] = gray
		}
	}
	return mat,nil
}	


// 將任意 image.Image 轉換成 image.Gray（灰階格式）
func ToGray(img image.Image) *image.Gray {
    bounds := img.Bounds()
    gray := image.NewGray(bounds)
    draw.Draw(gray, bounds, img, bounds.Min, draw.Src)
    return gray
}


// 輔助函數：保存灰階圖像文件（根據副檔名判斷格式）
func SaveImage(img *image.Gray, filename string) {
    file, err := os.Create(filename)
    if err != nil { panic(err) }
    defer file.Close()

    png.Encode(file, img)

}