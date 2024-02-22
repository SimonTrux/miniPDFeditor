package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/signintech/gopdf"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/upload", uploadHandler).Methods("POST")
	r.HandleFunc("/edit/{filename}", editHandler).Methods("GET")
	r.HandleFunc("/save/{filename}", saveHandler).Methods("POST")
	r.HandleFunc("/download/{filename}", downloadHandler).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.Handle("/", r)

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "page1.html", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error uploading file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Save the uploaded file
	filename := filepath.Join("uploads", handler.Filename)
	out, err := os.Create(filename)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Error copying file", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/edit/"+handler.Filename, http.StatusFound)
}

// editHandler function
func editHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	pdfPath := filepath.Join("uploads", filename)

	// Read the PDF content
	pdfContent, err := os.ReadFile(pdfPath)
	if err != nil {
		http.Error(w, "Error reading PDF file", http.StatusInternalServerError)
		return
	}

	// Convert PDF content to base64 for embedding in HTML
	pdfBase64 := base64.StdEncoding.EncodeToString(pdfContent)

	templates.ExecuteTemplate(w, "edit.html", map[string]interface{}{
		"Filename":  filename,
		"PdfBase64": pdfBase64,
	})
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	// Retrieve form data (text, signature image, etc.)
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	text := r.Form.Get("text")

	// Apply changes to the PDF using gopdf package
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	pdf.Text(text) // Adjust the coordinates as needed

	// Save the modified PDF
	outputFilename := filepath.Join("output", filename)
	err = pdf.WritePdf(outputFilename)
	if err != nil {
		http.Error(w, "Error saving PDF", http.StatusInternalServerError)
		return
	}

	// Provide download link for the modified PDF
	downloadLink := "/download/" + filename
	http.Redirect(w, r, downloadLink, http.StatusSeeOther)
}

// Add a new handler for serving the download link
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	// Provide download link for the modified PDF
	outputFilename := filepath.Join("output", filename)
	http.ServeFile(w, r, outputFilename)
}
