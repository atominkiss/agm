package pdfProcesser

import (
	"bytes"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"os/exec"
	"sync"
)

var pdfPool = sync.Pool{ //nolint:gochecknoglobals
	New: func() interface{} {
		return gofpdf.New("P", "mm", "A4", "")
	},
}

func GetPDFInstance() *gofpdf.Fpdf {
	return pdfPool.Get().(*gofpdf.Fpdf)
}

func ReleasePDFInstance(pdf *gofpdf.Fpdf) {
	pdfPool.Put(pdf)
}

func GenerateAndSavePDF(pdf *gofpdf.Fpdf, jsonData []byte, qrCodeImage *bytes.Buffer) error {
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "QR Code Content:")

	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 10, string(jsonData), "", "", false)

	pdf.ImageOptions(
		qrCodeImage.String(),
		0,
		0,
		200,
		0,
		true,
		gofpdf.ImageOptions{
			ImageType: "PNG",
		},
		0,
		"",
	)

	return pdf.OutputFileAndClose("output qr-code.pdf")
}

func PrintPDF(pdf *gofpdf.Fpdf) error {
	pdfBuf := new(bytes.Buffer)
	err := pdf.Output(pdfBuf)
	if err != nil {
		return fmt.Errorf("failed to generate PDF: %w", err)
	}

	printCmd := exec.Command("lpr", "-P", "your_printer_name", "-")
	printCmd.Stdin = pdfBuf
	err = printCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to print PDF: %w", err)
	}

	return nil
}
