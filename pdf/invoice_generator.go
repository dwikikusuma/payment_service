package pdf

import (
	// golang package
	"fmt"
	"payment_service/models"

	// external package
	"github.com/phpdave11/gofpdf"
)

// GenerateInvoicePDF generate invoice pdf by given Payment, and outputPath.
//
// It returns nil error when successful.
// Otherwise, error will be returned.
func GenerateInvoicePDF(payment models.Payment, outputPath string) error {
	pdf := gofpdf.New("P", "mm", "A4", "arial")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "[FC] Invoice Details")

	pdf.Ln(20)
	pdf.SetFont("Arial", "", 12)

	pdf.Cell(40, 10, fmt.Sprintf("Payment ID: #%d", payment.ID))
	pdf.Ln(10)

	pdf.Cell(40, 10, fmt.Sprintf("Order ID: #%d", payment.OrderID))
	pdf.Ln(10)

	pdf.Cell(40, 10, fmt.Sprintf("Total Amount: Rp%.2f", payment.Amount))
	pdf.Ln(10)

	pdf.Cell(40, 10, fmt.Sprintf("Status: %s", payment.Status))
	pdf.Ln(10)

	pdf.Cell(40, 10, fmt.Sprintf("External ID: %s", payment.ExternalID))

	return pdf.OutputFileAndClose(outputPath)
}
