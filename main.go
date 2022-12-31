package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public", http.FileServer(http.Dir("./public/"))))


	route.HandleFunc("/", homePage).Methods("GET")
	route.HandleFunc("/project", projectPage).Methods("GET")
	route.HandleFunc("/project/{id}", detailProject).Methods("GET")
	route.HandleFunc("/project", addProject).Methods("POST")
	route.HandleFunc("/contact", contactPage).Methods("GET")
	route.HandleFunc("/deleteProject/{id}", deleteProject).Methods("GET")
	route.HandleFunc("/editProject/{id}", editProject).Methods("GET")
	route.HandleFunc("/updateProject/{id}" , updateProject).Methods("POST")

	fmt.Println("Server running on port:8080")
	http.ListenAndServe("localhost:8080", route)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("view/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	dataCaller := map[string]interface{} {
		"Projects": dataSubmit,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, dataCaller)
}

func projectPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("view/myProject.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}


	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

// Struct buat menentukan variable sama tipe data nya
type dataReceive struct {
	ID int
	Projectname string
	Description string
	Technologies []string
	Startdate string
	Enddate string
	Duration string
}

// Nanti si variable dataSubmit ini bakal di isi sama value dari function di bawah
var dataSubmit = []dataReceive{

}

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10485760) // menggunakan r.ParseMultipartForm karena mengizinkan pengiriman file beda dengan r.ParseForm yang tidak mengizinkan pengiriman gambar
	// 10485760 itu parameter untuk ukuran batas file nya dalam satuan byte, jadi batas ukuran file yang diterima aku isi 10485760 byte atau 10 mb

	if err != nil {
		log.Fatal(err)
	}

	projectname := r.PostForm.Get("project-name")
	startDate := r.PostForm.Get("start-date")
	endDate := r.PostForm.Get("end-date")
	description := r.PostForm.Get("description")
	technologies := r.Form["technologies"] // pakai r.Form karena ingin menangkap query string


	// Duration Start
	const timeFormat = "2006-01-02" // Mendeklarasikan format tanggal
	timeStartDate, _:= time.Parse(timeFormat, startDate) //Mengubah format tanggal start date sesuai dengan const timeFormat
	timeEndDate, _:= time.Parse(timeFormat, endDate) //Mengubah format tanggal end date sesuai dengan const timeFormat

	// Hitung jarak antara start date dan end date hasilnya akan menjadi milisecond
	distance := timeEndDate.Sub(timeStartDate)

	//Ubah milisecond menjadi bulan, minggu dan hari
	monthDistance := int(distance.Hours() / 24 / 30)
	weekDistance := int(distance.Hours() / 24 / 7)
	daysDistance := int(distance.Hours() / 24)

	// variable buat menampung durasi yang sudah diolah
	var duration string
	// pengkondisian yang akan mengirimkan durasi yang sudah diolah
	if monthDistance >= 1 {
		duration = strconv.Itoa(monthDistance) + " months"
	} else if monthDistance < 1 && weekDistance >= 1 {
		duration = strconv.Itoa(weekDistance) + " weeks"
	} else if monthDistance < 1 && daysDistance >= 0 {
		duration = strconv.Itoa(daysDistance) + " days"
	} else {
		duration = "0 days"
	}
	// Duration End


	var newData = dataReceive{
		Projectname: projectname,
		Description: description,
		Technologies: technologies,
		Startdate: startDate,
		Enddate: endDate,
		Duration: duration,
	} 

	// fmt.Println("Project Name : " + projectname)
	// fmt.Println("Start-date : " + startDate)
	// fmt.Println("End-date : " + endDate)
	// fmt.Println("Description : " + description)
	// fmt.Println("Technologies : ", r.Form["technologies"] )
	// fmt.Println("Duration : " + duration)
	// fmt.Println("Image : " + imgname.Filename)

	dataSubmit = append(dataSubmit, newData)
	
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func contactPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("view/contact.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func detailProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("view/project-detail.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var dataProject = dataReceive{}

	for index, data := range dataSubmit {
		if index == id {
			dataProject = dataReceive{
				ID: id,
				Projectname: data.Projectname,
				Startdate: data.Startdate,
				Enddate: data.Enddate,
				Duration: data.Duration,
				Description: data.Description,
				Technologies: data.Technologies,
			}
		}
	}

	detailProject := map[string]interface{} {
		"Projects": dataProject,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, detailProject)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	dataSubmit = append(dataSubmit[:id], dataSubmit[id+1:]...)

	http.Redirect(w, r, "/", http.StatusFound)
}

func editProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("view/editProject.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var editData = dataReceive{}

	for index, data := range dataSubmit {
		if index == id{
			editData = dataReceive{
				ID: id,
				Projectname: data.Projectname,
				Startdate: data.Startdate,
				Enddate: data.Enddate,
				Duration: data.Duration,
				Description: data.Description,
				Technologies: data.Technologies,
			}
		}
	}

	dataEdit := map[string]interface{} {
		"Projects": editData,
	}

	tmpl.Execute(w, dataEdit)
}

func updateProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err := r.ParseMultipartForm(10485760) // menggunakan r.ParseMultipartForm karena mengizinkan pengiriman file beda dengan r.ParseForm yang tidak mengizinkan pengiriman gambar
	// 10485760 itu parameter untuk ukuran batas file nya dalam satuan byte, jadi batas ukuran file yang diterima aku isi 10485760 byte atau 10 mb

	if err != nil {
		log.Fatal(err)
	}

	projectname := r.PostForm.Get("project-name")
	startDate := r.PostForm.Get("start-date")
	endDate := r.PostForm.Get("end-date")
	description := r.PostForm.Get("description")
	technologies := r.Form["technologies"] // pakai r.Form karena ingin menangkap query string


	// Duration Start
	const timeFormat = "2006-01-02" // Mendeklarasikan format tanggal
	timeStartDate, _:= time.Parse(timeFormat, startDate) //Mengubah format tanggal start date sesuai dengan const timeFormat
	timeEndDate, _:= time.Parse(timeFormat, endDate) //Mengubah format tanggal end date sesuai dengan const timeFormat

	// Hitung jarak antara start date dan end date hasilnya akan menjadi milisecond
	distance := timeEndDate.Sub(timeStartDate)

	//Ubah milisecond menjadi bulan, minggu dan hari
	monthDistance := int(distance.Hours() / 24 / 30)
	weekDistance := int(distance.Hours() / 24 / 7)
	daysDistance := int(distance.Hours() / 24)

	// variable buat menampung durasi yang sudah diolah
	var duration string
	// pengkondisian yang akan mengirimkan durasi yang sudah diolah
	if monthDistance >= 1 {
		duration = strconv.Itoa(monthDistance) + " months"
	} else if monthDistance < 1 && weekDistance >= 1 {
		duration = strconv.Itoa(weekDistance) + " weeks"
	} else if monthDistance < 1 && daysDistance >= 0 {
		duration = strconv.Itoa(daysDistance) + " days"
	} else {
		duration = "0 days"
	}
	// Duration End


	var newData = dataReceive{
		Projectname: projectname,
		Description: description,
		Technologies: technologies,
		Startdate: startDate,
		Enddate: endDate,
		Duration: duration,
	} 

	dataSubmit[id] = newData
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}