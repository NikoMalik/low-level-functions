#include "textflag.h"
#include "go_asm.h"

// func Compare(a []byte, b []byte) int
TEXT ·Compare(SB),7,$0
    MOVQ    a+0(FP),SI      
    MOVQ    b+24(FP),DI     
    MOVQ    a_len+8(FP),AX   
    MOVQ    b_len+32(FP),BX  
    CMPQ    SI,DI
    JE      cmp_len          
    MOVQ    AX,CX            
    CMPQ    BX,CX            
    CMOVQLE BX,CX            
    CMPQ    CX,$32
    JA      cmpsb_elm
    INCQ    CX               

cmp_elm_loop:                
    DECQ    CX               
    JZ      cmp_len          
    MOVBQZX 0(SI),R8         
    MOVBQZX 0(DI),R9         
	CMPQ    R8,R9               
    JG      ret_greater      
    JL      ret_less         
	INCQ    SI                  
    INCQ    DI               
    JMP     cmp_elm_loop

cmpsb_elm:                   
    CLD                      
    REP; CMPSB               
    JE      cmp_len
    DECQ    SI               
    DECQ    DI               
    MOVBQZX 0(SI),R8         
    MOVBQZX 0(DI),R9         
    CMPQ    R8,R9            
    JG      ret_greater      
    JL      ret_less         
    JMP     ret_equal       

cmp_len:                    
    CMPQ    AX,BX           
    JG      ret_greater     
    JL      ret_less        

ret_equal:
    MOVQ    $0,ret+48(FP)   
    RET

ret_greater:
    MOVQ    $1,ret+48(FP)   
    RET

ret_less:
    MOVQ    $-1,ret+48(FP)  
    RET

