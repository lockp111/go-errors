package errors

import (
	"errors"
	"fmt"
	"io"
	"testing"
)

func TestParse(t *testing.T) {
	test1 := Register(404, "404 error")
	wrap1 := fmt.Errorf("wrap error: %w", test1)
	target1 := Parse(wrap1)
	if target1.Code() != test1.Code() {
		t.Fatal("p1-1")
	}

	withMsg := test1.WithMessage("with msg")
	wrap2 := fmt.Errorf("wrap error: %w", withMsg)
	target2 := Parse(wrap2)
	if target2.Code() != test1.Code() {
		t.Fatal("p2-1")
	}

	withErr := errors.New("with err")
	test3 := test1.WithError(withErr)
	wrap3 := fmt.Errorf("wrap error: %w", test3)
	target3 := Parse(wrap3)
	if target3.Code() != test1.Code() {
		t.Fatal("p3-1")
	}
}

func TestWithError(t *testing.T) {
	test1 := Register(403, "403 error")
	test2 := Register(404, "404 error")

	e1 := test1.WithError(io.ErrUnexpectedEOF)
	if !Is(e1, io.ErrUnexpectedEOF) {
		t.Fatal("e1-1")
	}
	if !Is(e1, test1) {
		t.Fatal("e1-2")
	}

	e2 := test2.WithError(e1)
	if !Is(e2, test1) {
		t.Fatal("e2-1")
	}
	if !Is(e2, io.ErrUnexpectedEOF) {
		t.Fatal("e2-2")
	}

	e3 := fmt.Errorf("wrap error: %w", e2)
	if !Is(e3, e2) {
		t.Fatal("e3-1")
	}
	if !Is(e3, test2) {
		t.Fatal("e3-2")
	}
	if !Is(e3, test1) {
		t.Fatal("e3-3")
	}
	if !Is(e3, io.ErrUnexpectedEOF) {
		t.Fatal("e3-4")
	}

	p3 := Parse(e3)
	if p3.Code() != test2.Code() {
		t.Fatal("e3-5")
	}
}

func TestWithMessage(t *testing.T) {
	test1 := Register(403, "403 error")

	m1 := test1.WithMessage("not found index")
	if !Is(m1, test1) {
		t.Fatal("m1-1")
	}

	m2 := fmt.Errorf("wrap error: %w", m1)
	if !Is(m2, test1) {
		t.Fatal("m2-1")
	}
}
