package main

import (
	"com.snippetbox.aitu/internal/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, http.StatusOK, "home.tmpl", &templateData{
		Snippets: snippets,
	})

	//files := []string{
	//	"./ui/html/base.tmpl",
	//	"./ui/html/partials/nav.tmpl",
	//	"./ui/html/pages/home.tmpl",
	//}
	//
	//ts, err := template.ParseFiles(files...)
	//if err != nil {
	//	app.serverError(w, err)
	//	return
	//}
	//
	//data := &templateData{
	//	Snippets: snippets,
	//}
	//
	//err = ts.ExecuteTemplate(w, "base", data)
	//if err != nil {
	//	app.serverError(w, err)
	//}

	//for _, snippet := range snippets {
	//	fmt.Fprintf(w, "%+v\n", snippet)
	//}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, http.StatusOK, "view.tmpl", &templateData{
		Snippet: snippet,
	})

	//files := []string{
	//	"./ui/html/base.tmpl",
	//	"./ui/html/partials/nav.tmpl",
	//	"./ui/html/pages/view.tmpl",
	//}
	//
	//ts, err := template.ParseFiles(files...)
	//if err != nil {
	//	app.serverError(w, err)
	//	return
	//}
	//
	//data := &templateData{
	//	Snippet: snippet,
	//}
	//
	//err = ts.ExecuteTemplate(w, "base", data)
	//if err != nil {
	//	app.serverError(w, err)
	//}
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	title := ""
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)

}
