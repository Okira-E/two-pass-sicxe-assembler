package utils

import (
	"bufio"
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"log"
	"os"
	"strconv"
	"strings"
)

func PanicIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Log(str ...interface{}) {
	count := 0
	fmt.Print("\n")
	for _, s := range str {
		if count != 0 {
			fmt.Print(" ")
		}
		count++
		fmt.Print(s)
	}
	fmt.Print("\n")
}

func HandleCodeInput() string {
	var asm string

	Log("Enter your code ending in an END instruction: ")

	for {
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		PanicIfError(err)

		if strings.Contains(strings.ToUpper(input), "END") {
			break
		}
		asm += input
	}

	return asm
}

func HandleFileInput() (string, error) {
	var asm string

	Log("Enter the path to your file: ")
	reader := bufio.NewReader(os.Stdin)
	filePath, err := reader.ReadString('\n')
	PanicIfError(err)

	file, err := os.Open(strings.TrimSpace(filePath))
	if err != nil {
		return "", err
	}

	defer func() {
		err := file.Close()
		PanicIfError(err)
	}()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		asm += scanner.Text() + "\n"
	}

	return asm, nil
}

func ToHexRepresentation(num int) string {
	return fmt.Sprintf("%X", num)
}

func HexToDecimal(num int) int {
	_num := strconv.Itoa(num)
	decimal, err := strconv.ParseInt(_num, 16, 64)
	PanicIfError(err)

	return int(decimal)
}

func PrintSymTable(symTable map[string]int) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"LABEL", "ADDR"})

	tableRows := []table.Row{}
	for label, addr := range symTable {
		tableRow := table.Row{
			label,
			fmt.Sprintf("%X", addr),
		}

		tableRows = append(tableRows, tableRow)
	}

	t.AppendRows(tableRows)
	t.Render()
}

func CheckLengthAfterChar(str string, char string) int {
	// Get the length of the last element in the slice returned by strings.Split(str, char).
	// Example: T234*21332*42342*433 -> 3
	arr := strings.Split(str, char)

	length := len(arr[len(arr)-1])

	return length
}

func Complement16For3Digits(hex string) string {
	complement := ""

	for i, val := range hex {
		char := string(val)

		if char != "-" {
			decimalChar, err := strconv.Atoi(char)
			PanicIfError(err)

			if i == len(hex)-1 {
				complement += fmt.Sprintf("%X", 16-decimalChar)
			} else {
				complement += fmt.Sprintf("%X", 15-decimalChar)
			}
		}
	}

	for len(complement) < 3 {
		complement = "F" + complement[:]
	}

	return complement
}

func InvalidCharOrSpace(char byte) bool {
	// Remove all extra spaces.
	return char < 33 || char > 126
}
