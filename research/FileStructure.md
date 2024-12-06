# Servo file structure

### This document describes the structure of a .svo or .servo file.
A servo file contains a single main class, 

## Scope
The scope is inherently set based on folder structure. Depending on the containing folder, the scope will only be within
the folder and its subfolders. This is useful for organizing files and keeping track of where you are in the file structure.

```servo
scope some.nested.folders;
```