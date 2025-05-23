package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"sort"

	"github.com/disintegration/imaging"
)

const (
	blockSize = 3
	mu        = 3.99 // 混沌參數 μ
)

var ClassCodeCount = map[uint8]int{0b00:0 , 0b01:0 , 0b10:0 , 0b11: 0}


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
func processBlock(block []uint8, bits []uint8) ([]uint8, int, uint8) {
	bestMSE := math.MaxFloat64 // 記錄目前找到的最佳 MSE 與對應區塊內容。
	best := make([]uint8, len(block))
	var usedBits int //這塊區塊最後用了幾個 bits
	var classCode uint8 //這塊對應的分類代碼（用於藏在中心像素）

	for bitCount := 2; bitCount <= 5; bitCount++ { //嘗試 2~5 bits 的嵌入方式
		tmp := make([]uint8, len(block))
		copy(tmp, block)
		bitIndex := 0
		for i := 0; i < len(block); i++ {
			if i == 4 {
				continue // 中央像素先略過
			}
			if bitIndex+bitCount > len(bits) { // 如果 bits 不夠用了，就停止嵌入
				break
			}
			data := uint8(0)
			for b := 0; b < bitCount; b++ { //將 bitCount 個 bits 組成一個 byte 值
				data = (data << 1) | bits[bitIndex] // 將 bits 中的位元逐個嵌入到 data 中
				bitIndex++
			}
			tmp[i] = embedLSB(tmp[i], data, bitCount)
		}
		mseBit := mse(block, tmp) // MSE 越小，代表修改越不明顯，視覺上越安全
		if mseBit < bestMSE { // 若這次比前一次更好，就更新「最佳解」
			bestMSE = mseBit
			copy(best, tmp)
			usedBits = bitCount
		}
	}
	classCode = map[int]uint8{2: 0b00, 3: 0b01, 4: 0b10, 5: 0b11}[usedBits] //選取 最佳bit 數量 對應的分類代碼
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
	copy(out.Pix, grayCover.Pix) //將主圖的像素值複製到輸出圖片

	blockCount := (width / blockSize) * (height / blockSize) //計算塊數

	order := logisticSequence(0.712, blockCount) // 生成混沌序列

	bitIndex := 0
	for _, idx := range order {
		if bitIndex >= len(secretBits) { // 若所有秘密位元都藏完，就不繼續處理後續區塊（提升效率）。
			break
		}
		bx := (idx % (width / blockSize)) * blockSize // 根據 Logistic 排序後的 idx，計算對應的 區塊左上角座標 (bx, by)。
		by := (idx / (width / blockSize)) * blockSize 

		block := make([]uint8, 9) //從主圖中擷取出 3×3 的像素資料到 block 陣列（共 9 pixels）。
		for y := 0; y < blockSize; y++ {
			for x := 0; x < blockSize; x++ {
				block[y*blockSize+x] = grayCover.GrayAt(bx+x, by+y).Y
			}
		}
		remainingBits := secretBits[bitIndex:] //把剩餘的秘密資料傳入 processBlock() 嘗試嵌入
		modified, usedBits, classCode := processBlock(block, remainingBits) // 嘗試嵌入
		ClassCodeCount[classCode]++
		bitIndex += (8 * (len(modified) - 1)) / usedBits // 更新 bitIndex 以指向下一個未嵌入的秘密位元
		// 藏入類別碼
		modified[4] = embedLSB(modified[4], classCode, 2) // 將類別碼嵌入到中心像素
		// 寫入回主圖
		for y := 0; y < blockSize; y++ {
			for x := 0; x < blockSize; x++ {
				out.SetGray(bx+x, by+y, color.Gray{Y: modified[y*blockSize+x]})
			}
		}
	}
	PSNR := computePSNR(grayCover, out)
	err = imaging.Save(out, "stego_output.png")
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
	fmt.Printf("總藏入量: %d\n", ClassCodeCount[0b00]*2 + ClassCodeCount[0b01]*3 + ClassCodeCount[0b10]*4 + ClassCodeCount[0b11]*5)
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