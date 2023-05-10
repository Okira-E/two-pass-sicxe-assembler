package assembler

import (
	"fmt"
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

	// Create the instructions array.
	var asmInstructions []types.AsmInstruction

	// Iterate over the asm string by lines.
	for _, line := range strings.Split(asm, "\n") {
		// Skip the line if it is empty.
		if line == "" {
			continue
		}
		// Skip the line if it is a comment.
		if strings.HasPrefix(line, ".") {
			continue
		}
		// Take the words inside a line and put them into a slice.
		// Each word is seperated by at least one space and/or a tab.
		words := strings.Fields(line)

		// Initialize the AsmInstruction struct.
		var asmInstruction types.AsmInstruction

		// Determine how many words are in the line.
		numberOfWords := len(words)
		// If there are more than 3 words, then the line is invalid.
		if numberOfWords > 3 {
			utils.Log("Invalid line: " + line)
			utils.Log("Too many words.")
			os.Exit(1)
		}
		// If there is only one word, then it is an opcode.
		if numberOfWords == 1 {
			asmInstruction.Label = "NIL"
			asmInstruction.OpCodeEn = words[0]
			asmInstruction.Operand = "NIL"
		} else if numberOfWords == 2 {
			// If there are two words, then the first word is an opcode and the second word is an operand.
			asmInstruction.Label = "NIL"
			asmInstruction.OpCodeEn = words[0]
			asmInstruction.Operand = words[1]
		} else {
			// If there are three words, then the first word is a label, the second word is an opcode, and the third word is an operand.
			asmInstruction.Label = words[0]
			asmInstruction.OpCodeEn = words[1]
			asmInstruction.Operand = words[2]
		}

		// Add the AsmInstruction struct to the instructions array.
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

func PrintAssemblerRules() {
	message := `
1. The assembler supports the instructions given on prompt 2.
2. The assembler supports the following addressing modes:
	- Immediate
	- Indirect
	- Simple
	- Indexed
3. The assembler supports the following pseudo-ops:
	- START
	- END
	- BYTE
	- WORD
	- RESB
	- RESW	
4. The assembler supports the following assembler directives:
	- BASE
5. Assembly rules:
	- The assembler assumes that the first instruction is the START instruction.
	- The assembler assumes that the last instruction is the END instruction.
	- Assembly code must be written in the following format:
		* <Label>	<Operation Code>	<Operand>
	- Assembly code must be written in a single line. symbols must be separated by at least one space and/or a tab.
	- Assembly code is not case-sensitive.
	- Comments are supported and must start with a '.' (dot).
	- The assembler assumes that the label field is optional. If it is present, it must be a valid label.
	- The assembler assumes that the operand field is optional. If it is present, it must be a valid operand.
	- The assembler assumes that the operation code field is mandatory. It must be a valid operation code.
	- The assembler assumes that the operand field is mandatory for the following operations:
		* BYTE
		* WORD
		* RESB
		* RESW
		* START
		* END
6. Examples can be found in the 'examples' directory.
		`

	utils.Log(message)
}

func PrintAsmWithObjectCodes(asmInstructions []types.AsmInstruction) {
	var tableRows []table.Row
	for _, val := range asmInstructions {
		label := val.Label
		operand := val.Operand

		if label == "NIL" {
			label = ""
		}
		if operand == "NIL" {
			operand = ""
		}

		tableRow := table.Row{fmt.Sprintf("%04X", val.Loc), label, val.OpCodeEn, operand, val.ObjCode}
		tableRows = append(tableRows, tableRow)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Loc", "Label", "OpCode", "Operand", "Object Code"})
	t.AppendRows(tableRows)
	t.Render()
}
