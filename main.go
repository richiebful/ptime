package main

import(
	"fmt"
	"flag"
	"time"
)

/*type struct ArgList{
	date, lat, long, tz, zip bool
}*/


func validLatitude(lat float64) bool{
	return (lat >= -90.0) && (lat <= 90.0) 
}

func validLongitude(long float64) bool{
	return (long >= -180) && (long <= 180.0)
}

func validZone(tz float64) bool{
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
	long := flag.Float64("long", -200.0 , "Longitude of position")
	lat := flag.Float64("lat", -200.0, "Latitude of position")
	zip := flag.Int("zip", -1, "Zip Code of location")

	flag.Parse()

	//TODO, error handling
	//argL := ArgList{0, 0, 0, 0, 0}
	date, err := time.Parse("01/02/2006 -0700", *dateString)
	fmt.Println(date, *zip)
	if (err != nil){
		fmt.Printf("Invalid date %s\n", *dateString)
		return
	}
	//argL.lat = validLatitude(lat)
	//argL.long = validLongitude(long)
	//argL.tz = validZone(tz)
	//argL.zip = validZip(zip)
	loc := Location{*lat, *long, *tz}
	fmt.Println(loc.lat, loc.long, loc.tz)
	times := genTimes(date, loc, "ISNA")
	dispTimes(times)
}
