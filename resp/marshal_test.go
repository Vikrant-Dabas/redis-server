package resp

import (
	"bytes"
	"testing"
)

func TestSimpleMarshal(t *testing.T) {
	tests := []struct {
		name string
		f    Format
		want string
	}{
		{
			name: "simple string",
			f: Format{
				Type:    typeSimple,
				Payload: []byte("Hello"),
			},
			want: "+Hello\r\n",
		},
		{
			name: "error string",
			f: Format{
				Type:    typeError,
				Payload: []byte("Error"),
			},
			want: "-Error\r\n",
		},
		{
			name: "integer",
			f: Format{
				Type:    typeInt,
				Payload: []byte("14"),
			},
			want: ":14\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.Marshal()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !bytes.Equal(got, []byte(tt.want)) {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBulkMarshal(t *testing.T) {
	f := Format{
		Type:    typeBulk,
		Payload: []byte("Hello World"),
	}

	want := "$11\r\nHello World\r\n"

	got, err := f.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(got, []byte(want)) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestArrayMarshal(t *testing.T) {
	fArr := Format{
		Type: typeArray,
		ArrayPayload: []Format{
			{Type: typeSimple, Payload: []byte("Hello")},
			{Type: typeError, Payload: []byte("Error")},
			{Type: typeInt, Payload: []byte("14")},
			{Type: typeBulk, Payload: []byte("Hello World")},
		},
	}

	want := "*4\r\n" +
		"+Hello\r\n" +
		"-Error\r\n" +
		":14\r\n" +
		"$11\r\nHello World\r\n"

	got, err := fArr.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(got, []byte(want)) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestNestedArrayMarshal(t *testing.T) {
	f := Format{
		Type: typeArray,
		ArrayPayload: []Format{
			{
				Type: typeArray,
				ArrayPayload: []Format{
					{Type: typeInt, Payload: []byte("42")},
				},
			},
		},
	}

	want := "*1\r\n*1\r\n:42\r\n"

	got, err := f.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(got, []byte(want)) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestEmptyArrayMarshal(t *testing.T) {
	f := Format{
		Type:         typeArray,
		ArrayPayload: []Format{},
	}

	want := "*0\r\n"

	got, err := f.Marshal()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(got, []byte(want)) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestInvalidInteger(t *testing.T) {
	f := Format{
		Type:    typeInt,
		Payload: []byte("abc"),
	}

	_, err := f.Marshal()
	if err == nil {
		t.Fatal("expected error for invalid integer payload")
	}
}

func TestUnsupportedType(t *testing.T) {
	f := Format{
		Type:    'x',
		Payload: []byte("oops"),
	}

	_, err := f.Marshal()
	if err == nil {
		t.Fatal("expected error for unsupported type")
	}
}

func TestTerminatorInPayload(t *testing.T) {
	f := Format{
		Type:    typeSimple,
		Payload: []byte("Hello\r\nWorld"),
	}

	_, err := f.Marshal()
	if err == nil {
		t.Fatal("expected error due to terminator in payload")
	}
}