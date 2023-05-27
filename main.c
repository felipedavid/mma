#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <assert.h>
#include <ctype.h>
#include <stddef.h>
#include <string.h>

#include "common.c"
#include "lex.c"

void instr_dec_error(const char *fmt, ...) {
	va_list args;
	printf("Instruction decoding error at line %d: ", token.line);
	va_start(args, fmt);
	vprintf(fmt, args);
	va_end(args);
	printf("\n");
}

bool assert_bit_size(u16 value, int n_bits) {
	u16 max_value = 0;
	for (int i = 0; i < n_bits; i++) {
		max_value |= (1 << n_bits);
	}
	return value <= max_value;
}

typedef enum {
	R_INSTRUCTION,
	I_INSTRUCTION,
	J_INSTRUCTION,
} InstructionKind;

typedef enum {
	MOD_NONE,
	MOD_ADD,
	MOD_ADDI,
	MOD_AND,
	MOD_BEQ,
	MOD_J,
	MOD_LW,
	MOD_NAND,
	MOD_NOR,
	MOD_OR,
	MOD_SLT,
	MOD_SW,
	MOD_SUB,
} InstructionMod;

typedef enum {
	R0,
	R1,
	R2,
	R3,
	R4,
	R5,
	R6,
	R7,
	REG_FILE_END,
} Register;

void validate_register(Register reg) {
	if (reg < R0 || reg >= REG_FILE_END) {
		instr_dec_error("Register '$r%d' is out of range. The architecture only supports registers from 1 to 7", reg);
	}
}

typedef struct {
	u16 rs;
	u16 rt;
	u16 rd;
} RFormat;

typedef struct {
	u16 rs;
	u16 rt;
	u16 immd;
} IFormat;

typedef struct {
	u16 address;
} JFormat;

typedef struct {
	InstructionKind kind;
	InstructionMod mod;
	union {
		RFormat r;
		IFormat i;
		JFormat j;
	};
} Instruction;

// This functions should be replaced by a hash table.. But I'm too lazy to do it right now
InstructionMod get_instr_mod(InstructionKind kind, const char *instr_name) {
	InstructionMod mod = MOD_NONE;
	if (kind == R_INSTRUCTION) {
		if (instr_name == add_keyword) {
			mod = MOD_ADD;
		} else if (instr_name == and_keyword) {
			mod = MOD_AND;
		} else if (instr_name == nand_keyword) {
			mod = MOD_NAND;
		} else if (instr_name == nor_keyword) {
			mod = MOD_NOR;
		} else if (instr_name == or_keyword) {
			mod = MOD_OR;
		} else if (instr_name == slt_keyword) {
			mod = MOD_SLT;
		} else if (instr_name == sub_keyword) {
			mod = MOD_SUB;
		}
	} else if (kind == I_INSTRUCTION) {
		if (instr_name == addi_keyword) {
			mod = MOD_ADDI;
		} else if (instr_name == beq_keyword) {
			mod = MOD_BEQ;
		} else if (instr_name == lw_keyword) {
			mod = MOD_LW;
		} else if (instr_name == sw_keyword) {
			mod = MOD_SW;
		}
	} else if (kind == J_INSTRUCTION && instr_name == j_keyword) {
		mod = MOD_J;
	} 

	if (mod == MOD_NONE) {
		instr_dec_error("There is no %d instruction called '%s'", kind, instr_name);
	}
	return mod;
}

Instruction *new_instruction(InstructionKind kind) {
	Instruction *instr = calloc(1, sizeof(Instruction));
	instr->kind = kind;
	return instr;
}

Instruction *new_r_instruction(const char *instr_name, u16 rs, u16 rt, u16 rd) {
	Register r_rs = (Register)rs;
	Register r_rt = (Register)rt;
	Register r_rd = (Register)rd;

	validate_register(r_rs);
	validate_register(r_rt);
	validate_register(r_rd);

	Instruction *instr = new_instruction(R_INSTRUCTION);
	instr->mod = get_instr_mod(R_INSTRUCTION, instr_name);
	instr->r.rs = r_rs;
	instr->r.rt = r_rt;
	instr->r.rd = r_rd;

	return instr;
}

Instruction *new_i_instruction(const char *instr_name, u16 rs, u16 rt, u16 immd) {
	Register r_rs = (Register)rs;
	Register r_rt = (Register)rt;

	validate_register(r_rs);
	validate_register(r_rt);

	if (!assert_bit_size(immd, 6)) {
		instr_dec_error("immediates for i-instructions cannot be longer than 6 bits");
	}

	Instruction *instr = new_instruction(I_INSTRUCTION);
	instr->mod = get_instr_mod(I_INSTRUCTION, instr_name);
	instr->i.rs = r_rs;
	instr->i.rt = r_rt;
	instr->i.immd = immd;

	return instr;
}

Instruction *new_j_instruction(const char *instr_name, u16 addr) {
	if (!assert_bit_size(addr, 12)) {
		instr_dec_error("addresses for jump instructions cannot be longer than 12 bits");
	}
	Instruction *instr = new_instruction(J_INSTRUCTION);
	instr->mod = get_instr_mod(J_INSTRUCTION, instr_name);
	instr->j.address = addr;
	return instr;
}

u16 encode_instruction(Instruction *instr) {
	return 0;
}

void assemble(const char *source) {
	init_stream(source);
	while (token.kind) {
		//print_token(token);
		next_token();
	}
}

void assemble_file(const char *file_name) {
	FILE *fp = fopen(file_name, "r");
	if (fp == NULL) {
		fatal("cannot read file");
	}

	fseek(fp, 0L, SEEK_END);
	size_t file_size = ftell(fp);
	rewind(fp);

	char *source = xmalloc(file_size+1);
	source[file_size] = '\0';

	size_t n_read = fread(source, file_size, 1, fp);
	if (n_read != 1) {
		fatal("error while reading file");
	}

	assemble(source);
}

void lex_test() {
	init_stream("add $1, $2, $3\nsw $3, addr");
	expect_token(TOKEN_IDENTIFIER);
	expect_token(TOKEN_REGISTER);
	expect_token(',');
	expect_token(TOKEN_REGISTER);
	expect_token(',');
	expect_token(TOKEN_REGISTER);
	expect_token(TOKEN_IDENTIFIER);
	expect_token(TOKEN_REGISTER);
	expect_token(',');
	expect_token(TOKEN_IDENTIFIER);
}

Instruction *parse_instr() {
	return NULL;
}

void parse_test() {
	init_stream("add $1, $2, $3\nsw $3, addr");
	parse_instr();
}

void print_instr(Instruction *instr) {
	switch (instr->kind) {
	case R_INSTRUCTION:
		printf("r-instruction, mod: %d, rs: %d, rt: %d, rd: %d", 
			instr->mod, instr->r.rs, instr->r.rt, instr->r.rd);
		break;
	case I_INSTRUCTION:
		break;
	case J_INSTRUCTION:
		printf("j-instruction, mod: %d, addr: 0x%x", 
			instr->mod, instr->j.address);
		break;
	}
	printf("\n");
}

void instr_test() {
	Instruction *instructions[] = {
		new_r_instruction(lw_keyword, R1, R2, R3),
		new_j_instruction(j_keyword, 0x100),
		NULL,
	};
	for (Instruction **instr = instructions; *instr != NULL; instr++) {
		print_instr(*instr);
	}
}

int main(int argc, char **argv) {
	if (argc != 2) {
		fatal("Usage: %s <source_file>\n", argv[0]);
	}

	init_keywords();
	//assemble_file(argv[1]);
	instr_test();
	
	return 0;
}
