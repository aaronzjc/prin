package iptable

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"prin/internal/util/logger"
	"prin/internal/util/req"
	"regexp"
	"strings"
)

type Form struct {
	T    string `json:"table"`
	Data string `json:"data"`
}

type Rule struct {
	Raw       string `json:"raw"`
	RawFormat string `json:"raw_format"`
	Seg       struct {
		Target string `json:"target"`
		Prot   string `json:"prot"`
		Opt    string `json:"opt"`
		Src    string `json:"source"`
		Dst    string `json:"destination"`
	} `json:"seg"`
}

type Chain struct {
	Summary string `json:"summary"`
	Rules   []Rule `json:"rules"`
}

type TableMap map[string]Chain

func convertTable(data string) (TableMap, error) {
	rows := strings.Split(data, "\n")
	tableMap := make(TableMap)
	r, _ := regexp.Compile("Chain ([0-9a-zA-Z\\-]+) .*")
	rr, _ := regexp.Compile("\\s+")
	chainName := ""
	for _, row := range rows {
		if row == "" {
			chainName = ""
			continue
		}
		if strings.HasPrefix(row, "Chain") {
			matches := r.FindStringSubmatch(row)
			chainName = matches[1]
			tableMap[chainName] = Chain{
				Summary: row,
				Rules:   []Rule{},
			}
			continue
		}
		if strings.HasPrefix(row, "target") {
			continue
		}

		rule := Rule{}
		rule.Raw = row
		segs := rr.Split(row, 5)
		if len(segs) != 5 {
			return nil, errors.New("unrecognize rule")
		}

		rule.Seg.Target = segs[0]
		rule.Seg.Prot = segs[1]
		rule.Seg.Opt = segs[2]
		rule.Seg.Src = segs[3]
		rule.Seg.Dst = segs[4]

		rule.RawFormat = fmt.Sprintf("Target=(%s) Prot=(%s) Opt=(%s) Src=(%s) Dst=(%s)", rule.Seg.Target, rule.Seg.Prot, rule.Seg.Opt, rule.Seg.Src, rule.Seg.Dst)

		chainData, ok := tableMap[chainName]
		if !ok {
			return nil, errors.New("no chain found before rules")
		}
		chainData.Rules = append(chainData.Rules, rule)
		tableMap[chainName] = chainData
	}
	return tableMap, nil
}

type Node struct {
	Name      string `json:"name"`
	Raw       string `json:"raw"`
	RawFormat string `json:"raw_format"`
	Children  []Node `json:"children"`
}

func BuildTree(nodes []Node, tableMap TableMap) {
	for k, c := range nodes {
		if c.Name == "" {
			return
		}
		chainData, ok := tableMap[c.Name]
		if !ok {
			continue
		}
		nodes[k].Raw = chainData.Summary
		for _, vv := range tableMap[c.Name].Rules {
			nodes[k].Children = append(nodes[k].Children, Node{
				Name:      vv.Seg.Target,
				RawFormat: vv.RawFormat,
			})
		}
		if len(nodes[k].Children) > 0 {
			BuildTree(nodes[k].Children, tableMap)
		}
	}
}

func Beauty(c *gin.Context) {
	var r Form
	var err error
	if err = c.ShouldBindJSON(&r); err != nil {
		req.JSON(c, req.CodeError, "参数异常", nil)
		return
	}

	tableMap, err := convertTable(r.Data)
	fmt.Println(tableMap)
	if err != nil {
		req.JSON(c, req.CodeError, err.Error(), nil)
		return
	}

	tbDefine := map[string][]string{
		"NAT": {"PREROUTING", "INPUT", "OUTPUT", "POSTROUTING"},
	}
	chains, ok := tbDefine[r.T]
	if !ok {
		req.JSON(c, req.CodeError, "未识别的表", nil)
		return
	}

	nodes := []Node{}
	for _, v := range chains {
		chainData, ok := tableMap[v]
		if !ok {
			continue
		}
		nodes = append(nodes, Node{
			Name:     v,
			Raw:      chainData.Summary,
			Children: nil,
		})
	}

	logger.Info("format iptable success !")

	BuildTree(nodes, tableMap)

	req.JSON(c, req.CodeSuccess, "success", nodes)
	return
}
