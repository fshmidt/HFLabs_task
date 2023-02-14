package main

import (
	"HFLabs/pkg/gdoc"
	"HFLabs/pkg/parsing"
)

const url = "https://confluence.hflabs.ru/pages/viewpage.action?pageId=1181220999"

func main() {

	// Получение таблицы в виде структуры
	tableData := parsing.ParseTable(url)

	// Перезапись таблицы в гугл доке
	gdoc.RewriteGdoc(tableData)
}
