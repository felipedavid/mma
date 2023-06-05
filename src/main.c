#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <assert.h>
#include <ctype.h>
#include <stddef.h>
#include <string.h>
#include <stdarg.h>

typedef struct {
    bool had_error;

    char *instr_img;
    char *data_img;
} Assembler;

Assembler assembler;

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
    }
}

void run_tests() {
    parse_test();
}

void assemble(const char *source) {
	init_stream(source);
    while (parse_line());

    fix_unresolved_symbols();

    print_symbols();

    enum { LINE_CHAR_LEN = 6 };
    char instr_line[LINE_CHAR_LEN];
    for (Instruction **i = instructions; i != buf_end(instructions); i++) {
        sprintf(instr_line, "%04x\n", encode_instruction(*i));     

        // This is the only place where I use stretchy buffers for characters,
        // so I did not see the necessity of a buf_push_all macro.
        for (int i = 0; i < LINE_CHAR_LEN-1; i++) {
            buf_push(assembler.instr_img, instr_line[i]);
        }
    }

    assembler.instr_img[buf_len(assembler.instr_img)-1] = '\0';
    printf("%s", assembler.instr_img);
}

void assemble_file(char *file_name) {
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

    if (assembler.had_error) {
        return;
    }

    size_t file_name_size = strlen(file_name);
    if (file_name[file_name_size-2] != '.' || file_name[file_name_size-1] != 'm') {
        fprintf(stderr, "File does not have '.m' extension\n");
        exit(1);
    }

    char *instr_img = assembler.instr_img;
    if (instr_img) {
        // WARNING: Changing memory that was not me that allocated. That can produce some
        // weird side effect probably?
        file_name[file_name_size-1] = 'i';

        FILE *instr_img_f = fopen(file_name, "w");
        if (instr_img_f == NULL) {
            fprintf(stderr, "Could not open '%s'.\n", file_name);
            exit(1);
        }

        size_t written = fwrite("v2.0 raw\n", 9, 1, instr_img_f); 
        if (!written) {
            fprintf(stderr, "Unable to write to '%s'.\n", file_name);
            exit(1);
        }

        fwrite(instr_img, buf_len(instr_img), 1, instr_img_f);

        fclose(instr_img_f);
    }

    char *data_img = assembler.data_img;
    if (data_img) {
        // WARNING: Changing memory that was not me that allocated. That can produce some
        // weird side effect probably?
        file_name[file_name_size-1] = 'd';

        FILE *data_img_f = fopen(file_name, "w");
        if (data_img_f == NULL) {
            fprintf(stderr, "Could not open '%s'.\n", file_name);
            exit(1);
        }

        size_t written = fwrite(data_img, buf_len(data_img), 1, data_img_f);
        if (!written) {
            fprintf(stderr, "Unable to write to '%s'.\n", file_name);
            exit(1);
        }
        fclose(data_img_f);
    }
}

int main(int argc, char **argv) {
	if (argc != 2) {
		fatal("Usage: %s <source_file>\n", argv[0]);
	}

	init_keywords();
	init_directives();
    assemble_file(argv[1]);
	
	return 0;
}
