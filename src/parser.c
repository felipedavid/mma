
void parsing_error_line(int line, const char *fmt, ...);
void parsing_error(const char *fmt, ...);

Register parse_register() {
	Register register_id = token.register_id;
	expect_token(TOKEN_REGISTER);	
	return (Register)register_id;
}

int symbols_get(const char *key);

Instruction *parse_instr() {
	if (is_token(0)) {
		return NULL;
	}

	Instruction *instr = NULL;

	const char *instr_name = token.name;
	Register rt = 0, rd = 0, rs = 0;

	// TODO: This should be aptimized by using a hashtable but I was too lazy to actually write one >.>
	if (is_name(add_keyword) || is_name(and_keyword) || is_name(nand_keyword) || is_name(nor_keyword) ||
		is_name(or_keyword)  || is_name(slt_keyword) || is_name(sub_keyword)) {
		expect_token(TOKEN_KEYWORD);
		rs = parse_register();
		expect_token(',');
		rd = parse_register();
		expect_token(',');
		rt = parse_register();
		instr = new_r_instruction(instr_name, rs, rt, rd);
	} else if (is_name(addi_keyword) || is_name(beq_keyword) || is_name(lw_keyword) || is_name(sw_keyword)) {
		u16 immd = 0;
		expect_token(TOKEN_KEYWORD);
		rt = parse_register();
		expect_token(',');
		if (is_token(TOKEN_REGISTER)) {
			rs = parse_register();
			expect_token(',');
			immd = token.int_value;
			expect_token(TOKEN_NUMBER);
		} else if (is_token(TOKEN_NUMBER)) {
			immd = token.int_value;
			next_token();
			expect_token('(');
			rs = parse_register();
			expect_token(')');
		}
		instr = new_i_instruction(instr_name, rs, rt, immd);
	} else if (is_name(j_keyword)) {
        instr = new_j_instruction(instr_name, 0);
		expect_token(TOKEN_KEYWORD);
		if (is_token(TOKEN_IDENTIFIER)) {
            int index = symbols_get(token.name);
            if (index == -1) {
                instr->needs_fix_label = true;
                instr->label_to_resolve = token.name;
            } else {
                instr->j.address = index * 16;
            }
			next_token();
		} else {
			instr->j.address = token.int_value;
			expect_token(TOKEN_NUMBER);
		}
	} else {
		fatal("Invalid instruction name %s", token.name);
	}

    instr->line = token.line;

	return instr;
}

// TODO: This should be reimplemented with a pointer hash table
typedef struct { const char *name; int index; } Symbol;
Symbol *symbols;

void print_symbols() {
    for (int i = 0; i < buf_len(symbols); i++) {
        printf("%s: %d\n", symbols[i].name, symbols[i].index);
    }
}

// OBS: key is expected to be a interned string, since we are doing pointer comparisons
int symbols_get(const char *key) {
    for (int i = 0; i < buf_len(symbols); i++) {
        if (symbols[i].name == key) {
            return symbols[i].index;
        }
    }
    return -1;
}

void fix_unresolved_symbols() {
    for (int i = 0; i < buf_len(instructions); i++) {
        Instruction *instr = instructions[i];
        if (instr->needs_fix_label) {
            int index = symbols_get(instr->label_to_resolve);
            if (index != -1) {
                instr->j.address = index * 16;
            } else {
                parsing_error_line(instr->line, "there is no label called '%s'", instr->label_to_resolve);
            }
        }
    }
}

typedef enum {
    PARSING_CODE_SECTION,
    PARSING_DATA_SECTION,
} ParsingMode;

ParsingMode parser_mode;

const char *text_directive;
const char *data_directive;

void init_directives() {
    text_directive = str_intern("text");
    data_directive = str_intern("data");
}

bool parse_line() {
	switch (token.kind) {
	case '\0':
		return false;
	case TOKEN_KEYWORD:
		buf_push(instructions, parse_instr());
		break;
	case TOKEN_IDENTIFIER_COLON:
        buf_push(symbols, ((Symbol){token.name, buf_len(instructions)}));
        next_token();
        parse_line();
		break;
	case TOKEN_DIRECTIVE:
        assert(text_directive);
        assert(data_directive);

        if (token.name == text_directive) {
            parser_mode = PARSING_CODE_SECTION;
        } else if (token.name == data_directive) {
            parser_mode = PARSING_DATA_SECTION;
        } else {
            parsing_error("there is no directive called '.%s'. Did you mean .text or .data", token.name);
        }
        next_token();
		break;
	default:
		parsing_error("a %s should not appear the begging of a line", token_kind_str(token.kind));
		return false;
	}
	return true;
}

void parsing_error_line(int line, const char *fmt, ...) {
	va_list args;
	printf("Parsing error at line %d: ", line);
	va_start(args, fmt);
	vprintf(fmt, args);
	va_end(args);
	printf("\n");

    assembler.had_error = true;
}

void parsing_error(const char *fmt, ...) {
	va_list args;
	printf("Parsing error at line %d: ", token.line);
	va_start(args, fmt);
	vprintf(fmt, args);
	va_end(args);
	printf("\n");

    assembler.had_error = true;
}

