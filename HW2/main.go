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
	pictureMatrix ,err := imageToDWT.ImageToMatrix("./resource/picture.png")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	LL,LH,HL,HH := imageToDWT.DWT2D(matrix)

	PLL,pLH,pHL,pHH := imageToDWT.DWT2D(pictureMatrix)
	//fmt.Print("LL: ", LL)
	//fmt.Print("LH: ", LH)
	//fmt.Print("HL: ", HL)
	//fmt.Print("HH: ", HH)


	outImage := imageToDWT.CombinDWTBandsToImage(LL,LH,HL,HH)
	file ,_ := os.Create("dwt_result.png")
	defer file.Close()
	bmp.Encode(file,outImage)


	pout := imageToDWT.CombinDWTBandsToImage(PLL,pLH,pHL,pHH)
	file,_ = os.Create("dwt_picture_result.png")
	defer file.Close()
	bmp.Encode(file,pout)




	reconstructed := imageToDWT.IDWT2D(LL,LH,HL,HH)


	imageToDWT.SaveMatrixAsImage(reconstructed,"reconstructed.png")



	


}