package qrcode

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	qr "github.com/skip2/go-qrcode"
	"prin/internal/util/logger"
	"prin/internal/util/req"
)

type Form struct {
	Content string `json:"content"`
}

func Generate(c *gin.Context) {
	var r Form
	var err error
	if err = c.ShouldBindJSON(&r); err != nil {
		req.JSON(c, req.CodeError, "参数异常", nil)
		return
	}

	var pic []byte
	pic, err = qr.Encode(r.Content, qr.Medium, 512)
	if err != nil {
		logger.Error("qrcode encode failed ! content = " + r.Content)
		req.JSON(c, req.CodeError, "二维码生成失败", nil)
		return
	}

	logger.Info("qrcode generate success ! content = " + r.Content)

	data := base64.StdEncoding.EncodeToString(pic)

	req.JSON(c, req.CodeSuccess, "success", map[string]string{
		"qrcode": data,
	})
	return
}
