package localserver

import (
	"log"
	"net/http"
	"text/template"
)

// Page : struct
type Page struct {
	Title string
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{Title: "Home Camera"}
	renderTemplate(w, "./view/index", &p)
}

// MyServer : Provide WEB GUI
func MyServer() {
	// Need to provide static files
	// "assets"内の静的ファイルを返すハンドラを作る
	staticHandler := http.FileServer(http.Dir("assets"))
	// ファイル名だけを得たいのでリクエストURIの中から "/assets/" を除外するハンドラを作成する
	splitHandler := http.StripPrefix("/assets/", staticHandler)
	// URLパス"/assets/"に対して、上記ハンドラを紐付ける
	// これにより localhost:8080/assets/test.jpg が呼ばれると、サーバ内のファイルからassets/test.jpgが呼ばれる
	http.Handle("/assets/", splitHandler)

	http.HandleFunc("/", viewHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
