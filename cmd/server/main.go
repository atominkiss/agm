package main

import (
	db "agm/internal/database"
	qrProcesser "agm/internal/qrProcesser"
	"encoding/json"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"gocv.io/x/gocv"
	"log"
	"sync"
)

var pdfPool = sync.Pool{
	New: func() interface{} {
		pdf := gofpdf.New("P", "mm", "A4", "")
		return pdf
	},
}

func main() {
	// Открытие камеры
	webcam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		log.Fatalf("Ошибка открытия видеокамеры: %v", err)
	}
	defer func(webcam *gocv.VideoCapture) {
		err := webcam.Close()
		if err != nil {
			log.Fatalf("Ошибка закрытия видеокамеры: %v", err)
		}
	}(webcam)

	// Подготовка матрицы изображения
	img := gocv.NewMat()
	defer func(img *gocv.Mat) {
		err := img.Close()
		if err != nil {
			log.Fatalf("Ошибка закрытия изображения: %v", err)
		}
	}(&img)

	window := gocv.NewWindow("QR Code Reader")
	defer func(window *gocv.Window) {
		err := window.Close()
		if err != nil {
			log.Fatalf("Ошибка закрытия окна: %v", err)
		}
	}(window)

	fmt.Println("Нажмите ESC для выхода")
	for {
		// Считывание изображения с камеры
		if ok := webcam.Read(&img); !ok {
			log.Println("Невозможно захватить изображение")
			return
		}
		if img.Empty() {
			continue
		}

		// Декодирование QR-кода с изображения
		qrCodeContent, err := qrProcesser.DecodeQRCode(img)
		if err == nil {
			fmt.Printf("QR-код обнаружен: %s\n", qrCodeContent)

			// Преобразование данных в JSON
			jsonData, err := json.Marshal(qrProcesser.QRData{Content: qrCodeContent})
			if err != nil {
				log.Fatalf("Ошибка кодирования JSON: %v", err)
			}

			// Сохранение JSON в базу данных DataGrip (PostgreSQL)
			db.SaveToDatabase(jsonData)

			// Извлечение JSON из базы данных
			extractedJSON, err := db.FetchJSONFromDB(qrCodeContent)
			if err != nil {
				log.Fatalf("Ошибка извлечения JSON из базы данных: %v", err)
			}
			// Генерация и печать QR-кода из извлеченного JSON
			qrProcesser.GenerateAndPrintQRCode(extractedJSON)
			break
		}

		// Отображение изображения в окне
		window.IMShow(img)
		if window.WaitKey(1) == 27 {
			break
		}
	}
}
