#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <assert.h>
#include <ctype.h>
#include <stddef.h>
#include <string.h>

#include "common.c"
#include "lex.c"
#include "instruction.c"
#include "parser.c"

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

void print_instr(Instruction *instr);

void parse_test() {
	init_stream("add $1, $2, $3\nsub $3, $2, $1");
	print_instr(parse_instr());
	print_instr(parse_instr());
}

void print_instr(Instruction *instr) {
	switch (instr->kind) {
	case R_INSTRUCTION:
		printf("r-instruction, mod: %d, rs: %d, rt: %d, rd: %d", 
			instr->mod, instr->r.rs, instr->r.rt, instr->r.rd);
		break;
	case I_INSTRUCTION:
		printf("i-instruction, mod: %d, rs: %d, rt: %d, immd: %d", 
			instr->mod, instr->i.rs, instr->i.rt, instr->i.immd);
		break;
	case J_INSTRUCTION:
		printf("j-instruction, mod: %d, addr: 0x%x", 
			instr->mod, instr->j.address);
		break;
	}
	printf("\n");
}

void instr_test() {
	token.line = 100;
	Instruction *instructions[] = {
		new_r_instruction(add_keyword, R1, R2, R3),
		new_i_instruction(addi_keyword, R1, R2, 65),
		new_j_instruction(j_keyword, 0x100),
		NULL,
	};
	for (Instruction **instr = instructions; *instr != NULL; instr++) {
		print_instr(*instr);
	}
}

void encode_test() {
	Instruction *instructions[] = {
		new_i_instruction(lw_keyword, R1, R2, 0),
		new_i_instruction(lw_keyword, R1, R3, 2),
		new_r_instruction(add_keyword, R2, R3, R4),
		new_i_instruction(sw_keyword, R1, R4, 4),
		NULL,
	};
	for (Instruction **instr = instructions; *instr != NULL; instr++) {
	u16 bin = encode_instruction(*instr);
		printf("%04x\n", bin);
	}
}

int main(int argc, char **argv) {
	init_keywords();
	parse_test();

	if (argc != 2) {
		fatal("Usage: %s <source_file>\n", argv[0]);
	}
	
	return 0;
}
