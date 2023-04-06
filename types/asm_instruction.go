package types

import (
	"fmt"
	"github.com/Okira-E/two-pass-sicxe-assembler/utils"
	"strconv"
	"strings"
)

type AsmInstruction struct {
	Loc      int
	Label    string
	OpCodeEn string
	Operand  string
	ObjCode  string
}

func (instruction AsmInstruction) String() string {
	loc := fmt.Sprintf("%X", instruction.Loc)

	return loc + " " + instruction.Label + " " + instruction.OpCodeEn + " " + instruction.Operand + " " + instruction.ObjCode
}

func (instruction AsmInstruction) CalculateInstructionLength() int {
	if instruction.OpCodeEn == "RESB" {
		operandInt, err := strconv.Atoi(instruction.Operand)
		utils.PanicIfError(err)

		val := operandInt
		utils.PanicIfError(err)

		return val
	} else if instruction.OpCodeEn == "RESW" {
		operandInt, err := strconv.Atoi(instruction.Operand)
		utils.PanicIfError(err)

		val := operandInt
		utils.PanicIfError(err)

		return val * 3
	} else if instruction.OpCodeEn == "BYTE" {
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
	opCode := strings.ReplaceAll(instruction.OpCodeEn, "+", "")
	opCode = strings.ReplaceAll(opCode, "#", "")
	opCode = strings.ReplaceAll(opCode, "@", "")

	if _, ok := OpTable[opCode]; !ok {
		return true
	}

	return false
}

func (instruction AsmInstruction) IsReserveInstruction() bool {
	return instruction.OpCodeEn == "RESW" || instruction.OpCodeEn == "RESB" || instruction.OpCodeEn == "BYTE"
}
