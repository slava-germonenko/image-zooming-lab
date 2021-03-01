package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
	"image"
	"image/jpeg"
	"log"
	"os"
	"time"
)

func main() {
	var source string
	var output string

	var scale int
	var zoom int

	// Считываем аргументы командной строки
	flag.StringVar(&source, "source", "", "Full path to an image to be zoomed.")
	flag.StringVar(&output, "output", "output.jpeg", "Output filename.")
	flag.IntVar(&scale, "scale", 1, "Scale value")
	flag.IntVar(&zoom, "zoom", 1, "Zooming value")
	flag.Parse()

	err := validateArgs(source, output, zoom, scale)
	if err != nil {
		log.Fatal(err.Error())
	}

	file, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}

	// Считываем картинку и представляем её как объект в Go
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	// Считываем информацию о ширине и высоте картинки
	_, _ = file.Seek(0, 0)
	imgConfig, _, err := image.DecodeConfig(file)
	if err != nil {
		log.Fatal(err)
	}
	_ = file.Close()

	outputHeight := uint(imgConfig.Width * scale)
	outputWidth := uint(imgConfig.Height * scale)

	start := time.Now()
	// Сначала увеличиваем картинку используя билинейную интерполяцию
	m := resize.Resize(
		outputHeight,
		outputWidth,
		img,
		resize.Bilinear,
		)

	// Обрезаем картинку
	cropped, err := cutter.Crop(m, cutter.Config{
		Width: int(outputHeight / uint(zoom)),
		Height: int(outputWidth / uint(zoom)),
		Mode: cutter.Centered,
	})

	fmt.Println("Time elapsed:", time.Now().Sub(start).Milliseconds())

	out, err := os.Create(output)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// Создаём новый файл картинки
	_ = jpeg.Encode(out, cropped, nil)
}

func validateArgs(source, output string, zoom, scale int) error {
	if len(source) == 0 {
		return errors.New("source image name is empty")
	}

	if len(output) == 0 {
		return errors.New("output filename can't be empty")
	}

	if scale < 0 {
		return errors.New("scale value must be greater than 0")
	}

	if zoom < 1 {
		return errors.New("zoom value must be greater than 0")
	}

	return nil
}
