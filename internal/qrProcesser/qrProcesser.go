package pdfProcesser

import (
	pdfProcesser "agm/internal/pdfProcesser"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tuotoo/qrcode"
	"gocv.io/x/gocv"
	"image/png"
	"log"
	"os"
	"rsc.io/qr"
)

type QRData struct {
	Content string `json:"content"`
}

// DecodeQRCode decodes a QR code from an image.
func DecodeQRCode(image gocv.Mat) (string, error) {
	imageBytes, err := gocv.IMEncode(gocv.PNGFileExt, image)
	if err != nil {
		return "", err
	}
	defer imageBytes.Close()

	tempFile, err := os.CreateTemp("", "qr-code-*.png")
	if err != nil {
		return "", err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Fatal(err)
		}
	}(tempFile.Name())

	_, err = tempFile.Write(imageBytes.GetBytes())
	if err != nil {
		return "", err
	}

	qrCode, err := qrcode.Decode(tempFile)
	if err != nil {
		return "", err
	}
	return qrCode.Content, nil
}

// GenerateAndPrintQRCode generates and prints a QR code from JSON data.
func GenerateAndPrintQRCode(jsonData []byte) error {
	qrData, err := decodeJSONData(jsonData)
	if err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	qrCode, err := generateQRCode(qrData.Content)
	if err != nil {
		return fmt.Errorf("failed to generate QR code: %w", err)
	}

	qrCodeImage, err := encodeQRCodeImage(qrCode)
	if err != nil {
		return fmt.Errorf("failed to encode QR code image: %w", err)
	}

	pdf := pdfProcesser.GetPDFInstance()
	defer pdfProcesser.ReleasePDFInstance(pdf)

	err = pdfProcesser.GenerateAndSavePDF(pdf, jsonData, qrCodeImage)
	if err != nil {
		return fmt.Errorf("failed to generate and save PDF: %w", err)
	}

	err = pdfProcesser.PrintPDF(pdf)
	if err != nil {
		return fmt.Errorf("failed to print PDF: %w", err)
	}

	return nil
}

func decodeJSONData(jsonData []byte) (QRData, error) {
	var qrData QRData
	err := json.Unmarshal(jsonData, &qrData)
	return qrData, err
}

func generateQRCode(content string) (*qr.Code, error) {
	return qr.Encode(content, qr.H)
}

func encodeQRCodeImage(qrCode *qr.Code) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, qrCode.Image())
	return buf, err
}
