package main

import (
	"HW2/imageToDWT"
	"fmt"
	"os"

	"golang.org/x/image/bmp"

)


func main(){
	fmt.Println("Hello World")
	matrix,err := imageToDWT.ImageToMatrix("./resource/elaine_512x512.png")

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	

	LL,LH,HL,HH := imageToDWT.DWT2D(matrix)


	//fmt.Print("LL: ", LL)
	//fmt.Print("LH: ", LH)
	//fmt.Print("HL: ", HL)
	//fmt.Print("HH: ", HH)


	outImage := imageToDWT.CombinDWTBandsToImage(LL,LH,HL,HH)
	file ,_ := os.Create("dwt_result.png")
	defer file.Close()
	bmp.Encode(file,outImage)



	alpha := 0.1
	imageToDWT.EmbedImageToHH(HH)



	reconstructed := imageToDWT.IDWT2D(LL,LH,HL,HH)


	imageToDWT.SaveMatrixAsImage(reconstructed,"reconstructed.png")



	


}