package server

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/diogenxs/dxs/models"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// go:embed static
var embedFS embed.FS

// start a web server on port 8080
type Server struct {
}

// type PageContent struct {
// 	Content  interface{}
// 	Metadata struct {
// 		Title       string
// 		Description string
// 		Navbar      []struct {
// 			Title string
// 			Link  string
// 		}
// 	}
// }

func NewServer() *Server {
	return &Server{}
}

func renderTemplate(w http.ResponseWriter, r *http.Request, name string, data any) {
	// Check if the header "HX-Request" is set to "true"
	if r.Header.Get("HX-Request") == "true" {
		// It is an HTMX request, only render the partial template
		tmpl := template.Must(template.ParseFiles("./server/templates/" + name))
		err := tmpl.ExecuteTemplate(w, "content", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		// indexData := PageContent{
		// 	Content: data,
		// 	Metadata: struct {
		// 		Title       string
		// 		Description string
		// 		Navbar      []struct {
		// 			Title string
		// 			Link  string
		// 		}
		// 	}{
		// 		Title:       "DXS",
		// 		Description: "DXS",
		// 		Navbar: []struct {
		// 			Title string
		// 			Link  string
		// 		}{
		// 			{
		// 				Title: "K8s Clusters",
		// 				Link:  "/k8s-clusters",
		// 			},
		// 			{
		// 				Title: "Alerts",
		// 				Link:  "/alerts",
		// 			},
		// 		},
		// 	},
		// }
		// Not an HTMX request, render the full page
		tmpl := template.Must(template.ParseFiles("./server/templates/index.html", "./server/templates/"+name))
		err := tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *Server) Start() {
	db, err := models.GetDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// create a new webserver on port 8080

	// handler := http.NewServeMux()
	handler := mux.NewRouter()

	// var staticFS = http.FS(embedFS)
	// handler.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(staticFS)))
	handler.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./server/static"))))

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./server/templates/index.html"))
		tmpl.Execute(w, nil)
	})

	handler.HandleFunc("/k8s-clusters", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, r, "k8s-clusters.html", nil)
	})

	handler.HandleFunc("/alerts", func(w http.ResponseWriter, r *http.Request) {
		alerts, err := models.ListAlerts(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		renderTemplate(w, r, "alerts.html", alerts)
	})

	handler.HandleFunc("/alerts/{alertid}", func(w http.ResponseWriter, r *http.Request) {
		// time.Sleep(time.Second * 2)
		vars := mux.Vars(r)
		alertid := vars["alertid"]
		err := models.AckAlert(db, alertid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}).Methods("DELETE")

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, handler)); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
