package watermark

import (
    "math"
)

// Flatten 矩陣轉成一維切片
func Flatten(matrix [][]float64) []float64 {
    flat := []float64{}
    for _, row := range matrix {
        flat = append(flat, row...)
    }
    return flat
}

// 機率分佈計算：以固定 bin 數為區間（例如 256 區間）
func HistogramProbs(data []float64, bins int) map[int]float64 {
    hist := make(map[int]int)
    total := len(data)

    for _, val := range data {
        bin := int(val + 256) // shift to positive
        if bin >= bins {
            bin = bins - 1
        } else if bin < 0 {
            bin = 0
        }
        hist[bin]++
    }

    probs := make(map[int]float64)
    for bin, count := range hist {
        probs[bin] = float64(count) / float64(total)
    }
    return probs
}

// 平均值（以機率分佈計算）
func Mean(data []float64, bins int) float64 {
    probs := HistogramProbs(data, bins)
    mean := 0.0
    for bin, p := range probs {
        mean += float64(bin-256) * p // shift back
    }
    return mean
}

// 變異數
func Variance(data []float64, bins int) float64 {
    probs := HistogramProbs(data, bins)
    mu := Mean(data, bins)
    var sigma2 float64
    for bin, p := range probs {
        x := float64(bin-256)
        sigma2 += math.Pow(x-mu, 2) * p
    }
    return sigma2
}

// Entropy 計算
func Entropy(data []float64, bins int) float64 {
    probs := HistogramProbs(data, bins)
    entropy := 0.0
    for _, p := range probs {
        if p > 0 {
            entropy -= p * math.Log2(p)
        }
    }
    return entropy
}


func DecideD(entropy float64) float64 {
    if entropy < 4.0 {
        return 0.01
    } else {
        return 0.02
    }
}