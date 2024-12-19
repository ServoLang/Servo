//
// Created by wylan on 12/17/24.
//

#ifndef COMPILER_H
#define COMPILER_H
#include <stdbool.h>

#include "chunk.h"

Chunk *currentChunk();

bool compile(const char *source, Chunk *chunk);

#endif //COMPILER_H
