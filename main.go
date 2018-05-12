package main

import (
	"io/ioutil"
	"log"
	"gopkg.in/gographics/imagick.v2/imagick"
	"github.com/otiai10/gosseract"
	"os"
	"bufio"
)

const FORMAT = "png"

func main() {
	srcFolder := os.Args[1]
	dstFolder := os.Args[2]
	files, err := ioutil.ReadDir(srcFolder)
	files[0].Name()
	if err != nil {
		log.Fatal("could not read source Directory", err)
	}

	for _, file := range files {
		if err := ConvertPdfToPng(srcFolder+"/"+file.Name(), dstFolder+"/.png/"+file.Name()); err != nil {
			log.Fatal(err)
		}
		if err := pngToTxt(srcFolder+"/.png/"+file.Name(), dstFolder+"/.txt/"+file.Name()); err != nil {
			log.Fatal(err)
		}

	}
}

func ConvertPdfToPng(src string, dst string) error {
	// Setup
	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	if err := mw.SetResolution(300, 300); err != nil {
		return err
	}
	if err := mw.ReadImage(src); err != nil {
		return err
	}
	if err := mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_FLATTEN); err != nil {
		return err
	}
	if err := mw.SetCompressionQuality(95); err != nil {
		return err
	}
	mw.SetIteratorIndex(0)

	if err := mw.SetFormat(FORMAT); err != nil {
		return err
	}
	return mw.WriteImage(dst)
}

func pngToTxt(src string, dst string) error {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage(src)
	text, _ := client.Text()
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	_, err = w.WriteString(text)
	if err != nil {
		return err
	}
	return nil
}
