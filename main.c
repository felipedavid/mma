#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <assert.h>
#include <ctype.h>
#include <stddef.h>
#include <string.h>

#include "common.c"
#include "lex.c"

void assert_bit_size(u16 value, int n_bits) {
	u16 max_value = 0;
	for (int i = 0; i < n_bits; i++) {
		max_value |= (1 << n_bits);
	}
	assert(value <= max_value);
}

typedef enum {
	R_INSTRUCTION,
	I_INSTRUCTION,
	J_INSTRUCTION,
} InstructionKind;

typedef struct {
	u16 rs;
	u16 rt;
	u16 rd;
	u16 funct;
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
	u16 op;
	union {
		RFormat r;
		IFormat i;
		JFormat j;
	};
} Instruction;

Instruction *new_instruction(InstructionKind kind) {
	Instruction *instr = calloc(1, sizeof(Instruction));
	instr->kind = kind;
	return instr;
}

Instruction *new_r_instruction(u16 op, u16 rs, u16 rt, u16 rd, u16 funct) {
	assert_bit_size(op, 4);
	assert_bit_size(rs, 3);
	assert_bit_size(rt, 3);
	assert_bit_size(rd, 3);
	assert_bit_size(funct, 3);

	Instruction *instr = new_instruction(R_INSTRUCTION);
	instr->op = op;
	instr->r.rs = rs;
	instr->r.rt = rt;
	instr->r.rd = rd;
	instr->r.funct = funct;
	return instr;
}

Instruction *new_i_instruction(u16 op, u16 rs, u16 rt, u16 immd) {
	assert_bit_size(op, 4);
	assert_bit_size(rs, 3);
	assert_bit_size(rt, 3);
	assert_bit_size(immd, 6);

	Instruction *instr = new_instruction(I_INSTRUCTION);
	instr->op = op;
	instr->i.rs = rs;
	instr->i.rt = rt;
	instr->i.immd = immd;
	return instr;
}

Instruction *new_j_instruction(u16 op, u16 addr) {
	assert_bit_size(op, 4);
	assert_bit_size(addr, 12);

	Instruction *instr = new_instruction(J_INSTRUCTION);
	instr->op = op;
	instr->j.address = addr;
	return instr;
}

u16 encode_instruction(Instruction *instr) {
	return 0;
}

void assemble(const char *source) {
	init_stream(source);
	while (token.kind) {
		print_token(token);
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

int main(int argc, char **argv) {
	if (argc != 2) {
		fatal("Usage: %s <source_file>\n", argv[0]);
	}

	init_keywords();
	assemble_file(argv[1]);
	
	return 0;
}
