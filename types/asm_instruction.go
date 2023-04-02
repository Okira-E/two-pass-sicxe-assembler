package types

import (
	"fmt"
	"github.com/Okira-E/two-pass-sicxe-assembler/utils"
	"strconv"
	"strings"
)

type AsmInstruction struct {
	Loc     int
	Label   string
	OpCode  string
	Operand string
}

func (instruction AsmInstruction) String() string {
	loc := fmt.Sprintf("%X", instruction.Loc)

	return loc + " " + instruction.Label + " " + instruction.OpCode + " " + instruction.Operand
}

func (instruction AsmInstruction) CalculateInstructionLength() int {
	if instruction.OpCode == "RESB" {
		operandInt, err := strconv.Atoi(instruction.Operand)
		utils.HandleError(err)

		val := operandInt
		utils.HandleError(err)

		return val
	} else if instruction.OpCode == "RESW" {
		operandInt, err := strconv.Atoi(instruction.Operand)
		utils.HandleError(err)

		val := operandInt
		utils.HandleError(err)

		return val * 3
	} else if instruction.OpCode == "BYTE" {
		if instruction.Operand[0] == 'C' {
			val := strings.Split(instruction.Operand, "`")[1]

			length := len(val)

			return length
		} else if instruction.Operand[0] == 'X' {
			val := strings.Split(instruction.Operand, "`")[1]

			return len(val) / 2
		}
	}

	return 0
}

func (instruction AsmInstruction) IsZeroLengthInstruction(OpTable map[string]CpuInstruction) bool {
	opCode := strings.ReplaceAll(instruction.OpCode, "+", "")
	opCode = strings.ReplaceAll(opCode, "#", "")
	opCode = strings.ReplaceAll(opCode, "@", "")

	if _, ok := OpTable[opCode]; !ok {
		return true
	}

	return false
}

func (instruction AsmInstruction) IsReserveInstruction() bool {
	return instruction.OpCode == "RESW" || instruction.OpCode == "RESB" || instruction.OpCode == "BYTE"
}
