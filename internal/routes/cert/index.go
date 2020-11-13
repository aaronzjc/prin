package cert

import (
	"prin/internal/util/logger"
	"prin/internal/util/req"

	"github.com/gin-gonic/gin"
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

	logger.Info("coder exec success ! content = " + r.Input)

	req.JSON(c, req.CodeSuccess, "success", map[string]string{
		"output": output,
	})
	return
}
