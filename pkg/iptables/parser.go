package iptables

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strings"
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
	Comment string
	Matches []*Match
	Target  *Target
}

func (r *Rule) AddMatch(match *Match) {
	// if it is comment, we dont treat it as a match
	if match.Name == "comment" {
		r.Comment = match.Flags["--comment"].Val
		return
	}
	r.Matches = append(r.Matches, match)
}

type Flag struct {
	IsNot bool
	Val   string
}

type Match struct {
	Name  string
	Flags map[string]Flag
}

func (m *Match) String() string {
	var str string
	for k, v := range m.Flags {
		if v.IsNot {
			str = str + fmt.Sprintf("(NOT %s=%s)", k, v.Val)
		} else {
			str = str + fmt.Sprintf("(%s=%s)", k, v.Val)
		}
	}
	return m.Name + ":" + str
}

func (m *Match) Parse(s *Scanner) error {
	m.Flags = make(map[string]Flag)
	m.Name, _ = s.ReadWord()
	for {
		try, err := s.Peek(2)
		if err != nil {
			return errors.New("invalid match")
		}
		if try == "" {
			return nil
		}
		if try != "--" {
			return nil
		}
		flag := Flag{}
		if try == "! " {
			flag.IsNot = true
		}
		key, _ := s.ReadWord()
		if key == "--comment" {
			flag.Val, _ = s.ReadComment()
		} else {
			flag.Val, _ = s.ReadWord()
		}
		m.Flags[key] = flag
	}
}

type Target struct {
	Name  string
	Flags map[string]string
}

func (t *Target) String() string {
	str := []string{t.Name}
	for k, v := range t.Flags {
		if v != "" {
			k = k + "=" + v
		}
		str = append(str, k)
	}
	return strings.Join(str, " ")
}

func (t *Target) Parse(s *Scanner) error {
	t.Flags = make(map[string]string)
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
		key, _ := s.ReadWord()
		val, _ := s.ReadWord()
		t.Flags[key] = val
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

	// 去除注释
	re := regexp.MustCompile(`-m comment --comment \".*\" `)
	rule.Text = re.ReplaceAllString(line, "")

	for {
		word, err := s.ReadWord()
		if err != nil {
			return nil, defaultErr
		}
		if word == "" {
			break
		}
		var isNot bool
		switch word {
		case "!":
			isNot = true
		case "-p":
			val, _ := s.ReadWord()
			rule.Proto = Flag{IsNot: isNot, Val: val}
			isNot = false
		case "-d":
			val, _ := s.ReadWord()
			rule.Dst = Flag{IsNot: isNot, Val: val}
			isNot = false
		case "-s":
			val, _ := s.ReadWord()
			rule.Src = Flag{IsNot: isNot, Val: val}
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

func (p *Parser) Render(t string) []*OutChain {
	traffic := map[string][]string{
		"in":      {"raw:PREROUTING", "mangle:PREROUTING", "nat:PREROUTING", "mangle:INPUT", "filter:INPUT", "nat:INPUT"},
		"forward": {"raw:PREROUTING", "mangle:PREROUTING", "nat:PREROUTING", "mangle:FORWARD", "filter:FORWARD", "mangle:POSTROUTING", "nat:POSTROUTING"},
		"out":     {"raw:OUTPUT", "mangle:OUTPUT", "nat:OUTPUT", "filter:output", "mangle:POSTROUTING", "nat:POSTROUTING"},
	}
	var formatChain func(string, string) *OutChain
	formatChain = func(tableName string, chainName string) (out *OutChain) {
		out = &OutChain{Name: tableName + ":" + chainName, Rules: []OutRule{}}
		table, ok := p.tm[tableName]
		if !ok {
			return
		}
		chain, err := table.GetChain(chainName)
		if err != nil {
			return
		}
		out.Policy = chain.Policy
		out.Rules = make([]OutRule, len(chain.Rules))
		for kk, rr := range chain.Rules {
			out.Rules[kk] = OutRule{Text: rr.Text, Comment: rr.Comment}
			matches := []string{}
			for _, mm := range rr.Matches {
				matches = append(matches, mm.String())
			}
			out.Rules[kk].Matches = matches
			out.Rules[kk].Target = rr.Target.String()
			if _, err := rr.Table.GetChain(rr.Target.Name); err == nil {
				out.Rules[kk].Chains = []*OutChain{formatChain(rr.Table.Name, rr.Target.Name)}
			}
		}
		return
	}

	conf := traffic[t]
	rootChains := make([]*OutChain, len(conf))
	for k, v := range conf {
		segs := strings.Split(v, ":")
		tableName, chainName := segs[0], segs[1]
		rootChains[k] = formatChain(tableName, chainName)
	}

	return rootChains
}

type OutChain struct {
	Name   string `json:"name"`
	Policy string `json:"policy"`

	Rules []OutRule `json:"rules"`
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
