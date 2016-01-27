package main

import(
	"fmt"
	"flag"
	"time"
	"math"
)

func validLatitude(lat float64) bool{
	return (lat >= -90.0) && (lat <= 90.0) 
}

func validLongitude(long float64) bool{
	return (long >= -180) && (long <= 180.0)
}

func validZone(tz int) bool{
	return (tz >= -12.0) && (tz <= 12.0)
}

func validZip(zip int) bool{
	return (zip > 0) && (zip < 100000)
}

func main(){
	todayDate, defZone := nowDate()
	defDate := todayDate.Format("01/02/2006 -0700")
	tz := flag.Int("tz", defZone, "Time zone of location")
	dateString := flag.String("date", defDate, "Date of calculation")
	long := flag.Float64("long", math.NaN() , "Longitude of position")
	lat := flag.Float64("lat", math.NaN(), "Latitude of position")
	zip := flag.Int("zip", -1, "Zip Code of location")

	flag.Parse()
	
	latF := validLatitude(*lat)
	longF := validLongitude(*long)
	zipF := validZip(*zip)
	zoneF := validZone(*tz)
	date, err := time.Parse("01/02/2006 -0700", *dateString)
	fmt.Println(date, *zip)
	if (err != nil){
		fmt.Printf("Invalid date %s\n", *dateString)
		return
	}else if (latF || longF) && zipF {
		fmt.Printf("Conflicting coordinates and zip\n")
		return
	}else if !latF && longF && *lat != math.NaN() {
		fmt.Printf("Invalid latitude, %f\n", lat)
		return
	}else if latF && !longF && *long != math.NaN(){
		fmt.Printf("Invalid longitude, %f\n", long)
		return
	}else if !latF && longF {
		fmt.Printf("Missing latitude\n")
		return
	}else if latF && !longF {
		fmt.Printf("Missing longitude\n")
		return
	}else if !zoneF {
		fmt.Printf("Invalid time zone, %i\n", tz)
		return
	}else if !zipF && *zip != -1 {
		fmt.Printf("Invalid zip, %i\n", zip)
		return
	}else if latF && longF {
		loc := Location{*lat, *long, *tz}
	}else if zipF {
		zipS := fmt.Sprintf("%d", *zip)
		latitude, longitude, error := getCoords("/usr/share/weather-util/zctas", zipS)
		loc := Location{latitude, longitude, *tz}
	}

	times := genTimes(date, loc, "ISNA")
	dispTimes(times)
}
