package main

import (
	"image"
	"image/draw"
	"log"
	"math"

	xdraw "golang.org/x/image/draw"
)

type CombinedImages struct {
	striped, compressed image.Image
}

func getCombined(leftImg, rightImg image.Image, angle float64) *CombinedImages {
	if angle > myConfiguration.MAX_ANGLE || angle < 0 {

		log.Fatalf("der Parameter angle muss 0>= angle <= {%.2f} sein: {%.2f}", myConfiguration.MAX_ANGLE, angle)
	}
	// fmt.Println("winkel", angle)
	compressLeft := angle
	factorLeft := math.Cos(compressLeft * math.Pi / 180)
	if factorLeft < 0.01 {
		factorLeft = 0
	}

	compressright := 120 - angle
	factorRight := math.Cos(compressright * math.Pi / 180)
	if factorRight < 0.01 {
		factorRight = 0
	}

	// fmt.Println("factorLeft,factorRight", factorLeft, factorRight)

	stripedImageBounds := leftImg.Bounds()
	stripedImageBounds.Max.X = leftImg.Bounds().Dx() + rightImg.Bounds().Dx()
	stripedImageBounds.Max.Y = max(leftImg.Bounds().Dy(), rightImg.Bounds().Dy())

	striped_img := image.NewRGBA(stripedImageBounds)
	// draw.Draw(striped_img, striped_img.Bounds(), leftImg, image.Pt(0, 0), draw.Src)
	// start_second := stripedImageBounds
	// start_second.Min = (image.Point{X: leftImg.Bounds().Max.X, Y: 0})
	// draw.Draw(striped_img, start_second, rightImg, image.Pt(0, 0), draw.Src)

	sliceWidth := leftImg.Bounds().Max.X / 19
	// fmt.Println(sliceWidth)
	targetStartX := 0
	sliceCount := 19
	for index := 0; index < sliceCount; index++ {
		startX := sliceWidth * index
		endX := startX + sliceWidth
		// fmt.Println(startX, endX)
		// copy left
		srcLeft := leftImg.Bounds().Intersect(image.Rectangle{Min: image.Pt(startX, 0), Max: image.Pt(endX, leftImg.Bounds().Max.Y)})
		// fmt.Println("index", index, "srcLeft", srcLeft)
		startPoint := image.Pt(targetStartX, 0)
		// fmt.Println(startPoint, srcLeft.Size())
		r := image.Rectangle{startPoint, startPoint.Add(srcLeft.Size())}
		startPointdraw := leftImg.Bounds().Min.Add(image.Pt(startX, 0))
		// fmt.Println("r", r, "startPointdraw", startPointdraw)
		draw.Draw(striped_img, r, leftImg, startPointdraw, draw.Src)
		// add sliceWidth
		targetStartX += sliceWidth
		// // copy img right
		// startPoint_2 := image.Pt(targetStartX, 0)
		// fmt.Println(startPoint_2, srcLeft.Size())
		startPointdraw_2 := rightImg.Bounds().Min.Add(image.Pt(startX, 0))
		// r_2 := image.Rectangle{startPoint_2, startPoint.Add(srcLeft.Size())}
		draw.Draw(striped_img, r.Add(image.Pt(sliceWidth, 0)), rightImg, startPointdraw_2, draw.Src)
		targetStartX += sliceWidth

	}
	compressedsliceWidthLeft := float64(leftImg.Bounds().Max.X/19) * factorLeft
	compressedsliceWidthRight := float64(leftImg.Bounds().Max.X/19) * factorRight

	compressedWidth := (compressedsliceWidthLeft + compressedsliceWidthRight) * float64(sliceCount)

	compressedImageBounds := leftImg.Bounds()
	compressedImageBounds.Max.X = int(compressedWidth)
	compressedImageBounds.Max.Y = max(leftImg.Bounds().Dy(), rightImg.Bounds().Dy())

	compressed_img := image.NewRGBA(compressedImageBounds)
	compressed_src_left := image.NewNRGBA(image.Rect(0, 0, int(compressedsliceWidthLeft), leftImg.Bounds().Max.Y))
	compressed_src_right := image.NewNRGBA(image.Rect(0, 0, int(compressedsliceWidthRight), rightImg.Bounds().Max.Y))
	// fmt.Println("compressedsliceWidthLeft, compressedsliceWidthRight, compressedImageBounds, compressed_src_left.Bounds(),  compressed_src_right.Bounds()")
	// fmt.Println(compressedsliceWidthLeft, compressedsliceWidthRight, compressedImageBounds, compressed_src_left.Bounds(), compressed_src_right.Bounds())
	targetStartX = 0
	for index := 0; index < sliceCount; index++ {
		startX := float64(sliceWidth) * float64(index)
		endX := startX + float64(sliceWidth)
		// fmt.Println("left startX, endX", startX, endX)
		// copy left
		srcLeft := leftImg.Bounds().Intersect(image.Rectangle{Min: image.Pt(int(startX), 0), Max: image.Pt(int(endX), leftImg.Bounds().Max.Y)})
		xdraw.ApproxBiLinear.Scale(compressed_src_left, compressed_src_left.Rect, leftImg, srcLeft.Bounds(), xdraw.Src, nil)
		// fmt.Println("index", index, "srcLeft", srcLeft)
		startPoint := image.Pt(targetStartX, 0)
		// fmt.Println(startPoint, srcLeft.Size())
		r := image.Rectangle{startPoint, startPoint.Add(compressed_src_left.Rect.Max)}
		// fmt.Println("r", r)
		// startPointdraw := leftImg.Bounds().Min.Add(image.Pt(int(compressed_src_left.Rect.Max.X), 0))
		startPointdraw := image.Pt(0, 0)
		// fmt.Println("r", r, "startPointdraw", startPointdraw)
		draw.Draw(compressed_img, r, compressed_src_left, startPointdraw, draw.Src)
		// add sliceWidth
		targetStartX += int(compressedsliceWidthLeft)
		// // copy img right
		srcRight := rightImg.Bounds().Intersect(image.Rectangle{Min: image.Pt(int(startX), 0), Max: image.Pt(int(endX), leftImg.Bounds().Max.Y)})
		xdraw.ApproxBiLinear.Scale(compressed_src_right, compressed_src_right.Rect, rightImg, srcRight.Bounds(), xdraw.Src, nil)

		startPoint = image.Pt(targetStartX, 0)
		r = image.Rectangle{startPoint, startPoint.Add(compressed_src_right.Rect.Max)}
		draw.Draw(compressed_img, r, compressed_src_right, startPointdraw, draw.Src)

		// fmt.Println(startPoint_2, srcLeft.Size())
		// startPointdraw_2 := rightImg.Bounds().Min.Add(image.Pt(startX, 0))
		// // r_2 := image.Rectangle{startPoint_2, startPoint.Add(srcLeft.Size())}
		// draw.Draw(compressed_img, r.Add(image.Pt(sliceWidth, 0)), rightImg, startPointdraw_2, draw.Src)
		targetStartX += int(compressedsliceWidthRight)

	}

	// fmt.Println(striped_img.Bounds(), compressedImageBounds, compressed_img.Bounds())
	// xdraw.ApproxBiLinear.Scale(compressed_img, compressed_img.Rect, striped_img, striped_img.Bounds(), xdraw.Src, nil)
	// fmt.Println(striped_img.Bounds(), compressedImageBounds, compressed_img.Bounds())
	combined := CombinedImages{
		striped:    striped_img,
		compressed: compressed_img,
	}

	return &combined

}
