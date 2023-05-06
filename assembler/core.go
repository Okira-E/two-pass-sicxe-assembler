package assembler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	. "github.com/Okira-E/two-pass-sicxe-assembler/types"
	"github.com/Okira-E/two-pass-sicxe-assembler/utils"
	"github.com/Okira-E/two-pass-sicxe-assembler/vars"
)

// FirstPass returns a map of the symbol table.
// The key is the location counter and the value is the label.
// It modifies the AsmInstruction struct by adding the location counter to each instruction.
func FirstPass(asmInstructions *[]AsmInstruction, baseRegister *BaseRegister) map[string]int {
	symTable := make(map[string]int)

	// Check missing START instruction.
	if (*asmInstructions)[0].OpCodeEn != "START" {
		utils.PanicIfError(errors.New("ERROR: program doesn't start with a START instruction"))
	}

	// Check missing END instruction.
	if (*asmInstructions)[len(*asmInstructions)-1].OpCodeEn != "END" {
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
			// If the OpCodeEn has + before it, we add 1 to its length.
			addedByteDueToExtendedFormat := 0
			if strings.Contains(asmInstructionRef.OpCodeEn, "+") {
				opCode = strings.ReplaceAll(asmInstructionRef.OpCodeEn, "+", "")
				addedByteDueToExtendedFormat = 1
			} else {
				opCode = asmInstructionRef.OpCodeEn
			}

			val := vars.OpTable[opCode].Format + addedByteDueToExtendedFormat

			newLoc += val
		} else if asmInstructionRef.IsReserveInstruction() {
			newLoc += asmInstructionRef.CalculateInstructionLength()
		} else if asmInstructionRef.OpCodeEn == "BASE" {
			baseOperand := asmInstructionRef.Operand
			baseOperand = strings.ReplaceAll(baseOperand, "#", "")
			baseOperand = strings.ReplaceAll(baseOperand, "@", "")

			_, err := strconv.Atoi(baseOperand)
			// The operand is a label.
			if err != nil {
				baseRegister.IsRef = true
				baseRegister.Value = baseOperand
			} else { // The operand is the numeric value.
				baseRegister.IsRef = false
				baseRegister.Value = baseOperand
			}
		}

		if strings.ToUpper(asmInstructionRef.Label) != "NIL" {
			symTable[asmInstructionRef.Label] = loc
		}

		loc = newLoc
	}

	return symTable
}

func SecondPass(asmInstructions *[]AsmInstruction, symTable map[string]int, baseRegister BaseRegister) string {
	objProgram := ""
	modificationRecords := ""

	startingAddress := 0
	for i, asmInstruction := range *asmInstructions {
		if asmInstruction.OpCodeEn == "START" {
			// If the operand is empty, we assume the starting address is 0.
			if asmInstruction.Operand != "" {
				operandInDecimal, err := strconv.ParseInt(asmInstruction.Operand, 16, 64)
				utils.PanicIfError(err)

				val := int(operandInDecimal)
				startingAddress = val
			}

			// This assumes END is the last line in the assembly.
			startingAddressInDec := utils.HexToDecimal(startingAddress)
			objProgram += "H" + asmInstruction.Label + " " + fmt.Sprintf("%06X", startingAddress) + " " + fmt.Sprintf("%06X", (*asmInstructions)[len(*asmInstructions)-1].Loc-startingAddressInDec) + "\n"
			continue
		} else if asmInstruction.OpCodeEn == "END" { // Gets the last line of the assembly. Is true last.
			objProgram += "\n" + modificationRecords
			objProgram += "\n" + "E" + fmt.Sprintf("%06X", startingAddress) + "\n"
			continue
		} else {
			isInExtendedFormat := asmInstruction.OpCodeEn[0] == '+'

			// Remove the '+' from the extended instruction for the OpTable to recognise it.
			if isInExtendedFormat {
				asmInstruction.OpCodeEn = asmInstruction.OpCodeEn[1:]
			}

			doesUseIndexReg := requiresTheIndexRegister(asmInstruction.Operand)

			if doesUseIndexReg {
				// Remove the ',X' from the operand.
				asmInstruction.Operand = asmInstruction.Operand[:len(asmInstruction.Operand)-2]
			}

			isImmediateAddrMode := asmInstruction.Operand[0] == '#'
			isIndirectAddrMode := asmInstruction.Operand[0] == '@'
			rawOperand := strings.ReplaceAll(asmInstruction.Operand, "#", "")
			rawOperand = strings.ReplaceAll(rawOperand, "@", "")

			// If the instruction doesn't belong to the opTable, it probably doesn't have object code.
			cpuInstruction, ok := vars.OpTable[asmInstruction.OpCodeEn]

			// Set the instructions that are not in the OpTable but have object code to true.
			if asmInstruction.OpCodeEn == "BYTE" || asmInstruction.OpCodeEn == "WORD" || asmInstruction.OpCodeEn == "RESW" || asmInstruction.OpCodeEn == "RESB" {
				ok = true
			}

			// If the instruction is not in the OpTable, we skip it.
			if !ok {
				continue
			}

			opCode := cpuInstruction.Opcode
			objCode := ""

			if asmInstruction.OpCodeEn == "BYTE" {
				operand := rawOperand
				if string(operand[0]) == "C" {
					// Remove the C' and ' from the string and convert it to hex.
					operand = fmt.Sprintf("%06X", rawOperand[2:len(rawOperand)-1])
				} else if string(operand[0]) == "X" {
					// Remove the X' and ' from the string.
					operand = rawOperand[2 : len(rawOperand)-1]
				} else {
					utils.PanicIfError(errors.New("ERROR: invalid BYTE instruction"))
				}
				// The operand is the object code.
				objCode += operand
			} else if asmInstruction.OpCodeEn == "WORD" {
				// Convert the operand to hex.
				operandInDecimal, err := strconv.ParseInt(rawOperand, 10, 64)
				if err != nil {
					utils.PanicIfError(err)

				}
				operand := fmt.Sprintf("%06X", operandInDecimal)
				// The operand is the object code.
				objCode += operand
			} else if asmInstruction.OpCodeEn == "RSUB" {
				// Set the NI flags to 11 and leave the operand empty.
				opCode += 3
				objCode = fmt.Sprintf("%02X", opCode) + "0000"
			} else if asmInstruction.OpCodeEn != "RESW" && asmInstruction.OpCodeEn != "RESB" {
				if cpuInstruction.Format == 2 {
					registers := strings.Split(rawOperand, ",")
					objCode += fmt.Sprintf("%02X", opCode)
					for _, register := range registers {
						objCode += fmt.Sprintf("%X", vars.RegisterTable[register])
					}
					if len(objCode) != 4 {
						objCode += "0"
					}
				} else {
					// Set the 7th and 8th bits from the left.
					if isImmediateAddrMode {
						opCode += 1
					} else if isIndirectAddrMode {
						opCode += 2
					} else {
						opCode += 3
					}

					firstTwoHalfBytes := fmt.Sprintf("%02X", opCode)
					// Set the first byte in the instruciton.
					objCode = firstTwoHalfBytes

					thirdHalfByte := 0

					// If the instruction has `,X` in the opCode, we set the 4th bit from the right.
					if doesUseIndexReg {
						thirdHalfByte += 8
					}

					if isInExtendedFormat { // 4-bit instruction
						// If the instruction is in extended format, we set the 1st bit from the right.
						// We also just put the operand in the object code in 2 bytes.
						thirdHalfByte += 1

						objCode += fmt.Sprintf("%01X", thirdHalfByte) + fmt.Sprintf("%05X", symTable[rawOperand])
						// We add a modification record for the operand because it is a relative address.
						if _, ok := symTable[rawOperand]; ok {
							modificationRecords += fmt.Sprintf("M%06X05\n", asmInstruction.Loc+1)
						}
					} else { // 3-bit instruction
						// See if we should use the Program Counter or the Base register to reach the operand by
						// displacement.

						// If the operand is in immediate addressing mode, we evaluate the object code immediately unless
						// it is a label that's value is relative to the starting address.
						if isImmediateAddrMode && operandIsAbsolute(rawOperand) {
							objCode += fmt.Sprintf("%01X", thirdHalfByte)

							// The value in the object code can come in two formats.
							var format string
							if isInExtendedFormat {
								format = "%04s"
							} else {
								format = "%03s"
							}

							objCode += fmt.Sprintf(format, rawOperand)
						} else {
							// If the operand is not in immediate addressing mode, we add it using displacement relative to
							// either the program counter or the value in the base register.

							// Calculate the Program Counter displacement
							programCounterDisp := symTable[rawOperand] - (*asmInstructions)[i+1].Loc
							// Calculate the Base register displacement.
							baseRegisterExists := baseRegister.Value != ""
							var baseRegisterDisp int
							// IsRef tells us if the value provided in the base register is a label itself or not.
							// "Not" being the numeric value provided immediately.
							if baseRegisterExists {
								if baseRegister.IsRef {
									baseRegisterDisp = symTable[rawOperand] - symTable[baseRegister.Value]
								} else {
									baseRegisterValueInDecimal, err := strconv.ParseInt(baseRegister.Value, 16, 64)
									baseValue := int(baseRegisterValueInDecimal)

									utils.PanicIfError(err)
									baseRegisterDisp = symTable[rawOperand] - baseValue
								}
							}

							// Decide which displacement should be used.
							var dispTechnique string
							if programCounterDisp >= -2048 && programCounterDisp <= 2047 {
								dispTechnique = "pc"
								thirdHalfByte += 2
							} else if baseRegisterExists && baseRegisterDisp >= 0 && baseRegisterDisp <= 4095 {
								dispTechnique = "base"
								thirdHalfByte += 4
							}

							objCode += fmt.Sprintf("%01X", thirdHalfByte)

							if dispTechnique == "pc" {
								// If it is in a negative value. Apply 16s complement to it first.
								if programCounterDisp < 0 {
									objCode += utils.Complement16For3Digits(fmt.Sprintf("%03X", programCounterDisp))
								} else {
									objCode += fmt.Sprintf("%03X", programCounterDisp)
								}
							} else if dispTechnique == "base" {
								objCode += fmt.Sprintf("%03X", baseRegisterDisp)
							}
						}
					}
				}
			}

			// Modify the object code of the current asmInstruction.
			(*asmInstructions)[i].ObjCode = objCode

			// Always check if we are on a new line. If we are, start a new `T` record. Else, check if we need to start a new
			// line.
			newLine := uint8(10)
			if objProgram[len(objProgram)-1] == newLine && (asmInstruction.OpCodeEn != "RESW" && asmInstruction.OpCodeEn != "RESB") {
				objProgram += "T" + fmt.Sprintf("%06X", asmInstruction.Loc) + "*" + objCode
			} else {
				currentTextRecordLength := utils.CheckLengthAfterChar(objProgram, "*")
				// If the current text record length is less than 30, we can add the object code to the current line.
				textRecordLengthAfterCurrentOpCode := currentTextRecordLength/2 + len(objCode)/2

				if textRecordLengthAfterCurrentOpCode > 30 && asmInstruction.OpCodeEn != "RESW" && asmInstruction.OpCodeEn != "RESB" {
					// If the current text record length is more than 30, we need to start a new line.
					objProgram += "\n" + "T" + fmt.Sprintf("%06X", asmInstruction.Loc) + "*" + objCode
				} else if asmInstruction.OpCodeEn == "RESW" || asmInstruction.OpCodeEn == "RESB" {
					objProgram += "\n"
				} else {
					objProgram += objCode
				}
			}
		}
	}

	// Clean the object program from any extra spaces.
	finalObjectProgram := ""
	lines := strings.Split(objProgram, "\n")
	for _, line := range lines {
		// Check if the current line in the object program is an empty line.
		cleanedLine := strings.ReplaceAll(line, " ", "")
		if len(cleanedLine) != 0 {
			currentTextRecordLength := fmt.Sprintf("%02X", utils.CheckLengthAfterChar(line, "*")/2)

			line = strings.ReplaceAll(line, "*", currentTextRecordLength)

			finalObjectProgram += "\n" + line
		}
	}
	//finalObjectProgram += "\n" + modificationRecords

	// Replace '*' characters with the length of the line.

	return finalObjectProgram
}
