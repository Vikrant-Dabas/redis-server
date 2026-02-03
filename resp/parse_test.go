package resp

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func newReader(s string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(s))
}

func TestParseSimple(t *testing.T) {
	r := newReader("+OK\r\n")

	out, err := Parse(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "+ OK"
	if string(out) != want {
		t.Fatalf("got %q, want %q", out, want)
	}
}

func TestParseError(t *testing.T) {
	r := newReader("-ERR\r\n")

	out, err := Parse(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "- ERR"
	if string(out) != want {
		t.Fatalf("got %q, want %q", out, want)
	}
}

func TestParseInt(t *testing.T) {
	r := newReader(":42\r\n")

	out, err := Parse(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := ": 42"
	if string(out) != want {
		t.Fatalf("got %q, want %q", out, want)
	}
}

func TestParseBulk(t *testing.T) {
	r := newReader("$5\r\nHello\r\n")

	out, err := Parse(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "$ Hello"
	if string(out) != want {
		t.Fatalf("got %q, want %q", out, want)
	}
}

func TestParseArray(t *testing.T) {
	input :=
		"*3\r\n" +
			"+OK\r\n" +
			":1\r\n" +
			"-ERR\r\n"

	r := newReader(input)

	out, err := Parse(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want :=
		"*\n" +
			"+ - OK\n" +
			": - 1\n" +
			"- - ERR\n"

	if string(out) != want {
		t.Fatalf("got:\n%q\nwant:\n%q", out, want)
	}
}

func TestNestedArray(t *testing.T) {
	input :=
		"*1\r\n" +
			"*1\r\n" +
			":5\r\n"

	r := newReader(input)

	out, err := Parse(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want :=
		"*\n" +
			"* - \n"

	if string(out) != want {
		t.Fatalf("got %q, want %q", out, want)
	}
}

func TestInvalidTerminator(t *testing.T) {
	r := newReader("+OK\n")

	_, err := Parse(r)
	if err == nil {
		t.Fatal("expected error due to invalid terminator")
	}
}

func TestUnexpectedTypeInsidePayload(t *testing.T) {
	r := newReader("+O*K\r\n")

	_, err := Parse(r)
	if err == nil {
		t.Fatal("expected error due to unexpected message type inside payload")
	}
}

func TestUnknownMessageType(t *testing.T) {
	r := newReader("?WHAT\r\n")

	_, err := Parse(r)
	if err == nil {
		t.Fatal("expected error for unknown message type")
	}
}

func TestReadValidateInput(t *testing.T) {
	r := newReader("+PING\r\n")

	payload, err := readValidateInput(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(payload, []byte("+PING\r\n")) {
		t.Fatalf("got %q, want %q", payload, "+PING\r\n")
	}
}