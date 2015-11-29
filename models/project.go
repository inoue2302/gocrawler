package models

type Project struct {
	Id                 int    `db:"id"`
	CompanyName        string `db:"company_name"`
	Title              string `db:"title"`
	Price              string `db:"price"`
	Keitai             string `db:"keitai"`
	Nearest            string `db:"nearest"`
	Term               string `db:"term"`
	Content            string `db:"content"`
	Environment        string `db:"environment"`
	NeedExperience     string `db:"needExperience"`
	ReceptionExprience string `db:"receptionExprience"`
	Iquidation         string `db:"iquidation"`
	Negotiations       string `db:"negotiations"`
	CreatedAt          string `db:"created_at"`
	UpdatedAt          string `db:"updated_at"`
}
