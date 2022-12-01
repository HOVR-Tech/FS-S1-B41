package main

import (
	"context"
	"day_09/connection"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Project struct {
	ID 						int
	ProjectName	  			string
	ProjectStartDate 		time.Time
	ProjectEndDate 			time.Time
	ProjectStartDateString	string
	ProjectEndDateString	string
	ProjectDuration			string
	ProjectDescription		string
	ProjectTechnologies		[]string
}

var ProjectList = []Project{}

func main() {
	route := mux.NewRouter()

	// DATABASE CONNECTION
	connection.DatabaseConnect()

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
	} else {
		var renderData []Project
		var item = Project{}

		rows, _ := connection.Conn.Query(context.Background(), `SELECT "ID", "ProjectName", "ProjectStartDate", "ProjectEndDate", "ProjectDescription", "ProjectTechnologies" FROM public.project`)
		for rows.Next() {

		err := rows.Scan(&item.ID, &item.ProjectName, &item.ProjectStartDate, &item.ProjectEndDate, &item.ProjectDescription, &item.ProjectTechnologies)
		
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
	ProjectName 		 := r.PostForm.Get("project-name")
	ProjectStartDate	 := r.PostForm.Get("start-date")
	ProjectEndDate		 := r.PostForm.Get("end-date")
	ProjectDescription	 := r.PostForm.Get("project-description")
	ProjectTechnologies  := []string{r.PostForm.Get("nodejs"), r.PostForm.Get("reactjs"), r.PostForm.Get("golang"), r.PostForm.Get("typescript")}
		
		// var newProject = Project{
		// 	ProjectName 		: projectName,
		// 	ProjectStartDate 	: FormatDate(projectStartDate),
		// 	ProjectEndDate		: FormatDate(projectEndDate),
		// 	ProjectDuration 	: GetDuration(projectStartDate, projectEndDate),
		// 	ProjectDescription	: projectDescription,
		// 	ProjectTechnologies : []string{projectUseNodeJS, projectUseReactJS, projectUseGolang, projectUseTypescript},
		// }

		_, err = connection.Conn.Exec(context.Background(), `INSERT INTO public.project("ProjectName", "ProjectStartDate", "ProjectEndDate", "ProjectDescription", "ProjectTechnologies") VALUES ($1, $2, $3, $4, $5)`, ProjectName, ProjectStartDate, ProjectEndDate, ProjectDescription, ProjectTechnologies)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("message : " + err.Error()))
			return
		}
	
			// ProjectList = append(ProjectList, newProject)
			
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
	err := r.ParseForm()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	} else {
			ID, _			   	:= strconv.Atoi(mux.Vars(r)["id"])
			ProjectName		   	:= r.PostForm.Get("ProjectName")
			ProjectStartDate   	:= r.PostForm.Get("start-date")
			ProjectEndDate	   	:= r.PostForm.Get("end-date")
			ProjectDescription 	:= r.PostForm.Get("project-description")
			ProjectTechnologies := []string{r.PostForm.Get("nodejs"), r.PostForm.Get("reactjs"), r.PostForm.Get("golang"), r.PostForm.Get("typescript")}
		
		// DATABASE
			_, err := connection.Conn.Exec(context.Background(), `UPDATE public.project SET "ProjectName"=$1, "ProjectStartDate"=$2, "ProjectEndDate"=$3, "ProjectDescription"=$4, "ProjectTechnologies"=$5 WHERE "ID"=$6`, ProjectName, ProjectStartDate, ProjectEndDate, ProjectDescription, ProjectTechnologies, ID)
		
		// ERROR
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("message :" + err.Error()))
			return
		}

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func DeleteProject(w http.ResponseWriter, r *http.Request) {

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	ProjectList = append(ProjectList[:index], ProjectList[index+1:]...)

	http.Redirect(w, r, "/", http.StatusFound)
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