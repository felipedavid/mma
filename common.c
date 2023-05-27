typedef unsigned char u8;
typedef unsigned short u16;
typedef unsigned int u32;
typedef unsigned long long int u64;

void fatal(const char *fmt, ...) {
	va_list args;
	va_start(args, fmt);
	vprintf(fmt, args);
	va_end(args);
	exit(1);
}

void *xmalloc(size_t size) {
	void *ptr = malloc(size);
	if (ptr == NULL) {
		fatal("xmalloc failed");
	}
	return ptr;
}

void *xrealloc(void *ptr, size_t n_bytes) {
	ptr = realloc(ptr, n_bytes);
	if (ptr == NULL) {
		fprintf(stderr, "xrealloc failed");
		exit(1);
	}
	return ptr;
}

bool is_space(char ch) {
	return ch == ' ' || ch == '\r' || ch == '\t';
}

bool is_digit(char ch) {
	return (ch >= '0') && (ch <= '9');
}

bool is_alpha(char ch) {
	return ((ch >= 'a') && (ch <= 'z')) || ((ch >= 'A') && (ch <= 'Z'));
}

bool is_alphanum(char ch) {
	return is_digit(ch) || is_alpha(ch);
}

// STRETCHY BUFFERS 
typedef struct {
	size_t len;
	size_t cap;
	char content[0];
} Buf_Header;

#define BUF(x) x
#define buf__header(b) ((b) ? (Buf_Header *)((char*)(b) - offsetof(Buf_Header, content)) : 0)
#define buf__fits(b, n)
#define buf__make_fit(b, n) (((buf_len(b)+n) < buf_cap(b)) ? 0 : ((b) = buf__grow((b), (buf_len(b)+n), sizeof(*b))))

#define buf_len(b) ((b) ? buf__header(b)->len : 0)
#define buf_cap(b) ((b) ? buf__header(b)->cap : 0)
#define buf_end(b) ((b) + buf_len(b))
#define buf_push(b, ...) (buf__make_fit((b), 1), (b)[buf__header(b)->len++] = (__VA_ARGS__))
#define buf_free(b) ((b) ? (free(buf__header(b)), (b) = NULL) : 0)

#define MAX(x, y) (((x) > (y)) ? (x) : (y))

void *buf__grow(const void *buf, size_t min_cap, size_t elem_size) {
	size_t new_cap = MAX(buf_cap(buf) * 2 + 1, min_cap);
	assert(min_cap <= new_cap);

	size_t new_size = sizeof(Buf_Header) + (elem_size * new_cap);

	Buf_Header *new_buffer = NULL;
	if (buf == NULL) {
		new_buffer = xmalloc(new_size);
		new_buffer->len = 0;
	}
	else {
		new_buffer = xrealloc(buf__header(buf), new_size);
	}
	new_buffer->cap = new_cap;

	return (void *)new_buffer->content;
}

// String interning
typedef struct {
	size_t len;
	const char *str;
} Intern;

Intern *interns;

const char *str_intern_range(const char *start, const char *end) {
	size_t len = (end - start);
	for (Intern * it = interns; it != buf_end(interns); it++) {
		if (it->len == len && strncmp(it->str, start, len) == 0) {
			return it->str;
		}
	}

	char *str = malloc(len + 1);
	memcpy(str, start, len);
	str[len] = '\0';

	buf_push(interns, ((Intern) {len, str}));

	return str;
}

const char *str_intern(const char *str) {
	return str_intern_range(str, (str + strlen(str)));
}
