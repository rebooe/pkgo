package util

import (
	"testing"
)

type st struct {
	Aa int    `json:",omitempty"`
	Bb string `json:",omitempty"`
	Cc *int   `json:",omitempty"`
}

func TestStructToMap(t *testing.T) {
	c := 3
	type args struct {
		dest map[string]any
		src  any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestStructToMap",
			args: args{
				dest: map[string]any{},
				src: &st{
					Aa: 1,
					Bb: "a",
					Cc: &c,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StructToMap(tt.args.dest, tt.args.src); (err != nil) != tt.wantErr {
				t.Errorf("StructToMap() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Logf("%v", tt.args.dest)
		})
	}
}

func TestMapToStruct(t *testing.T) {
	c := 3
	type args struct {
		dest any
		src  map[string]any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestMapToStruct",
			args: args{
				dest: &st{},
				src: map[string]any{
					"Aa": 1,
					"Bb": "A",
					"Cc": &c,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MapToStruct(tt.args.dest, tt.args.src); (err != nil) != tt.wantErr {
				t.Errorf("MapToStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Logf("%+v", tt.args.dest)
		})
	}
}

func TestMapMerge(t *testing.T) {
	type args struct {
		dest map[string]any
		src  map[string]any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestMapMerge",
			args: args{
				dest: map[string]any{
					"Aa": 1,
					"Bb": "a",
				},
				src: map[string]any{
					"Aa": 2,
					"Bb": "b",
					"Cc": "c",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MapMerge(tt.args.dest, tt.args.src)
			t.Logf("%+v", tt.args.dest)
		})
	}
}

// func TestStructMerge(t *testing.T) {
// 	c := 3
// 	type args struct {
// 		dest any
// 		src  any
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "TestStructMerge",
// 			args: args{
// 				dest: &st{},
// 				src: &st{
// 					Aa: 1,
// 					Bb: "a",
// 					Cc: &c,
// 				},
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if err := StructMerge(tt.args.dest, tt.args.src); (err != nil) != tt.wantErr {
// 				t.Errorf("StructMerge() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 			t.Logf("%+v", tt.args.dest)
// 		})
// 	}
// }

func TestCaller(t *testing.T) {
	got, _ := Caller(0)
	t.Logf("skip:0, %s", got)

	func() {
		got, _ := Caller(0)
		t.Logf("skip:0, %s", got)
	}()
}
