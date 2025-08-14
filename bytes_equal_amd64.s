#include "textflag.h"




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
    JZ not_found
    TESTQ CX, CX
    JZ empty_value

    MOVQ (DI), R8                
    MOVQ 8(DI), CX                
    
    XORQ R9, R9                    

element_loop:
    //  R14 = R9 * 24
    MOVQ R9, R14
    SHLQ $3, R14                  // R14 = R9 * 8
    LEAQ (R14)(R14*2), R14        // R14 = R9 * 24
    
    LEAQ (SI)(R14*1), R15       
    MOVQ (R15), R10              
    MOVQ 8(R15), R11              
    
    TESTQ R10, R10
    JZ next_element
    
    CMPQ R11, CX
    JNE next_element
    
    MOVQ R11, R12                 
    
    CMPQ R12, $64
    JAE avx_compare

scalar_compare:
    XORQ AX, AX                
scalar_loop:
    MOVB (R10)(AX*1), BL        
    CMPB BL, (R8)(AX*1)          
    JNE next_element
    INCQ AX
    CMPQ AX, R12
    JB scalar_loop
    JMP found

avx_compare:
    MOVQ R12, R13
    SHRQ $6, R13                  // 64 bytes blocks
avx_loop:
    VMOVDQU (R10), Y0
    VMOVDQU (R8), Y1
    VPCMPEQB Y0, Y1, Y2
    VPMOVMSKB Y2, AX
    CMPL AX, $0xFFFFFFFF
    JNE next_element
    
    VMOVDQU 32(R10), Y0
    VMOVDQU 32(R8), Y1
    VPCMPEQB Y0, Y1, Y2
    VPMOVMSKB Y2, AX
    CMPL AX, $0xFFFFFFFF
    JNE next_element
    
    ADDQ $64, R10
    ADDQ $64, R8
    DECQ R13
    JNZ avx_loop
    
    // 0-63
    MOVQ R12, R13
    ANDQ $63, R13
    JZ found
    
    // 32
    CMPQ R13, $32
    JB less_than_32
    VMOVDQU (R10), Y0
    VMOVDQU (R8), Y1
    VPCMPEQB Y0, Y1, Y2
    VPMOVMSKB Y2, AX
    CMPL AX, $0xFFFFFFFF
    JNE next_element
    ADDQ $32, R10
    ADDQ $32, R8
    SUBQ $32, R13

less_than_32:
    // SCALAR OSTATOK TRASH
    TESTQ R13, R13
    JZ found
    XORQ AX, AX
scalar_tail:
    MOVB (R10)(AX*1), BL
    CMPB BL, (R8)(AX*1)
    JNE next_element
    INCQ AX
    CMPQ AX, R13
    JB scalar_tail
    JMP found

next_element:
    // revive pointer value
    MOVQ valuePtr+16(FP), DI
    MOVQ (DI), R8                 // R8 = pointer
    
    // NEXT
    INCQ R9
    CMPQ R9, DX
    JB element_loop
    JMP not_found

found:
    MOVB $1, ret+32(FP)
    JMP exit

empty_value:
    //try to find this empty shit
    XORQ R9, R9
empty_loop:
    MOVQ R9, R14
    SHLQ $3, R14
    LEAQ (R14)(R14*2), R14        // R14 = R9 * 24
    LEAQ (SI)(R14*1), R15         // R15 = &slice[R9]
    MOVQ 8(R15), R11              // R11 = len element
    CMPQ R11, $0
    JE found_empty
    INCQ R9
    CMPQ R9, DX
    JB empty_loop
    JMP not_found

found_empty:
    MOVB $1, ret+32(FP)
    JMP exit

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
