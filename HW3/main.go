package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"

	"sort"
	"github.com/disintegration/imaging"
)

const (
	blockSize = 3
	mu        = 3.99 // 混沌參數 μ
)

type BlockPosition struct {
	x int
	y int
}

type PixelState struct {
	value uint8
	valid bool
}

var ClassCodeCount = map[uint8]int{0b00: 0, 0b01: 0, 0b10: 0, 0b11: 0}




// Logistic Map 混沌排序產生器
func logisticSequence(x0 float64, total int) []int {

	seq := make([]float64, total) //創建一個長度為total (512*512 為 29241 )的float64切片
	seq[0] = x0
	for i := 1; i < total; i++ {
		seq[i] = mu * seq[i-1] * (1 - seq[i-1])
	}

	indexes := make([]int, total)
	for i := range indexes {
		indexes[i] = i
	}
	
		sort.Slice(indexes, func(i, j int) bool {
			return seq[indexes[i]] < seq[indexes[j]] // 根據混沌數值大小進行排序
		})
	return indexes
}

// 灰階像素陣列轉 bit stream
func imageToBitStream(img *image.Gray) []uint8 {
	bits := []uint8{}
	for _, px := range img.Pix {
		for i := 7; i >= 0; i-- {
			bits = append(bits, (px>>i)&1)
		}
	}
	return bits
}

// LSB 嵌入函數
func embedLSB(value uint8, data uint8, bitCount int) uint8 {
	mask := ^uint8((1 << bitCount) - 1)
	return (value & mask) | (data & ((1 << bitCount) - 1))
}

// MSE 計算
func mse(original, modified []uint8) float64 {
	var sum float64
	for i := range original {
		diff := float64(original[i]) - float64(modified[i])
		sum += diff * diff
	}
	return sum / float64(len(original))
}

// 嘗試不同藏入量，選擇最佳方案
func processBlock(block []PixelState, bits []uint8) ([]PixelState, int, uint8) {
	bestMSE := math.MaxFloat64
	best := make([]PixelState, len(block))
	var usedBits int
	var classCode uint8

	for bitCount := 2; bitCount <= 5; bitCount++ {
		tmp := make([]PixelState, len(block))
		copy(tmp, block)
		bitIndex := 0
		for i := 0; i < len(block); i++ {
			if i == 6 {
				continue
			}
			if bitIndex+bitCount > len(bits) {
				break
			}

			if !block[i].valid {
				continue
			}

			data := uint8(0)
			for b := 0; b < bitCount; b++ {
				data = (data << 1) | bits[bitIndex]
				bitIndex++
			}
			tmp[i].value = embedLSB(tmp[i].value, data, bitCount)
		}

		// Calculate MSE only for valid pixels
		var sum float64
		var count int
		for i := range block {
			if block[i].valid && tmp[i].valid {
				diff := float64(block[i].value) - float64(tmp[i].value)
				sum += diff * diff
				count++
			}
		}
		mseBit := sum / float64(count)

		if mseBit < bestMSE {
			bestMSE = mseBit
			copy(best, tmp)
			usedBits = bitCount
		}
	}
	classCode = map[int]uint8{2: 0b00, 3: 0b01, 4: 0b10, 5: 0b11}[usedBits]
	return best, usedBits, classCode
}

func toGray(img image.Image) *image.Gray {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray.Set(x, y, img.At(x, y))
		}
	}
	return gray
}

func padToMultipleOf3(img *image.Gray) *image.Gray {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	newW := ((w + 2) / 3) * 3
	newH := ((h + 2) / 3) * 3
	if newW == w && newH == h {
		return img
	}

	padded := image.NewGray(image.Rect(0, 0, newW, newH))
	for y := 0; y < h; y++ {
		copy(padded.Pix[y*newW:y*newW+w], img.Pix[y*w:(y+1)*w])
	}
	return padded
}

func main() {
	file, err := os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal("無法開啟日誌檔案:", err)
	}
	defer file.Close()

	log.SetOutput(file)

	logger := log.New(os.Stdout, "[HW3]", log.Ldate|log.Ltime|log.Lshortfile)
	log.Println("開始執行")

	coverImg, err := imaging.Open("image/cameraman_512x512.bmp")
	if err != nil {
		log.Fatal("無法讀取主圖片:", err)
	}
	secretImg, err := imaging.Open("image/elaine_512x512.bmp")
	if err != nil {
		log.Fatal("無法讀取秘密圖片:", err)
	}

	grayCover := toGray(imaging.Grayscale(coverImg))
	graySecret := toGray(imaging.Grayscale(secretImg))

	width := grayCover.Bounds().Dx()
	height := grayCover.Bounds().Dy()

	if width%3 != 0 || height%3 != 0 {
		grayCover = padToMultipleOf3(grayCover) // 擴展為3的倍數
		width = grayCover.Bounds().Dx()
		height = grayCover.Bounds().Dy()
	}

	secretWidth := graySecret.Bounds().Dx()
	secretHeight := graySecret.Bounds().Dy()

	if secretWidth%3 != 0 || secretHeight%3 != 0 {
		graySecret = padToMultipleOf3(graySecret) // 擴展為3的倍數
	}

	secretBits := imageToBitStream(graySecret)
	fmt.Printf("📦 準備嵌入 %d bits 資料\n", len(secretBits))

	out := image.NewGray(grayCover.Bounds()) //創建輸出圖片 大小為512*512
	copy(out.Pix, grayCover.Pix)             //將主圖的像素值複製到輸出圖片

	blockCount := (width / blockSize) * (height / blockSize) //計算塊數

	order := logisticSequence(0.712, blockCount) // 生成混沌序列

	bitIndex := 0
	for _, idx := range order {
		if bitIndex >= len(secretBits) { // 若所有秘密位元都藏完，就不繼續處理後續區塊（提升效率）。
			break
		}
		bx := (idx % (width / blockSize)) * blockSize // 根據 Logistic 排序後的 idx，計算對應的 區塊左上角座標 (bx, by)。
		by := (idx / (width / blockSize)) * blockSize

		totalFixBlock := 9

		top := BlockPosition{x: bx + 1, y: by}
		bottom := BlockPosition{x: bx + 1, y: by + 2}
		left := BlockPosition{x: bx, y: by + 1}
		right := BlockPosition{x: bx + 2, y: by + 1}

		topPixel := grayCover.GrayAt(top.x, top.y).Y
		bottomPixel := grayCover.GrayAt(bottom.x, bottom.y).Y
		leftPixel := grayCover.GrayAt(left.x, left.y).Y
		rightPixel := grayCover.GrayAt(right.x, right.y).Y

		compareTop := grayCover.GrayAt(top.x, top.y-1).Y
		compareBottom := grayCover.GrayAt(bottom.x, bottom.y+1).Y
		compareLeft := grayCover.GrayAt(left.x-1, left.y).Y
		compareRight := grayCover.GrayAt(right.x+1, right.y).Y

		block := make([]PixelState, 13)
		// Initialize all pixels as valid
		for i := range block {
			block[i] = PixelState{valid: true}
		}

		if topPixel > compareTop {

			if top.y > 0 {
				block[0] = PixelState{value: compareTop, valid: true}
				totalFixBlock++
			} else {
				block[0].valid = false
			}
		} else if compareTop > topPixel {
			block[2] = PixelState{valid: false}
			block[0] = PixelState{valid: false}
			totalFixBlock--
		} else if compareTop == topPixel {
			block[0] = PixelState{valid: false}
		} 

		if bottomPixel > compareBottom {
			if bottom.y < height {
				block[12] = PixelState{value: compareBottom, valid: true}
				totalFixBlock++
			} else {
				block[12] = PixelState{valid: false}
			}
		} else if compareBottom > bottomPixel {

			block[10] = PixelState{valid: false}
			block[12] = PixelState{valid: false}
			totalFixBlock--
		} else if compareBottom == bottomPixel {
			block[12] = PixelState{valid: false}
		}

		if leftPixel > compareLeft {
			if left.x > 0 {
				block[4] = PixelState{value: compareLeft, valid: true}
				totalFixBlock++
			} else {
				block[4] = PixelState{valid: false}
			}
		} else if compareLeft > leftPixel {
			block[4] = PixelState{valid: false}
			block[5] = PixelState{valid: false}
			totalFixBlock--
		} else if compareLeft == leftPixel {
			block[4] = PixelState{valid: false}
		}

		if rightPixel > compareRight {
			if right.x < width {

				block[8] = PixelState{value: compareRight, valid: true}
				totalFixBlock++
			} else {
				block[8] = PixelState{valid: false}
			}
		} else if compareRight > rightPixel {

			block[7] = PixelState{valid: false}
			block[8] = PixelState{valid: false}
			totalFixBlock--
		} else if compareRight == rightPixel {
			block[8] = PixelState{valid: false}
		}

		for i := 0; i < blockSize; i++ {
			if block[i+1].valid {
				block[i+1] = PixelState{value: grayCover.GrayAt(bx+i, by).Y, valid: true}
			}
		}
		for i := 0; i < blockSize; i++ {
			if block[i+5].valid {
				block[i+5] = PixelState{value: grayCover.GrayAt(bx+i, by+1).Y, valid: true}
			}
		}
		for i := 0; i < blockSize; i++ {
			if block[i+9].valid {
				block[i+9] = PixelState{value: grayCover.GrayAt(bx+i, by+2).Y, valid: true}
			}
		}


		log.Printf("bx,by: %d,%d", bx, by)
		log.Printf("top,bottom,left,right: %d,%d,%d,%d", top.x, bottom.x, left.x, right.x)
		log.Printf("top.y,bottom.y,left.y,right.y: %d,%d,%d,%d", top.y, bottom.y, left.y, right.y)
		log.Printf("topPixel,bottomPixel,leftPixel,rightPixel: %d,%d,%d,%d", topPixel, bottomPixel, leftPixel, rightPixel)
		log.Printf("compareTop,compareBottom,compareLeft,compareRight: %d,%d,%d,%d", compareTop, compareBottom, compareLeft, compareRight)

		log.Printf(
			"original pixel : [%d,%d,%d,%d,%d,%d,%d,%d,%d]",
			grayCover.GrayAt(bx, by).Y,
			grayCover.GrayAt(bx+1, by).Y,
			grayCover.GrayAt(bx+2, by).Y,
			grayCover.GrayAt(bx, by+1).Y,
			grayCover.GrayAt(bx+1, by+1).Y,
			grayCover.GrayAt(bx+2, by+1).Y,
			grayCover.GrayAt(bx, by+2).Y,
			grayCover.GrayAt(bx+1, by+2).Y,
			grayCover.GrayAt(bx+2, by+2).Y,
		)
		log.Printf("out pixel : [%d,%d,%d,%d,%d,%d,%d,%d,%d]", block[1].value, block[2].value, block[3].value, block[5].value, block[6].value, block[7].value, block[9].value, block[10].value, block[11].value)

		log.Printf("block: %v ", block)

		remainingBits := secretBits[bitIndex:]                              //把剩餘的秘密資料傳入 processBlock() 嘗試嵌入
		modified, usedBits, classCode := processBlock(block, remainingBits) // 嘗試嵌入
		log.Printf("modified: %v", modified)
		ClassCodeCount[classCode]++
		bitIndex += totalFixBlock * usedBits // 更新 bitIndex 以指向下一個未嵌入的秘密位元

		// 藏入類別碼
		modified[6] = PixelState{
			value: embedLSB(modified[6].value, classCode, 2),
			valid: true,
		}

		// 寫入回主圖

		if modified[0].valid {
			out.SetGray(bx+1, by-1, color.Gray{Y: modified[0].value}) //修改 top
		}
		if modified[4].valid {
			out.SetGray(bx-1, by+1, color.Gray{Y: modified[4].value}) //修改 left
		}
		if modified[8].valid {
			out.SetGray(bx+3, by+1, color.Gray{Y: modified[8].value}) //修改 right
		}
		if modified[12].valid {
			out.SetGray(bx+1, by+3, color.Gray{Y: modified[12].value}) //修改 bottom
		}

		for i := 0; i < blockSize; i++ {
			if modified[i+1].valid {
				out.SetGray(bx+i, by, color.Gray{Y: modified[i+1].value})
			}
		}
		for i := 0; i < blockSize; i++ {
			if modified[i+5].valid {
				out.SetGray(bx+i, by+1, color.Gray{Y: modified[i+5].value})
			}
		}
		for i := 0; i < blockSize; i++ {
			if modified[i+9].valid {
				out.SetGray(bx+i, by+2, color.Gray{Y: modified[i+9].value})
			}
		}
	}


	PSNR := computePSNR(grayCover, out)
	cropImage := imaging.CropCenter(out, 512, 512)

	err = imaging.Save(cropImage, "stego_output.png")

	if err != nil {
		log.Fatal("無法儲存結果圖片:", err)
	}
	fmt.Println("✅ 藏圖完成，結果為 stego_output.png")
	fmt.Printf("PSNR: %f\n", PSNR)

	fmt.Println("ClassCodeCount:", ClassCodeCount)
	fmt.Println("ClassCodeCount[0b00]:", ClassCodeCount[0b00])
	fmt.Println("ClassCodeCount[0b01]:", ClassCodeCount[0b01])
	fmt.Println("ClassCodeCount[0b10]:", ClassCodeCount[0b10])
	fmt.Println("ClassCodeCount[0b11]:", ClassCodeCount[0b11])

	fmt.Printf("總藏入量: %d\n", ClassCodeCount[0b00]*2+ClassCodeCount[0b01]*3+ClassCodeCount[0b10]*4+ClassCodeCount[0b11]*5)

	logger.Println("執行完成")

}

// 計算 MSE（Mean Squared Error）
func computeMSE(img1, img2 *image.Gray) float64 {
	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()

	if bounds1.Dx() != bounds2.Dx() || bounds1.Dy() != bounds2.Dy() {
		panic("圖片尺寸不一致，無法計算 MSE")
	}

	var sum float64
	for y := 0; y < bounds1.Dy(); y++ {
		for x := 0; x < bounds1.Dx(); x++ {
			v1 := float64(img1.GrayAt(x, y).Y)
			v2 := float64(img2.GrayAt(x, y).Y)
			diff := v1 - v2
			sum += diff * diff
		}
	}
	total := float64(bounds1.Dx() * bounds1.Dy())
	return sum / total
}

// 計算 PSNR（Peak Signal-to-Noise Ratio）
func computePSNR(img1, img2 *image.Gray) float64 {
	mse := computeMSE(img1, img2)
	if mse == 0 {
		return math.Inf(1) // 完全一樣，PSNR 無限大
	}
	return 10 * math.Log10((255*255)/mse)
}




// 計算 MSE（Mean Squared Error）
func computeMSE(img1, img2 *image.Gray) float64 {
	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()

	if bounds1.Dx() != bounds2.Dx() || bounds1.Dy() != bounds2.Dy() {
		panic("圖片尺寸不一致，無法計算 MSE")
	}

	var sum float64
	for y := 0; y < bounds1.Dy(); y++ {
		for x := 0; x < bounds1.Dx(); x++ {
			v1 := float64(img1.GrayAt(x, y).Y)
			v2 := float64(img2.GrayAt(x, y).Y)
			diff := v1 - v2
			sum += diff * diff
		}
	}
	total := float64(bounds1.Dx() * bounds1.Dy())
	return sum / total
}

// 計算 PSNR（Peak Signal-to-Noise Ratio）
func computePSNR(img1, img2 *image.Gray) float64 {
	mse := computeMSE(img1, img2)
	if mse == 0 {
		return math.Inf(1) // 完全一樣，PSNR 無限大
	}
	return 10 * math.Log10((255 * 255) / mse)
}