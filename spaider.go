package main

import (
	"./crawler"
	"time"
)

func main() {
	url := "https://xxxxxxx/"
	spaider := crawler.NewSpaider(url)

	pageNums := spaider.GetPageNum()
	for _, value := range pageNums {
		projectNumbers := spaider.GetProjectNumber(value)
		for _, number := range projectNumbers {
			spaider.AddNumberToQueue(number, "https://xxxxx/")
		}
		time.Sleep(time.Second)
	}
}
