// ================================================================
// Go ASM helpers (plan9 syntax, (x86-64))
// ================================================================
//
//
//
// I FUCKING HATE THIS DOT  (U+00B7)
#define DOT ·

// RAX, RBX, RCX, RDX, RDI, RSI, R8, R9, R10, R11  (caller-saved)
// 64bit
#define rax AX
#define rbx BX
#define rcx CX
#define rdx DX
#define rsi SI
#define rdi DI
#define rbp BP
#define rsp SP
#define r8 R8
#define r9 R9
#define r10 R10
#define r11 R11
#define r12 R12
#define r13 R13
#define r14 R14
#define r15 R15

// 32-bit
#define eax AX
#define ebx BX
#define ecx CX
#define edx DX
#define esi SI
#define edi DI
#define ebp BP
#define esp SP
#define r8d R8L
#define r9d R9L
#define r10d R10L
#define r11d R11L
#define r12d R12L
#define r13d R13L
#define r14d R14L
#define r15d R15L

// 16-bit
#define ax AX
#define bx BX
#define cx CX
#define dx DX
#define si SI
#define di DI
#define bp BP
#define sp SP
#define r8w R8W
#define r9w R9W
#define r10w R10W
#define r11w R11W
#define r12w R12W
#define r13w R13W
#define r14w R14W
#define r15w R15W

// 8-bit low
#define al AL
#define bl BL
#define cl CL
#define dl DL
#define sil SIL
#define dil DIL
#define bpl BPL
#define spl SPL
#define r8b R8B
#define r9b R9B
#define r10b R10B
#define r11b R11B
#define r12b R12B
#define r13b R13B
#define r14b R14B
#define r15b R15B

// 8-bit high (only for AX..DX)
#define ah AH
#define bh BH
#define ch CH
#define dh DH

// SIMD  (AVX-512)
// #define Z0 Z0
// #define Z1 Z1
// #define Z2 Z2
// #define Z3 Z3
// #define Z4 Z4
// #define Z5 Z5
// #define Z6 Z6
// #define Z7 Z7
// #define Z8 Z8
// #define Z9 Z9
// #define Z10 Z10
// #define Z11 Z11
// #define Z12 Z12
// #define Z13 Z13
// #define Z14 Z14
// #define Z15 Z15
// #define Z16 Z16
// #define Z17 Z17
// #define Z18 Z18
// #define Z19 Z19
// #define Z20 Z20
// #define Z21 Z21
// #define Z22 Z22
// #define Z23 Z23
// #define Z24 Z24
// #define Z25 Z25
// #define Z26 Z26
// #define Z27 Z27
// #define Z28 Z28
// #define Z29 Z29
// #define Z30 Z30
// #define Z31 Z31

// SIMD  (AVX2/YMM)
#define y0 Y0
#define y1 Y1
#define y2 Y2
#define y3 Y3
#define y4 Y4
#define y5 Y5
#define y6 Y6
#define y7 Y7
#define y8 Y8
#define y9 Y9
#define y10 Y10
#define y11 Y11
#define y12 Y12
#define y13 Y13
#define y14 Y14
#define y15 Y15

// SIMD  (SSE/XMM)
#define x0 X0
#define x1 X1
#define x2 X2
#define x3 X3
#define x4 X4
#define x5 X5
#define x6 X6
#define x7 X7
#define x8 X8
#define x9 X9
#define x10 X10
#define x11 X11
#define x12 X12
#define x13 X13
#define x14 X14
#define x15 X15

// segment registers
#define FSREG FS
#define GSREG GS

// math
#define ADD(dst, src) ADDQ src, dst
#define SUB(dst, src) SUBQ src, dst
#define MUL(reg) MULQ reg
#define DIV(reg) DIVQ reg
#define INC(reg) INCQ reg
// INC:: //for 64bit
// MOVQ $0, AX   // rax = 0
// INCQ AX       // rax = 1
// INCQ AX       // rax = 2
#define DEC(reg) DECQ reg
// MOVQ $5, AX   // rax = 5
// DECQ AX       // rax = 4
// DECQ AX       // rax = 3
#define NEG(reg) NEGQ reg

#define DIV_SHRQ(req, number) SHRQ number, req

#define MUL_SHLQ(req, number) SHLQ number, req

#define MUL2(reg) SHLQ $1, reg // mul  2
#define DIV2(reg) SHRQ $1, reg // div 2

// reg = reg * 8
#define MUL8(reg) SHLQ $3, reg

// reg = reg / 8
#define DIV8(reg) SHRQ $3, re

#define AND(dst, src) ANDQ src, dst
#define OR(dst, src) ORQ src, dst
#define XOR(dst, src) XORQ src, dst
#define NOT(reg) NOTQ reg

// moves
#define MOV(src, dst) MOVQ src, dst
#define MOVZX(src, dst) MOVLQZX src, dst
#define MOVSX(src, dst) MOVLQSX src, dst
#define LEA(addr, reg) LEAQ addr, reg

// compare
#define CMP(a, b) CMPQ a, b
#define TEST(a, b) TESTQ a, b

// jumps
// #define JMP(label) JMP label
// #define JE(label) JE label
// #define JNE(label) JNE label
// #define JL(label) JL label
// #define JZ(label) JZ label
// #define JG(label) JG label
// #define JLE(label) JLE label
// #define JGE(label) JGE label
// #define JNZ(label) JNZ label
// #define JC(label) JC label
// #define JNC(label) JNC label

// #define CALL(addr) CALL addr
// #define RET RET

// stack
#define PUSH(reg) PUSHQ reg
#define POP(reg) POPQ reg

// shlq $n, reg → reg << n (mul 2ⁿ)
//
// shrq $n, reg → reg >> n (div 2ⁿ)

// SIMD
#define VMOV(src, dst) VMOVDQU src, dst
#define VADD(src, dst) VPADDQ src, dst, dst
#define VEXTRACTI128_HI(src, dst) VEXTRACTI128 $1, src, dst
#define VPSHUFD_SWAPQ(src, dst) VPSHUFD $0x4E, src, dst
#define VMOVQ_XMM_TO_GPR(src, dst) VMOVQ src, dst
#define VZERO(dst) VPXOR dst, dst, dst

#define VSUB(src, dst) VPSUBQ src, dst, dst
#define VMUL(src, dst) VPMULLQ src, dst, dst
#define VAND(src, dst) VPAND src, dst, dst
#define VOR(src, dst) VPOR src, dst, dst
#define VXOR(src, dst) VPXOR src, dst, dst
#define VSHL(count, dst) VPSLLQ count, dst, dst
#define VSHR(count, dst) VPSRLQ count, dst, dst
#define VPADDQ3(src1, src2, dst) VPADDQ src1, src2, dst

#define MEMCPY_16(src, dst) \
    MOVOU(src), X0;         \
    MOVOU X0, (dst)

#define MEMSET_256(value, dst) \
    MOVQ value, RAX;           \
    MOVQ dst, RDI;             \
    MOVQ $32, RCX;             \
    REP;                       \
    STOSQ

#define PRINT_REG(reg) \
    MOVQ reg, 0(SP);   \
    CALL ·printInt(SB)

#define GO_CALL(func, args) \
    MOVQ func, AX;          \
    MOVQ args, DI;          \
    CALL AX

// registers flags
#define FLAGS EFLAGS

#define CENTER_POINT(label) \
    label:                  \
    BYTE $0x90;             \
    BYTE $0x90;             \
    BYTE $0x90;             \
    BYTE $0x90;             \
    BYTE $0x90;             \
    BYTE $0x90;             \
    BYTE $0x90;             \
    BYTE $0x90

#define LOAD_PARAM(reg, offset) MOVQ offset(FP), reg

// save register to stack
#define SAVE_REG(reg) MOVQ reg, saved_##reg(SP)

// restore register from stack
#define RESTORE_REG(reg) MOVQ saved_##reg(SP), reg

// swap value in registers
#define SWAP(reg1, reg2) XCHGQ reg1, reg2

// load 256-bit vector without cache
#define VMOVNT(addr, reg) VMOVNTDQ addr, reg

// compare (YMM)
#define VPCMPEQ(reg1, reg2, dst) VPCMPEQB reg1, reg2, dst

// create mask from vector compare
#define VMOVMSK(vec, reg) VPMOVMSKB vec, reg

// check (all 1 in mask)
#define ALL_EQ(reg)        \
    CMPL reg, $0xFFFFFFFF; \
    JEQ

#define IF_ZERO(req, label) \
    TESTQ req, req;         \
    JZ label;

#define IF_NOT_ZERO(reg, label) \
    TESTQ reg, reg;             \
    JNZ label

// check any incompare
#define ANY_NEQ(reg, label) \
    CMPL reg, $0xFFFFFFFF;  \
    JNE label

#define RETURN_FALSE           \
    MOVB $0, ret + offset(FP); \
    RET

#define RETURN_TRUE            \
    MOVB $1, ret + offset(FP); \
    RET

// clear high bits from  AVX registers
// #define VZEROUPPER VZEROUPPER

#define ALIGN_32 \
    BYTE $0x90;  \
    BYTE $0x90;  \
    BYTE $0x90;  \
    BYTE $0x90

// get len slice
// len_reg = *(slice_ptr + 8)
#define SLICE_LEN(slice_ptr, len_reg) \
    MOVQ 8(slice_ptr), len_reg

// get pointer to array data
#define SLICE_DATA(slice_ptr, data_reg) \
    MOVQ(slice_ptr), data_reg

// loop for elements
#define FOR_SLICE(idx_reg, len_reg, loop_label, end_label) \
    XORQ idx_reg, idx_reg;                                 \
    loop_label:                                            \
    CMPQ idx_reg, len_reg;                                 \
    JGE end_label

// lock
#define MEMORY_BARRIER MFENCE

#define PAUSE   \
    BYTE $0xF3; \
    BYTE $0x90

// Load immediate value
#define MOV_IMM(value, reg) MOVQ $value, reg

// Atomic operations
#define ATOMIC_ADD(dst, value) \
    LOCK;                      \
    ADDQ value, dst
#define ATOMIC_XCHG(src, dst) \
    LOCK;                     \
    XCHGQ src, dst
#define ATOMIC_CMPXCHG(old, new, dst) \
    LOCK;                             \
    CMPXCHGQ new, dst

// Prefetch data
#define PREFETCH(addr) PREFETCHT0(addr)

// Bit manipulation
#define BSF(src, dst) BSFQ src, dst
#define BSR(src, dst) BSRQ src, dst
#define TZCNT(src, dst) TZCNTQ src, dst
#define LZCNT(src, dst) LZCNTQ src, dst

// Conditional move
#define CMOVNE(src, dst) CMOVQNE src, dst
#define CMOVEQ(src, dst) CMOVQEQ src, dst

// name — name array
// size — count array
// val — value (byte)
#define DEFINE_BYTE_ARRAY(name, size, val) \
    GLOBL name(SB), RODATA, $size;         \
    name:                                  \
    REPT size                              \
        BYTE $val                          \
            END

// array 64-bit word
#define DEFINE_QWORD_ARRAY(name, size, val) \
    GLOBL name(SB), RODATA, $(size * 8);    \
    name:                                   \
    REPT size                               \
        DQ val                              \
            END

#define GARR(base, idx, sz) base + ((idx) * (sz))(SB)
