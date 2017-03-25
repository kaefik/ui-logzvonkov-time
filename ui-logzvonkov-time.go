package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/auth"
	"github.com/martini-contrib/render"

	//	"image"
	//	_ "image/gif"
	//	"image/jpeg"
	//	_ "image/png"
)

////------------ Объявление типов и глобальных переменных

var (
	hd             string
	user           string
	nameConfigFile string // имя конфигурационного файла
)

var (
	tekuser     string // текущий пользователь который задает условия на срабатывания
	pathcfg     string // адрес где находятся папки пользователей, если пустая строка, то текущая папка
	pathcfguser string
)

type page struct {
	Title  string
	Msg    string
	Msg2   string
	TekUsr string
}

// структура справочника телефонов менеджеров
type DataConfigFile struct {
	Numtel          string // номер телефона
	Fio_man         string // ФИО менеджера
	Fio_rg          string // ФИО РГ
	Planresultkolzv string // плановое кол-во результативных звоноков
	Plankolvstrech  string // плановое количество встреч
	Nstr            int    // номер строки который редактируется или удаляется, если -1 то новое условие
}

var (
	dataConfigFile []DataConfigFile // структура данных конфиг файла
)

//------------ END Объявление типов и глобальных переменных

// сохранить в новый файл
func SaveNewstrtofile(namef string, str string) int {
	file, err := os.Create(namef)
	if err != nil {
		// handle the error here
		return -1
	}
	defer file.Close()

	file.WriteString(str)
	return 0
}

// сохранить файл перезаписью
func Savestrtofile(namef string, str string) int {
	file, err := os.OpenFile(namef, os.O_CREATE|os.O_WRONLY, 0776)
	if err != nil {
		// handle the error here
		return -1
	}
	defer file.Close()

	file.WriteString(str)
	return 0
}

// сохранение данных из dataConfigFile в файл с именем namef
func saveToFileCfg(namef string, tt []DataConfigFile) {
	str := ""
	for _, t := range tt {
		str += t.Numtel + ";" + t.Fio_man + ";" + t.Fio_rg + ";" + t.Planresultkolzv + ";" + t.Plankolvstrech + "\n"
	}

	fmt.Println(tt)
	Savestrtofile(namef, str)
}

// чтение файла с именем namef и возвращение содержимое файла, иначе текст ошибки
func readfiletxt(namef string) string {
	file, err := os.Open(namef)
	if err != nil {
		return "handle the error here"
	}
	defer file.Close()
	// get the file size
	stat, err := file.Stat()
	if err != nil {
		return "error here"
	}
	// read the file
	bs := make([]byte, stat.Size())
	_, err = file.Read(bs)
	if err != nil {
		return "error here"
	}
	return string(bs)
}

// чтение файла конфига логов звонков
func readConfigFile(nameFile string) []DataConfigFile {
	dataConfigFile := make([]DataConfigFile, 0)
	s := readfiletxt(nameFile)
	ss := strings.Split(s, "\n")
	for k, v := range ss {
		ts := strings.Split(v, ";")

		if len(ts) >= 5 {
			dataConfigFile = append(dataConfigFile, DataConfigFile{Numtel: ts[0], Fio_man: ts[1], Fio_rg: ts[2], Planresultkolzv: ts[3], Plankolvstrech: ts[4], Nstr: k})
		}
	}
	return dataConfigFile
}

// обработчик начальной страницы
func indexHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request) {
	rr.HTML(200, "index", &page{Title: "Работа с конфиг файлом лога звонков", Msg: "Начальная страница", TekUsr: "Текущий пользователь: " + string(user)})
}

// обработка редактирования строки конфига
func EditConfigFileHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request, params martini.Params) {
	nstr, _ := strconv.Atoi(params["nstr"])

	tt := dataConfigFile[nstr]
	//	fmt.Println(tt)
	rr.HTML(200, "edit", &tt)
}

// просмотр конфига
func ViewConfigFileHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request) {

	rr.HTML(200, "view", &dataConfigFile)
}

// сохранение измененных данных перезаписью файла
func SaveConfigFileHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request, params martini.Params) {
	nstr, _ := strconv.Atoi(params["nstr"])
	planresultkolzv := r.FormValue("planresultkolzv")
	plankolvstrech := r.FormValue("plankolvstrech")
	dataConfigFile[nstr].Plankolvstrech = plankolvstrech
	dataConfigFile[nstr].Planresultkolzv = planresultkolzv

	saveToFileCfg(nameConfigFile, dataConfigFile)

	rr.HTML(200, "save", "")
}

func authFunc(username, password string) bool {
	return (auth.SecureCompare(username, "admin") && auth.SecureCompare(password, "qwe123!!"))
}

// функция парсинга аргументов программы
func parse_args() bool {
	flag.StringVar(&hd, "hd", "", "Рабочая папка где нах-ся папки пользователей для сохранения ")
	flag.StringVar(&user, "user", "", "Рабочая папка где нах-ся папки пользователей для сохранения ")
	flag.Parse()
	pathcfg = hd
	if user == "" {
		tekuser = "testuser"
	} else {
		tekuser = user
	}
	return true
}

func main() {

	nameConfigFile = "list-num-tel.cfg"
	dataConfigFile = readConfigFile(nameConfigFile)

	fmt.Println(nameConfigFile)
	m := martini.Classic()

	if !parse_args() {
		return
	}

	if pathcfg == "" {
		pathcfguser = ""
	} else {
		pathcfguser = pathcfg + string(os.PathSeparator)
	}

	m.Use(render.Renderer(render.Options{
		Directory:  "templates", // Specify what path to load the templates from.
		Layout:     "layout",    // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"}}))

	m.Use(auth.BasicFunc(authFunc))
	m.Get("/edit/:nstr", EditConfigFileHandler)
	m.Post("/save/:nstr", SaveConfigFileHandler)
	m.Get("/view", ViewConfigFileHandler)
	m.Get("/", indexHandler)
	m.RunOnAddr(":9999")

}
