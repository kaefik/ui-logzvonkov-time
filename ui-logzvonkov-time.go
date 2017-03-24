package main

import (
	"flag"
	"fmt"
	//	"html/template"
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
	Nstr            string // номер строки который редактируется или удаляется, если -1 то новое условие
}

var (
	dataConfigFile []DataConfigFile
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

// сохранить файл
func Savestrtofile(namef string, str string) int {
	file, err := os.OpenFile(namef, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0776)
	if err != nil {
		// handle the error here
		return -1
	}
	defer file.Close()

	file.WriteString(str)
	return 0
}

//сохранение данных из TTasker товара в файл с именем namef
//func savetofilecfg(namef string, t TTasker) {
//	str := t.Url + ";" + t.Uslovie + ";" + t.Price + ";" + "\n"
//	Savestrtofile(namef, str)
//}

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

func indexHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request) {
	rr.HTML(200, "index", &page{Title: "Йоу Начало", Msg: "Начальная страница", TekUsr: "Текущий пользователь: " + string(user)})
}

// обработка редактирования задания
func EditConfigFileHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request, params martini.Params) {
	nstr, _ := strconv.Atoi(params["nstr"])
	fmt.Println("nstr = ", nstr)
	fmt.Println(dataConfigFile)
	//	tt := dataConfigFile[nstr]
	//	rr.HTML(200, "edit", &tt)
}

// просмотр заданий выбранного магазина
func ViewConfigFileHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request) {
	s := readfiletxt(nameConfigFile)
	ss := strings.Split(s, "\n")
	for _, v := range ss {
		ts := strings.Split(v, ";")

		if len(ts) >= 5 {
			dataConfigFile = append(dataConfigFile, DataConfigFile{Numtel: ts[0], Fio_man: ts[1], Fio_rg: ts[2], Planresultkolzv: ts[3], Plankolvstrech: ts[4]})
		}
	}
	rr.HTML(200, "view", &dataConfigFile)
}

func ExecHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request, params martini.Params) {
	//	var tt TTasker
	//	nstr, _ := strconv.Atoi(params["nstr"])
	//	//	shop := params["shop"]
	//	shop := r.FormValue("shop")

	//	tt.Shop = r.FormValue("shop")
	//	tt.Url = r.FormValue("surl")
	//	tt.Uslovie = r.FormValue("uslovie")
	//	tt.Price = r.FormValue("schislo")

	//	if _, err := os.Stat(pathcfguser + string(user)); os.IsNotExist(err) {
	//		os.Mkdir(pathcfguser+string(user), 0776)
	//	}
	//	namef := pathcfguser + string(user) + string(os.PathSeparator) + shop + "-url.cfg"

	//	if nstr == -1 {
	//		savetofilecfg(namef, tt)
	//	} else {
	//		s := readfiletxt(namef)
	//		ss := strings.Split(s, "\n")
	//		if nstr <= (len(ss) - 1) {
	//			ss[nstr] = tt.Url + ";" + tt.Uslovie + ";" + tt.Price + ";"
	//		}
	//		str := ""
	//		for _, v := range ss {
	//			if v != "" {
	//				str += v + "\n"
	//			}
	//		}
	//		SaveNewstrtofile(namef, str)
	//	}

	//	ss1 := "Введенное условие для магазина " + shop
	//	ss := tt.Url + "   " + tt.Uslovie + " " + tt.Price
	//	rr.HTML(200, "exec", &page{Title: "Введенное условие для магазина " + shop, Msg: ss, Msg2: ss1})
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

	dataConfigFile = make([]DataConfigFile, 0)

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
	m.Get("/", ViewConfigFileHandler)
	m.RunOnAddr(":9999")

}

//	m.Get("/addtask", AddTaskHandler)
//	m.Get("/edit/:shop/:nstr", EditTaskHandler)
//	m.Get("/del/:shop/:nstr", DelTaskHandler)
//	m.Post("/exec/:shop/:nstr", ExecHandler)
//	m.Get("/clickview", clickViewTaskHandler)
