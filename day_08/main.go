package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Project struct {
	ProjectName	  		string
	ProjectStartDate 	string
	ProjectEndDate 		string
	ProjectDuration		string
	ProjectDescription	string
	ProjectTechnologies	[4]string
	ProjectImage		string
}

var ProjectList = []Project{}

func main() {
	route := mux.NewRouter()

	// ROUTING PATH TO PUBLIC
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// INDEX
	route.HandleFunc("/", Home).Methods("GET")
	
	// CONTACT
	route.HandleFunc("/contact", Contact).Methods("GET")
	
	// CREATE PROJECT 
	route.HandleFunc("/form-add-project", FormAddProject).Methods("GET")
	route.HandleFunc("/add-project", AddProject).Methods("POST")
	route.HandleFunc("/project-details/{index}", ProjectDetails).Methods("GET")
	
	// UPDATE PROJECT
	route.HandleFunc("/edit-project/{index}", EditProject).Methods("GET")
	route.HandleFunc("/delete-project/{index}", DeleteProject).Methods("GET")
	
	// PORT HANDLING
	fmt.Println("Server running on port 5000")
	http.ListenAndServe("localhost:5000", route)
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	response := map[string]interface{}{
		"ProjectList": ProjectList,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, response)
}

func Contact(w http.ResponseWriter, r *http.Request) {
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

// CREATE PROJECT
func FormAddProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/form-add-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func AddProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	} else {
	projectName 		 := r.PostForm.Get("project-name")
	projectStartDate	 := r.PostForm.Get("start-date")
	projectEndDate		 := r.PostForm.Get("end-date")
	projectDescription	 := r.PostForm.Get("project-description")
	projectUseNodeJS 	 := r.PostForm.Get("nodejs")
	projectUseReactJS	 := r.PostForm.Get("reactjs")
	projectUseGolang	 := r.PostForm.Get("golang")
	projectUseTypescript := r.PostForm.Get("typescript")
	ProjectImage		 := r.PostForm.Get("upload-image")

		var newProject = Project{
			ProjectName 		: projectName,
			ProjectStartDate 	: FormatDate(projectStartDate),
			ProjectEndDate		: FormatDate(projectEndDate),
			ProjectDuration 	: GetDuration(projectStartDate, projectEndDate),
			ProjectDescription	: projectDescription,
			ProjectTechnologies : [4]string{projectUseNodeJS, projectUseReactJS, projectUseGolang, projectUseTypescript},
			ProjectImage		: ProjectImage,
		}
	
			ProjectList = append(ProjectList, newProject)
			
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		}
}

func ProjectDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/project-details.html")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
		return
	} else {
		var renderDetails = Project{}
		index, _ := strconv.Atoi(mux.Vars(r)["index"])

		for i, data := range ProjectList {
			if index == i {
				renderDetails = Project{
					ProjectName		   : data.ProjectName,
					ProjectStartDate   : data.ProjectStartDate,
					ProjectEndDate	   : data.ProjectEndDate,
					ProjectDuration	   : data.ProjectDuration,
					ProjectDescription : data.ProjectDescription,
					ProjectTechnologies: data.ProjectTechnologies,
				}
			}
		}
		
		data := map[string]interface{}{
			"renderDetails": renderDetails,
		}
	
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
	}
}


// UPDATE PROJECT
func EditProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/edit-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	} else {
		var updateData = Project{}
		index, _ := strconv.Atoi(mux.Vars(r)["index"])

			for i, data := range ProjectList{
				if index == i {
					updateData = Project{
						ProjectName		   : data.ProjectName,
						ProjectStartDate   : ReturnDate(data.ProjectStartDate),
						ProjectEndDate	   : ReturnDate(data.ProjectEndDate),
						ProjectDescription : data.ProjectDescription,
						ProjectTechnologies: data.ProjectTechnologies,
				}

				ProjectList = append(ProjectList[:index], ProjectList[index+1:]...)
			}
		}
		data := map[string]interface{}{
			"updateData" : updateData,
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, data)
	}
}

func DeleteProject(w http.ResponseWriter, r *http.Request) {

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	ProjectList = append(ProjectList[:index], ProjectList[index+1:]...)

	http.Redirect(w, r, "/", http.StatusFound)
}


// ADDITIONAL FUNCTION
// GET DURATION
func GetDuration(startDate string, endDate string) string {

	layout := "2006-01-02"

	date1, _ := time.Parse(layout, startDate)
	date2, _ := time.Parse(layout, endDate)

	count := date2.Sub(date1).Hours() / 23
	var duration string

	if count > 30 {
		if (count / 30) <= 1 {
			duration = "1 Month"
		} else {
			duration = strconv.Itoa(int(count)/30) + " Months"
		}
	} else {
		if count <= 1 {
			duration = "1 Day"
		} else {
			duration = strconv.Itoa(int(count)) + "Days"
		}
	}
	return duration
}

// CHANGE DATE FORMAT
func FormatDate(InputDate string) string {

	layout := "2006-01-06"
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