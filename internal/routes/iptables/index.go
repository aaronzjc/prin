package iptables

import (
	"bufio"
	"prin/internal/util/req"
	ipt "prin/pkg/iptables"
	"strings"

	"github.com/gin-gonic/gin"
)

type Form struct {
	T    string `json:"type"`
	Data string `json:"data"`
}

func Beauty(c *gin.Context) {
	var r Form
	var err error
	if err = c.ShouldBindJSON(&r); err != nil {
		req.JSON(c, req.CodeError, "参数异常", nil)
		return
	}
	s := ipt.NewScanner(bufio.NewReader(strings.NewReader(r.Data)))
	parser := ipt.NewParser(s)
	err = parser.Parse()
	if err != nil {
		req.JSON(c, req.CodeError, "解析失败", nil)
		return
	}
	result, err := parser.Render(r.T)
	if err != nil {
		req.JSON(c, req.CodeError, err.Error(), nil)
		return
	}

	req.JSON(c, req.CodeSuccess, "success", result)
	return
}
