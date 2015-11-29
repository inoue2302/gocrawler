package crawler

import (
	"../aws"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

/*
 * スパイダー構造体
 */
type Spaider struct {
	url string
}

/**
 * 初期化処理
 * @param  url string
 * @return void
 */
func (self *Spaider) initialize(url string) {
	self.url = url
}

/**
 * コンストラクタ
 * @param url string
 * @return *spaider
 */
func NewSpaider(url string) *Spaider {
	spaider := new(Spaider)
	spaider.initialize(url)
	return spaider
}

/**
 * ページャーの番号取得
 * @return pageNum slice
 */
func (self *Spaider) GetPageNum() []string {
	var pageNums []string
	doc, _ := goquery.NewDocument(self.url)
	doc.Find(".pagination").Each(func(_ int, s *goquery.Selection) {
		s.Find(".active").Each(func(_ int, s *goquery.Selection) {
			_, err := strconv.Atoi(s.Text())
			if err == nil {
				pageNums = append(pageNums, s.Text())
			}
		})
		s.Find("li>a").Each(func(_ int, s *goquery.Selection) {
			_, err := strconv.Atoi(s.Text())
			if err == nil {
				pageNums = append(pageNums, s.Text())
			}
		})
	})
	return pageNums
}

/**
 * プロジェクト採番番号を取得
 * @pram pageNum string
 * @rerturn projectUrls
 */
func (self *Spaider) GetProjectNumber(pageNum string) []string {
	var projectUrls []string
	pageUrl := self.url + "&page=" + pageNum
	doc, _ := goquery.NewDocument(pageUrl)
	doc.Find(".pro_title__title").Each(func(_ int, s *goquery.Selection) {
		url, exists := s.Find("a").Attr("href")
		if exists {
			projectUrls = append(projectUrls, strings.Split(url, "/")[6])
		}
	})
	return projectUrls
}

func (self *Spaider) AddNumberToQueue(projectNumber string, queueUrl string) {
	sqsr := aws.NewSqsWrapper(queueUrl)
	sqsr.CreateSendParam(projectNumber)
	sqsr.SendMesage()
}
