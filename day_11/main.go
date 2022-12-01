package main

import (
	"context"
	"day_11/connection"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

// ==================================================
// STRUCT TEMPLATE
// ==================================================

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
	UserID				   int
}

var ProjectList = []Project{}

// ACCOUNT STRUCT
type User struct {
	ID 		 int
	Name 	 string
	Email 	 string
	Password string
}

// var UserList = []User{}

// SESSION STRUCT 
type Session struct {
	IsLogin   bool
	UserID    int
	UserName  string
	FlashData string
} 

var Data = Session{}

// ==================================================
// MAIN, HANDLEFUNC
// ==================================================

func main() {
	route := mux.NewRouter()

	// DATABASE CONNECT
	connection.DatabaseConnect()
	
	// ROUTING
	// ROUTE PATH PUBLIC
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// LOGIN 
	route.HandleFunc("/login", FormLogin).Methods("GET")
	route.HandleFunc("/logged-in", LoggedIn).Methods("POST")

	// REGISTER 
	route.HandleFunc("/form-registration", FormRegistration).Methods("GET")
	route.HandleFunc("/register", Register).Methods("POST")
	
	// LOGOUT
	route.HandleFunc("/logout", Logout).Methods("GET")

	// HOME
	route.HandleFunc("/", Home).Methods("GET")
	
	// CONTACT
	route.HandleFunc("/contact", Contact).Methods("GET")
	
	
	// PROJECT
	route.HandleFunc("/form-add-project", FormAddProject).Methods("GET")
	route.HandleFunc("/edit-project/{id}", EditForm).Methods("GET") 
	route.HandleFunc("/edited-project/{id}", EditProject).Methods("POST")
	route.HandleFunc("/project-details/{id}", ProjectDetails).Methods("GET")
	route.HandleFunc("/add-project", AddProject).Methods("POST")
	route.HandleFunc("/delete-project/{id}", DeleteProject).Methods("GET")



	fmt.Println("Server running on port 5000")
	http.ListenAndServe("localhost:5000", route)
}

// ==================================================
// HANDLERS
// ==================================================

// LOGIN
func FormLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/form-login.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		
		// GET COOKIES DATA
		store := sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")
		// CHECK LOGIN STATUS
		if session.Values["IsLogin"] != true {
			Data.IsLogin = false
		} else {
			Data.IsLogin  = session.Values["IsLogin"].(bool)
			Data.UserName = session.Values["UserName"].(string)
			Data.UserID   = session.Values["UserID"].(int)
		}
		

		response := map[string]interface{}{
			"Data": Data,
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, response)
	}
}

func LoggedIn(w http.ResponseWriter, r *http.Request){
	// SETUP COOKIE STORE
	store := sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	} else {
			Email 	 := r.PostForm.Get("input-email")
			Password := r.PostForm.Get("input-password")

			LoginUser := User{}
			// GET ID FROM DATABASE
			err := connection.Conn.QueryRow(context.Background(), `SELECT * FROM public.user WHERE "Email" = $1`, Email).Scan(&LoginUser.ID, &LoginUser.Name, &LoginUser.Email, &LoginUser.Password)

			// ERROR HANDLING
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("message : " + err.Error()))
				return
			} else {
				// CHECK PASSWORD
				err = bcrypt.CompareHashAndPassword([]byte(LoginUser.Password), []byte(Password))
				// ERROR HANDLING
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("message : " + err.Error()))
					return
			} else {
				// CREATE SESSION CACHE
				session.Values["IsLogin"] = true
				session.Values["UserName"] = LoginUser.Name
				session.Values["UserID"] = LoginUser.ID
				// LOGGED IN DURATION
				session.Options.MaxAge = 10800 // 10800 SECONDS = 3 HOURS

				session.AddFlash("Login Success!", "message")
				session.Save(r, w)

				http.Redirect(w, r, "/", http.StatusMovedPermanently)
			}
		}
	}
}

// REGISTRATION
func FormRegistration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/form-registration.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {

			// GET COOKIES DATA
			store := sessions.NewCookieStore([]byte("SESSION_KEY"))
			session, _ := store.Get(r, "SESSION_KEY")
			// CHECK LOGIN STATUS
			if session.Values["IsLogin"] != true {
				Data.IsLogin = false
			} else {
				Data.IsLogin  = session.Values["IsLogin"].(bool)
				Data.UserName = session.Values["UserName"].(string)
				Data.UserID   = session.Values["UserID"].(int)
			}

			response := map[string]interface{}{
				"Data": Data,
			}

			w.WriteHeader(http.StatusOK)
			tmpl.Execute(w, response)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		Name := r.PostForm.Get("input-name")
		Email := r.PostForm.Get("input-email")
		Password := r.PostForm.Get("input-password")
		
		// ENCRYPTING PASSWORD WITH BCRYPT
		PasswordHash, _ := bcrypt.GenerateFromPassword([]byte(Password), 10)

		_, err := connection.Conn.Exec(context.Background(), `INSERT INTO public.user("Name", "Email", "Password") VALUES ($1, $2, $3)`, Name, Email, PasswordHash)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("message : " + err.Error()))
			return
		}

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

// LOGOUT
func Logout(w http.ResponseWriter, r *http.Request){
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")
	session.Options.MaxAge = -1
	session.Save(r, w)
	
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// INDEX
func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {

		// GET COOKIES DATA
		store := sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")

		// CHECK LOGIN STATUS
		if session.Values["IsLogin"] != true {
			Data.IsLogin = false
		} else {
			fm := session.Flashes("message")
			
			var flashes []string
			if len(fm) > 0 {
				session.Save(r, w)
				for _, fl := range fm {
					flashes = append(flashes, fl.(string))
				}
			}

			Data.FlashData = strings.Join(flashes, "")
			Data.IsLogin   = session.Values["IsLogin"].(bool)
			Data.UserName  = session.Values["Username"].(string)
			Data.UserID	   = session.Values["UserID"].(int)
		}


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
					UserID: item.UserID,
				}
				renderData = append(renderData, item)
			}
		}
		response := map[string]interface{}{
			"renderData": renderData,
			"Data"		: Data,
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

		// GET COOKIES DATA
		store := sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")
		// CHECK LOGIN STATUS
		if session.Values["IsLogin"] != true {
			Data.IsLogin = false
		} else {
			Data.IsLogin = session.Values["IsLogin"].(bool)
			Data.UserName = session.Values["UserName"].(string)
			Data.UserID = session.Values["UserID"].(int)
		}

		response := map[string]interface{}{
			"Data": Data,
		}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, response)
}

// PROJECT
func FormAddProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/form-add-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}
	// GET COOKIES DATA
	store := sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")
	
	// CHECK LOGIN STATUS
	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin  = session.Values["IsLogin"].(bool)
		Data.UserName = session.Values["UserName"].(string)
		Data.UserID   = session.Values["UserID"].(int)
	}

	response := map[string]interface{}{
		"Data": Data,
	}
	
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, response)
}

func AddProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	} else {
		ProjectName := r.PostForm.Get("project-name")
		ProjectStartDate := r.PostForm.Get("date-start")
		ProjectEndDate := r.PostForm.Get("date-end")
		ProjectDescription := r.PostForm.Get("project-description")
		ProjectTechnologies := []string{r.PostForm.Get("nodejs"), r.PostForm.Get("reactjs"), r.PostForm.Get("golang"), r.PostForm.Get("typescript")}

		_, err = connection.Conn.Exec(context.Background(), `INSERT INTO public.project("ProjectName", "ProjectStartDate", "ProjectEndDate", "ProjectDescription", "ProjectTechnologies") VALUES ($1, $2, $3, $4, $5, $6)`, ProjectName, ProjectStartDate, ProjectEndDate, ProjectDescription, ProjectTechnologies)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("message : " + err.Error()))
			return
		}

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func ProjectDetails(w http.ResponseWriter, r *http.Request) {
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

				// GET COOKIES DATA
				store := sessions.NewCookieStore([]byte("SESSION_KEY"))
				session, _ := store.Get(r, "SESSION_KEY")
				// CHECK LOGIN STATUS
				if session.Values["IsLogin"] != true {
						Data.IsLogin = false
				} else {
						Data.IsLogin = session.Values["IsLogin"].(bool)
						Data.UserName = session.Values["UserName"].(string)
						Data.UserID = session.Values["UserID"].(int)
				}
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
					"Data"		   : Data,
				}
				w.WriteHeader(http.StatusOK)
				tmpl.Execute(w, response)
		}
	}
}

func EditForm(w http.ResponseWriter, r *http.Request) {
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

func EditProject(w http.ResponseWriter, r *http.Request) {
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

func DeleteProject(w http.ResponseWriter, r *http.Request) {
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



// ==================================================
// TIME HANDLERS
// ==================================================

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


