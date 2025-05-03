package imageToDWT

import (
	"fmt"
	"image"
	"math"
	//log "github.com/sirupsen/logrus"
	"github.com/disintegration/imaging"
)



func RotatePatch(rotatedImage *image.Gray,originalTemp [][]int) [][]float64 {
	bestScore := -math.MaxFloat64
	var bestAngle float64
	var score float64

	for Theta:= -5.0 ;Theta<=5.0;Theta+=1{
		score = 0.0
		rotated := imaging.Rotate(rotatedImage,Theta,image.Transparent)
		rotated = imaging.CropCenter(rotated, 512, 512)
		//imaging.Save(rotated,fmt.Sprintf("./temp/rotated %f.png",Theta))
		rotatedMatrix := NRGBAtoGrayMatrix(rotated)

		rotateDWT3 := DWT3Level(rotatedMatrix)
		template := rotateDWT3.HL3
		bitMat := ExtractLSBMatrix(template)
		//bitMat = AlignMatrix(bitMat,0,0,512,512)
		
		score,_,_ = MatchTemplateScore(bitMat, originalTemp)
		fmt.Print("theta: ",Theta, "  score: ",score, "\n")
		if score > bestScore {
			bestScore = score
			bestAngle = Theta
		}

	} 

	fmt.Print("bestAngle: ", bestAngle, "  bestScore: ", bestScore, "\n")
	bestRotated := imaging.Rotate(rotatedImage,bestAngle,image.Transparent)
	bestRotated = imaging.CropCenter(bestRotated, 512, 512)
	imaging.Save(bestRotated,"./temp/bestRotated.png")
	bestRotatedMatrix := NRGBAtoGrayMatrix(bestRotated)
	return bestRotatedMatrix


}

func NRGBAtoGrayMatrix(img *image.NRGBA) [][]float64 {
    bounds := img.Bounds()
    width := bounds.Dx()
    height := bounds.Dy()

    matrix := make([][]float64, height)
    for y := 0; y < height; y++ {
        row := make([]float64, width)
        for x := 0; x < width; x++ {
            offset := y*img.Stride + x*4
            r := float64(img.Pix[offset+0])

            gray := r
            row[x] = gray
        }
        matrix[y] = row
    }
    return matrix
}

// 比對 HL3 的 LSB bit matrix 與 template，同步碼要鋪滿 LL3 區域
func MatchTemplateScore(bitMat [][]int, templateBlock [][]int) (float64, int, int) {

    h := 64
    w := 64
	fmt.Println(h," ",w)
    blockSize := len(templateBlock)




    // 1. 平鋪 template → fullTemplate
    fullTemplate := make([][]int, h)
    for j := range fullTemplate {
        fullTemplate[j] = make([]int, w)
        for i := range fullTemplate[j] {
            fullTemplate[j][i] = templateBlock[j%blockSize][i%blockSize]
        }
    }

    // 2. 搜尋最大對齊位置（滑動視窗）
    maxOffset := 100 // 支援裁切 ±32px
    bestScore := -math.MaxFloat64
    bestY, bestX := 0, 0

    for offY := 0; offY <= maxOffset; offY++ {
        for offX := 0; offX <= maxOffset; offX++ {
            score := 0.0
            for j := 0; j < h && j+offY < len(bitMat); j++ {
                for i := 0; i < w && i+offX < len(bitMat[0]); i++ {
                    t := fullTemplate[j][i]
                    b := bitMat[j+offY][i+offX]
                    if t == b {
                        score += 1.0
                    } else {
                        score -= 1.0
                    }
                }
            }
            if score > bestScore {
                bestScore = score
                bestY = offY
                bestX = offX
            }
        }
    }
    return bestScore, bestY, bestX
}



