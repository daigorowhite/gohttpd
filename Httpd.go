package main

import ("net/http"
        // "fmt"
        "io/ioutil"
        "text/template"
        "os"
        "bytes"
        "strings"
         )

type Page struct{
	Title string
	Body []byte
}

type All struct{
	Pages string
}

const PagePostFix = ".txt"
const DataFolder = "wiki/"
const ViewFolder = "view/"

func viewHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[6:]
	p, _ := loadPage(title)
	t, _ := template.ParseFiles(ViewFolder + "view.html")
	t.Execute(w, p)
}

func editHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[6:]
	p, _ := loadPage(title)
	t, _ := template.ParseFiles(ViewFolder + "edit.html")
	t.Execute(w, p)
}

func toEditHandler(w http.ResponseWriter, r *http.Request){
	title := r.FormValue("title")
	http.Redirect(w, r, "/edit/" + title , http.StatusFound)
}

func saveHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[6:]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/" + title , http.StatusFound)
}

func allPageHandler(w http.ResponseWriter, r *http.Request){
	pages, _ := ioutil.ReadDir(DataFolder)	
	t, _ := template.ParseFiles(ViewFolder + "all.html")
	html_pages := translateHtml(pages)
	t.Execute(w, html_pages)
}

func translateHtml(pages []os.FileInfo) string {
	html_pages :=  ""
	t, _ := template.ParseFiles(ViewFolder + "_wiki_link.html")
	for _, page := range pages{
		buffer := bytes.NewBufferString("")
		t.Execute(buffer, strings.Replace(page.Name(),".txt","",-1))
		html_pages = html_pages + string(buffer.Bytes())
	}
	return html_pages
}

func (p *Page) save() error{
	filename := DataFolder + p.Title + PagePostFix
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error){
	filename := DataFolder + title + PagePostFix
	body,err := ioutil.ReadFile(filename)
	if err != nil {
		return &Page{Title: title}, err
	}
	return &Page{Title:title, Body: []byte(template.HTMLEscapeString(string(body)))}, nil
}

func main(){
	http.HandleFunc("/view/",viewHandler)
	http.HandleFunc("/edit/",editHandler)
	http.HandleFunc("/save/",saveHandler)
	http.HandleFunc("/toedit/",toEditHandler)
	http.HandleFunc("/",allPageHandler)
	http.ListenAndServe(":8080",nil)
}
