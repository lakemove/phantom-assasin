package main

import (
	"fmt"
	"log"
	"github.com/mitchellh/cli"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("hello gin")
	c := cli.NewCLI("leo", "1.0.0")
	c.Run()
	r := gin.Default()
	r.GET("/", handle1)
	r.Run(":8080")
}

func handle1(c *gin.Context) {
	c.String(200, "hello gin\n")
}
