package main

import(
	"fmt"
	"time"
	"math"
)

/*type struct Coords{
	long float64
	lat float64
}*/

//https://en.wikipedia.org/wiki/Julian_day
func julianCent(jul float64) float64{
	return (jul - 2451545.0)/36525.0
}

func julian(dt time.Time) float64{
	year := float64(dt.Year())
	month := float64(dt.Month())
	day := float64(dt.Day())
	hour := float64(dt.Hour())
	minute := float64(dt.Minute())
	second := float64(dt.Second())
	_, zone := dt.Zone()

	if month <= 2 {
                year -= 1
                month += 12
	}

	a := math.Floor(year/100.0)
	b := 2 - a + math.Floor(a / 4.0)
	jd := math.Floor(365.25 * (year + 4716.0)) + math.Floor(30.6001 * (month + 1)) + day + b - 1524.5
	mins := hour * 60.0 + minute + second / 60.0
	return jd + mins/1444.0 - float64(zone)/24.0
}

func fixAngle(deg float64) float64{
	return math.Mod(deg, 360.0)
}

func toRadians(deg float64) float64{
	return deg*math.Pi/180.0
}

func toDegrees(rad float64) float64{
	return rad*180.0/math.Pi
}

func meanObliquityEcliptic(jul float64) float64{
	seconds := 21.448 - jul*(46.8150 + jul*(0.00059 - jul*(0.001813)))
	e0 := 23.0 + (26.0 + (seconds/60.0))/60.0
	return e0
}

func obliquityCorrection(jul float64) float64{
	e0 := meanObliquityEcliptic(jul)
	omega := 125.04 - 1934.136 * jul
	e := e0 + 0.00256 * math.Cos(toRadians(omega))
	//fmt.Println("obliquity correction ", e)
	return e
}

func geomMeanAnomalySun(jul float64) float64{
	a := 357.52911 + jul * (35999.05029 - 0.0001537 * jul)
	//fmt.Println("Mean Anomaly: ", a)
	return a
}

func geomMeanLongSun(jul float64) float64{
	L0 := 280.46646 + jul * (36000.76983 + jul*(0.0003032))
	L0 = fixAngle(L0)
	//fmt.Println("Mean Long: ", L0)
	return L0
}

func eccentricityEarthOrbit(jul float64) float64{
	e := 0.016708634 - jul * (0.000042037 + 0.0000001267 * jul)
	//fmt.Println("Eccentricity: ", e)
	return e
}

func sunEqOfCenter(jul float64) float64{
	m := geomMeanAnomalySun(jul)
	mrad := toRadians(m)
	sinm := math.Sin(mrad)
	sin2m := math.Sin(2.0 * mrad)
	sin3m := math.Sin(3.0 * mrad)
	C := sinm * (1.914602 - jul * (0.004817 + 0.000014 * jul)) + sin2m * (0.019993 - 0.000101 * jul) + sin3m * 0.000289
	return C
}
	

func sunTrueLong(jul float64) float64{
	l0 := geomMeanLongSun(jul)
	c := sunEqOfCenter(jul)
	o := l0 + c
	return o
}

func sunApparentLong(jul float64) float64{
	o := sunTrueLong(jul)
	omega := 125.04 - 1934.136 * jul
	lambda := o - 0.00569 - 0.00478 * math.Sin(toRadians(omega))
	return lambda
}


func solarDeclination(jul float64) float64{
	e := obliquityCorrection(jul)
	lambda := sunApparentLong(jul)

	sint := math.Sin(toRadians(e)) * math.Sin(toRadians(lambda))
	return toDegrees(math.Asin(sint))
}

func equationOfTime(jul float64) float64{
	epsilon := obliquityCorrection(jul);
	l0 := geomMeanLongSun(jul);
	e := eccentricityEarthOrbit(jul);
	m := geomMeanAnomalySun(jul);

	y := math.Tan(toRadians(epsilon)/2.0);
	y *= y;
	
	sin2l0 := math.Sin(2.0 * toRadians(l0));
	sinm   := math.Sin(toRadians(m));
	cos2l0 := math.Cos(2.0 * toRadians(l0));
	sin4l0 := math.Sin(4.0 * toRadians(l0));
	sin2m  := math.Sin(2.0 * toRadians(m));
	
	eTime := y * sin2l0 - 2.0 * e * sinm + 4.0 * e * y * sinm * cos2l0 - 0.5 * y * y * sin4l0 - 1.25 * e * e * sin2m;
	
	return toDegrees(eTime) * 4.0 / 60.0
}

func dhuhr(zone int, longitude float64, clock_offset float64) float64{
	return 12.0 + float64(zone) - longitude/15.0 - clock_offset
}

func sunrise(lat float64, solDec float64, dhuhr float64) float64{
	return dhuhr - timeAngle(lat, solDec, toRadians(0.833))
}

func sunset(lat float64, solDec float64, dhuhr float64) float64{
	return dhuhr + timeAngle(lat, solDec, toRadians(0.833))
}

func fajr(lat float64, solDec float64, dhuhr float64) float64{
	angle := toRadians(15.0)
	return dhuhr - timeAngle(lat, solDec, angle)
}

func isha(lat float64, solDec float64, dhuhr float64) float64{
	angle := toRadians(15.0)
	return dhuhr + timeAngle(lat, solDec, angle)
}

func maghrib(lat float64, solDec float64, dhuhr float64) float64{
	angle := toRadians(4.0)
	return dhuhr + timeAngle(lat, solDec, angle)
}

func asr(lat float64, long float64, jul float64, factor float64) float64{
	dec := solarDeclination(jul)
	eqt := equationOfTime(jul)
	zone := 0
	//_, zone := time.Zone(
	dhuhrT := dhuhr(zone, long, eqt)
	angle := -1.0/math.Atan(factor + math.Tan(math.Abs(lat - dec)))
	return dhuhrT + timeAngle(lat, dec, angle)
}

func timeAngle(latitude float64, solDec float64, angle float64) float64{
	latitude = toRadians(latitude)
	solDec = toRadians(solDec)
	return 1/15.0*toDegrees(math.Acos(-math.Sin(angle) - math.Sin(latitude)*math.Sin(solDec)/(math.Cos(latitude) * math.Cos(solDec))))
}

func main(){
	date := time.Date(2015, time.December, 29, 12, 0, 0, 0, time.UTC)
	jul := julianCent(julian(date))
	eqTime := equationOfTime(jul)
	solDec := solarDeclination(jul)
	dh := dhuhr(0.0, 0.0, eqTime)
	fr := fajr(0.0, solDec, dh)
	fmt.Println("Fajr: ", fr)
	fmt.Println("Dhuhr: ", dh)
	is := isha(0.0, solDec, dh)
	as := asr(0.0, 0.0, jul, 1.0)
	fmt.Println("Asr: ", as)
	mr := maghrib(0.0, solDec, dh)
	fmt.Println("Maghrib: ", mr)
	fmt.Println("Isha: ", is)

}
