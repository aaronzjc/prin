package handler

import (
	"encoding/base64"
	"net/url"
	"prin/internal/util/logger"
	"prin/internal/util/req"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Coder struct{}

type CoderForm struct {
	Input string `json:"input"`
	Type  string `json:"type"`
}

func (ctr *Coder) Decode(c *gin.Context) {
	var r CoderForm
	var err error
	if err = c.ShouldBindJSON(&r); err != nil {
		req.JSON(c, req.CodeError, "参数异常", nil)
		return
	}

	var output string

	switch r.Type {
	case "urlencode":
		output = url.QueryEscape(r.Input)
	case "urldecode":
		output, err = url.QueryUnescape(r.Input)
		if err != nil {
			req.JSON(c, req.CodeError, "转换失败", nil)
		}
	case "unicodeencode":
		conv := strconv.QuoteToASCII(r.Input)
		output = conv[1 : len(conv)-1]
	case "unicodedecode":
		output, err = strconv.Unquote(strings.Replace(strconv.Quote(r.Input), `\\u`, `\u`, -1))
		if err != nil {
			req.JSON(c, req.CodeError, "转换失败", nil)
		}
	case "base64encode":
		output = base64.StdEncoding.EncodeToString([]byte(r.Input))
	case "base64decode":
		conv, err := base64.StdEncoding.DecodeString(r.Input)
		if err != nil {
			req.JSON(c, req.CodeError, "转换失败", nil)
		}
		output = string(conv)
	default:
		req.JSON(c, req.CodeError, "未识别的操作", nil)
	}

	logger.Info("coder exec success ! content = " + r.Input)

	req.JSON(c, req.CodeSuccess, "success", map[string]string{
		"output": output,
	})
}
