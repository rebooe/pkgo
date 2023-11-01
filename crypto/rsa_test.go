package crypto

import (
	"encoding/base64"
	"testing"
)

var (
	pubKey  = `MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAMr5GLnZwrLvdurBrQxBDQDxIyNvLmHmTGvS8FQdzwChWQT+mdm6ApuTAK69yg3ETmjvjtzSobRO5f8MLEtGnTkCAwEAAQ==`
	privKey = `MIIBVQIBADANBgkqhkiG9w0BAQEFAASCAT8wggE7AgEAAkEAyvkYudnCsu926sGtDEENAPEjI28uYeZMa9LwVB3PAKFZBP6Z2boCm5MArr3KDcROaO+O3NKhtE7l/wwsS0adOQIDAQABAkA7fEuR1E8qb+HzJTXZHIt6FjHNJb17Nap6A0Up8d6D+T9XvV9LgLH2zENkIZ/7sFYdOTdq2e90zlmYJH9b6JYxAiEA+diqgSzf295lc+TEilL71mpjVq1prnnuxC6qB9ad1O0CIQDP+OJnb1CRj2Q/W8zdFNbduR2EmK/XFOzSdojvVhkL/QIgAXLcKjuUYLX9aJqe+R5aD3g2cz42KqjSVZjfq4P3DlECIQCfGj9SdCVGBlXh5r/2TlAGteywGQNE3vxCEm618r8cnQIhAPMntI+EySA+MJBTPvD5eYxrfpZfrkqcgZdvc/wQNbiw`
)

func TestParsePublicKey(t *testing.T) {
	parse := ParseRsaKey{}
	key, err := parse.DecodeBase64(pubKey).ToPublicKey()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(key)
}

func TestParsePrivateKey(t *testing.T) {
	parse := ParseRsaKey{}
	key, err := parse.DecodeBase64(privKey).ToPrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(key)
}

func TestEncrypt(t *testing.T) {
	parse := ParseRsaKey{}
	key, err := parse.DecodeBase64(pubKey).ToPublicKey()
	if err != nil {
		t.Fatal(err)
	}

	res, err := RSAEncrypt(key, []byte(`abc`))
	if err != nil {
		t.Fatal(err)
	}
	resStr := base64.StdEncoding.EncodeToString(res)
	t.Logf("%s", resStr)
}

func TestDecrypt(t *testing.T) {
	text, err := base64.StdEncoding.DecodeString("KgR9ypep0XmfU8XLSCYTzp3nPh7N6+gGkKga6vhviCDb1xMMBzdIwOiM4NRFWrLXbEshtz9mvIZYRG3tLop4eQ==")
	if err != nil {
		t.Fatal(err)
	}

	parse := ParseRsaKey{}
	key, err := parse.DecodeBase64(privKey).ToPrivateKey()
	if err != nil {
		t.Fatal(err)
	}

	res, err := RSADecrypt(key, text)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", res)
}
