#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <assert.h>
#include <ctype.h>
#include <stddef.h>
#include <string.h>
#include <stdarg.h>

#include "common.c"
#include "lexer.c"
#include "instruction.c"
#include "parser.c"

void parse_test() {
    init_stream("add $1, $2, $3\nsub $3, $2, $1\naddi $1, $2, 23\nj 123\n");
	for (;;) {
		Instruction *instr = parse_instr();
		if (instr == NULL) {
			break;
		}
		print_instr(instr);
	}
}

void run_tests() {
    parse_test();
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

int main(int argc, char **argv) {
	init_keywords();

#if 1
    run_tests();
#endif

	if (argc != 2) {
		fatal("Usage: %s <source_file>\n", argv[0]);
	}
	
	return 0;
}
