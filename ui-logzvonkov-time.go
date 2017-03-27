package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
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

type sortDataConfigFile []DataConfigFile

//------------ END Объявление типов и глобальных переменных

// ---- функции сортировки по полю ФИО РГ
func (s sortDataConfigFile) Less(i, j int) bool {
	return s[i].Fio_rg < s[j].Fio_rg
}

func (s sortDataConfigFile) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortDataConfigFile) Len() int {
	return len(s)
}

// ---- END функции сортировки по полю ФИО РГ

// удаляем в массиве позицию i
func removeItemFromDataConfigFile(s []DataConfigFile, i int) []DataConfigFile {
	//	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	//	return s[:len(s)-1]
	res := make([]DataConfigFile, 0)
	for k, v := range s {
		if k != i {
			res = append(res, v)
		}
	}
	return res
}

// сортировка по ФИО РГ
func sortFioRgDataConfigFile(data []DataConfigFile) []DataConfigFile {
	sort.Sort(sortDataConfigFile(data))
	data = refreshIndexInNstr(data)
	return data
}

// обновление индекса по полю строки Nstr
func refreshIndexInNstr(data []DataConfigFile) []DataConfigFile {
	for k, _ := range data {
		data[k].Nstr = k
	}
	return data
}

// сохранить файл перезаписью
func SaveStrToFile(namef string, str string) int {
	file, err := os.Create(namef)
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
	//	fmt.Println(tt)
	SaveStrToFile(namef, str)
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
	rr.HTML(200, "edit", &tt)
}

// просмотр конфига
func ViewConfigFileHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request) {
	rr.HTML(200, "view", &dataConfigFile)
}

func printDataConfigFile(d []DataConfigFile) {
	fmt.Println("Длина массива: ", len(d))
	for _, v := range d {
		fmt.Println(v)
	}
}

// сохранение измененных данных перезаписью файла
func SaveConfigFileHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request, params martini.Params) {
	nstr, _ := strconv.Atoi(params["nstr"])
	numtel := r.FormValue("numtel")
	fioman := r.FormValue("fioman")
	fiorg := r.FormValue("fiorg")
	planresultkolzv := r.FormValue("planresultkolzv")
	plankolvstrech := r.FormValue("plankolvstrech")

	fmt.Println("Исходный масив")
	printDataConfigFile(dataConfigFile)

	if nstr != -1 {
		dataConfigFile[nstr].Numtel = numtel
		dataConfigFile[nstr].Fio_man = fioman
		dataConfigFile[nstr].Fio_rg = fiorg
		dataConfigFile[nstr].Plankolvstrech = plankolvstrech
		dataConfigFile[nstr].Planresultkolzv = planresultkolzv
	} else {
		dataConfigFile = append(dataConfigFile, DataConfigFile{Numtel: numtel, Fio_man: fioman, Fio_rg: fiorg, Plankolvstrech: plankolvstrech, Planresultkolzv: planresultkolzv})
	}

	fmt.Println("Измененный масив")
	printDataConfigFile(dataConfigFile)

	dataConfigFile = sortFioRgDataConfigFile(dataConfigFile)

	fmt.Println("Отсортированный масив")
	printDataConfigFile(dataConfigFile)

	saveToFileCfg(nameConfigFile, dataConfigFile)
	rr.HTML(200, "save", "")
}

// удаление выбранногоэлемента
func DelConfigFileHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request, params martini.Params) {
	nstr, _ := strconv.Atoi(params["nstr"])
	dataConfigFile = removeItemFromDataConfigFile(dataConfigFile, nstr)
	//	fmt.Println("измененные:  ", dataConfigFile)
	saveToFileCfg(nameConfigFile, dataConfigFile)
	rr.HTML(200, "save", "")
}

// добавление нового элемента
func AddConfigFileHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request, params martini.Params) {
	rr.HTML(200, "add", "")
}

func authFunc(username, password string) bool {
	return (auth.SecureCompare(username, "admin") && auth.SecureCompare(password, "qwe123!!"))
}

// функция парсинга аргументов программы
func parseArgs() bool {
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

	if !parseArgs() {
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
	m.Get("/del/:nstr", DelConfigFileHandler)
	m.Post("/save/:nstr", SaveConfigFileHandler)
	m.Get("/add", AddConfigFileHandler)
	m.Get("/view", ViewConfigFileHandler)
	m.Get("/", indexHandler)
	m.RunOnAddr(":9999")

}
