package mm

import (
	"encoding/csv"
	"fmt"
	"os"
)

type CSVWriter struct {
	f *os.File
	w *csv.Writer
}

func NewCSVWriter(path string, header []string) (*CSVWriter, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	w := csv.NewWriter(f)
	if err := w.Write(header); err != nil {
		_ = f.Close()
		return nil, err
	}
	return &CSVWriter{f: f, w: w}, nil
}

func (c *CSVWriter) WriteRow(row []string) {
	_ = c.w.Write(row)
}

func (c *CSVWriter) Close() error {
	c.w.Flush()
	if err := c.w.Error(); err != nil {
		_ = c.f.Close()
		return err
	}
	return c.f.Close()
}

func fmtF(x float64) string { return fmt.Sprintf("%.6f", x) }
func fmtInt(x int) string   { return fmt.Sprintf("%d", x) }
