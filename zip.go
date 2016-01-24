package main

import(
	"bufio"
	"os"
	"strings"
	"strconv"
	"errors"
)

func ripCoords(line string) (float64, float64){
	leftParen := strings.Index(line, "(")
	rightParen := strings.Index(line, ")")
	comma := strings.Index(line, ",")
	radLat, err := strconv.ParseFloat(line[leftParen+1:comma], 64)
	radLong, err := strconv.ParseFloat(line[comma+2:rightParen], 64)
	if err != nil{
		return 0.0, 0.0
	}
	lat := deg(float64(radLat))
	long := deg(float64(radLong))
	return lat, long
}

func isZipcode(line string, zipcode string) bool{
	substring := "["+zipcode+"]" 
	return strings.Contains(line, substring)
}

func getCoords(dataPath string, zipcode string) (float64, float64, error){
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
			exitF++
		}
	}

	if (exitF == 0){
		return lat, long, errors.New("zip does not exist")
	}
	
	return lat, long, scanner.Err()
}
