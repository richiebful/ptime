package main

import(
	"fmt"
	"flag"
	"time"
	"math"
	"github.com/user/ptime"
)

func nowDate() (time.Time, int){
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0.0, 0.0, 0.0, 0.0, now.Location())
	_, offset := now.Zone()
	return now, offset/3600
}

func validLatitude(lat float64) bool{
	return (lat >= -90.0) && (lat <= 90.0) 
}

func validLongitude(long float64) bool{
	return (long >= -180) && (long <= 180.0)
}

func validZone(tz int) bool{
	return (tz >= -12.0) && (tz <= 12.0)
}

func validZip(zip string) bool{
	return len(zip) == 5
}

func main(){
	todayDate, defZone := nowDate()
	defDate := todayDate.Format("01/02/2006")
	tz := flag.Int("tz", defZone, "Time zone of location")
	dateString := flag.String("date", defDate, "Date of calculation")
	long := flag.Float64("long", math.NaN() , "Longitude of position")
	lat := flag.Float64("lat", math.NaN(), "Latitude of position")
	zip := flag.String("zip", "-1", "Zip Code of location")

	flag.Parse()

        loc := ptime.Location{}
	
	latF := validLatitude(*lat)
	longF := validLongitude(*long)
	zipF := validZip(*zip)
	zoneF := validZone(*tz)
	date, err := time.Parse("01/02/2006", *dateString)
	
	if (err != nil){
		fmt.Printf("Invalid date %s\n", *dateString)
		return
	}else if math.IsNaN(*lat) && math.IsNaN(*long) && !zipF {
		fmt.Printf("No location defined\n")
		return
	}else if (latF || longF) && zipF {
		fmt.Printf("Conflicting coordinates and zip\n")
		return
	}else if !latF && longF {
		fmt.Printf("Invalid latitude, %f\n", lat)
		return
	}else if latF && !longF {
		fmt.Printf("Invalid longitude, %f\n", long)
		return
	}else if math.IsNaN(*lat) && !math.IsNaN(*long) {
		fmt.Printf("Missing latitude\n")
		return
	}else if !math.IsNaN(*lat) && math.IsNaN(*long) {
		fmt.Printf("Missing longitude\n")
		return
	}else if !zoneF {
		fmt.Printf("Invalid time zone, %d\n", tz)
		return
	}else if !zipF && *zip != "-1" {
		fmt.Printf("Invalid zip, %d\n", zip)
		return
	}else if latF && longF {
		loc.Lat = *lat
                loc.Long = *long
                loc.Tz = *tz
	}else if zipF {
		latitude, longitude, err := ptime.GetCoords("/usr/share/weather-util/zctas", *zip)
                if (err != nil){
			fmt.Printf("Zip code %s is not available\n", *zip)
                        return
                }
		loc.Lat = latitude
                loc.Long = longitude
                loc.Tz = *tz
	}
	fmt.Println(date, loc)
	times := ptime.GenTimes(date, loc, "ISNA")
	ptime.DispTimes(times)
}
