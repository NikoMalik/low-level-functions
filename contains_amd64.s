#include "textflag.h"
#include "go_asm.h"

// func Contains(a []byte, b []byte) bool

TEXT ·Contains(SB),7,$0
    MOVQ    a+0(FP), SI       
    MOVQ    b+24(FP), DI      
    MOVQ    a_len+8(FP), AX   
    MOVQ    b_len+32(FP), BX  

    TESTQ   BX, BX            
    JE      ret_true          

    CMPQ    AX, BX            
    JB      ret_false         

contains_loop:
    MOVQ    AX, CX            
    SUBQ    BX, CX            
    INCQ    CX                

search_loop:
    DECQ    CX                
    JS      ret_false         

    MOVBQZX 0(SI), R8         
    MOVBQZX 0(DI), R9         
    CMPQ    R8, R9            
    JNE     skip              

    MOVQ    BX, DX            
    LEAQ    1(SI), SI         
    LEAQ    1(DI), DI         
    REP; CMPSB                
    JE      ret_true          

skip:
    LEAQ    1(SI), SI         
    LEAQ    -1(DI), DI        
    JMP     search_loop       

ret_false:
    MOVQ    $0, ret+40(FP)      
    RET

ret_true:
    MOVQ    $1, ret+40(FP)     
    RET
	