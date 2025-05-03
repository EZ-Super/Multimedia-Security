package watermark


import (
	"HW2/imageToDWT"
	//"fmt"
	//log "github.com/sirupsen/logrus"
)

// EmbedWaterMarkBand embeds a watermark band into a target band using the specified alpha and beta values.
//
// The formula used for embedding is as follows:
//
// alpha * host band + beta * watermark band = target band
func EmbedWaterMarkBand(
	hostDW3 imageToDWT.DWTResult, // Host image band
	wmLL [][]float64, // Watermark band
	waterMark [][]float64, // Watermark image band
	alpha,beta []float64, // Alpha and Beta values
){

	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			hostDW3.LL3[y][x] = alpha[1] * hostDW3.LL3[y][x] + beta[1] * wmLL[y][x]
		}
	}

}


