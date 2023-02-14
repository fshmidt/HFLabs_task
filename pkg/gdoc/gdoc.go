package gdoc

import (
	"HFLabs/pkg/parsing"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/tanaikech/go-gdoctableapp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/docs/v1"
	"log"
	"net/http"
	"os"
)

func ServiceAccount(credentialFile string) *http.Client {
	b, err := os.ReadFile(credentialFile)
	if err != nil {
		log.Fatal(err)
	}
	var c = struct {
		Email      string `json:"client_email"`
		PrivateKey string `json:"private_key"`
	}{}
	json.Unmarshal(b, &c)
	config := &jwt.Config{
		Email:      c.Email,
		PrivateKey: []byte(c.PrivateKey),
		Scopes: []string{
			docs.DocumentsScope,
		},
		TokenURL: google.JWTTokenURL,
	}
	client := config.Client(oauth2.NoContext)
	return client
}

func objFromStruct(tabledata []parsing.TableRow) *gdoctableapp.CreateTableRequest {

	rows := len(tabledata)
	columns := len(tabledata[0].Cells)
	values := make([][]interface{}, rows)

	for i := 0; i < rows; i++ {
		cells := make([]interface{}, columns)
		for j := 0; j < columns; j++ {
			cells[j] = tabledata[i].Cells[j]
		}
		values[i] = cells
	}

	obj := &gdoctableapp.CreateTableRequest{
		Rows:    int64(rows),
		Columns: int64(columns),
		Index:   1,
		Values:  values,
	}
	return obj
}
func RewriteGdoc(tableData []parsing.TableRow) {

	// Чтение приватных значений из .env файла и загрузка их в окружение
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env var: %s", err.Error())
	}

	// Подключение к Google Documents
	documentID := os.Getenv("DOCUMENT_ID")
	client := ServiceAccount(os.Getenv("SERVICE_ACCOUNT_FILE"))

	// Очищаем доку от предыдущей версии таблицы
	g := gdoctableapp.New()
	tableIndex := 0
	res, err := g.Docs(documentID).TableIndex(tableIndex).DeleteTable().Do(client)
	if err != nil {
		log.Fatal(err)
	}

	// Создаем table request из структуры
	obj := objFromStruct(tableData)

	// Запись таблицы в гугл доку
	g = gdoctableapp.New()
	res, err = g.Docs(documentID).CreateTable(obj).Do(client)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Your table is in the document.", res.Message)
	fmt.Printf("You can check it at https://docs.google.com/document/d/%s\n", os.Getenv("DOCUMENT_ID"))
}
