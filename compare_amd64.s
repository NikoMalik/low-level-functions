#include "textflag.h"

// func CompareImpl(a, b unsafe.Pointer, len int) bool
TEXT ·CompareImpl(SB), NOSPLIT, $0-25
    MOVQ a+0(FP), SI
    MOVQ b+8(FP), DI
    MOVQ len+16(FP), CX
    XORQ AX, AX              // AX = 0 
    TESTQ CX, CX
    // JZ equal                

    MOVQ CX, DX
    SHRQ $5, DX              
    JZ less_than_32         

loop32:
    VMOVDQU (SI), Y0
    VMOVDQU (DI), Y1
    VPCMPEQB Y0, Y1, Y2
    VPMOVMSKB Y2, BX
    CMPL BX, $0xFFFFFFFF
    JNE not_equal
    ADDQ $32, SI
    ADDQ $32, DI
    DECQ DX
    JNZ loop32
    ANDQ $31, CX             

//16-31
less_than_32:
    MOVQ CX, DX
    SHRQ $4, DX              
    JZ less_than_16
    VMOVDQU (SI), X0
    VMOVDQU (DI), X1
    VPCMPEQB X0, X1, X2
    VPMOVMSKB X2, BX
    CMPW BX, $0xFFFF
    JNE not_equal
    ADDQ $16, SI
    ADDQ $16, DI
    ANDQ $15, CX            

less_than_16:
    TESTQ $8, CX
    JZ less_than_8
    MOVQ (SI), R8
    MOVQ (DI), R9
    CMPQ R8, R9
    JNE not_equal
    ADDQ $8, SI
    ADDQ $8, DI

less_than_8:
    TESTQ $4, CX
    JZ less_than_4
    MOVL (SI), R8
    MOVL (DI), R9
    CMPL R8, R9
    JNE not_equal
    ADDQ $4, SI
    ADDQ $4, DI

less_than_4:
    TESTQ $3, CX
    JZ equal

    TESTQ $2, CX            
    JZ last_byte
    MOVW (SI), R8
    MOVW (DI), R9
    CMPW R8, R9
    JNE not_equal
    ADDQ $2, SI
    ADDQ $2, DI

last_byte:
    TESTQ $1, CX           
    JZ equal
    MOVB (SI), AL
    MOVB (DI), BL
    CMPB AL, BL
    JNE not_equal
    JMP equal

not_equal:
    MOVB $0, ret+24(FP)
    VZEROUPPER
    RET

equal:
    MOVB $1, ret+24(FP)
    VZEROUPPER
    RET
