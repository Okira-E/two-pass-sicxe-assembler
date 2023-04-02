package assembler

import (
	"errors"
	"fmt"
	. "github.com/Okira-E/two-pass-sicxe-assembler/types"
	"github.com/Okira-E/two-pass-sicxe-assembler/utils"
	"github.com/Okira-E/two-pass-sicxe-assembler/vars"
	"strconv"
	"strings"
)

// FirstPass returns a map of the symbol table.
// The key is the location counter and the value is the label.
// It modifies the AsmInstruction struct by adding the location counter to each instruction.
func FirstPass(asmInstructions *[]AsmInstruction) map[string]int {
	symTable := make(map[string]int)

	// Check missing START instruction.
	if (*asmInstructions)[0].OpCode != "START" {
		utils.PanicIfError(errors.New("ERROR: program doesn't start with a START instruction"))
	}

	// Check missing END instruction.
	if (*asmInstructions)[len(*asmInstructions)-1].OpCode != "END" {
		utils.PanicIfError(errors.New("ERROR: program doesn't end with an END instruction"))
	}

	startingAddress := (*asmInstructions)[0].Operand
	if startingAddress == "nil" {
		startingAddress = "0"
	}

	startingAddressInt, err := strconv.ParseInt(startingAddress, 16, 64)
	utils.PanicIfError(err)

	loc := int(startingAddressInt)

	for i := 0; i < len(*asmInstructions); i++ {
		asmInstructionRef := &(*asmInstructions)[i]
		asmInstructionRef.Loc = loc
		newLoc := loc

		if !asmInstructionRef.IsZeroLengthInstruction(vars.OpTable) {
			opCode := ""
			// If the OpCode has + before it, we add 1 to its length.
			addedByteDueToExtendedFormat := 0
			if strings.Contains(asmInstructionRef.OpCode, "+") {
				opCode = strings.ReplaceAll(asmInstructionRef.OpCode, "+", "")
				addedByteDueToExtendedFormat = 1
			} else {
				opCode = asmInstructionRef.OpCode
			}

			val := vars.OpTable[opCode].Format + addedByteDueToExtendedFormat

			newLoc += val
		} else if asmInstructionRef.IsReserveInstruction() {
			newLoc += asmInstructionRef.CalculateInstructionLength()
		}

		if strings.ToUpper(asmInstructionRef.Label) != "NIL" {
			symTable[asmInstructionRef.Label] = loc
		}

		loc = newLoc
	}

	return symTable
}

func SecondPass(asmInstructions *[]AsmInstruction, symTable map[string]int) string {
	objProgram := ""

	startingAddress := 0
	for _, asmInstruction := range *asmInstructions {
		if asmInstruction.OpCode == "START" {
			// If the operand is empty, we assume the starting address is 0.
			if asmInstruction.Operand != "" {
				val, err := strconv.Atoi(asmInstruction.Operand)
				utils.PanicIfError(err)
				startingAddress = val
			}

			// This assumes END is the last line in the assembly.
			startingAddressInDec := utils.HexToDecimal(startingAddress)
			objProgram += "H" + asmInstruction.Label + " " + fmt.Sprintf("%06d", startingAddress) + " " + fmt.Sprintf("%06X", (*asmInstructions)[len(*asmInstructions)-1].Loc-startingAddressInDec) + "\n"
			continue
		} else if asmInstruction.OpCode == "END" {
			objProgram += "E" + fmt.Sprintf("%06s", strconv.Itoa(startingAddress)) + "\n"
			continue
		} else {
			isInExtendedFormat := asmInstruction.OpCode[0] == '+'
			isImmediateAddrMode := asmInstruction.OpCode[0] == '#'
			isIndirectAddrMode := asmInstruction.OpCode[0] == '@'

			opCode := vars.OpTable[asmInstruction.OpCode].Opcode

			if isImmediateAddrMode {
				opCode += 1
			} else if isIndirectAddrMode {
				opCode += 2
			} else if isInExtendedFormat {

			}
		}
	}

	return objProgram
}
