package main

import(
	"fmt"
	"time"
	"math"
	"sort"
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

type PrayerTimes []PrayerTime

func (p PrayerTimes) Len() int{ return len(p) }
func (p PrayerTimes) Swap(i, j int){ p[i], p[j] = p[j], p[i] }
func (p PrayerTimes) Less(i, j int) bool{return p[i].time < p[j].time}

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

func dhuhrTime(loc Location, julT float64) float64{
	eqt, _ := sunPosition(julT)
	return 12.0 + float64(loc.tz) - loc.long/15.0 - eqt
}

func acot(n float64) float64{
	return math.Atan(1/n)
}

func asrTime(loc Location, julT, dhuhr, factor float64) float64{
	_, dec := sunPosition(julT)
	angle := -math.Atan(1/(factor + math.Tan(rad(math.Abs(loc.lat - dec)))))
	aTime := timeAngle(loc.lat, julT, dhuhr, angle, 1)
	fmt.Println(angle, aTime)
	return aTime
}

func timeAngle(lat, julT, dhuhr, angle, dir float64) float64{
	_, decl := sunPosition(julT)
	lat = rad(lat)
	decl = rad(decl)
	tmAngle := 1/15.0*deg(math.Acos((-math.Sin(angle) - math.Sin(lat)*math.Sin(decl))/(math.Cos(lat) * math.Cos(decl))))
	return dhuhr + dir*tmAngle
}

func formatTime(angularT float64) (int, int){
	hour := math.Floor(angularT)
	minute := math.Floor((angularT - hour) * 60)
	return int(hour), int(minute)
}

func initTimes(method string) (*PrayerTimes, error){
	calcMethods := map[string][]Method{
		"ISNA":	{
			{"dhuhr", 0.0},
			{"-angle", 15.0},
			{"fajr", -10.0},
			{"-angle", 0.833},
			{"asr", 1.0},
			{"angle", 0.833},
			{"angle", 0.833},
			{"angle", 15.0},
		},
		"MWL":	{
			{"dhuhr", 0.0},
			{"-angle", 18.0},
			{"fajr", -10.0},
			{"-angle", 0.833},
			{"asr", 1.0},
			{"angle", 0.833},
			{"angle", 0.833},
			{"angle", 17.0},
		},
		"EGAS":	{
			{"dhuhr", 0.0},
			{"-angle", 19.5},
			{"fajr", -10.0},
			{"-angle", 0.833},
			{"asr", 1.0},
			{"angle", 0.833},
			{"angle", 0.833},
			{"angle", 17.5},
		},
		"Makkah":	{
			{"dhuhr", 0.0},
			{"-angle", 18.5},
			{"fajr", -10.0},
			{"-angle", 0.833},
			{"asr", 1.0},
			{"angle", 0.833},
			{"angle", 0.833},
			{"maghrib", 90.0},
		},
		"Karachi":	{
			{"dhuhr", 0.0},
			{"-angle", 18.0},
			{"fajr", -10.0},
			{"-angle", 0.833},
			{"asr", 1.0},
			{"angle", 0.833},
			{"angle", 0.833},
			{"angle", 18.0},
		},
		"Tehran":	{
			{"dhuhr", 0.0},
			{"-angle", 17.7},
			{"fajr", -10.0},
			{"-angle", 0.833},
			{"asr", 1.0},
			{"angle", 0.833},
			{"angle", 4.5},
			{"angle", 14.0},
		},
		"Jafari":	{
			{"dhuhr", 0.0},
			{"-angle", 16.0},
			{"fajr", -10.0},
			{"-angle", 0.833},
			{"asr", 1.0},
			{"angle", 0.833},
			{"angle", 4.5},
			{"angle", 14.0},
		},
	}
	
	ptimes := PrayerTimes{
		{"dhuhr", 12.0, Method{"", 0.0}},
		{"fajr" , 5.0, Method{"", 0.0}},
		{"imsak", 5.0, Method{"", 0.0}},
		{"sunrise", 6.0, Method{"", 0.0}},
		{"asr", 13.0, Method{"", 0.0}},
		{"sunset", 18.0, Method{"", 0.0}},
		{"maghrib", 18.0, Method{"", 0.0}},
		{"isha", 18.0, Method{"", 0.0}},
	}

	for i := 0; i < 8; i++{
		ptimes[i].method = calcMethods[method][i]
	}

	return &ptimes, nil
}

func dispTimes(ptimes PrayerTimes){
	sort.Sort(ptimes)
	for i := 0; i < 8; i++ {
		time := ptimes[i].time
		label := ptimes[i].label
		//fmt.Println(ptimes[i].label, time)
		time = fixHour(time + 0.5/60.0);
		hr, min := formatTime(time)
		fmt.Printf("%s\t%.2d:%.2d\n", label, hr, min)
	}
}

func calculateTimes(ptimes PrayerTimes, jul float64, loc Location){
	pre := map[string]float64{"dhuhr": 0.0, "fajr": 0.0, "maghrib": 0.0}
	for i := 0; i < 8; i++{
		adjT := jul + ptimes[i].time/24.0
		switch ptimes[i].method.name {
		case "dhuhr":
			pre["dhuhr"] = dhuhrTime(loc, adjT)
			ptimes[i].time = pre["dhuhr"]
		case "angle":
			angle := rad(ptimes[i].method.number)
			ptimes[i].time = timeAngle(loc.lat, adjT, pre["dhuhr"], angle, 1)
			pre[ptimes[i].label] = ptimes[i].time
		case "asr":
			factor := ptimes[i].method.number
			ptimes[i].time = asrTime(loc, adjT, pre["dhuhr"], factor)
		case "fajr":
			ptimes[i].time = pre["fajr"] + ptimes[i].method.number/60.0
		case "maghrib":
			ptimes[i].time = pre["maghrib"] + ptimes[i].method.number/60.0
		case "-angle":
			angle := rad(ptimes[i].method.number)
			ptimes[i].time = timeAngle(loc.lat, adjT, pre["dhuhr"], angle, -1)
			pre[ptimes[i].label] = ptimes[i].time
		}
	}
}

func main(){
	date, _ := time.Parse("01/02/2006 15:04:05 -0700", "01/22/2016 00:00:00 -0500")
	_, offset := date.Zone()
	offset /= 3600.0
	loc := Location{40.0, -80.0, offset}
	timeRef, _ := initTimes("ISNA")
	ptimes := *timeRef
		
	jul := adjJulian(julian(date), loc)

	calculateTimes(ptimes, jul, loc)
	dispTimes(ptimes)
}
