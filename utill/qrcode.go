package utill

import (
	"encoding/base64"
	"fmt"
	"library/Entity"
	"log"

	"github.com/skip2/go-qrcode"
)

func GenerateQRCode(book Entity.BookInventory) string {
	var s string
	if book.AvailableCopies > 0 {
		s = "Available"

	} else {
		s = "Not Available"

	}
	bookData := fmt.Sprintf("Title: %s\nAuthor: %s\nPublisher: %s\nVersion: %s\nAvailablity: %s", book.Title, book.Authors, book.Publisher, book.Version, s)
	qrCodeBytes, err := qrcode.Encode(bookData, qrcode.Medium, 256)
	if err != nil {
		log.Fatal("Error generating QR code:", err)
	}
	base64ImageData := base64.StdEncoding.EncodeToString(qrCodeBytes)

	return base64ImageData
}
