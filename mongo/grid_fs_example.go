package mongo

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ExampleWriteExcelToGridFS() {
	fmt.Println("write excel to grid fs")
	orders := []map[string]interface{}{
		map[string]interface{}{
			"sn":       "aaa",
			"pay_time": "2019-01-01",
		},
		map[string]interface{}{
			"sn":       "bbb",
			"pay_time": "2019-02-01",
		},
		map[string]interface{}{
			"sn":       "ccc",
			"pay_time": "2019-03-01",
		},
	}

	err := WriteExcelToGridFSWithMap(orders, "ps_order", "download")
	if err != nil {
		panic(err)
	}
}

func ExampleReadGridFSToExcel() {
	fmt.Println("read excel to grid fs")
	err := Init(DefaultConfig())
	if err != nil {
		panic(err)
	}

	server := gin.Default()
	server.GET("/download", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusOK)
		c.Header("Content-Disposition", "attachment; filename=order.xlsx")
		c.Header("Content-Type", "application/text/plain")
		err = GridFSRead("ps_order", "download", c.Writer)
		if err != nil {
			panic(err)
		}
	})
	server.Run(":8899")
}
