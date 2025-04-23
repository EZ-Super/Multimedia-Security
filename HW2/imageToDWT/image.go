package imageToDWT

import (
	"fmt"

	"image/png"
	"image"
	"image/color"
	"os"
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


func SaveMatrixAsImage(matrix [][]float64, path string) error {
	h := len(matrix)
	w := len(matrix[0])
	img := image.NewGray(image.Rect(0, 0, w, h))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			val := matrix[y][x]
			if val < 0 {
				val = 0
			}
			if val > 255 {
				val = 255
			}
			img.SetGray(x, y, color.Gray{Y: uint8(val)})
		}
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	return png.Encode(out, img)
}