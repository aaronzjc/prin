package iptables

import (
	"strings"
	"testing"
)

func TestReadWord(t *testing.T) {
	cases := []struct {
		Input  string
		Expect string
	}{
		{"hello world", "hello"},
		{"", ""},
		{"hello", "hello"},
		{"\nhello", "\nhello"},
		{"-- hello", "--"},
		{"-p ", "-p"},
		{"! -p", "!"},
	}

	for _, v := range cases {
		s := NewScanner(strings.NewReader(v.Input))
		word, err := s.ReadWord()
		if err != nil {
			t.Fatal(v.Input, word, err)
		}
		if word != v.Expect {
			t.Fatal("input", v.Input, "result", word, err)
		}
	}

	t.Log("read word passed")
}

func TestPeek(t *testing.T) {
	cases := []struct {
		Input  string
		Expect string
		Err    bool
	}{
		{"-- hello", "--", false},
		{"-p ", "-p", false},
		{"! -p", "! ", false},
		{"", "", false},
		{"a", "a", false},
	}

	for _, v := range cases {
		s := NewScanner(strings.NewReader(v.Input))
		word, err := s.Peek(2)
		if v.Err != (err != nil) {
			t.Fatal("input", v.Input, "result", err)
		}
		if word != v.Expect {
			t.Fatal("input", v.Input, "expect", v.Expect, "result", word)
		}
	}
	t.Log("read peek passed")
}

func TestReadComment(t *testing.T) {
	cases := []struct {
		Input  string
		Expect string
		Err    bool
	}{
		{`hello "world"`, "world", false},
		{`hello "world`, "", true},
		{`follow me "i love \"u\" woo" end`, `i love \"u\" woo`, false},
	}

	for _, v := range cases {
		s := NewScanner(strings.NewReader(v.Input))
		word, err := s.ReadComment()
		if v.Err != (err != nil) {
			t.Fatal("input", v.Input, "expect", v.Err, "result", err)
		}
		if word != v.Expect {
			t.Fatal("input", v.Input, "expect", v.Expect, "result", word)
		}
	}
	t.Log("read comment passed")
}

func TestReadLine(t *testing.T) {
	cases := []struct {
		Input  string
		Expect string
	}{
		{"hello world", "hello world"},
		{"hello\nworld", "hello"},
		{"\nhello", ""},
	}

	for _, v := range cases {
		s := NewScanner(strings.NewReader(v.Input))
		line, _ := s.ReadLine()
		if line != v.Expect {
			t.Fatal("input", v.Input, "expect", v.Expect, "result", line)
		}
	}

	t.Log("read line passed")
}
