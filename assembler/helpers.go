package assembler

import (
	"github.com/Okira-E/two-pass-sicxe-assembler/types"
	"github.com/Okira-E/two-pass-sicxe-assembler/utils"
	"github.com/Okira-E/two-pass-sicxe-assembler/vars"
	"github.com/jedib0t/go-pretty/table"
	"os"
	"strconv"
	"strings"
)

func ParseCode(asm string) []types.AsmInstruction {
	asm = strings.ToUpper(asm)

	// Replaces consecutive spaces with a single space.
	for i := range asm {
		char := asm[i]
		if i == 0 {
			if utils.InvalidCharOrSpace(char) {
				// Remove the character.
				asm = asm[:i] + asm[i+1:]
			}
			continue
		}

		prevChar := asm[i-1]
		if utils.InvalidCharOrSpace(char) {
			if utils.InvalidCharOrSpace(prevChar) {
				// Check if it's the last character. If so, remove it.
				if i == len(asm)-1 {
					asm = asm[:i]
					break
				}
				// Remove the character.
				asm = asm[:i] + asm[i+1:]
			} else {
				// Check if it's the last character. If so, replace it with a space.
				if i == len(asm)-1 {
					asm = asm[:i] + " "
					break
				}
				// Replace the character with a space.
				asm = asm[:i] + " " + asm[i+1:]
			}
		}

		// Since we parse by spaces, we need to remove the spaces after commas.
		if char == ',' && i != len(asm)-2 && asm[i+1] == ' ' && asm[i+2] != ' ' {
			asm = asm[:i+1] + asm[i+2:]
		}
	}

	// Create the instructions array.
	var asmInstructions []types.AsmInstruction

	// Split the code into lines.
	lines := strings.Split(asm, " ")

	for i := 0; i < len(lines)-3; i += 3 {
		var asmInstruction types.AsmInstruction

		asmInstruction.Label = lines[i]
		asmInstruction.OpCodeEn = lines[i+1]
		asmInstruction.Operand = lines[i+2]
		asmInstructions = append(asmInstructions, asmInstruction)
	}

	return asmInstructions
}

func PrintInstructionSet() {
	var tableRows []table.Row
	for opName, instructionProperties := range vars.OpTable {
		var allowedBits string

		if instructionProperties.Format == 3 {
			allowedBits = "3/4"
		} else {
			allowedBits = strconv.Itoa(instructionProperties.Format)
		}
		tableRow := table.Row{opName, allowedBits, utils.ToHexRepresentation(instructionProperties.Opcode)}
		tableRows = append(tableRows, tableRow)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Operation Name", "Format", "Opcode"})
	t.AppendRows(tableRows)
	t.Render()
}

func requiresTheIndexRegister(opCode string) bool {
	return strings.Contains(opCode, ",X")
}

func operandIsAbsolute(operand string) bool {
	// Note:- The following code assumes that a label field CANNOT be an absolute value. I.e. the EQU is not supported
	// If you want to support symbol-defining instructions, make sure if the value evaluated in the operand field is
	// indeed absolute. For now, any label is supposed to hold an address.

	// Check if the operand is a number or not. Number -> absolute value. NaN -> relative address.
	_, err := strconv.Atoi(operand)
	if err != nil {
		return false
	} else {
		return true
	}
}
