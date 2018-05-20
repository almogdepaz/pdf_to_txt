package main

import (
	"gopkg.in/gographics/imagick.v2/imagick"
	"github.com/otiai10/gosseract"
	"os"
	"bufio"
	"io/ioutil"
	"log"
	"fmt"
)

const FORMAT = "png"
const LANG = "heb"
const DEFAULT_SERC_FOLDER = "pdf"

func main() {
	pdfFolder := DEFAULT_SERC_FOLDER
	imgFolder := pdfFolder + "/images"
	txtFolder := pdfFolder + "/text"

	pdfs, err := ioutil.ReadDir(pdfFolder)
	if err != nil {
		log.Fatal("could not read source Directory", err)
	}

	for _, file := range pdfs {
		if err := ConvertPdfToPng(pdfFolder+"/"+file.Name(), imgFolder+"/"+file.Name(), FORMAT); err != nil {
			log.Fatal(err)
		}
	}

	log.Println("done converting to png")

	images, err := ioutil.ReadDir(imgFolder)

	if err != nil {
		log.Fatal("could not read source Directory", err)
	}

	for _, file := range images {
		if err := pngToTxt(imgFolder+"/"+file.Name(), txtFolder+"/"+file.Name()+".txt"); err != nil {
			log.Fatal(err)
		}
	}
}

func ConvertPdfToPng(src string, dst string, format string) error {
	// Setup
	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	if err := mw.SetResolution(800, 600); err != nil {
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
	if err := mw.SetFormat(format); err != nil {
		return err
	}
	numberOfPages := int(mw.GetNumberImages())
	for i := 0; i < numberOfPages; i++ {
		mw.SetIteratorIndex(i)
		path := fmt.Sprintf("(%d).", i)
		err := mw.WriteImage(dst + path + format)
		if err != nil {
			return err
		}
	}

	return nil
}

func pngToTxt(src string, dst string) error {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetLanguage(LANG)
	client.SetImage(src)
	text, _ := client.Text()
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()
	_, err = w.WriteString(text)
	if err != nil {
		return err
	}
	return nil
}
