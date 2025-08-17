#include "asm.h"
#include "textflag.h"

TEXT ·SumUnroll8(SB), NOSPLIT, $0-32
	// Load slice parameters
	LOAD_PARAM(rax,arr_data+0(FP))
	LOAD_PARAM(rcx,arr_len+8(FP))

	// Check length and handle small cases
	CMP(rcx,$8)     // if len(array) < 8 jmp to small_vector
	JL small_vector

	// Initialize accumulators
	VZERO(y0)
	VZERO(y1)
	XOR(rsi,rsi) // 0

	// Main loop: process 8 elements per iteration
	// copy rcx data to rdx
	MOV(rcx,rdx)    // using rdx for math operations
	DIV_SHRQ(rdx,$3)// DX = len / 8
	MUL_SHLQ(rdx,$3)// DX = number of elements to process (multiple of 8)
	JE small_vector // Jump if no full groups

loop:
	// Load 8 elements
	// rax + rsi*8 //0*8 first 8 elements
	VMOVDQU (rax)(rsi*8), y2 // int64 is 8 bytes load to y2     // LOAD=4elements

	VMOVDQU 32(rax)(rsi*8), y3 // LOAD=4elements

	// (rax)(rsi*8) load arr[0..3] (Y2), // because ymm registers only can hold 256bit
	// 32(rax)(rsi*8) load arr[4..7] (Y3).
	//  arr[rsi..rsi+3] in Y2
	//  arr[rsi+4..rsi+7] in Y3
	// address = rax + rsi*8 + 32

	// Accumulate
	VADD(y2,y0) // sum 4 elements and save  to y0
	VADD(y3,y1) // sum 4 elements and save to y1

	// Next group
	// rsi +8 every loop cycle, so in first iteration rsi is 0,in next 8
	ADD(rsi,$8)// rsi +=8
	CMP(rsi,rdx)// check rsi index with len array
	JB loop // if rsi < rdx

	// Combine accumulators
	// SUM ALL RESULTS
	VADD(y1,y0)

small_vector:
	// Handle remaining elements (0-7)
	MOV(rcx,rdx)
	SUB(rdx,rsi) // rdx = remaining elements

	// rdx = len(array) - processed_index | rdx = len(array)
	// if rdx == 0 jmp to store_result
	JE store_result // Jump if no elements left

	// Process 4 elements if available
	CMP(rdx,$4)

	// if rdx < 4 jmp small_cleanup
	JL      small_cleanup
	VMOVDQU (AX)(SI*8), y1 // load 4 elements
	VADD(y1,y0)
	ADD(rsi,$4)            // processing 4 elements   rsi+=4
	SUB(rdx,$4)            // rdx-=4

small_cleanup:
	// Horizontal reduction
	// drobim 256bit to 128 bit and 128 bit
	VEXTRACTI128 $1, y0, x1    // 128 bit from y0 to x1 //sse
	VADD(x1,x0)                // 128bit number + 128bit number
	VPSHUFD      $0x4E, x0, x1 // shuffle  one to another
	VADD(x1,x0)
	VMOVQ_XMM_TO_GPR(x0,r8)    // save result to basic register 64bit

	// Process remaining elements (0-3)
	TESTQ rdx, rdx // rdx & rdx //if rdx == 0

	// if rdx == 0
	JZ store_result // jump if zero

	// LEA = Load Effective Address.
	LEA((rax)(rsi*8),r9)

	// (rax)(rsi*8) = rax + rsi*8 → pointer arr[rsi].

	// Efficiently handle 1-3 elements
	MOV((r9),r10) // r10= arr[rsi]
	DECQ rdx      // rdx--
	JZ   one_done
	ADD(8(r9),r10)// arr[rsi+1] //offset8
	DECQ rdx      // rdx--
	JZ   one_done
	ADD(16(r9),r10)// arr[rsi+2] //offset8

one_done:
	// added ostatok to result
	ADD(r10,r8)

store_result:
	MOV(r8,ret+24(FP)) // save result to stack and return
	VZEROUPPER         // clean ymm with operations avx,sse
	RET                // return
