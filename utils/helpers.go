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

func HandleFileInput() string {
	var asm string

	Log("Enter the path to your file: ")
	reader := bufio.NewReader(os.Stdin)
	filePath, err := reader.ReadString('\n')
	HandleError(err)

	file, err := os.Open(strings.TrimSpace(filePath))
	HandleError(err)
	defer func() {
		err := file.Close()
		HandleError(err)
	}()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		asm += scanner.Text() + "\n"
	}

	return asm
}

func ToHexRepresentation(num int) string {
	return fmt.Sprintf("%X", num)
}

func DecimalToHex(num int) int {
	hex, err := strconv.ParseInt(fmt.Sprintf("%X", num), 16, 64)
	PanicIfError(err)

	return int(hex)
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

func DoNothing(param ...interface{}) {
}
