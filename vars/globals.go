package vars

import . "github.com/Okira-E/two-pass-sicxe-assembler/types"

var OpTable = map[string]CpuInstruction{
	"ADD": {
		Format: 3,
		Opcode: 0x18,
	},
	"ADDF": {
		Format: 3,
		Opcode: 0x58,
	},
	"ADDR": {
		Format: 2,
		Opcode: 0x90,
	},
	"AND": {
		Format: 3,
		Opcode: 0x40,
	},
	"CLEAR": {
		Format: 2,
		Opcode: 0xB4,
	},
	"COMP": {
		Format: 3,
		Opcode: 0x28,
	},
	"COMPF": {
		Format: 3,
		Opcode: 0x88,
	},
	"COMPR": {
		Format: 2,
		Opcode: 0xA0,
	},
	"DIV": {
		Format: 3,
		Opcode: 0x24,
	},
	"DIVF": {
		Format: 3,
		Opcode: 0x64,
	},
	"DIVR": {
		Format: 2,
		Opcode: 0x9C,
	},
	"FIX": {
		Format: 1,
		Opcode: 0xC4,
	},
	"FLOAT": {
		Format: 1,
		Opcode: 0xC0,
	},

	"HIO": {
		Format: 1,
		Opcode: 0xF4,
	},
	"J": {
		Format: 3,
		Opcode: 0x3C,
	},
	"JEQ": {
		Format: 3,
		Opcode: 0x30,
	},
	"JGT": {
		Format: 3,
		Opcode: 0x34,
	},
	"JLT": {
		Format: 3,
		Opcode: 0x38,
	},
	"JSUB": {
		Format: 3,
		Opcode: 0x48,
	},
	"LDA": {
		Format: 3,
		Opcode: 0x00,
	},
	"LDB": {
		Format: 3,
		Opcode: 0x68,
	},
	"LDCH": {
		Format: 3,
		Opcode: 0x50,
	},
	"LDF": {
		Format: 3,
		Opcode: 0x70,
	},
	"LDL": {
		Format: 3,
		Opcode: 0x08,
	},
	"LDS": {
		Format: 3,
		Opcode: 0x6C,
	},
	"LDT": {
		Format: 3,
		Opcode: 0x74,
	},
	"LDX": {
		Format: 3,
		Opcode: 0x04,
	},
	"LPS": {
		Format: 3,
		Opcode: 0xD0,
	},
	"MUL": {
		Format: 3,
		Opcode: 0x20,
	},
	"MULF": {
		Format: 3,
		Opcode: 0x60,
	},
	"MULR": {
		Format: 2,
		Opcode: 0x98,
	},
	"NORM": {
		Format: 1,
		Opcode: 0xC8,
	},
	"OR": {
		Format: 3,
		Opcode: 0x44,
	},
	"RD": {
		Format: 3,
		Opcode: 0xD8,
	},
	"RMO": {
		Format: 2,
		Opcode: 0xAC,
	},
	"RSUB": {
		Format: 3,
		Opcode: 0x4C,
	},
	"SHIFTL": {
		Format: 2,
		Opcode: 0xA4,
	},
	"SHIFTR": {
		Format: 2,

		Opcode: 0xA8,
	},
	"SIO": {
		Format: 1,
		Opcode: 0xF0,
	},
	"SSK": {
		Format: 3,
		Opcode: 0xEC,
	},

	"STA": {
		Format: 3,
		Opcode: 0x0C,
	},
	"STB": {
		Format: 3,
		Opcode: 0x78,
	},
	"STCH": {
		Format: 3,
		Opcode: 0x54,
	},
	"STF": {
		Format: 3,
		Opcode: 0x80,
	},
	"STI": {
		Format: 3,
		Opcode: 0xD4,
	},
	"STL": {
		Format: 3,
		Opcode: 0x14,
	},
	"STS": {
		Format: 3,
		Opcode: 0x7C,
	},
	"STSW": {
		Format: 3,
		Opcode: 0xE8,
	},
	"STT": {
		Format: 3,
		Opcode: 0x84,
	},
	"STX": {
		Format: 3,
		Opcode: 0x10,
	},
	"SUB": {
		Format: 3,
		Opcode: 0x1C,
	},
	"SUBF": {
		Format: 3,
		Opcode: 0x5C,
	},
	"SUBR": {
		Format: 2,
		Opcode: 0x94,
	},
	"SVC": {
		Format: 2,
		Opcode: 0xB0,
	},
	"TD": {
		Format: 3,
		Opcode: 0xE0,
	},
	"TIO": {
		Format: 1,
		Opcode: 0xF8,
	},

	"TIX": {
		Format: 3,
		Opcode: 0x2C,
	},
	"TIXR": {
		Format: 2,
		Opcode: 0xB8,
	},
	"WD": {
		Format: 3,
		Opcode: 0xDC,
	},
}
