package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	val := &atomic.Int64{}

	e.GET("/", func(c echo.Context) error {
		valSync := val.Add(1)

		return c.String(http.StatusOK, fmt.Sprintf("Hello World %d", valSync))
	})

	log.Fatalln(e.Start(":1337"))
}
