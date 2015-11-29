package crawler

import (
	"../aws"
	"../models"
	"database/sql"
	"github.com/PuerkitoBio/goquery"
	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
	"time"
)

/*
 * スクレーパー構造体
 */
type Scraiper struct {
	url     string
	project models.Project
	doc     *goquery.Document
}

/**
 * 初期化処理
 * @param  url string
 * @return void
 */
func (self *Scraiper) initialize(url string) {
	self.url = url
	self.project = models.Project{}
}

/**
 * コンストラクタ
 * @param  url string
 * @return *Scraiper
 */
func NewScraiper(url string) *Scraiper {
	scraiper := new(Scraiper)
	scraiper.initialize(url)
	return scraiper

}

/**
 * ジョブナンバーを取得
 * @param  queueUrl string
 * @return string message
 */
func (self *Scraiper) getJobNumber(queueUrl string) string {
	sqsr := aws.NewSqsWrapper(queueUrl)
	sqsr.CreateReceiveParam()
	message := sqsr.ReceiveMessage()
	if message == "err" {
		panic("No Queue Message")
	}
	return message
}

/**
 * ドキュメントを取得
 * @param  number string
 * @return void
 */
func (self *Scraiper) getDetailDoc(number string) *goquery.Document {
	jobUrl := self.url + number
	doc, _ := goquery.NewDocument(jobUrl)
	return doc
}

/**
 * ドキュメントを構造体へセット
 * @param  queueUrl string
 * @return void
 */
func (self *Scraiper) SetDoc(queueUrl string) {
	number := self.getJobNumber(queueUrl)
	self.doc = self.getDetailDoc(number)
	self.project.Id, _ = strconv.Atoi(number)
}

/**
 * 会社名を構造体へセット
 * @return void
 */
func (self *Scraiper) SetCompanyName() {
	self.project.CompanyName = strings.TrimSpace(self.doc.Find(".company_an").Text())
}

/**
 * 案件名を構造体へセット
 * @return void
 */
func (self *Scraiper) SetTitle() {
	self.project.Title = strings.TrimSpace(self.doc.Find(".anken_syosai_top").Text())
}

/**
 * 単価を構造体へセット
 * @return void
 */
func (self *Scraiper) SetUnitPrice() {
	tankin := strings.TrimSpace(self.doc.Find(".tankin").Text())
	keitai := strings.TrimSpace(self.doc.Find(".keitai").Text())
	tankin = strings.TrimSpace(strings.Trim(tankin, keitai))
	self.project.Keitai = keitai
	self.project.Price = tankin
}

/**
 * 企業サブ情報をセット
 * @return void
 */
func (self *Scraiper) SetWorkSubInfo() {
	self.doc.Find(".table_cell_in_2").Each(func(index int, s *goquery.Selection) {
		if index == 0 {
			self.project.Nearest = strings.TrimSpace(s.Text())
		}
		if index == 1 {
			self.project.Term = strings.TrimSpace(s.Text())
		}
	})
}

/**
 * 企業情報をセット
 * @return void
 */
func (self *Scraiper) SetWorkInfo() {
	self.doc.Find(".table_cell_in").Each(func(index int, s *goquery.Selection) {
		if index == 0 {
			self.project.Content = strings.TrimSpace(s.Text())
		}
		if index == 1 {
			self.project.Environment = strings.TrimSpace(s.Text())
		}
		if index == 2 {
			self.project.NeedExperience = strings.TrimSpace(s.Text())
		}
		if index == 3 {
			self.project.ReceptionExprience = strings.TrimSpace(s.Text())
		}
		if index == 4 {
			self.project.Iquidation = strings.TrimSpace(s.Text())
		}
		if index == 5 {
			self.project.Negotiations = strings.TrimSpace(s.Text())
		}
	})
}

/**
 * DBへ登録もしくは更新
 * @return void
 */
func (self *Scraiper) AddProject() {
	nowTime := time.Now()
	self.project.CreatedAt = nowTime.Format("2006-01-02 15:04:05")
	self.project.UpdatedAt = nowTime.Format("2006-01-02 15:04:05")
	//コネクション作成
	db, err := sql.Open("mysql", "job_user:XXXXXXX@/freelance_job")
	if err != nil {
		panic(err.Error())
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	// 構造体とテーブルの紐付け
	dbmap.AddTableWithName(models.Project{}, "project")
	count, err := dbmap.SelectInt("select count(*) from project where id = " + strconv.Itoa(self.project.Id))

	var err2 error
	if count == 0 {
		err2 = dbmap.Insert(&self.project)
	} else {
		_, err2 = dbmap.Update(&self.project)
	}

	if err2 != nil {
		panic(err2.Error())
	}
}
