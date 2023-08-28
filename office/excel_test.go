package office

import (
	"os"
	"testing"
)

func Test_excel(t *testing.T) {
	excel := NewSimpExcel()
	if err := excel.SetRow(&[]string{"标题1", "标题2"}); err != nil {
		t.Error(err)
		return
	}
	if _, err := excel.WriteTo(os.Stdout); err != nil {
		t.Error(err)
		return
	}
}
