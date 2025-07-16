// compare.s
#include "textflag.h"

// func CompareImpl(a, b unsafe.Pointer, len int) bool
TEXT ·CompareImpl(SB), NOSPLIT, $0-25
    MOVQ a+0(FP), SI    
    MOVQ b+8(FP), DI   
    MOVQ len+16(FP), CX 

    TESTQ CX, CX
    JE    equal

    MOVQ CX, AX      
    SHRQ $5, AX       
    TESTQ AX, AX       

    JE    remainder     

avx_loop:
    VMOVUPS (SI), Y0  
    VMOVUPS (DI), Y1   
    VPCMPEQB Y0, Y1, Y2 
    VPMOVMSKB Y2, BX   
    CMPL BX, $0xFFFFFFFF 
    JNE   not_equal  
    ADDQ $32, SI     
    ADDQ $32, DI      
    DECQ AX            
    JNZ   avx_loop      

remainder:
    MOVQ CX, AX         
    ANDQ $31, AX        
    TESTQ AX, AX  
    JE    equal    
    MOVQ AX, CX     
    MOVQ CX, DX
    SHRQ $3, DX      
    TESTQ DX, DX
    JE    byte_remainder
remainder_loop:
    MOVQ (SI), R8
    MOVQ (DI), R9
    CMPQ R8, R9
    JNE   not_equal
    ADDQ $8, SI
    ADDQ $8, DI
    DECQ DX
    JNZ   remainder_loop
byte_remainder:
    MOVQ CX, AX
    ANDQ $7, AX       
    TESTQ AX, AX
    JE    equal
    MOVQ AX, CX
    REP; CMPSB
    JE    equal

not_equal:
    MOVB $0, ret+24(FP)
    RET

equal:
    MOVB $1, ret+24(FP) 
    RET
