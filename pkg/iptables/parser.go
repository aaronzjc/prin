package iptables

import (
	"bufio"
	"context"
	"errors"
	"log"
	"regexp"
	"strings"
	"time"
)

var (
	traffic = map[string][]string{
		"in":      {"raw:PREROUTING", "mangle:PREROUTING", "nat:PREROUTING", "mangle:INPUT", "filter:INPUT", "nat:INPUT"},
		"forward": {"raw:PREROUTING", "mangle:PREROUTING", "nat:PREROUTING", "mangle:FORWARD", "filter:FORWARD", "mangle:POSTROUTING", "nat:POSTROUTING"},
		"out":     {"raw:OUTPUT", "mangle:OUTPUT", "nat:OUTPUT", "filter:output", "mangle:POSTROUTING", "nat:POSTROUTING"},
	}
)

type Table struct {
	Name   string
	Chains []*Chain
}

func (t *Table) Add(ch *Chain) {
	t.Chains = append(t.Chains, ch)
}

func (t *Table) GetChain(name string) (*Chain, error) {
	for _, v := range t.Chains {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, errors.New("not found")
}

type Chain struct {
	Name   string
	Policy string
	Stat   string
	Rules  []*Rule
}

func (c *Chain) Add(r *Rule) {
	c.Rules = append(c.Rules, r)
}

type Rule struct {
	Text    string
	Table   *Table
	Chain   *Chain
	Dst     Flag
	Src     Flag
	Proto   Flag
	Ifce    Flag
	Comment string
	Matches []*Match
	Target  *Target
}

func (r *Rule) AddMatch(match *Match) {
	// if it is comment, we dont treat it as a match
	if match.Name == "comment" {
		r.Comment = match.Flags[0].Val
		return
	}
	r.Matches = append(r.Matches, match)
}

type Flag struct {
	IsNot bool
	Key   string
	Val   string
}

func (f Flag) String() (res string) {
	if res = f.Key + f.Val; res == "" {
		return
	}
	if f.Key != "" && f.Val != "" {
		if f.IsNot {
			return f.Key + " != " + f.Val
		}
		return f.Key + " = " + f.Val
	}
	if f.IsNot {
		res = "! " + res
	}
	return
}

type Match struct {
	Name  string
	Flags []Flag
}

func (m *Match) String() string {
	var str string
	for _, v := range m.Flags {
		str = str + "(" + v.String() + ")"
	}
	if str != "" {
		return m.Name + ":" + str
	}
	return m.Name
}

func (m *Match) Parse(s *Scanner) error {
	m.Flags = []Flag{}
	m.Name, _ = s.ReadWord()
	for {
		var flag Flag
		try, err := s.Peek(2)
		if err != nil {
			return errors.New("invalid match")
		}
		if try == "" || (try != "! " && try != "--") {
			return nil
		}
		if try == "! " {
			flag.IsNot = true
			s.ReadWord() // skip "!"
		}
		key, _ := s.ReadWord()
		if key == "--comment" {
			flag.Val, _ = s.ReadComment()
		} else {
			flag.Val, _ = s.ReadWord()
		}
		flag.Key = key
		m.Flags = append(m.Flags, flag)
	}
}

type Target struct {
	Name string
	Flag Flag
}

func (t *Target) String() string {
	return t.Name + " " + t.Flag.String()
}

func (t *Target) Parse(s *Scanner) error {
	t.Flag = Flag{}
	t.Name, _ = s.ReadWord()
	for {
		try, err := s.Peek(2)
		if err != nil {
			return errors.New("invalid target")
		}
		if try == "" {
			return nil
		}
		if try != "--" {
			return nil
		}
		t.Flag.Key, _ = s.ReadWord()
		t.Flag.Val, _ = s.ReadWord()
	}
}

type Parser struct {
	s  *Scanner
	ct *Table            // 当前处理的Table
	tm map[string]*Table // 处理后的结构数据
}

func (p *Parser) GetFormatData() map[string]*Table {
	return p.tm
}

func (p *Parser) Parse() error {
	for {
		line, err := p.s.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				return nil
			}
			return err
		}
		if len(line) < 3 {
			continue
		}
		switch {
		case line[0] == '*':
			// table define
			tbName := line[1:]
			tb := Table{Name: tbName}
			p.tm[tbName] = &tb
			p.ct = &tb
		case line[0] == ':':
			// chain define
			segs := strings.Split(line, " ")
			ch := Chain{Name: segs[0][1:], Policy: segs[1], Stat: segs[2]}
			p.ct.Add(&ch)
		case line[:2] == "-A":
			// rule define
			rule, err := p.parseRule(line)
			if err != nil {
				log.Printf("parse rule err = %v, skip [%s]", err, line)
				continue
			}
			rule.Chain.Add(rule)
		case line[0] == '#': // comment
		default:
		}
	}
}

func (p *Parser) parseRule(line string) (*Rule, error) {
	var err error
	var rule Rule

	var (
		defaultErr       = errors.New("invalid rule")
		invalidMatchErr  = errors.New("invalid match")
		invalidTargetErr = errors.New("invalid target")
	)

	s := NewScanner(bufio.NewReader(strings.NewReader(line)))
	// should be -A
	if word, err := s.ReadWord(); err != nil || word != "-A" {
		return nil, defaultErr
	}
	// should be chain name
	chainName, err := s.ReadWord()
	if err != nil {
		return nil, defaultErr
	}
	chain, _ := p.ct.GetChain(chainName)
	rule.Table = p.ct
	rule.Chain = chain

	// remove comment and save
	re := regexp.MustCompile(`-m comment --comment \".*\" `)
	rule.Text = re.ReplaceAllString(line, "")

	var isNot bool
	for {
		word, err := s.ReadWord()
		if err != nil {
			return nil, defaultErr
		}
		if word == "" {
			break
		}
		switch word {
		case "!":
			isNot = true
		case "-p":
			val, _ := s.ReadWord()
			rule.Proto = Flag{IsNot: isNot, Key: "--protocol", Val: val}
			isNot = false
		case "-d":
			val, _ := s.ReadWord()
			rule.Dst = Flag{IsNot: isNot, Key: "--destination", Val: val}
			isNot = false
		case "-i":
			val, _ := s.ReadWord()
			rule.Ifce = Flag{IsNot: isNot, Key: "--interface", Val: val}
			isNot = false
		case "-s":
			val, _ := s.ReadWord()
			rule.Src = Flag{IsNot: isNot, Key: "--source", Val: val}
			isNot = false
		case "-m":
			match := Match{}
			if err := match.Parse(s); err != nil {
				return nil, invalidMatchErr
			}
			rule.AddMatch(&match)
		case "-j":
			rule.Target = &Target{}
			if err := rule.Target.Parse(s); err != nil {
				return nil, invalidTargetErr
			}
		}
	}
	return &rule, nil
}

func (p *Parser) Render(t string) ([]*OutChain, error) {
	// 预防前端使用循环引用链，导致递归死循环，这里设置一个超时时间
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)

	var formatChain func(string, string) (*OutChain, error)
	formatChain = func(tableName string, chainName string) (out *OutChain, err error) {
		if ctx.Err() != nil {
			return nil, errors.New("处理超时")
		}
		out = &OutChain{Name: tableName + ":" + chainName, Rules: []*OutRule{}}
		table, ok := p.tm[tableName]
		if !ok {
			return
		}
		chain, err := table.GetChain(chainName)
		if err != nil {
			return
		}
		out.Policy = chain.Policy
		out.Rules = make([]*OutRule, len(chain.Rules))
		for kk, rr := range chain.Rules {
			or := OutRule{Text: rr.Text, Comment: rr.Comment}
			out.Rules[kk] = &or
			matches := []string{}
			// -i, -d, -s, -p
			if rr.Ifce.Val != "" {
				matches = append(matches, rr.Ifce.String())
			}
			if rr.Dst.Val != "" {
				matches = append(matches, rr.Dst.String())
			}
			if rr.Src.Val != "" {
				matches = append(matches, rr.Src.String())
			}
			if rr.Proto.Val != "" {
				matches = append(matches, rr.Proto.String())
			}
			for _, mm := range rr.Matches {
				matches = append(matches, mm.String())
			}
			or.Matches = matches
			or.Target = rr.Target.String()
			if _, err := rr.Table.GetChain(rr.Target.Name); err == nil {
				fc, err := formatChain(rr.Table.Name, rr.Target.Name)
				if err != nil {
					return nil, err
				}
				or.Chains = []*OutChain{fc}
			}
		}
		return
	}

	conf := traffic[t]
	rootChains := make([]*OutChain, len(conf))
	for k, v := range conf {
		segs := strings.Split(v, ":")
		tableName, chainName := segs[0], segs[1]
		fc, err := formatChain(tableName, chainName)
		if err != nil {
			return nil, err
		}
		rootChains[k] = fc
	}

	return rootChains, nil
}

type OutChain struct {
	Name   string `json:"name"`
	Policy string `json:"policy"`

	Rules []*OutRule `json:"rules"`
}

type OutRule struct {
	Text    string      `json:"text"`
	Comment string      `json:"comment"`
	Matches []string    `json:"matches"`
	Target  string      `json:"target"`
	Chains  []*OutChain `json:"chain,omitempty"`
}

func NewParser(s *Scanner) *Parser {
	return &Parser{s: s, tm: make(map[string]*Table)}
}
