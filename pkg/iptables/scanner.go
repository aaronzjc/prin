package iptables

import (
	"bufio"
	"io"
	"strings"
)

type Scanner struct {
	r *bufio.Reader
}

func (s *Scanner) ReadWord() (string, error) {
	word, err := s.r.ReadString(' ')
	if err != nil {
		if err.Error() == "EOF" {
			return strings.TrimRight(word, " "), nil
		}
	}
	return strings.TrimRight(word, " "), err
}

func (s *Scanner) Peek(n int) (string, error) {
	bs, err := s.r.Peek(n)
	if err != nil {
		if err.Error() == "EOF" {
			return string(bs), nil
		}
	}
	return string(bs), err
}

func (s *Scanner) ReadLine() (string, error) {
	str, _, err := s.r.ReadLine()
	return string(str), err
}

func (s *Scanner) ReadComment() (string, error) {
	counter := 2
	var err error
	var str []byte
	var cur, pre byte
	for counter > 0 {
		if cur, err = s.r.ReadByte(); err != nil {
			return "", err
		}
		if cur == '"' {
			if pre != '\\' {
				counter--
				continue
			}
		}
		// not comment
		if counter == 2 {
			continue
		}
		str = append(str, cur)
		pre = cur
	}
	s.r.ReadByte() // 多读取一个空格
	return string(str), nil
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}
