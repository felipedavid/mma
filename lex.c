const char **keywords;

const char *add_keyword;
const char *addi_keyword;
const char *and_keyword;
const char *beq_keyword;
const char *j_keyword;
const char *lw_keyword;
const char *nand_keyword;
const char *nor_keyword;
const char *or_keyword;
const char *slt_keyword;
const char *sw_keyword;
const char *sub_keyword;

#define KEYWORD(name) name##_keyword = str_intern(#name); buf_push(keywords, name##_keyword)

void init_keywords() {
	KEYWORD(add);
	KEYWORD(addi);
	KEYWORD(and);
	KEYWORD(beq);
	KEYWORD(j);
	KEYWORD(lw);
	KEYWORD(nand);
	KEYWORD(nor);
	KEYWORD(or);
	KEYWORD(slt);
	KEYWORD(sw);
	KEYWORD(sub);
}

bool is_keyword_name(const char *name) {
	for (const char **it = keywords; it != buf_end(keywords); it++) {
		if (*it == name) {
			return true;
		}
	}
	return false;
}

void syntax_error(const char *fmt, ...);

typedef enum {
	TOKEN_NONE = 0,
	TOKEN_IDENTIFIER = 128,
	TOKEN_NUMBER,
	TOKEN_DIRECTIVE,
	TOKEN_IDENTIFIER_COLON,
	TOKEN_REGISTER,
	TOKEN_KEYWORD,
} TokenKind;

typedef struct {
	TokenKind kind;
	const char *start;
	const char *end;
	int line;
	union {
		int int_value;
		const char *name;
		int register_val;
	};
} Token;

Token token;
const char *stream;

const char *token_kind_str(TokenKind kind) {
	static char buf[256]; // OBS: will be overwritten by next call
	switch (kind) {
	case TOKEN_NONE:
		return "none";
	case TOKEN_IDENTIFIER:
		return "identifier";
	case TOKEN_IDENTIFIER_COLON:
		return "identifier followed by a colon";
	case TOKEN_DIRECTIVE:
		return "directive";
	case TOKEN_NUMBER:
		return "number";
	case TOKEN_REGISTER:
		return "register";
	case TOKEN_KEYWORD:
		return "keyword";
	default:
		if (kind < 128 && isprint(kind)) {
			snprintf(buf, sizeof(buf), "'%c'", kind);
		} else {
			snprintf(buf, sizeof(buf), "<ASCII %d>", kind);
		}
	}
	return buf;
}

void next_token() {
BEGIN:
	token.start = stream;
	switch (*stream) {
	case '#':
		stream++;
		while (*stream != '\n') {
			stream++;
		}
		goto BEGIN;
	case '\n':
		stream++;
		token.line++;
		goto BEGIN;
	case ' ':
		while (is_space(*stream)) {
			stream++;
		}
		goto BEGIN;
	case '0' ... '9': {
		int val = 0;
		while (is_digit(*stream)) {
			val *= 10;
			val += *stream - '0';
			stream++;
		}
		token.kind = TOKEN_NUMBER;
		token.int_value = val;
	} break;
	case 'a' ... 'z':
	case 'A' ... 'Z':
	case '_': {
		while (is_alphanum(*stream) || *stream == '_') {
			stream++;
		}

		if (*stream == ':') {
			stream++;
			token.kind = TOKEN_IDENTIFIER_COLON;
			token.name = str_intern_range(token.start, stream-1);
		} else {
			token.kind = TOKEN_IDENTIFIER;
			token.name = str_intern_range(token.start, stream);

			if (is_keyword_name(token.name)) {
				token.kind = TOKEN_KEYWORD;
			}
		}
	} break;
	case '.': {
		stream++; 
		bool ok = false;
		while (is_alphanum(*stream) || *stream == '_') {
			ok = true;
			stream++;
		}
		if (!ok) {
			syntax_error("invalid directive '%.*s'", (int)(stream - token.start), token.start);
		}

		token.kind = TOKEN_DIRECTIVE;
		token.name = str_intern_range(token.start+1, stream);
	} break;
	case '$': {
		stream++;
		if (*stream == 'r' || *stream == 'R') {
			stream++;
		}

		if (is_digit(*stream)) {
			int val = 0;
			while (is_digit(*stream)) {
				val *= 10;
				val += *stream - '0';
				stream++;
			}

			assert(val >= 0);
			if (val > 7) {
				syntax_error("register '%.*s' does not exist", (int)(stream - token.start), token.start);
			}
			token.kind = TOKEN_REGISTER;
			token.register_val = val;
		} else {
			syntax_error("invalid register");
		}
    } break;
	default:
		token.kind = *stream++;
		break;
	}
	token.end = stream;
};

void print_token(Token token) {
	printf("TOKEN{ type: %s, lexeme: \"%.*s\", line: %d, ",
			token_kind_str(token.kind), (int)(token.end - token.start), token.start, token.line);
	switch (token.kind) {
	case TOKEN_IDENTIFIER_COLON:
		printf("val: \"%s\"", token.name);
		break;
	case TOKEN_IDENTIFIER:
		printf("val: \"%s\"", token.name);
		break;
	case TOKEN_NUMBER:
		printf("val: %d", token.int_value);
		break;
	case TOKEN_REGISTER:
		printf("val: %d", token.register_val);
		break;
	default:
		break;
	}
	printf("}\n");
}

void init_stream(const char *source) {
	stream = source;
	token.line = 1;
	next_token();
}

bool is_token(TokenKind kind) {
	return token.kind == kind;
}


bool match_token(TokenKind kind) {
	if (token.kind == kind) {
		next_token();
		return true;
	}
	return false;
}

bool expect_token(TokenKind kind) {
	if (token.kind == kind) {
		next_token();
		return true;
	}
	syntax_error("expected %s, got %s", token_kind_str(kind), token_kind_str(token.kind));
	return false;
}

void syntax_error(const char *fmt, ...) {
	va_list args;
	printf("Syntax error at line %d: ", token.line);
	va_start(args, fmt);
	vprintf(fmt, args);
	va_end(args);
	printf("\n");
}

