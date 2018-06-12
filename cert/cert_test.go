package cert_test

import (
	"os"
	"testing"
	"time"

	"gitlab.com/covveco/special-needs/cert"
)

func TestCertificate_Generate(t *testing.T) {
	f, err := os.Create("test.pdf")
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()

	c := &cert.Certificate{
		Name: "John Lim Boon Hock",
		Date: time.Now(),
	}
	if err := c.Generate(f); err != nil {
		t.Error(err)
		return
	}
}
