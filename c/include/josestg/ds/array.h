#ifndef JOSESTG_DS_ARRAY_H
#define JOSESTG_DS_ARRAY_H

#include <stddef.h>
#include <stdlib.h>

// Fixed-size array.
typedef struct {
  void *head;
  size_t length;
  size_t elem_size;
} array_t;

typedef enum {
  S_OK = 0,
  S_ERR_SELF_IS_NULL = 1,
  S_ERR_RETURN_PARAMS_IS_NULL = 2,
  S_ERR_OUT_OF_MEMORY = 4,
  S_ERR_OUT_OF_RANGE = 5,
  S_ERR_ELEMENT_SIZE_MISMATCH = 6,
  S_ERR_INVALID_ARGUMENTS = 7,
} status_t;

// Initialize array with length elements of elem_size bytes.
// Memory is zeroed. Returns S_ERR_OUT_OF_MEMORY on allocation failure.
status_t array_init(array_t *self, size_t length, size_t elem_size);

// Free array memory. Array must have been initialized with array_init.
status_t array_deinit(array_t *self);

// Allocate and initialize a new array.
array_t *array_new(size_t length, size_t elem_size);

// Free array allocated with array_new. Safe to call multiple times.
status_t array_free(array_t **self);

// Write value_in to array[index]. value_size must match elem_size.
status_t array_set(const array_t *self, size_t index, const void *value_in,
                   size_t value_size);

// Read array[index] into value_out. value_size must match elem_size.
status_t array_get(const array_t *self, size_t index, void *value_out,
                   size_t value_size);

// Get array length.
status_t array_len(const array_t *self, size_t *length);

#endif // JOSESTG_DS_ARRAY_H
