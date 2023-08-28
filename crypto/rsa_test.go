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
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				key: pubKey,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePublicKey(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%v", got)
		})
	}
}

func TestParsePrivateKey(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				key: privKey,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePrivateKey(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%v", got)
		})
	}
}

func TestEncrypt(t *testing.T) {
	key, err := ParsePublicKey(pubKey)
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
	key, err := ParsePrivateKey(privKey)
	if err != nil {
		t.Fatal(err)
	}
	res, err := RSADecrypt(key, text)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", res)
}
