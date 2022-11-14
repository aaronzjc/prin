package handler

import (
	"encoding/base64"
	"prin/internal/util/logger"
	"prin/internal/util/req"

	"github.com/gin-gonic/gin"
	qr "github.com/skip2/go-qrcode"
)

type Qrcode struct{}

type QrcodeForm struct {
	Content string `json:"content"`
}

func (ctr *Qrcode) Generate(c *gin.Context) {
	var r QrcodeForm
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
}
