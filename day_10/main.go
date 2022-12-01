package main

import (
	"context"
	"day_10/connection"
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
}

var ProjectList = []Project{}

func main() {
	route := mux.NewRouter()

	// Database Connect
	connection.DatabaseConnect()

	// route path folder untuk public
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	//routing
	route.HandleFunc("/", home).Methods("GET")
	
	
	// Project
	route.HandleFunc("/form-add-project", formAddProject).Methods("GET")
	route.HandleFunc("/edit-project/{id}", editForm).Methods("GET") 
	route.HandleFunc("/edited-project/{id}", editProject).Methods("POST")
	route.HandleFunc("/project-details/{id}", projectDetails).Methods("GET")
	route.HandleFunc("/add-project", addProject).Methods("POST")
	route.HandleFunc("/delete-project/{id}", deleteProject).Methods("GET")


	// Contact
	route.HandleFunc("/contact", contact).Methods("GET")
	

	fmt.Println("Server running on port 5000")
	http.ListenAndServe("localhost:5000", route)
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
		rows, _ := connection.Conn.Query(context.Background(), `SELECT "ID", "ProjectName", "ProjectStartDate", "ProjectEndDate", "ProjectDescription", "ProjectTechnologies" FROM public.project`)
		
		for rows.Next() {
			// Connect Struct
			err := rows.Scan(&item.ID, &item.ProjectName, &item.ProjectStartDate, &item.ProjectEndDate, &item.ProjectDescription, &item.ProjectTechnologies)
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

// PROJECT
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

		// var newProject = Project{
		// 	ProjectName:         projectName,
		// 	ProjectStartDate:    FormatDate(projectStartDate),
		// 	ProjectEndDate:      FormatDate(projectEndDate),
		// 	ProjectDuration:     GetDuration(projectStartDate, projectEndDate),
		// 	ProjectDescription:  projectDescription,
		// 	ProjectTechnologies: []string{projectUseNodeJS, projectUseReactJS, projectUseGolang, projectUseTypeScript},
		// }

		_, err = connection.Conn.Exec(context.Background(), `INSERT INTO public.project("ProjectName", "ProjectStartDate", "ProjectEndDate", "ProjectDescription", "ProjectTechnologies", "ProjectImage") VALUES ($1, $2, $3, $4, $5, $6)`, ProjectName, ProjectStartDate, ProjectEndDate, ProjectDescription, ProjectTechnologies)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("message : " + err.Error()))
			return
		}

		// ProjectList = append(ProjectList, newProject)

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

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
			err := connection.Conn.QueryRow(context.Background(), `SELECT "ID", "ProjectName", "ProjectStartDate", "ProjectEndDate", "ProjectDescription", "ProjectTechnologies" FROM public.project WHERE "ID" = $1`, ID).Scan(&renderDetails.ID, &renderDetails.ProjectName, &renderDetails.ProjectStartDate, &renderDetails.ProjectEndDate, &renderDetails.ProjectDescription, &renderDetails.ProjectTechnologies)
			
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
				}
				response := map[string]interface{}{
					"renderDetails": renderDetails,
				}
				w.WriteHeader(http.StatusOK)
				tmpl.Execute(w, response)
				}
			}
		}

func editForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/edit-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		ID, _ := strconv.Atoi(mux.Vars(r)["id"])
		updateData := Project{}

		err = connection.Conn.QueryRow(context.Background(), `SELECT "ID", "ProjectName", "ProjectStartDate", "ProjectEndDate", "ProjectDescription", "ProjectTechnologies" FROM public.project WHERE "ID" = $1`, ID).Scan(&updateData.ID, &updateData.ProjectName, &updateData.ProjectStartDate, &updateData.ProjectEndDate, &updateData.ProjectDescription, &updateData.ProjectTechnologies)
		
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("message : " + err.Error()))
			return
		}
		updateData = Project{
			ID:                     updateData.ID,
			ProjectName:            updateData.ProjectName,
			ProjectStartDateString: ReturnDate(updateData.ProjectStartDate),
			ProjectEndDateString:   ReturnDate(updateData.ProjectEndDate),
			ProjectDescription:     updateData.ProjectDescription,
			ProjectTechnologies:    updateData.ProjectTechnologies,
		}

		response := map[string]interface{}{
			"updateData": updateData,
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, response)
	}
}

func editProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	} else {
		ID, _ := strconv.Atoi(mux.Vars(r)["id"])
		ProjectName := r.PostForm.Get("project-name")
		ProjectStartDate := r.PostForm.Get("start-date")
		ProjectEndDate := r.PostForm.Get("end-date")
		ProjectDescription := r.PostForm.Get("project-description")
		ProjectTechnologies := []string{r.PostForm.Get("nodejs"), r.PostForm.Get("reactjs"), r.PostForm.Get("golang"), r.PostForm.Get("typescript")}
		
		// DATABASE CONNECTION
		_, err = connection.Conn.Exec(context.Background(), `UPDATE public.project SET "ProjectName"=$1, "ProjectStartDate"=$2, "ProjectEndDate"=$3, "ProjectDescription"=$4, "ProjectTechnologies"=$5 WHERE "ID"=$6`, ProjectName,ProjectStartDate, ProjectEndDate, ProjectDescription, ProjectTechnologies, ID)

		
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


