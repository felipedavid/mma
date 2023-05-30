void instr_dec_error(const char *fmt, ...) {
	va_list args;
	printf("Instruction decoding error at line %d: ", token.line);
	va_start(args, fmt);
	vprintf(fmt, args);
	va_end(args);
	printf("\n");
}

bool assert_bit_size(u16 value, int n_bits) {
	u16 max_value = 0;
	for (int i = 0; i < n_bits; i++) {
		max_value |= (1 << n_bits);
	}
	return value <= max_value;
}

typedef enum {
	R_INSTRUCTION,
	I_INSTRUCTION,
	J_INSTRUCTION,
} InstructionKind;

typedef enum {
	MOD_NONE,
	MOD_ADD,
	MOD_ADDI,
	MOD_AND,
	MOD_BEQ,
	MOD_J,
	MOD_LW,
	MOD_NAND,
	MOD_NOR,
	MOD_OR,
	MOD_SLT,
	MOD_SW,
	MOD_SUB,
} InstructionMod;

typedef enum {
	R0,
	R1,
	R2,
	R3,
	R4,
	R5,
	R6,
	R7,
	REG_FILE_END,
} Register;

void validate_register(Register reg) {
	if (reg < R0 || reg >= REG_FILE_END) {
		instr_dec_error("Register '$r%d' is out of range. The architecture only supports registers from 1 to 7", reg);
	}
}

typedef struct {
	u16 rs;
	u16 rt;
	u16 rd;
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
	InstructionMod mod;
	union {
		RFormat r;
		IFormat i;
		JFormat j;
	};
} Instruction;

// This functions should be replaced by a hash table.. But I'm too lazy to do it right now
InstructionMod get_instr_mod(InstructionKind kind, const char *instr_name) {
	InstructionMod mod = MOD_NONE;
	if (kind == R_INSTRUCTION) {
		if (instr_name == add_keyword) {
			mod = MOD_ADD;
		} else if (instr_name == and_keyword) {
			mod = MOD_AND;
		} else if (instr_name == nand_keyword) {
			mod = MOD_NAND;
		} else if (instr_name == nor_keyword) {
			mod = MOD_NOR;
		} else if (instr_name == or_keyword) {
			mod = MOD_OR;
		} else if (instr_name == slt_keyword) {
			mod = MOD_SLT;
		} else if (instr_name == sub_keyword) {
			mod = MOD_SUB;
		}
	} else if (kind == I_INSTRUCTION) {
		if (instr_name == addi_keyword) {
			mod = MOD_ADDI;
		} else if (instr_name == beq_keyword) {
			mod = MOD_BEQ;
		} else if (instr_name == lw_keyword) {
			mod = MOD_LW;
		} else if (instr_name == sw_keyword) {
			mod = MOD_SW;
		}
	} else if (kind == J_INSTRUCTION && instr_name == j_keyword) {
		mod = MOD_J;
	} 

	if (mod == MOD_NONE) {
		instr_dec_error("There is no %d instruction called '%s'", kind, instr_name);
	}
	return mod;
}

Instruction *new_instruction(InstructionKind kind) {
	Instruction *instr = calloc(1, sizeof(Instruction));
	instr->kind = kind;
	return instr;
}

Instruction *new_r_instruction(const char *instr_name, u16 rs, u16 rt, u16 rd) {
	Register r_rs = (Register)rs;
	Register r_rt = (Register)rt;
	Register r_rd = (Register)rd;

	validate_register(r_rs);
	validate_register(r_rt);
	validate_register(r_rd);

	Instruction *instr = new_instruction(R_INSTRUCTION);
	instr->mod = get_instr_mod(R_INSTRUCTION, instr_name);
	instr->r.rs = r_rs;
	instr->r.rt = r_rt;
	instr->r.rd = r_rd;

	return instr;
}

Instruction *new_i_instruction(const char *instr_name, u16 rs, u16 rt, u16 immd) {
	Register r_rs = (Register)rs;
	Register r_rt = (Register)rt;

	validate_register(r_rs);
	validate_register(r_rt);

	if (!assert_bit_size(immd, 6)) {
		instr_dec_error("immediates for i-instructions cannot be longer than 6 bits");
	}

	Instruction *instr = new_instruction(I_INSTRUCTION);
	instr->mod = get_instr_mod(I_INSTRUCTION, instr_name);
	instr->i.rs = r_rs;
	instr->i.rt = r_rt;
	instr->i.immd = immd;

	return instr;
}

Instruction *new_j_instruction(const char *instr_name, u16 addr) {
	if (!assert_bit_size(addr, 12)) {
		instr_dec_error("addresses for jump instructions cannot be longer than 12 bits");
	}
	Instruction *instr = new_instruction(J_INSTRUCTION);
	instr->mod = get_instr_mod(J_INSTRUCTION, instr_name);
	instr->j.address = addr;
	return instr;
}


u16 encode_r_instruction(Instruction *instr) {
	u16 val = 0;
	val |= (instr->r.rd << 3);
	val |= (instr->r.rt << 6);
	val |= (instr->r.rs << 9);

	// sets the op/funct fields
	switch (instr->mod) {
	case MOD_ADD:
		break;
	case MOD_AND:
		val |= 4;
		break;
	case MOD_NAND:
		val |= 1;
		break;
	case MOD_NOR:
		val |= 7;
		break;
	case MOD_OR:
		val |= 5;
		break;
	case MOD_SLT:
		val |= 2;
		break;
	case MOD_SUB:
		val |= 3;
		break;
	default:
		instr_dec_error("Invalid R-type instruction");
		return 0;
	}

	return val;
}

u16 encode_i_instruction(Instruction *instr) {
	u16 val = 0;
	val |= (instr->i.immd);
	val |= (instr->i.rt << 6);
	val |= (instr->i.rs << 9);

	switch (instr->mod) {
	case MOD_ADDI:
		val |= (8 << 12);
		break;
	case MOD_BEQ:
		val |= (4 << 12);
		break;
	case MOD_LW:
		val |= (3 << 12);
		break;
	case MOD_SW:
		val |= (0xa << 12);
		break;
	default:
		instr_dec_error("Invalid I-type instruction");
		break;
	}

	return val;
}

u16 encode_j_instruction(Instruction *instr) {
	u16 val = 0;
	if (instr->mod != MOD_J) {
		instr_dec_error("Invalid J-type instruction");
	}

	val |= (instr->j.address);
	val |= (2 << 12);	

	return val;
}

u16 encode_instruction(Instruction *instr) {
	switch (instr->kind) {
	case R_INSTRUCTION:
		return encode_r_instruction(instr);
	case I_INSTRUCTION:
		return encode_i_instruction(instr);
	case J_INSTRUCTION:
		return encode_j_instruction(instr);
	default:
		instr_dec_error("Invalid instruction type");
	}
}
