package office

import (
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"
)

type SimpExcel struct {
	file   *excelize.File
	curCol int
}

func NewSimpExcel() *SimpExcel {
	return &SimpExcel{
		file: excelize.NewFile(),
	}
}

func (exc *SimpExcel) SetRow(slice interface{}) error {
	exc.curCol++
	cell := fmt.Sprintf("A%d", exc.curCol)
	return exc.file.SetSheetRow("Sheet1", cell, slice)
}

func (exc *SimpExcel) WriteTo(w io.Writer, opts ...excelize.Options) (int64, error) {
	n, err := exc.file.WriteTo(w, opts...)
	if err != nil {
		return n, err
	}
	err = exc.file.Close()
	return n, err
}
