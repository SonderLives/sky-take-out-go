package req

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// ParsePage 从 Query 参数中解析分页参数，page 最小为 1，size 范围 [1, 100]
func ParsePage(c *gin.Context) (page, size int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ = strconv.Atoi(c.DefaultQuery("size", "10"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}
	return
}
