
package main

import (
	"fmt"
	"log"
	"net/http"
	"html/template"
	"strconv"
	"time"
	"github.com/user/ptime"
	"encoding/json"
)

type LocationForm struct {
	Latitude, Zipcode, Longitude string
}

type OutputData struct{
	fajr, dhuhr, asr, maghrib, isha string
}

func initPage(rw http.ResponseWriter, req *http.Request){
	index := template.Must(template.ParseFiles(
			"www/html/main.html",
			"www/css/main.css",
			"www/js/main.js"))
	err := index.Execute(rw, nil)
	if err != nil{
		log.Println("Err: ",err)
	}
}

func mapifyTimes(times ptime.PrayerTimes) map[string]string {
	log.Println(times)
	result := make(map[string]string, len(times))
	for _, t := range times{		
		time := ptime.FixHour(t.Time + 0.5/60.0);
		hr, min := ptime.FormatTime(time)
		output := fmt.Sprintf("%.2d:%.2d", hr, min)
		result[t.Label] = output
	}
	return result
}

func updatePage(rw http.ResponseWriter, req *http.Request, times ptime.PrayerTimes){
	jsonTimes, err := json.Marshal(mapifyTimes(times))
	if err != nil{
		panic(nil)
	}
	log.Println("jsonTimes:", string(jsonTimes))
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write(jsonTimes)
	log.Println(rw)
}

func validZip(zip string) (int, bool){
	zipValue, err := strconv.Atoi(zip)
	if err != nil {
		return zipValue, false
	} else{
		return zipValue,  len(zip) == 5
	}
}

func validCoords(lat, long string) (float64, float64, bool) {
	latValue, err := strconv.ParseFloat(lat, 64)
	longValue, err2 := strconv.ParseFloat(long, 64)
	if (err != nil) || (err2 != nil) {
		return 0, 0, false
	} else{
		return latValue, longValue,
		(latValue >= -90.0) && (latValue <= 90.0) &&
			(longValue >= -180) && (longValue <= 180.0)
	}
}

func validZone(zone string) (int, bool){
	tz, err := strconv.Atoi(zone)
	if err != nil{
		return tz, false
	} else{
		return tz, true
	}
}

func processForm(rw http.ResponseWriter, req *http.Request){
	req.ParseForm()
	log.Println("req.Form" , req.Form)
	date, err := time.Parse("2006/1/2", req.Form["Date"][0])
	tz, tzFlag := validZone(req.Form["Zone"][0])
	_, zipFlag := validZip(req.Form["Zipcode"][0])
	lat, long, coordFlag := validCoords(req.Form["Latitude"][0], req.Form["Longitude"][0])
	if err != nil {
		log.Println("Invalid date")
		panic(nil)
	}
	if !tzFlag {
		log.Println("Invalid time zone")
		panic(nil)
	}

	log.Println(zipFlag, coordFlag)
	if zipFlag {
		zipLat, zipLong, err := ptime.GetCoords("/usr/share/weather-util/zctas", req.Form["Zipcode"][0])
		if err != nil{
			log.Println("Invalid zip code")
			panic(nil)
		}
		loc := ptime.Location{zipLat, zipLong, tz}
		log.Println(date, loc)
		times := ptime.GenTimes(date, loc, "ISNA")
		updatePage(rw, req, times)
		//ptime.DispTimes(times)
		log.Println("Valid zip")
		return
	} else if coordFlag {
		loc := ptime.Location{lat, long, tz}
		times := ptime.GenTimes(date, loc, "ISNA")
		updatePage(rw, req, times)
		log.Println("Valid coordinates")
		return
	} else if !coordFlag{
		fmt.Fprintf(rw, "Invalid coordinates")
		panic(nil)
	}
}

func showMain(rw http.ResponseWriter, req *http.Request){
	if (req.Method == "GET"){
		initPage(rw, req);
	} else if (req.Method == "POST"){
		processForm(rw, req);
	}
}

func main(){
	http.HandleFunc("/", showMain)
	http.Handle("/www/css/", http.StripPrefix("/www/css/", http.FileServer(http.Dir("www/css"))))
	http.Handle("/www/js/", http.StripPrefix("/www/js/", http.FileServer(http.Dir("www/js"))))
	http.Handle("/www/img/", http.StripPrefix("/www/img/", http.FileServer(http.Dir("www/img"))))
	log.Println("Server running on localhost:8080")
	log.Fatal(http.ListenAndServe(":4040", nil))
}
