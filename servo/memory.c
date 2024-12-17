//
// Created by wylan on 12/17/24.
//

#include <stdlib.h>
#include "memory.h"

void *reallocate(void *pointer, size_t oldSize, size_t newSize) {
    if (newSize == 0) {
        free(pointer);
        return NULL;
    }

    void *result = realloc(pointer, newSize);
    if (result == NULL) exit(EXIT_FAILURE); // EXIT_FAILURE = 1
    return result;
}

