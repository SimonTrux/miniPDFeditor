package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"io"
	"path/filepath"
	"encoding/base64"

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
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.Handle("/", r)

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
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
	signatureImage, _, err := r.FormFile("signature")
	if err != nil {
		http.Error(w, "Error retrieving signature image", http.StatusInternalServerError)
		return
	}
	defer signatureImage.Close()

	// Apply changes to the PDF using gopdf package
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	pdf.Text(text)  // Adjust the coordinates as needed
	// Add signature image logic here

	// Save the modified PDF
	outputFilename := filepath.Join("output", filename)
	err = pdf.WritePdf(outputFilename)
	if err != nil {
		http.Error(w, "Error saving PDF", http.StatusInternalServerError)
		return
	}

	// Provide download link for the modified PDF
	http.ServeFile(w, r, outputFilename)
}

//
//// saveHandler function
//func saveHandler(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	filename := vars["filename"]
//
//	// Retrieve form data (text, signature image, etc.)
//	err := r.ParseMultipartForm(10 << 20) // Set max memory to 10 MB for the entire form
//	if err != nil {
//		http.Error(w, "Error parsing form", http.StatusInternalServerError)
//		return
//	}
//
//	text := r.FormValue("text")
//	signatureImage, _, err := r.FormFile("signature")
//	if err != nil {
//		http.Error(w, "Error retrieving signature image", http.StatusInternalServerError)
//		return
//	}
//	defer signatureImage.Close()
//
//	// Read the existing PDF content
//	pdfPath := filepath.Join("uploads", filename)
//	pdfContent, err := os.ReadFile(pdfPath)
//	if err != nil {
//		http.Error(w, "Error reading existing PDF file", http.StatusInternalServerError)
//		return
//	}
//
//	// Create a new PDF instance and load the existing content
//	pdf := gopdf.GoPdf{}
//	err = pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
//	if err != nil {
//		http.Error(w, "Error initializing PDF", http.StatusInternalServerError)
//		return
//	}
//	err = pdf.AddPage()
//	if err != nil {
//		http.Error(w, "Error adding page to PDF", http.StatusInternalServerError)
//		return
//	}
//	err = pdf.AddTTFFont("Arial", "assets/fonts/arial.ttf")
//	if err != nil {
//		http.Error(w, "Error adding font to PDF", http.StatusInternalServerError)
//		return
//	}
//	err = pdf.SetFont("Arial", "", 14)
//	if err != nil {
//		http.Error(w, "Error setting font in PDF", http.StatusInternalServerError)
//		return
//	}
//
//	// Write the existing content to the new PDF
//	err = pdf.WritePdfStream(pdfContent)
//	if err != nil {
//		http.Error(w, "Error adding existing content to PDF", http.StatusInternalServerError)
//		return
//	}
//
//	// Add text to the PDF
//	err = pdf.SetXY(50, 50)
//	if err != nil {
//		http.Error(w, "Error setting X coordinate in PDF", http.StatusInternalServerError)
//		return
//	}
//	err = pdf.Cell(nil, text)
//	if err != nil {
//		http.Error(w, "Error adding text to PDF", http.StatusInternalServerError)
//		return
//	}
//
//	// Add logic to handle the signature image here
//
//	// Save the modified PDF
//	outputFilename := filepath.Join("output", filename)
//	err = pdf.WritePdf(outputFilename)
//	if err != nil {
//		http.Error(w, "Error saving PDF", http.StatusInternalServerError)
//		return
//	}
//
//	// Provide download link for the modified PDF
//	http.ServeFile(w, r, outputFilename)
//}
