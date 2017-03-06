package main

import (
	"fmt"
	"regexp"
)

const (
	TOP_PRIORITY = "\u2191"
	MID_PRIORITY = "\u2014"
	LOW_PRIORITY = "\u2193"
	NO_PRIORITY = "\u22C5"

	WIP_STATUS = "\u2941"
	DONE_STATUS = "\u2713"
	NOT_STARTED_STATUS = "\u2717"

	FULL_SHORT_PRINT_FORM = "%v %v %v%v%v\t%v%v\t%v"
	NEWLINE_TYPE_SHORT_PRINT_FORM = "\t\t\t  %v"
	NEWLINE_TASK_SHORT_PRINT_FORM = "\t\t\t\t\t%v"
	NEWLINE_TASK_AFTER_TYPE_SHORT_PRINT_FORM = "\t\t%v"

	FULL_LONG_PRINT_FORM = "%v %v %v%v%v%v\t%v%v\t%v"
	NEWLINE_TYPE_LONG_PRINT_FORM = "\t\t\t\t  %v"
	NEWLINE_TASK_LONG_PRINT_FORM = "\t\t\t\t\t\t%v"
	NEWLINE_TASK_AFTER_TYPE_LONG_PRINT_FORM = "\t\t%v"
)

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

func SplitTextByNths(str string, wordCount int) []string {
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
	spaces := " "
	for i := 1; i < spaceCount; i++ {
		spaces += " "
	}
	return spaces
}

func Round(x, unit float64) float64 {
    return float64(int64(x/unit+0.5)) * unit
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