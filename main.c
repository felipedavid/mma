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

int main(int argc, char **argv) {
	if (argc != 2) {
		fatal("Usage: %s <source_file>\n", argv[0]);
	}

	init_keywords();
	//assemble_file(argv[1]);
	instr_test();
	
	return 0;
}
