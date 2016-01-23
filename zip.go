package main

import(
	"fmt"
	"bufio"
	"os"
	"strings"
	"strconv"
)

func ripCoords(line string) (float64, float64){
	leftParen := strings.Index(line, "(")
	rightParen := strings.Index(line, ")")
	comma := strings.Index(line, ",")
	fmt.Println(line[leftParen+1:comma], line[comma+3:rightParen])
	radLat, err := strconv.ParseFloat(line[leftParen+1:comma], 64)
	radLong, err := strconv.ParseFloat(line[comma+3:rightParen], 64)
	if err != nil{
		return 0.0, 0.0
	}
	lat := deg(float64(radLat))
	long := deg(float64(radLong))
	return lat, long
}

func isZipcode(line string, zipcode int) bool{
	components := []string{"[", strconv.Itoa(zipcode), "]"}
	substring := strings.Join(components, "") 
	//fmt.Println(substring, line)
	return strings.Contains(line, substring)
}

func getCoords(dataPath string, zipcode int) (float64, float64, error){
	file, err := os.Open(dataPath)
	lat, long := 0.0, 0.0
	if err != nil {
		return lat, long, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	exitF := 0
	for scanner.Scan(){
		line := scanner.Text()
		if (exitF == 1){
			lat, long = ripCoords(line)
			break
		}
		if (isZipcode(line, zipcode)){
			fmt.Println("triggered")
			exitF++
		}
	}
	
	return lat, long, scanner.Err()
}

func main(){
	lat, long, err := getCoords("/usr/share/weather-util/zctas", 19610)
	fmt.Println(lat, long, err)
}
