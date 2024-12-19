//
// Created by wylan on 12/16/24.
//

#ifndef chunk_h
#define chunk_h

#include <stdint.h>

#include "common.h"
#include "value.h"

typedef enum {
    OP_CONSTANT,
    OP_NULL,
    OP_TRUE,
    OP_FALSE,

    OP_ADD,
    OP_SUBTRACT,
    OP_MULTIPLY,
    OP_DIVIDE,
    OP_MODULO,
    OP_POWER,

    OP_CONSTANT_LONG,
    OP_NEGATE,
    OP_RETURN,
} OpCode;

typedef struct {
    int offset;
    int line;
} LineStart;

typedef struct {
    int count;
    int capacity;
    uint8_t* code;
    ValueArray constants;
    int lineCount;
    int lineCapacity;
    LineStart *lines;
} Chunk;

void initChunk(Chunk* chunk);
void freeChunk(Chunk* chunk);
void writeChunk(Chunk* chunk, uint32_t byte, int line);
void writeConstant(Chunk *chunk, Value value, int line);
int addConstant(Chunk *chunk, Value value);

int getLine(Chunk *chunk, int instruction);

#endif
