package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"prin/internal/util/req"

	"github.com/gin-gonic/gin"
)

type Stat struct{}

func (ctr *Stat) Online(c *gin.Context) {
	svc := os.Getenv("ONLINE_SVC")
	if svc == "" {
		req.JSON(c, req.CodeSuccess, "success", map[string]string{
			"count": "",
		})
		return
	}
	url := fmt.Sprintf("%s/online/%s", svc, "prin")
	resp, err := http.Get(url)
	if err != nil {
		req.JSON(c, req.CodeSuccess, "success", map[string]string{
			"count": "",
		})
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		req.JSON(c, req.CodeSuccess, "success", map[string]string{
			"count": "",
		})
		return
	}
	req.JSON(c, req.CodeSuccess, "success", map[string]string{
		"count": string(body),
	})
}
