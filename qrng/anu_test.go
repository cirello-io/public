package qrng_test

import (
	"testing"

	"cirello.io/qrng"
)

func TestRead(t *testing.T) {
	t.Run("test small", func(t *testing.T) {
		buf := make([]byte, 1024)
		n, err := qrng.Read(buf)
		if err != nil {
			t.Fatal("unexpected error found:", err)
		} else if n != len(buf) {
			t.Fatal("unexpected partial read:", n, len(buf))
		}
	})
	t.Run("test large", func(t *testing.T) {
		buf := make([]byte, 2048)
		n, err := qrng.Read(buf)
		if err != nil {
			t.Fatal("unexpected error found:", err)
		} else if n != len(buf) {
			t.Fatal("unexpected partial read:", n, len(buf))
		}
	})
}

func TestUint8(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"happy case", args{1}, false},
		{"too small", args{-1}, true},
		{"too large", args{1025}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := qrng.Uint8(tt.args.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint8() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && len(got) != tt.args.length {
				t.Errorf("Uint8() = %v, want %v", len(got), tt.args.length)
			}
		})
	}
}

func TestUint16(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"happy case", args{1}, false},
		{"too small", args{-1}, true},
		{"too large", args{1025}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := qrng.Uint16(tt.args.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint16() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && len(got) != tt.args.length {
				t.Errorf("Uint16() = %v, want %v", len(got), tt.args.length)
			}
		})
	}
}

func TestHex16(t *testing.T) {
	type args struct {
		length    int
		blockSize int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"happy case", args{1, 1}, false},
		{"length too small", args{-1, 1}, true},
		{"length too large", args{1025, 1}, true},
		{"blockSize too small", args{1, -1}, true},
		{"blockSize too large", args{1, 1025}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := qrng.Hex16(tt.args.length, tt.args.blockSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hex16() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && len(got) != tt.args.length {
				t.Errorf("Hex16() = %v, want %v", len(got), tt.args.length)
			}
		})
	}
}
