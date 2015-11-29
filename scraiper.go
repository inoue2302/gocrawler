package main

import (
	"./crawler"
)

func main() {
	url := "https://xxxx/"
	scraiper := crawler.NewScraiper(url)
	queueUrl := "https://xxxxx/"
	scraiper.SetDoc(queueUrl)
	scraiper.SetCompanyName()
	scraiper.SetTitle()
	scraiper.SetUnitPrice()
	scraiper.SetWorkSubInfo()
	scraiper.SetWorkInfo()
	scraiper.AddProject()
}
