package main

import (
	"HW2/imageToDWT"
	"HW2/waterMark"
	"fmt"
	"github.com/disintegration/imaging"
	"os"
	"image"
	log "github.com/sirupsen/logrus"
)

func main(){

    file, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        log.Fatalf("無法打開日誌文件: %v", err)
    }
    defer file.Close()

	log.SetOutput(file)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.InfoLevel)
	log.Info("Start program")
	matrix,err := imageToDWT.ImageToMatrix("./resource/elaine_512x512.png")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	

	
	waterMark ,err := imageToDWT.ImageToMatrix("./resource/picture 64.png")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	
	hostDW3 := imageToDWT.DWT3Level(matrix) 
	originalDW3 := hostDW3.Clone()


	var alpha ,beta []float64 = make([]float64,0),make([]float64,0)


	for i:=0 ;i<4;i++{
		D := 0.0
		flat := make([]float64,0)
		if i == 0{
			flat = watermark.Flatten(hostDW3.LH3)
		}else if i == 1{
			flat = watermark.Flatten(hostDW3.HL3)
		}else if i == 2{
			flat = watermark.Flatten(hostDW3.LH2)
		}else if i == 3{
			flat = watermark.Flatten(hostDW3.HL2)
		}
		enctropy := watermark.Entropy(flat,512)
		D = watermark.DecideD(enctropy)
		alphaVal ,betaVal := imageToDWT.AlphaAndBeta(i,D)
		alpha = append(alpha, alphaVal)
		beta = append(beta, betaVal)

	}


	templateBlock := imageToDWT.GenerateTemplateBlock(16,1234)
	imageToDWT.EmbedTemplateLSB(hostDW3.HL3,templateBlock)
	watermark.EmbedWaterMarkBand(hostDW3,waterMark,waterMark,alpha,beta)



	result := imageToDWT.IDWT3D(
		hostDW3.LL3, hostDW3.HL3, hostDW3.LH3, 
		hostDW3.HH3,hostDW3.HL2, hostDW3.LH2, hostDW3.HH2,
		hostDW3.HL1, hostDW3.LH1, hostDW3.HH1,
	)

	log.Info(fmt.Sprintf("psnr: %f \n" , imageToDWT.PSNR(matrix,result)))

	embedImg := imageToDWT.MatrixToImage(result)
	imageToDWT.SaveImage(embedImg,"./temp/elaine_512x512_watermark.png")


	
	embedDW3 := imageToDWT.DWT3Level(result)
	wmHL :=watermark.ExtractWaterMarkBand(embedDW3,*originalDW3,alpha,beta)



	wmHLImg := imageToDWT.MatrixToImage(wmHL)
	imageToDWT.SaveImage(wmHLImg,"./result/elaine_512x512_watermark_wmHL.png")



	rotate,err := imaging.Open("./temp/elaine_512x512_watermark.png")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	
	rotate = imaging.Rotate(rotate, -3, image.Transparent)
	imaging.Save(rotate,"./temp/rotated_elaine_512x512_watermark.png")
	rotate = imageToDWT.ToGray(rotate)
	
	patchImage := imageToDWT.RotatePatch(rotate.(*image.Gray),templateBlock)
	rotateDW3 := imageToDWT.DWT3Level(patchImage)

	rwm := watermark.ExtractWaterMarkBand(rotateDW3,*originalDW3,alpha,beta)

	log.Info(fmt.Sprintf("Rotated NC : %f",imageToDWT.CalcNC(waterMark,rwm)))

	rotateHL := imageToDWT.MatrixToImage(rwm)
	imageToDWT.SaveImage(rotateHL,"./result/elaine_512x512_watermark_rotate_HL.png")


	imageToDWT.CropPatch(embedImg,64,64,32,32);
	imageToDWT.SaveImage(embedImg,"./temp/elaine_512x512_watermark_crop.png")

	cropMatrix,_ := imageToDWT.ImageToMatrix("./temp/elaine_512x512_watermark_crop.png")
	cropDW3 := imageToDWT.DWT3Level(cropMatrix)
	cropWM := watermark.ExtractWaterMarkBand(cropDW3,*originalDW3,alpha,beta)
	cropWMImg := imageToDWT.MatrixToImage(cropWM)
	imageToDWT.SaveImage(cropWMImg,"./result/elaine_512x512_watermark_crop_wmHL.png")

	log.Info(fmt.Sprintf("crop NC : %f",imageToDWT.CalcNC(waterMark,cropWM)))


}


