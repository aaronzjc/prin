package handler

import (
	"os"
	"prin/internal/app"
	"prin/internal/util/logger"
	"prin/internal/util/req"
	"prin/internal/util/tool"
	"strings"

	"github.com/gin-gonic/gin"
)

type Cert struct{}

type CertForm struct {
	Domains string `json:"domains"`
}

func (ctr *Cert) Generate(c *gin.Context) {
	var r CertForm
	var err error
	if err = c.ShouldBindJSON(&r); err != nil {
		req.JSON(c, req.CodeError, "参数异常", nil)
		return
	}

	if len(r.Domains) == 0 {
		req.JSON(c, req.CodeError, "域名为空", nil)
		return
	}

	domains := strings.Split(r.Domains, ",")

	caCertData, _ := os.ReadFile(app.RootPath + "/assets/ca.pem")
	caKeyData, _ := os.ReadFile(app.RootPath + "/assets/ca.key")
	caCerts, err := tool.ParseCertsPEM(caCertData)
	if err != nil {
		req.JSON(c, req.CodeError, err.Error(), nil)
		return
	}
	if len(caCerts) == 0 {
		req.JSON(c, req.CodeError, "无效的根证书", nil)
		return
	}
	caKey, err := tool.ParsePrivKey(caKeyData)
	if err != nil {
		req.JSON(c, req.CodeError, err.Error(), nil)
		return
	}
	dstKey, dstCert, err := tool.GenCertAndPrivkey(caKey, caCerts[0], domains)
	if err != nil {
		req.JSON(c, req.CodeError, err.Error(), nil)
		return
	}

	logger.Info("generate certificate success ! content = " + r.Domains)

	req.JSON(c, req.CodeSuccess, "success", map[string]string{
		"cert": string(dstCert),
		"key":  string(dstKey),
	})
}
