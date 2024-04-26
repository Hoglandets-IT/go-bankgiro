package seal

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"testing"
)

func TestHmacSealer_SetHashFunction(t *testing.T) {
	type fields struct {
		Key            []byte
		KeyVer         []byte
		Hash           func() hash.Hash
		Mac            []byte
		OriginalData   []byte
		FormattedData  []byte
		NormalizedData []byte
	}
	type args struct {
		hashFunc func() hash.Hash
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "SHA 1",
			fields: fields{
				Hash: sha1.New,
			},
			args: args{
				hashFunc: sha1.New,
			},
		},
		{
			name: "SHA 256",
			fields: fields{
				Hash: sha256.New,
			},
			args: args{
				hashFunc: sha256.New,
			},
		},
		{
			name: "SHA 512",
			fields: fields{
				Hash: sha512.New,
			},
			args: args{
				hashFunc: sha512.New,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hm := &HmacSealer{
				Key:            tt.fields.Key,
				KeyVer:         tt.fields.KeyVer,
				Hash:           tt.fields.Hash,
				Mac:            tt.fields.Mac,
				OriginalData:   tt.fields.OriginalData,
				FormattedData:  tt.fields.FormattedData,
				NormalizedData: tt.fields.NormalizedData,
			}
			hm.SetHashFunction(tt.args.hashFunc)
		})
	}
}

func TestHmacSealer_SetKey(t *testing.T) {
	type fields struct {
		Key            []byte
		KeyVer         []byte
		Hash           func() hash.Hash
		Mac            []byte
		OriginalData   []byte
		FormattedData  []byte
		NormalizedData []byte
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Valid Key",
			fields: fields{
				Key: []byte("1234567890ABCDEF1234567890ABCDEF"),
			},
			args: args{
				key: "1234567890ABCDEF1234567890ABCDEF",
			},
		},
		{
			name: "Invalid Key: Odd Length",
			fields: fields{
				Key: []byte("1234567890ABCDEF1234567890ABCDE"),
			},
			args: args{
				key: "1234567890ABCDEF1234567890ABCDE",
			},
			wantErr: true,
		},
		{
			name: "Invalid Key: Too Short",
			fields: fields{
				Key: []byte("123456"),
			},
			args: args{
				key: "123456",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hm := &HmacSealer{
				Key:            tt.fields.Key,
				KeyVer:         tt.fields.KeyVer,
				Hash:           tt.fields.Hash,
				Mac:            tt.fields.Mac,
				OriginalData:   tt.fields.OriginalData,
				FormattedData:  tt.fields.FormattedData,
				NormalizedData: tt.fields.NormalizedData,
			}
			if err := hm.SetKey(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("HmacSealer.SetKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHmacSealer_SetKeyBytes(t *testing.T) {
	type fields struct {
		Key            []byte
		KeyVer         []byte
		Hash           func() hash.Hash
		Mac            []byte
		OriginalData   []byte
		FormattedData  []byte
		NormalizedData []byte
	}
	type args struct {
		key []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hm := &HmacSealer{
				Key:            tt.fields.Key,
				KeyVer:         tt.fields.KeyVer,
				Hash:           tt.fields.Hash,
				Mac:            tt.fields.Mac,
				OriginalData:   tt.fields.OriginalData,
				FormattedData:  tt.fields.FormattedData,
				NormalizedData: tt.fields.NormalizedData,
			}
			if err := hm.SetKeyBytes(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("HmacSealer.SetKeyBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHmacSealer_GenerateKvv(t *testing.T) {
	type fields struct {
		Key            []byte
		KeyVer         []byte
		Hash           func() hash.Hash
		Mac            []byte
		OriginalData   []byte
		FormattedData  []byte
		NormalizedData []byte
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hm := &HmacSealer{
				Key:            tt.fields.Key,
				KeyVer:         tt.fields.KeyVer,
				Hash:           tt.fields.Hash,
				Mac:            tt.fields.Mac,
				OriginalData:   tt.fields.OriginalData,
				FormattedData:  tt.fields.FormattedData,
				NormalizedData: tt.fields.NormalizedData,
			}
			if err := hm.GenerateKvv(); (err != nil) != tt.wantErr {
				t.Errorf("HmacSealer.GenerateKvv() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHmacSealer_CheckKvv(t *testing.T) {
	type fields struct {
		Key            []byte
		KeyVer         []byte
		Hash           func() hash.Hash
		Mac            []byte
		OriginalData   []byte
		FormattedData  []byte
		NormalizedData []byte
	}
	type args struct {
		kvv string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hm := &HmacSealer{
				Key:            tt.fields.Key,
				KeyVer:         tt.fields.KeyVer,
				Hash:           tt.fields.Hash,
				Mac:            tt.fields.Mac,
				OriginalData:   tt.fields.OriginalData,
				FormattedData:  tt.fields.FormattedData,
				NormalizedData: tt.fields.NormalizedData,
			}
			if err := hm.CheckKvv(tt.args.kvv); (err != nil) != tt.wantErr {
				t.Errorf("HmacSealer.CheckKvv() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHmacSealer_UpdateFormatted(t *testing.T) {
	type fields struct {
		Key            []byte
		KeyVer         []byte
		Hash           func() hash.Hash
		Mac            []byte
		OriginalData   []byte
		FormattedData  []byte
		NormalizedData []byte
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hm := &HmacSealer{
				Key:            tt.fields.Key,
				KeyVer:         tt.fields.KeyVer,
				Hash:           tt.fields.Hash,
				Mac:            tt.fields.Mac,
				OriginalData:   tt.fields.OriginalData,
				FormattedData:  tt.fields.FormattedData,
				NormalizedData: tt.fields.NormalizedData,
			}
			hm.UpdateFormatted()
		})
	}
}

func TestHmacSealer_SetData(t *testing.T) {
	type fields struct {
		Key            []byte
		KeyVer         []byte
		Hash           func() hash.Hash
		Mac            []byte
		OriginalData   []byte
		FormattedData  []byte
		NormalizedData []byte
	}
	type args struct {
		data string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hm := &HmacSealer{
				Key:            tt.fields.Key,
				KeyVer:         tt.fields.KeyVer,
				Hash:           tt.fields.Hash,
				Mac:            tt.fields.Mac,
				OriginalData:   tt.fields.OriginalData,
				FormattedData:  tt.fields.FormattedData,
				NormalizedData: tt.fields.NormalizedData,
			}
			hm.SetData(tt.args.data)
		})
	}
}

func TestHmacSealer_AddData(t *testing.T) {
	type fields struct {
		Key            []byte
		KeyVer         []byte
		Hash           func() hash.Hash
		Mac            []byte
		OriginalData   []byte
		FormattedData  []byte
		NormalizedData []byte
	}
	type args struct {
		data string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hm := &HmacSealer{
				Key:            tt.fields.Key,
				KeyVer:         tt.fields.KeyVer,
				Hash:           tt.fields.Hash,
				Mac:            tt.fields.Mac,
				OriginalData:   tt.fields.OriginalData,
				FormattedData:  tt.fields.FormattedData,
				NormalizedData: tt.fields.NormalizedData,
			}
			hm.AddData(tt.args.data)
		})
	}
}

func TestHmacSealer_SetDataBytes(t *testing.T) {
	type fields struct {
		Key            []byte
		KeyVer         []byte
		Hash           func() hash.Hash
		Mac            []byte
		OriginalData   []byte
		FormattedData  []byte
		NormalizedData []byte
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hm := &HmacSealer{
				Key:            tt.fields.Key,
				KeyVer:         tt.fields.KeyVer,
				Hash:           tt.fields.Hash,
				Mac:            tt.fields.Mac,
				OriginalData:   tt.fields.OriginalData,
				FormattedData:  tt.fields.FormattedData,
				NormalizedData: tt.fields.NormalizedData,
			}
			hm.SetDataBytes(tt.args.data)
		})
	}
}

func TestHmacSealer_AddDataBytes(t *testing.T) {
	type fields struct {
		Key            []byte
		KeyVer         []byte
		Hash           func() hash.Hash
		Mac            []byte
		OriginalData   []byte
		FormattedData  []byte
		NormalizedData []byte
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hm := &HmacSealer{
				Key:            tt.fields.Key,
				KeyVer:         tt.fields.KeyVer,
				Hash:           tt.fields.Hash,
				Mac:            tt.fields.Mac,
				OriginalData:   tt.fields.OriginalData,
				FormattedData:  tt.fields.FormattedData,
				NormalizedData: tt.fields.NormalizedData,
			}
			hm.AddDataBytes(tt.args.data)
		})
	}
}

func TestHmacSealer_Validate(t *testing.T) {
	type fields struct {
		Key            []byte
		KeyVer         []byte
		Hash           func() hash.Hash
		Mac            []byte
		OriginalData   []byte
		FormattedData  []byte
		NormalizedData []byte
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hm := &HmacSealer{
				Key:            tt.fields.Key,
				KeyVer:         tt.fields.KeyVer,
				Hash:           tt.fields.Hash,
				Mac:            tt.fields.Mac,
				OriginalData:   tt.fields.OriginalData,
				FormattedData:  tt.fields.FormattedData,
				NormalizedData: tt.fields.NormalizedData,
			}
			if err := hm.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("HmacSealer.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHmacSealer_Calculate(t *testing.T) {
	type fields struct {
		Key            []byte
		KeyVer         []byte
		Hash           func() hash.Hash
		Mac            []byte
		OriginalData   []byte
		FormattedData  []byte
		NormalizedData []byte
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hm := &HmacSealer{
				Key:            tt.fields.Key,
				KeyVer:         tt.fields.KeyVer,
				Hash:           tt.fields.Hash,
				Mac:            tt.fields.Mac,
				OriginalData:   tt.fields.OriginalData,
				FormattedData:  tt.fields.FormattedData,
				NormalizedData: tt.fields.NormalizedData,
			}
			if err := hm.Calculate(); (err != nil) != tt.wantErr {
				t.Errorf("HmacSealer.Calculate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
