//
// Created by wylan on 12/17/24.
//

#ifndef VM_H
#define VM_H

#include "chunk.h"
#include "value.h"

typedef struct {
    Chunk *chunk;
    uint8_t *ip;
    Value *stack;
    int stackCount;
    int stackCapacity;
} VM;

typedef enum {
    INTERPRET_OK,
    INTERPRET_COMPILE_ERROR,
    INTERPRET_RUNTIME_ERROR
} InterpretResult;

void initVM();
void freeVM();
InterpretResult interpret(const char *source);
void push(Value value);
Value pop();

#endif //VM_H
