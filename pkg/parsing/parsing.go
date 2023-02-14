package parsing

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

type TableRow struct {
	Cells []string
}

func ParseTable(url string) []TableRow {

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Поиск таблицы на странице
	table := doc.Find("table")

	var tableData []TableRow

	// вычитывание таблицы из selection в структуру
	table.Find("tr").Each(func(i int, tr *goquery.Selection) {
		var row TableRow

		// Перебор ячеек заголовков
		tr.Find("th").Each(func(x int, th *goquery.Selection) {
			cell := th.Text()
			row.Cells = append(row.Cells, cell)
		})

		// Перебор ячеек строки
		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			cell := td.Text()
			row.Cells = append(row.Cells, cell)
		})

		// Добавление строки в массив
		tableData = append(tableData, row)
	})
	return tableData
}
