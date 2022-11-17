package main

import (
	"Day10/connection"
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
	// {
	// Title: "Blog Satu",
	// Author: "Hydrilla Fragrant",
	// Content: "Berita Blog Satu",
	// },
	// {
	// Title: "Blog Dua",
	// Author: "Hydrilla Fragrant",
	// Content: "Berita Blog Dua",
	// },
}

// Struct Project
type Project struct {
	ID                     int
	ProjectName            string
	ProjectStartDate       time.Time
	ProjectEndDate         time.Time
	ProjectStartDateString string
	ProjectEndDateString   string
	ProjectDuration        string
	ProjectDescription     string
	ProjectTechnologies    []string
	ProjectImage           string
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
	route.HandleFunc("/blog-details/{id}", blogDetails).Methods("GET")
	route.HandleFunc("/form-blog", formAddBlog).Methods("GET")
	route.HandleFunc("/add-blog", addBlog).Methods("POST")
	route.HandleFunc("/delete-blog/{id}", deleteBlog).Methods("GET")
	
	// Project
	route.HandleFunc("/project-details/{id}", projectDetails).Methods("GET")
	route.HandleFunc("/form-add-project", formAddProject).Methods("GET")
	route.HandleFunc("/add-project", addProject).Methods("POST")
	route.HandleFunc("/update-project/{id}", updateProject).Methods("POST")
	route.HandleFunc("/delete-project/{id}", deleteProject).Methods("GET")


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

	tmpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		var renderData []Project
		item := Project{}
		// Database Connection
		rows, _ := connection.Conn.Query(context.Background(), `SELECT "ID", "ProjectName", "ProjectStartDate", "ProjectEndDate", "ProjectDescription", "ProjectTechnologies", "ProjectImage" FROM public.project`)
		
		for rows.Next() {
			// Connect Struct
			err := rows.Scan(&item.ID, &item.ProjectName, &item.ProjectStartDate, &item.ProjectEndDate, &item.ProjectDescription, &item.ProjectTechnologies, &item.ProjectImage)
			// ERROR HANDLING
			if err != nil {
				fmt.Println(err.Error())
				return
			} else {
				// PARSING DATE
				item := Project{
					ID:          		 item.ID,
					ProjectName:  		 item.ProjectName,
					ProjectDuration:     GetDuration(item.ProjectStartDate, item.ProjectEndDate),
					ProjectDescription:  item.ProjectDescription,
					ProjectTechnologies: item.ProjectTechnologies,
					ProjectImage:        item.ProjectImage,
				}
				renderData = append(renderData, item)
			}
		}
		response := map[string]interface{}{
			"renderData": renderData,
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, response)
	}
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

	rows, _ := connection.Conn.Query(context.Background(), `SELECT id, title, content, post_date, author FROM blog`)

	var result []Blog // array data

	for rows.Next() {
		var each = Blog{} //untuk manggil struct

		err := rows.Scan(&each.Id, &each.Title, &each.Content, &each.Post_date, &each.Author)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// each.Author = "Hydrilla Fragrant"
		each.Format_date = each.Post_date.Format("2 January 2006")

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

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var tmpl, err = template.ParseFiles("views/blog-details.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	var BlogDetails = Blog{}

	// index, _ := strconv.Atoi(mux.Vars(r)["index"])

	// for i, data := range Blogs {
	// 		if index == i {
	// 				BlogDetails = Blog{
	// 					Title	 : data.Title,
	// 					Content	 : data.Content,
	// 					Post_date: data.Post_date,
	// 					Author	 : data.Author,
	// 				}
	// 		}		
	// }

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, title, content, post_date, author FROM blog WHERE id=$1", id).Scan(
		&BlogDetails.Id, &BlogDetails.Title, &BlogDetails.Content, &BlogDetails.Post_date, &BlogDetails.Author,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}

	// BlogDetails.Author = "Hydrilla Fragrant"
	BlogDetails.Format_date = BlogDetails.Post_date.Format("2 January 2006")

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
	var author = r.PostForm.Get("inputAuthor")

	// var newBlog = Blog{
	// 	Title:   title,
	// 	Content: content,
	// 	Author:  "Hydrilla Fragrant",
	// }

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO blog(title, content, image, author) VALUES ($1, $2, 'images.png', $3)", title, content, author)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("message : " + err.Error()))
			return
		}

	// Blogs = append(Blogs, newBlog)
	

	http.Redirect(w, r, "/blog", http.StatusMovedPermanently)
}

func deleteBlog(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// Blogs = append(Blogs[:index], Blogs[index+1:]...)

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM blog WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/blog", http.StatusFound)
}

// PROJECT
func projectDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	
	tmpl, err := template.ParseFiles("views/project-details.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
		} else {
			ID, _ := strconv.Atoi(mux.Vars(r)["id"])
			renderDetails := Project{}

			// GET ID FROM DATABASE
			err := connection.Conn.QueryRow(context.Background(), `SELECT "ID", "ProjectName", "ProjectStartDate", "ProjectEndDate", "ProjectDescription", "ProjectTechnologies", "ProjectImage"
			FROM public.project WHERE "ID" = $1`, ID).Scan(&renderDetails.ID, &renderDetails.ProjectName, &renderDetails.ProjectStartDate, &renderDetails.ProjectEndDate, &renderDetails.ProjectDescription, &renderDetails.ProjectTechnologies, &renderDetails.ProjectImage)
			
			// ERROR HANDLING
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("message : " + err.Error()))
			} else {
				// PARSING DATE
				renderDetails := Project{
					ID: 				 	renderDetails.ID, 	
					ProjectName:  		 	renderDetails.ProjectName,
					ProjectStartDateString: FormatDate(renderDetails.ProjectStartDate),
					ProjectEndDateString:   FormatDate(renderDetails.ProjectEndDate),
					ProjectDuration:     	GetDuration(renderDetails.ProjectStartDate, renderDetails.ProjectEndDate),
					ProjectDescription:  	renderDetails.ProjectDescription,
					ProjectTechnologies: 	renderDetails.ProjectTechnologies,
					ProjectImage: 		 	renderDetails.ProjectImage,
				}
				response := map[string]interface{}{
					"renderDetails": renderDetails,
				}
				w.WriteHeader(http.StatusOK)
				tmpl.Execute(w, response)
				}
			}
		}
		
func formAddProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/form-add-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}
	
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	} else {
		ProjectName := r.PostForm.Get("project-name")
		ProjectStartDate := r.PostForm.Get("date-start")
		ProjectEndDate := r.PostForm.Get("date-end")
		ProjectDescription := r.PostForm.Get("project-description")
		ProjectTechnologies := []string{r.PostForm.Get("nodejs"), r.PostForm.Get("reactjs"), r.PostForm.Get("golang"), r.PostForm.Get("typescript")}
		ProjectImage := r.PostForm.Get("upload-image")

		// var newProject = Project{
		// 	ProjectName:         projectName,
		// 	ProjectStartDate:    FormatDate(projectStartDate),
		// 	ProjectEndDate:      FormatDate(projectEndDate),
		// 	ProjectDuration:     GetDuration(projectStartDate, projectEndDate),
		// 	ProjectDescription:  projectDescription,
		// 	ProjectTechnologies: []string{projectUseNodeJS, projectUseReactJS, projectUseGolang, projectUseTypeScript},
		// }

		_, err = connection.Conn.Exec(context.Background(), `INSERT INTO public.project("ProjectName", "ProjectStartDate", "ProjectEndDate", "ProjectDescription", "ProjectTechnologies", "ProjectImage") VALUES ($1, $2, $3, $4, $5, $6)`, ProjectName, ProjectStartDate, ProjectEndDate, ProjectDescription, ProjectTechnologies, ProjectImage)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("message : " + err.Error()))
			return
		}

		// ProjectList = append(ProjectList, newProject)

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func updateProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	} else {
		ID, _ := strconv.Atoi(mux.Vars(r)["id"])
		ProjectName := r.PostForm.Get("project-name")
		ProjectStartDate := r.PostForm.Get("date-start")
		ProjectEndDate := r.PostForm.Get("date-end")
		ProjectDescription := r.PostForm.Get("project-description")
		ProjectTechnologies := []string{r.PostForm.Get("nodejs"), r.PostForm.Get("reactjs"), r.PostForm.Get("golang"), r.PostForm.Get("typescript")}
		ProjectImage := r.PostForm.Get("upload-image")

		// DATABASE CONNECTION
		_, err = connection.Conn.Exec(context.Background(), `UPDATE public.project SET "ProjectName"=$1, "ProjectStartDate"=$2, "ProjectEndDate"=$3, "ProjectDescription"=$4, "ProjectTechnologies"=$5, "ProjectImage"=$6 WHERE "ID"=$7`, ProjectName,ProjectStartDate, ProjectEndDate, ProjectDescription, ProjectTechnologies, ProjectImage, ID)

		// ERROR HANDLING
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("message : " + err.Error()))
			return
		}
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		
	}
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	ID, _ := strconv.Atoi(mux.Vars(r)["id"])

	// DELETE PROJECT BY ID
	_, err := connection.Conn.Exec(context.Background(), `DELETE FROM public.project WHERE "ID" = $1`, ID)
	
	// ERROR HANDLING
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}



// ADDITIONAL FUNCTION
//DURATION
func GetDuration(startDate time.Time, endDate time.Time) string {

	// layout := "2006-01-02"

	// date1, _ := time.Parse(layout, startDate)
	// date2, _ := time.Parse(layout, endDate)

	margin := endDate.Sub(startDate).Hours() / 24
	var duration string

	if margin >= 30 {
		if (margin / 30) == 1 {
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
func FormatDate(InputDate time.Time) string {

	// layout := "2006-01-02"
	// t, _ := time.Parse(layout, InputDate)

	Formated := InputDate.Format("02 January 2006")

	return Formated
}

// RETURN DATE FORMAT
func ReturnDate(InputDate time.Time) string {

	// layout := "02 January 2006"
	// t, _ := time.Parse(layout, InputDate)

	Formated := InputDate.Format("2006-01-02")

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