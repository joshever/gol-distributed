package gol

import "fmt"

func Handle(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
