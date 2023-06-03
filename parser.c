
void parsing_error(const char *fmt, ...) {
	va_list args;
	printf("Parsing error at line %d: ", token.line);
	va_start(args, fmt);
	vprintf(fmt, args);
	va_end(args);
	printf("\n");
}

Register parse_register() {
	Register register_id = token.register_id;
	expect_token(TOKEN_REGISTER);	
	return (Register)register_id;
}

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
		rt = parse_register();
		expect_token(',');
		rd = parse_register();
		instr = new_r_instruction(instr_name, rs, rt, rd);
	} else if (is_name(addi_keyword) || is_name(beq_keyword) || is_name(lw_keyword) || is_name(sw_keyword)) {
		u16 immd = 0;
		expect_token(TOKEN_KEYWORD);
		rs = parse_register();
		expect_token(',');
		if (is_token(TOKEN_REGISTER)) {
			rt = parse_register();
			expect_token(',');
			immd = token.int_value;
			expect_token(TOKEN_NUMBER);
		} else if (is_token(TOKEN_NUMBER)) {
			immd = token.int_value;
			next_token();
			expect_token('(');
			rt = parse_register();
			expect_token(')');
		}
		instr = new_i_instruction(instr_name, rs, rt, immd);
	} else if (is_name(j_keyword)) {
		expect_token(TOKEN_KEYWORD);
		u16 addr = token.int_value;
		expect_token(TOKEN_NUMBER);
		instr = new_j_instruction(instr_name, addr);
	} else {
		fatal("Invalid instruction name %s", token.name);
	}

	return instr;
}

void parse_line() {
	switch (token.kind) {
	case TOKEN_KEYWORD:
		parse_instr();
		break;
	case TOKEN_IDENTIFIER_COLON:
		fatal("Not implemented yet!");
		break;
	case TOKEN_DIRECTIVE:
		fatal("Not implemented yet!");
		break;
	default:
		parsing_error("a %s should not appear the begging of a line", token_kind_str(token.kind));
	}
}
