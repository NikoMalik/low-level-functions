#include "asm.h"
#include "textflag.h"

TEXT ·containsStringAVX2(SB), NOSPLIT, $40-33
	LOAD_PARAM(rsi, slicePtr+0(FP))
	LOAD_PARAM(rdx, sliceLen+8(FP))
	LOAD_PARAM(rdi,valuePtr+16(FP))
	LOAD_PARAM(rcx,valueLen+24(FP))

	// Save registers
	MOV(rbx, 0(rsp))
	MOV(r12, 8(rsp))
	MOV(r13,16(rsp))
	MOV(r14,24(rsp))
	MOV(r15,32(rsp))

	// Handle empty cases
	IF_ZERO(rdx,not_found)
	IF_ZERO(rcx,empty_value)

	MOV((rdi),r8) // needle data pointer
	XOR(r9,r9)    // index

element_loop:
	// Calculate element offset (16 bytes per element)
	MOV(r9,r14)

	// R14 = index * 16
	MUL_SHLQ(r14,$4)

	// Load string element
	// Element address
	LEA((rsi)(r14*1),r15)

	// String data pointer
	MOV((r15),r10)

	// String length
	MOV(8(r15),r11)

	// Skip zero-length elements
	IF_ZERO(r10,next_element)

	// Compare lengths
	CMP(r11,rcx)
	JNE next_element

	// Prepare for comparison
	// Comparison length
	MOV(r11,r12)

	// AVX comparison for >=64 byte blocks
	CMP(r12,$64)
	JAE avx_compare

scalar_compare:
	XOR(rax,rax)

scalar_loop:
	MOVB (r10)(rax*1), bl
	MOVB (r8)(rax*1), r13b
	CMPB bl, r13b
	JNE  next_element
	INC(rax)
	CMP(rax,r12)
	JB   scalar_loop       // rax < r12
	JMP  found

avx_compare:
	MOV(r12,r13)

	// 64-byte blocks
	DIV_SHRQ(r13,$6)

	// Save needle pointer
	MOV(r8,r15)

	// Save element pointer
	MOV(r10,r14)

avx_loop:
	VMOV((r10),y0)
	VMOV((r8),y1)
	VPCMPEQ(y0,y1,y2)
	VMOVMSK(y2,rax)
	ANY_NEQ(rax,avx_mismatch)

	VMOV(32(r10),y0)
	VMOV(32(r8),y1)
	VPCMPEQ(y0,y1,y2)
	VMOVMSK(y2,rax)
	ANY_NEQ(rax,avx_mismatch)
	ADD(r10,$64)
	ADD(r8,$64)
	DEC(r13)
	JNZ avx_loop

	// Handle remaining bytes (0-63)
	MOV(r12,r13)
	AND(r13,$63)
	JZ avx_match

	// Handle 32-63 bytes
	CMP(r13,$32)
	JB less_than_32 // if r13<32
	VMOV((r10),y0)
	VMOV((r8),y1)
	VPCMPEQ(y0,y1,y2)
	VMOVMSK(y2,rax)
	ANY_NEQ(rax,avx_mismatch)
	ADD(r10,$32)
	ADD(r8,$32)
	SUB(r13,$32)

less_than_32:
	IF_ZERO(r13,avx_match)
	XOR(rax,rax)

scalar_tail:
	MOVB (r10)(rax*1), bl
	MOVB (r8)(rax*1), r13b
	CMPB bl, r13b
	JNE  avx_mismatch
	INC(rax)
	CMP(rax,r13)
	JB   scalar_tail

avx_match:
	JMP found

avx_mismatch:
	// Restore needle pointer
	MOV(r15,r8)

	// Restore element pointer
	MOV(r14,r10)

next_element:
	INC(r9)
	CMP(r9,rdx)
	JB  element_loop
	JMP not_found

found:
	MOVB $1, ret+32(FP)
	JMP  exit

empty_value:
	XOR(r9,r9)

empty_loop:
	MOV(r9,r14)

	// 16 bytes per element (string)
	MUL_SHLQ(r14,$4)

	// R15 = &slice[R9]
	LEA((rsi)(r14*1),r15)

	// String length
	MOV(8(r15),r11)
	IF_ZERO(r11,found_empty)
	INC(r9)
	CMP(r9,rdx)
	JB  empty_loop
	JMP not_found

found_empty:
	MOVB $1, ret+32(FP)
	JMP  exit

not_found:
	MOVB $0, ret+32(FP)

exit:
	VZEROUPPER
	MOV(0(rsp),rbx)
	MOV(8(rsp),r12)
	MOV(16(rsp),r13)
	MOV(24(rsp),r14)
	MOV(32(rsp),r15)
	RET

// func containsByteSliceAVX2(slicePtr unsafe.Pointer, sliceLen int, valuePtr unsafe.Pointer, valueLen int) bool
TEXT ·containsByteSliceAVX2(SB), NOSPLIT, $40-33
	MOVQ slicePtr+0(FP), SI
	MOVQ sliceLen+8(FP), DX
	MOVQ valuePtr+16(FP), DI
	MOVQ valueLen+24(FP), CX

	MOVQ BX, 0(SP)
	MOVQ R12, 8(SP)
	MOVQ R13, 16(SP)
	MOVQ R14, 24(SP)
	MOVQ R15, 32(SP)

	TESTQ DX, DX
	JZ    not_found
	TESTQ CX, CX
	JZ    empty_value

	MOVQ (DI), R8
	MOVQ 8(DI), CX

	XORQ R9, R9

element_loop:
	//  R14 = R9 * 24
	MOVQ R9, R14
	SHLQ $3, R14           // R14 = R9 * 8
	LEAQ (R14)(R14*2), R14 // R14 = R9 * 24

	LEAQ (SI)(R14*1), R15
	MOVQ (R15), R10
	MOVQ 8(R15), R11

	TESTQ R10, R10
	JZ    next_element

	CMPQ R11, CX
	JNE  next_element

	MOVQ R11, R12

	CMPQ R12, $64
	JAE  avx_compare

scalar_compare:
	XORQ AX, AX

scalar_loop:
	MOVB (R10)(AX*1), BL
	CMPB BL, (R8)(AX*1)
	JNE  next_element
	INCQ AX
	CMPQ AX, R12
	JB   scalar_loop
	JMP  found

avx_compare:
	MOVQ R12, R13
	SHRQ $6, R13  // 64 bytes blocks

avx_loop:
	VMOVDQU   (R10), Y0
	VMOVDQU   (R8), Y1
	VPCMPEQB  Y0, Y1, Y2
	VPMOVMSKB Y2, AX
	CMPL      AX, $0xFFFFFFFF
	JNE       next_element

	VMOVDQU   32(R10), Y0
	VMOVDQU   32(R8), Y1
	VPCMPEQB  Y0, Y1, Y2
	VPMOVMSKB Y2, AX
	CMPL      AX, $0xFFFFFFFF
	JNE       next_element

	ADDQ $64, R10
	ADDQ $64, R8
	DECQ R13
	JNZ  avx_loop

	// 0-63
	MOVQ R12, R13
	ANDQ $63, R13
	JZ   found

	// 32
	CMPQ      R13, $32
	JB        less_than_32
	VMOVDQU   (R10), Y0
	VMOVDQU   (R8), Y1
	VPCMPEQB  Y0, Y1, Y2
	VPMOVMSKB Y2, AX
	CMPL      AX, $0xFFFFFFFF
	JNE       next_element
	ADDQ      $32, R10
	ADDQ      $32, R8
	SUBQ      $32, R13

less_than_32:
	// SCALAR OSTATOK TRASH
	TESTQ R13, R13
	JZ    found
	XORQ  AX, AX

scalar_tail:
	MOVB (R10)(AX*1), BL
	CMPB BL, (R8)(AX*1)
	JNE  next_element
	INCQ AX
	CMPQ AX, R13
	JB   scalar_tail
	JMP  found

next_element:
	// revive pointer value
	MOVQ valuePtr+16(FP), DI
	MOVQ (DI), R8            // R8 = pointer

	// NEXT
	INCQ R9
	CMPQ R9, DX
	JB   element_loop
	JMP  not_found

found:
	MOVB $1, ret+32(FP)
	JMP  exit

empty_value:
	// try to find this empty shit
	XORQ R9, R9

empty_loop:
	MOVQ R9, R14
	SHLQ $3, R14
	LEAQ (R14)(R14*2), R14 // R14 = R9 * 24
	LEAQ (SI)(R14*1), R15  // R15 = &slice[R9]
	MOVQ 8(R15), R11       // R11 = len element
	CMPQ R11, $0
	JE   found_empty
	INCQ R9
	CMPQ R9, DX
	JB   empty_loop
	JMP  not_found

found_empty:
	MOVB $1, ret+32(FP)
	JMP  exit

not_found:
	MOVB $0, ret+32(FP)

exit:
	// RETURN REGISTERS
	VZEROUPPER
	MOVQ 0(SP), BX
	MOVQ 8(SP), R12
	MOVQ 16(SP), R13
	MOVQ 24(SP), R14
	MOVQ 32(SP), R15
	RET
