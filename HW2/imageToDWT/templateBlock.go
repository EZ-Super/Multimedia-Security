package imageToDWT

import (
	"math/rand"
	"math"
)

// 產生 n × n 的偽隨機同步模板（元素為 0 或 1）
func GenerateTemplateBlock(n int, seed int64) [][]int {
    rand.Seed(seed) // 讓嵌入與提取都能用同樣的 seed 產生一致模板
    block := make([][]int, n)
    for y := 0; y < n; y++ {
        block[y] = make([]int, n)
        for x := 0; x < n; x++ {
            block[y][x] = rand.Intn(2) // 0 或 1
        }
    }
    return block
}


func EmbedTemplateLSB(HL3 [][]float64, template [][]int) {
    h := len(HL3)
    w := len(HL3[0])
    n := len(template) // 假設 template 是 n × n
    for y := 0; y < h; y += n {
        for x := 0; x < w; x += n {
            for dy := 0; dy < n && y+dy < h; dy++ {
                for dx := 0; dx < n && x+dx < w; dx++ {
                    bit := template[dy][dx]
                    HL3[y+dy][x+dx] = SetIntLSB(HL3[y+dy][x+dx], bit)
                }
            }
        }
    }
}

func ExtractLSBMatrix(HL3 [][]float64) [][]int {
	h := len(HL3)
	w := len(HL3[0])
	bitMat := make([][]int, h)
	for y := 0; y < h; y++ {
		bitMat[y] = make([]int, w)
		for x := 0; x < w; x++ {
			bitMat[y][x] = GetIntLSB(HL3[y][x]) // return 0 or 1
		}
	}
	return bitMat
}


// 將 float64 的 LSB 設為 bit（0 或 1）
func SetIntLSB(value float64, bit int) float64 {
    // 轉為整數
    intVal := int(math.Round(value))
    // 清除原本的 LSB
    intVal &= ^1
    // 設定新的 LSB
    intVal |= (bit & 1)
    return float64(intVal)
}
// 從 float64 取得 LSB 位元（回傳 0 或 1）
func GetIntLSB(value float64) int {
    intVal := int(math.Round(value))
    return intVal & 1
}
