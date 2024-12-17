cmake_minimum_required(VERSION 3.30)
project(Servo)

set(CMAKE_CXX_STANDARD 20)

add_executable(Servo servo/main.c
        servo/chunk.h
        servo/common.h
        servo/chunk.c
        servo/memory.h
        servo/memory.c
        servo/debug.h
        servo/debug.c
        servo/value.c
        servo/value.h
        servo/vm.c
        servo/vm.h)
