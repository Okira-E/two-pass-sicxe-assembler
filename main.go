package main

import (
	"fmt"
	"github.com/Okira-E/two-pass-sicxe-assembler/assembler"
	"github.com/Okira-E/two-pass-sicxe-assembler/types"
	"github.com/Okira-E/two-pass-sicxe-assembler/utils"
)

func main() {
	const choicesPrompt = "Enter (1) to read from a file, (2) to view the instruction set, or (3) to view the rules of the assembly: "

	utils.Log("Welcome to my SIC/XE Assembler!")
	utils.Log(choicesPrompt)

	var asm string
	counter := 0
	for {
		if counter > 0 {
			utils.Log(choicesPrompt)
		}
		counter++

		var choice int
		_, err := fmt.Scan(&choice)
		utils.PanicIfError(err)

		if choice == 1 {
			for {
				asm, err = utils.HandleFileInput()
				if err != nil {
					utils.Log(err.Error())
					continue
				} else {
					break
				}
			}
			break
		} else if choice == 2 {
			assembler.PrintInstructionSet()
			continue
		} else if choice == 3 {
			assembler.PrintAssemblerRules()
		} else {
			utils.Log("Invalid choice.")
		}
	}

	// Return a slice of AsmInstruction structs.
	asmInstructions := assembler.ParseCode(asm)

	utils.Log("-- First Pass")
	// Creates the Symbol Table, as well as, assigns memory locations for each line in the assembly.
	baseRegister := types.BaseRegister{}
	symTable := assembler.FirstPass(&asmInstructions, &baseRegister)
	// Print the Symbol table.
	utils.Log("Symbol Table:")
	utils.PrintSymTable(symTable)

	// Returns the object program, as well as, creates an object code for each line in the assembly.
	objProgram := assembler.SecondPass(&asmInstructions, symTable, baseRegister)
	utils.Log("-- Second Pass")

	assembler.PrintAsmWithObjectCodes(asmInstructions)

	utils.Log(objProgram)
}
