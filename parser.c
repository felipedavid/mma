
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
	Instruction *instr = NULL;

	// TODO: This should be aptimized by using a hashtable but I was too lazy to actually write one >.>
	if (is_name(add_keyword) || is_name(and_keyword) || is_name(nand_keyword) || is_name(nor_keyword) ||
		is_name(or_keyword)  || is_name(slt_keyword) || is_name(sub_keyword)) {
		expect_token(TOKEN_KEYWORD);
		Register rs = parse_register();
		expect_token(',');
		Register rt = parse_register();
		expect_token(',');
		Register rd = parse_register();
		instr = new_r_instruction(add_keyword, rs, rt, rd);
	} else if (is_name(addi_keyword) || is_name(beq_keyword) || is_name(lw_keyword) || is_name(sw_keyword)) {
		expect_token(TOKEN_KEYWORD);
	} else if (is_name(j_keyword)) {
		expect_token(TOKEN_KEYWORD);
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
