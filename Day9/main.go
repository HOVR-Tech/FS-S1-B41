package main

import (
	"Day9/connection"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var Data = map[string]interface{}{
	"Title": "Personal Web",
}

// Struct Blog
type Blog struct {
	Id		  	int
	Title	  	string
	Post_date 	time.Time
	Format_date string
	Author		string
	Content		string
}

var Blogs = []Blog{
	{
	Title: "Blog Satu",
	Author: "Hydrilla Fragrant",
	Content: "Berita Blog Satu",
	},
	{
	Title: "Blog Dua",
	Author: "Hydrilla Fragrant",
	Content: "Berita Blog Dua",
	},
}

// Struct Project
type Project struct {
	ID                  int
	ProjectName         string
	ProjectStartDate    string
	ProjectEndDate      string
	ProjectDuration     string
	ProjectDescription  string
	ProjectTechnologies []string
}

var ProjectList = []Project{}

func main() {
	route := mux.NewRouter()

	// Database Connect
	connection.DatabaseConnect()

	// route path folder untuk public
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	//routing
	route.HandleFunc("/hello", helloWorld).Methods("GET")
	route.HandleFunc("/", home).Methods("GET")
	
	//BLog 
	route.HandleFunc("/blog", blog).Methods("GET")
	route.HandleFunc("/blog-details/{index}", blogDetails).Methods("GET")
	route.HandleFunc("/form-blog", formAddBlog).Methods("GET")
	route.HandleFunc("/add-blog", addBlog).Methods("POST")
	route.HandleFunc("/delete-blog/{index}", deleteBlog).Methods("GET")
	
	// Project
	route.HandleFunc("/project", project).Methods("GET")
	route.HandleFunc("/project-details/{index}", projectDetails).Methods("GET")
	route.HandleFunc("/project/create", CreateProject).Methods("POST")
	route.HandleFunc("/update-project/{index}", updateProject).Methods("GET")
	route.HandleFunc("/delete-project/{index}", deleteProject).Methods("GET")


	// Contact
	route.HandleFunc("/contact", contact).Methods("GET")
	

	fmt.Println("Server running on port 5000")
	http.ListenAndServe("localhost:5000", route)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}


// INDEX
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


// BLOG
func blog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/blog.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	// var query = "SELECT id, title, content FROM blog"

	rows, _ := connection.Conn.Query(context.Background(), "SELECT id, title, content FROM blog")

	var result []Blog // array data

	for rows.Next() {
		var each = Blog{} //untuk manggil struct

		err := rows.Scan(&each.Id, &each.Title, &each.Content)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		each.Author = "Hydrilla Fragrant"
		// each.Format_date = each.Post_date.Format("2 January 2006")

		result = append(result, each)
	}

	fmt.Println(result)

	respData := map[string]interface{}{
		"Blogs": result,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, respData)
}

func blogDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/blog-details.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	var BlogDetails = Blog{}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	for i, data := range Blogs {
			if index == i {
					BlogDetails = Blog{
						Title	 : data.Title,
						Content	 : data.Content,
						Post_date: data.Post_date,
						Author	 : data.Author,
					}
			}		
	}


	data := map[string]interface{}{
		"Blog": BlogDetails,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

func formAddBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/add-blog.html")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func addBlog(w http.ResponseWriter, r *http.Request) {
	
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Title : " + r.PostForm.Get("inputTitle"))
	fmt.Println("Content : " + r.PostForm.Get("inputContent"))

	var title = r.PostForm.Get("inputTitle")
	var content = r.PostForm.Get("inputContent")

	var newBlog = Blog{
		Title:   title,
		Content: content,
		Author:  "Hydrilla Fragrant",
	}

	Blogs = append(Blogs, newBlog)
	

	http.Redirect(w, r, "/blog", http.StatusMovedPermanently)
}

func deleteBlog(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	Blogs = append(Blogs[:index], Blogs[index+1:]...)

	http.Redirect(w, r, "/blog", http.StatusFound)
}


// PROJECT
func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}
	
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func projectDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/project-details.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		var renderDetail = Project{}
		index, _ := strconv.Atoi(mux.Vars(r)["index"])

		for i, data := range ProjectList {
			if index == i {
				renderDetail = Project{
					ProjectName:         data.ProjectName,
					ProjectStartDate:    data.ProjectStartDate,
					ProjectEndDate:      data.ProjectEndDate,
					ProjectDuration:     data.ProjectDuration,
					ProjectDescription:  data.ProjectDescription,
					ProjectTechnologies: data.ProjectTechnologies,
				}
			}
		}
		data := map[string]interface{}{
			"renderDetail": renderDetail,
		}
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, data)
	}
}

func CreateProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	} else {
		projectName := r.PostForm.Get("project-name")
		projectStartDate := r.PostForm.Get("date-start")
		projectEndDate := r.PostForm.Get("date-end")
		projectDescription := r.PostForm.Get("project-description")
		projectUseNodeJS := r.PostForm.Get("nodejs")
		projectUseReactJS := r.PostForm.Get("reactjs")
		projectUseGolang := r.PostForm.Get("golang")
		projectUseTypeScript := r.PostForm.Get("typescript")

		var newProject = Project{
			ProjectName:         projectName,
			ProjectStartDate:    FormatDate(projectStartDate),
			ProjectEndDate:      FormatDate(projectEndDate),
			ProjectDuration:     GetDuration(projectStartDate, projectEndDate),
			ProjectDescription:  projectDescription,
			ProjectTechnologies: []string{projectUseNodeJS, projectUseReactJS, projectUseGolang, projectUseTypeScript},
		}

		fmt.Println(newProject)

		ProjectList = append(ProjectList, newProject)

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func updateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/update-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		var updateData = Project{}
		index, _ := strconv.Atoi(mux.Vars(r)["index"])

		for i, data := range ProjectList {
			if index == i {
				updateData = Project{
					ProjectName:         data.ProjectName,
					ProjectStartDate:    ReturnDate(data.ProjectStartDate),
					ProjectEndDate:      ReturnDate(data.ProjectEndDate),
					ProjectDescription:  data.ProjectDescription,
					ProjectTechnologies: data.ProjectTechnologies,
				}
				ProjectList = append(ProjectList[:index], ProjectList[index+1:]...)
			}
		}
		data := map[string]interface{}{
			"updateData": updateData,
		}
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, data)
	}
}

func deleteProject(w http.ResponseWriter, r *http.Request) {

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	ProjectList = append(ProjectList[:index], ProjectList[index+1:]...)

	http.Redirect(w, r, "/", http.StatusFound)
}



// ADDITIONAL FUNCTION
//DURATION
func GetDuration(startDate string, endDate string) string {

	layout := "2006-01-02"

	date1, _ := time.Parse(layout, startDate)
	date2, _ := time.Parse(layout, endDate)

	margin := date2.Sub(date1).Hours() / 24
	var duration string

	if margin > 30 {
		if (margin / 30) <= 1 {
			duration = "1 Month"
		} else {
			duration = strconv.Itoa(int(margin)/30) + " Months"
		}
	} else {
		if margin <= 1 {
			duration = "1 Day"
		} else {
			duration = strconv.Itoa(int(margin)) + " Days"
		}
	}

	return duration
}

// DATE
func FormatDate(InputDate string) string {

	layout := "2006-01-02"
	t, _ := time.Parse(layout, InputDate)

	Formated := t.Format("02 January 2006")

	return Formated
}

// RETURN DATE FORMAT
func ReturnDate(InputDate string) string {

	layout := "02 January 2006"
	t, _ := time.Parse(layout, InputDate)

	Formated := t.Format("2006-01-02")

	return Formated
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