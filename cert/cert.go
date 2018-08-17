package cert

import (
	"bytes"
	"io"
	"time"

	"github.com/covveco/micro-learning/cert/assets"
	"github.com/jung-kurt/gofpdf"
)

const (
	fontFamily string = "Arial"
)

type Certificate struct {
	Name string
	Date time.Time
}

func (c *Certificate) Generate(w io.Writer) error {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetTextColor(54, 85, 157)

	// Background image
	certBytes, err := assets.Asset("cert.png")
	if err != nil {
		return err
	}
	pdf.RegisterImageOptionsReader("cert", gofpdf.ImageOptions{ImageType: "png"}, bytes.NewReader(certBytes))
	pdf.ImageOptions("cert", 0, 0, 297, 210, false, gofpdf.ImageOptions{}, 0, "")

	// Name
	pdf.SetFont(fontFamily, "B", 28)
	pdf.MoveTo(18, 95) // 17.5, not 18
	pdf.MultiCell(260, 10, c.Name, "", "C", false)

	// Date
	pdf.SetFont(fontFamily, "B", 20)
	pdf.MoveTo(15, 173)
	pdf.MultiCell(100, 10, formatDate(c.Date), "", "C", false)

	return pdf.Output(w)
}

func formatDate(d time.Time) string {
	return d.Format("02.01.2006")
}
