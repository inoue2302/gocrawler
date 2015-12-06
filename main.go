package main

import (
	"./crawler"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

/**
 * スレイパー実行
 * @return void
 */
func scraiper() {
	url := "https://xxxx.xx/"
	scraiper := crawler.NewScraiper(url)
	queueUrl := "https://sqs.ap-northeast-1.amazonaws.com/xxxx/xxxx"
	scraiper.SetDoc(queueUrl)
	scraiper.SetCompanyName()
	scraiper.SetTitle()
	scraiper.SetUnitPrice()
	scraiper.SetWorkSubInfo()
	scraiper.SetWorkInfo()
	scraiper.AddProject()
}

func main() {
	send := make(chan int)
	quit := make(chan bool)
	workerquit := make(chan bool)

	//メモリ監視
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	go func() {
	loop:
		for {
			select {
			case <-quit:
				workerquit <- true
				break loop
			case <-send:
				scraiper()
			}
		}
	}()

	go func() {
		url := "xxxx"
		spaider := crawler.NewSpaider(url)

		pageNums := spaider.GetPageNum()
		for _, value := range pageNums {
			projectNumbers := spaider.GetProjectNumber(value)
			for _, number := range projectNumbers {
				send <- spaider.AddNumberToQueue(number, "https://sqs.ap-northeast-1.amazonaws.com/xxx/xxxxx")
				time.Sleep(time.Second)
			}
		}
		fmt.Println("finish")
		quit <- true
	}()
	<-workerquit
}
