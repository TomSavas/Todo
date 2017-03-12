package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	TOP_PRIORITY = "\u2191"
	MID_PRIORITY = "\u2014"
	LOW_PRIORITY = "\u2193"
	NO_PRIORITY = "\u22C5"

	WIP_STATUS = "\u2941"
	DONE_STATUS = "\u2713"
	NOT_STARTED_STATUS = "\u2717"

	INFINITY = "\u221E"
)
var	NEWLINE_TYPE_LONG_PRINT_FORM string = TerribleIndentationHack(49) + "%v"
var	NEWLINE_TASK_LONG_PRINT_FORM string = TerribleIndentationHack(70) + "%v"
var NEWLINE_TYPE_SHORT_PRINT_FORM string = TerribleIndentationHack(34) + "%v"
var NEWLINE_TASK_SHORT_PRINT_FORM string = TerribleIndentationHack(55) + "%v"

func StartBold() {
	fmt.Print("\033[1m")
}

func StartFaint() {
	fmt.Print("\033[2m")
}

func StartItalic() {
	fmt.Print("\033[3m")
}

func StartUnderline() {
	fmt.Print("\033[4m")
}

func EndModifiers() {
	fmt.Print("\033[0m")
}

func ToBold(str string) string {
	return "\033[1m" + str + "\033[0m"
}

func ToFaint(str string) string {
	return "\033[2m" + str + "\033[0m"
}

func ToItalic(str string) string {
	return "\033[3m" + str + "\033[0m"
}

func ToUnderline(str string) string {
	return "\033[4m" + str + "\033[0m"
}

func CheckAndFaint(isColored bool) {
	if isColored {
		StartFaint()
	}
}

func FitStrings(types []string, maxLineLength int, putBarsAtStart bool) []string {
	splitType := func(typee string, maxLineLength int, putBarsAtStart bool) []string {
		splitTypes := []string {""}
		lineNumber := 0

		if putBarsAtStart {
			typee = "|" + typee
		}

		for _, word := range(strings.Split(typee, " ")) {
			for {
				lineLength := len(splitTypes[lineNumber])
				//TODO: only first line is exactly maxLineLength length, other ones are maxLineLength - 1 
				if lineLength + len(word) > maxLineLength {
					if maxLineLength - lineLength > maxLineLength / 5 {
						splitTypes[lineNumber] += word[:maxLineLength - lineLength - 1] + "-" // + FloatToString(float64(maxLineLength)) + " " + FloatToString(float64(lineLength)) + " " + FloatToString(float64(maxLineLength-lineLength))
						word = word[maxLineLength - lineLength - 1:]
						if putBarsAtStart {
							word = " " + word
						}
					}
					lineNumber++
					splitTypes = append(splitTypes, "")
				} else {
					break
				}
			}
			splitTypes[lineNumber] += word + " "
		} 
		return splitTypes
	}

	uniteSlices := func (firstSlice, secondSlice []string) []string {
		returnableSlice := make([]string, len(firstSlice) + len(secondSlice))
		i := 0
		for ; i < len(firstSlice); i++ {
			returnableSlice[i] = firstSlice[i]
		}

		for ; i - len(firstSlice) < len(secondSlice); i++ {
			returnableSlice[i] = secondSlice[i - len(firstSlice)]
		}

		return returnableSlice
	}

	fittedTypes := []string{}
	for i, _ := range(types) {
		fittedTypes = uniteSlices(fittedTypes, splitType(types[i], maxLineLength, putBarsAtStart))
	}

	return fittedTypes
}

func SplitTextByNWords(str string, wordCount int) []string {
	strings := []string{}
	lowerBound, upperBound, spaceCount := 0, 0, 0
	for i := 0; i < len(str); i++ {
		if str[i] == ' ' {
			spaceCount++
		}
		if spaceCount % wordCount == 0 && spaceCount != 0 || i == len(str) - 1 {
			spaceCount = 0
			upperBound = i + 1
			strings = append(strings, str[lowerBound:upperBound])
			lowerBound = upperBound 
		}
	} 
	return strings
}

func TerribleIndentationHack(spaceCount int) string {
	var spaces string
	for i := 0; i < spaceCount; i++ {
		spaces += " "
	}
	return spaces
}

func FloatToString(num float64) string {
	return fmt.Sprintf("%.1f", num)
}

func Round(x, unit float64) float64 {
    num := strconv.FormatFloat(float64(int64(x/unit+0.5)) * unit, 'g', 1, 64)
    if len(num) > 3 {
    	num = num[0:3]
    }
    returnableNum, _ := strconv.ParseFloat(num, 64) 
    return returnableNum
}

func SplitBySemicolons(src string) []string {
	if len(src) == 0 {
		return []string{}
	}

	query, _ := regexp.Compile("\\s*?;+?\\s*")
	return query.Split(src, -1)
}

func ValidateID(id string) bool {
	if found, _ := regexp.MatchString("[^\\d\\s]+", id); found {
		fmt.Println(ID_CANT_CONTAIN_LETTERS, HINT_FOR_HELP)
		return false
	} else if found, _ = regexp.MatchString("\\s+", id); found{
		fmt.Println(ID_CANT_CONTAIN_SPACES, HINT_FOR_HELP)
	}
	return true
}

func ValidateIDs (ids []string) bool {
	for _, id := range ids {
		if !ValidateID(id){
			return false
		}
	}
	return true
}