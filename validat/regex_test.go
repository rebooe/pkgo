package validat

import (
	"testing"
)

type validatorTest struct {
	Mobile string `binding:"required,len=3"`
}

func Test_regex(t *testing.T) {
	users := validatorTest{
		Mobile: "123123",
	}
	if err := Struct(&users); err != nil {
		t.Errorf("%v", err)
	}
}
