package main

import(
	"fmt"
	"time"
	"math"
)

type Location struct{
	lat float64
	long float64
	tz int
}

type Method struct{
	name string
	number float64
}

type PrayerTime struct{
	label string
	time float64
	method Method	
}

func adjJulian(jul float64, loc Location) float64{
	jul = jul - loc.long/(24.0*15.0)
	return jul
}

func julian(dt time.Time) float64{
	year := float64(dt.Year())
	month := float64(dt.Month())
	day := float64(dt.Day())
	hour := float64(dt.Hour())
	minute := float64(dt.Minute())
	second := float64(dt.Second())

	if month <= 2 {
                year -= 1
                month += 12
	}

	a := math.Floor(year/100.0)
	b := 2 - a + math.Floor(a / 4.0)
	jd := math.Floor(365.25 * (year + 4716.0)) + math.Floor(30.6001 * (month + 1)) + day + b - 1524.5
	mins := hour * 60.0 + minute + second / 60.0
	return jd + mins/1440.0
}

func fixAngle(deg float64) float64{
	return fix(deg, 360.0)
}

func fixHour(hr float64) float64{
	return fix(hr, 24.0)
}

func fix(num, den float64) float64{ 
	num = num - den * (math.Floor(num/den))
	if num < 0{
		return num + den
	}else{
		return num
	}
}

func rad(deg float64) float64{
	return deg*math.Pi/180.0
}

func deg(rad float64) float64{
	return rad*180.0/math.Pi
}

func sunPosition(jd float64) (float64, float64) {
	D := jd - 2451545.0
	g := fixAngle(357.529 + 0.98560028 * D)
	q := fixAngle(280.459 + 0.98564736 * D)
	L := fixAngle(q + 1.915 * math.Sin(rad(g)) + 0.020* math.Sin(2*rad(g)))

	//R := 1.00014 - 0.01671* math.Cos(rad(g)) - 0.00014* math.Cos(2*rad(g));
	e := 23.439 - 0.00000036* D

	RA := deg(math.Atan2(math.Cos(rad(e))* math.Sin(rad(L)), math.Cos(rad(L))))/ 15.0
	eqt := q/15 - fixHour(RA)
	decl := deg(math.Asin(math.Sin(rad(e))* math.Sin(rad(L))))

	return eqt, decl
}

func dhuhr(loc Location, clock_offset float64) float64{
	return 12.0 + float64(loc.tz) - loc.long/15.0 - clock_offset
}

func sunrise(lat float64, solDec float64, dhuhr float64) float64{
	return dhuhr - timeAngle(lat, solDec, rad(0.833))
}

func sunset(lat float64, solDec float64, dhuhr float64) float64{
	return dhuhr + timeAngle(lat, solDec, rad(0.833))
}

func fajr(lat float64, solDec float64, dhuhr float64) float64{
	angle := rad(15.0)
	return dhuhr - timeAngle(lat, solDec, angle)
}

func isha(lat float64, solDec float64, dhuhr float64) float64{
	angle := rad(15.0)
	return dhuhr + timeAngle(lat, solDec, angle)
}

func maghrib(lat float64, solDec float64, dhuhr float64) float64{
	angle := rad(4.0)
	return dhuhr + timeAngle(lat, solDec, angle)
}

func asr(loc Location, jul float64, factor float64) float64{
	eqt, dec := sunPosition(jul)
	dhuhrT := dhuhr(loc, eqt)
	angle := -math.Atan(1/(factor + math.Tan(loc.lat - dec)))
	return dhuhrT + timeAngle(loc.lat, dec, angle)
}

func timeAngle(lat float64, decl float64, angle float64) float64{
	lat = rad(lat)
	decl = rad(decl)
	return 1/15.0*deg(math.Acos(-math.Sin(angle) - math.Sin(lat)*math.Sin(decl)/(math.Cos(lat) * math.Cos(decl))))
}

func formatTime(angularT float64) (int, int){
	hour := math.Floor(angularT)
	minute := math.Floor((angularT - hour) * 60)
	return int(hour), int(minute)
}

func newPrayerTimes(method string) (*[]PrayerTime, error){
	ptimes := []PrayerTime{
		{"imsak", 5.0, Method{"", 0.0}},
		{"fajr" , 5.0, Method{"", 0.0}},
		{"sunrise", 6.0, Method{"", 0.0}},
		{"dhuhr", 12.0, Method{"", 0.0}},
		{"asr", 13.0, Method{"", 0.0}},
		{"sunset", 18.0, Method{"", 0.0}},
		{"maghrib", 18.0, Method{"", 0.0}},
		{"isha", 18.0, Method{"", 0.0}},
	}

	if (method == "ISNA"){
		fmt.Println("Satan");
		
	}else{
		return &ptimes, fmt.Errorf("Invalid method")
	}
	
	return &ptimes, nil
}

func main(){
	date, _ := time.Parse("01/02/2006 15:04:05 -0700", "01/20/2016 00:00:00 -0500")
	_, offset := date.Zone()
	loc := Location{40.0, -80.0, offset}
	ptimes, _ := newPrayerTimes("ISNA")
		
	jul := adjJulian(julian(date), loc)
	//eqTime := equationOfTime(jul)
	//solDec := solarDeclination(jul)
	eqTime, solDec := sunPosition(jul)
	fmt.Println(eqTime*60.0, solDec)
	//fmt.Println(eqTime, eqt)
	//fmt.Println(solDec, decl)
	
	dhuhrT := dhuhr(loc, eqTime)
	
	fajrT := fajr(loc.lat, solDec, dhuhrT)
	hr, min := formatTime(fajrT)
	fmt.Printf("Fajr, %d:%.2d\n", hr, min)
	
	hr, min = formatTime(dhuhrT)
	fmt.Printf("Dhuhr, %d:%.2d\n", hr, min)
	
	asrT := asr(loc, jul, 1.0)
	hr, min = formatTime(asrT)
	fmt.Printf("Asr, %d:%.2d\n", hr, min)
	
	maghribT := maghrib(loc.lat, solDec, dhuhrT)
	hr, min = formatTime(maghribT)
	fmt.Printf("Maghrib, %d:%.2d\n", hr, min)

	ishaT := isha(loc.lat, solDec, dhuhrT)
	hr, min = formatTime(ishaT)
	fmt.Printf("Isha, %d:%.2d\n", hr, min)
}
