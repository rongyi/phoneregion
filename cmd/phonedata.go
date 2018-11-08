package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/rongyi/phoneregion"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Print("请输入您要查询的手机号码，不需要+86")
		os.Exit(1)
	}
	b, err := Asset("phone.dat")
	if err != nil {
		fmt.Println("no data found")
		os.Exit(1)
	}
	parser, err := phoneregion.NewParser(bytes.NewReader(b))
	pr, err := parser.Find(os.Args[1])
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	fmt.Println(pr.String())
}
