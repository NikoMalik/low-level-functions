#include "textflag.h"



#define SYS_CLONE 56
#define SYS_WRITE 1
#define SYS_EXIT 60

// signature: func rawClone(flags uintptr, stack unsafe.Pointer) int
TEXT ·rawClone(SB), NOSPLIT, $0-16
    MOVQ    flags+0(FP), DI     //load flags
    MOVQ    stack+8(FP), SI      // pointer to stack
    LEAQ    threadFunc<>(SB), DX // pointer to function

    XORQ    R10, R10              // parent_tid
    XORQ    R8, R8                // child_tid
    XORQ    R9, R9                // TLS

    MOVQ    $SYS_CLONE, AX
    SYSCALL

    TESTQ   AX, AX
    JZ      new_thread

    MOVQ    AX, ret+16(FP)
    RET

new_thread:
    MOVQ    DX, AX
    CALL    AX                   
    XORQ    DI, DI
    MOVQ    $SYS_EXIT, AX
    SYSCALL

// function pointer to thread
TEXT threadFunc<>(SB), NOSPLIT, $0-0
    MOVQ    $SYS_WRITE, AX
    MOVQ    $1, DI           // stdout
    LEAQ    msg<>(SB), SI    // pointer to message
    MOVQ    msg_len<>(SB), DX // len message
    SYSCALL

    MOVQ    $SYS_EXIT, AX
    XORQ    DI, DI           // exit code 0
    SYSCALL


DATA msg<>+0(SB)/8, $"Hello fr"
DATA msg<>+8(SB)/8, $"om asm "
DATA msg<>+16(SB)/8, $"thread\n"
GLOBL msg<>(SB), RODATA, $24

DATA msg_len<>(SB)/8, $24
GLOBL msg_len<>(SB), RODATA, $8

