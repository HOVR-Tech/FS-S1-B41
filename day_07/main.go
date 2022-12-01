package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	// ROUTING PATH TO PUBLIC
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// ROUTING HTML
	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")

	// PROJECT
	route.HandleFunc("/form-add-project", FormAddProject).Methods("GET")
	// route.HandleFunc("/add-project", AddProject).Methods("POST")

	fmt.Println("Server running on port 5000")
	http.ListenAndServe("localhost:5000", route)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/contact.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func FormAddProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/form-add-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	tmpl.Execute(w, nil)
}

func AddProject(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}  
	
	// GET DATA FROM INPUT FORM
<<<<<<< HEAD
	fmt.Println("Project Name : " + r.PostForm.Get("project-name"))
	fmt.Println("Start Date : " + r.PostForm.Get("start-date"))
	fmt.Println("End Date : " + r.PostForm.Get("end-date"))
	fmt.Println("Description : " + r.PostForm.Get("project-description"))
	fmt.Println("Use Node JS " + r.PostForm.Get("nodejs"))
	fmt.Println("Use React JS " + r.PostForm.Get("reactjs"))
	fmt.Println("Use Next JS " + r.PostForm.Get("nextjs"))
	fmt.Println("Use Typescript " + r.PostForm.Get("typescript"))
=======
	fmt.Println("Project Name : " + r.ParseForm.Get("project-name"))
	fmt.Println("Start Date : " + r.ParseForm.Get("start-date"))
	fmt.Println("End Date : " + r.ParseForm.Get("end-date"))
	fmt.Println("Description : " + r.ParseForm.Get("project-description"))
	fmt.Println("Use Node JS " + r.ParseForm.Get("nodejs"))
	fmt.Println("Use React JS " + r.ParseForm.Get("reactjs"))
	fmt.Println("Use Next JS " + r.ParseForm.Get("nextjs"))
	fmt.Println("Use Typescript " + r.ParseForm.Get("typescript"))
>>>>>>> 182be16d2ed89514c218774ce24ab23d7f98cccc


	http.Redirect(w, r, "/form-add-project", http.StatusMovedPermanently)
}

