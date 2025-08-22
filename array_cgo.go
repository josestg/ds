package ds

// #include "c/src/array.c"
import "C"

// CGO only compile C files only if they are in the same directory with the Go files who uses it.
// So this trick, will solve it.
