#include <greatest.h>
#include <josestg/ds/array.h>

const size_t int_elem_size = sizeof(int);

TEST test_array_init_and_free(void) {
  array_t arr;
  ASSERT_EQ(array_init(&arr, 10, int_elem_size), S_OK);
  ASSERT(arr.head != NULL);
  ASSERT_EQ(arr.length, 10);

  ASSERT_EQ(array_deinit(&arr), S_OK);
  ASSERT_EQ(arr.head, NULL);
  ASSERT_EQ(arr.length, 0);
  PASS();
}

TEST test_array_zero_values(void) {
  array_t arr;
  ASSERT_EQ(array_init(&arr, 10, int_elem_size), S_OK);
  for (int i = 0; i < 10; i++) {
    int val = -1;
    ASSERT_EQ(array_get(&arr, i, &val, int_elem_size), S_OK);
    ASSERT_EQ(val, 0);
  }

  ASSERT_EQ(array_deinit(&arr), S_OK);
  ASSERT_EQ(arr.head, NULL);
  ASSERT_EQ(arr.length, 0);
  PASS();
}

TEST test_array_init_invalid_inputs(void) {
  ASSERT_EQ(array_init(NULL, 10, int_elem_size), S_ERR_SELF_IS_NULL);

  array_t arr;
  ASSERT_EQ(array_init(&arr, 0, int_elem_size), S_ERR_OUT_OF_RANGE);
  PASS();
}

TEST test_array_new_and_free(void) {
  array_t *arr = array_new(5, int_elem_size);
  ASSERT(arr != NULL);
  ASSERT(arr->head != NULL);
  ASSERT_EQ(arr->length, 5);
  ASSERT_EQ(arr->elem_size, int_elem_size);

  int val = 99;
  ASSERT_EQ(array_set(arr, 0, &val, int_elem_size), S_OK);

  int result = 0;
  ASSERT_EQ(array_get(arr, 0, &result, int_elem_size), S_OK);
  ASSERT_EQ(result, 99);

  ASSERT_EQ(array_free(&arr), S_OK);
  ASSERT_EQ(arr, NULL);
  PASS();
}

TEST test_array_new_zero_length(void) {
  array_t *arr = array_new(0, int_elem_size);
  ASSERT_EQ(arr, NULL);
  PASS();
}

TEST test_array_free_null_pointer(void) {
  ASSERT_EQ(array_free(NULL), S_ERR_SELF_IS_NULL);
  PASS();
}

TEST test_array_free_double_free(void) {
  array_t *arr = array_new(3, int_elem_size);
  ASSERT(arr != NULL);

  ASSERT_EQ(array_free(&arr), S_OK);
  ASSERT_EQ(arr, NULL);

  ASSERT_EQ(array_free(&arr), S_OK);
  PASS();
}

TEST test_array_new_large_allocation(void) {
  array_t *arr = array_new(1000, sizeof(int));
  ASSERT(arr != NULL);
  ASSERT_EQ(arr->length, 1000);

  for (size_t i = 0; i < 1000; i++) {
    int val = (int)i;
    ASSERT_EQ(array_set(arr, i, &val, sizeof(int)), S_OK);
  }

  for (size_t i = 0; i < 1000; i++) {
    int val = 0;
    ASSERT_EQ(array_get(arr, i, &val, sizeof(int)), S_OK);
    ASSERT_EQ(val, (int)i);
  }

  ASSERT_EQ(array_free(&arr), S_OK);
  PASS();
}

TEST test_array_set_and_get(void) {
  array_t arr;
  array_init(&arr, 5, int_elem_size);

  int val = 42;
  ASSERT_EQ(array_set(&arr, 2, &val, int_elem_size), S_OK);

  int value = 0;
  ASSERT_EQ(array_get(&arr, 2, &value, int_elem_size), S_OK);
  ASSERT_EQ(value, val);

  array_deinit(&arr);
  PASS();
}

TEST test_array_bounds_check(void) {
  array_t arr;
  array_init(&arr, 3, int_elem_size);

  ASSERT_EQ(array_set(&arr, 3, (void *)1, int_elem_size), S_ERR_OUT_OF_RANGE);
  ASSERT_EQ(array_get(&arr, 3, NULL, int_elem_size),
            S_ERR_RETURN_PARAMS_IS_NULL);
  ASSERT_EQ(array_get(NULL, 1, NULL, int_elem_size), S_ERR_SELF_IS_NULL);

  array_deinit(&arr);
  PASS();
}

TEST test_array_len(void) {
  array_t arr;
  array_init(&arr, 7, int_elem_size);

  size_t len = 0;
  ASSERT_EQ(array_len(&arr, &len), S_OK);
  ASSERT_EQ(len, 7);

  ASSERT_EQ(array_len(NULL, &len), S_ERR_SELF_IS_NULL);
  ASSERT_EQ(array_len(&arr, NULL), S_ERR_RETURN_PARAMS_IS_NULL);

  array_deinit(&arr);
  PASS();
}

TEST test_array_int_simulation(void) {
  array_t arr;
  const size_t length = 100;

  printf("creating array of length %zu...\n", length);
  ASSERT_EQ(array_init(&arr, length, int_elem_size), S_OK);

  printf("filling array with i * 2...\n");
  for (size_t i = 0; i < length; i++) {
    int val = (int)i * 2;
    ASSERT_EQ(array_set(&arr, i, &val, int_elem_size), S_OK);
  }

  printf("verifying values...\n");
  for (size_t i = 0; i < length; i++) {
    int val = 0;
    ASSERT_EQ(array_get(&arr, i, &val, int_elem_size), S_OK);
    ASSERT_EQ(val, (i * 2));
  }

  printf("adding 1 to values at index 10 to 19...\n");
  for (size_t i = 10; i < 20; i++) {
    int val = 0;
    ASSERT_EQ(array_get(&arr, i, &val, int_elem_size), S_OK);
    val++;
    ASSERT_EQ(array_set(&arr, i, &val, int_elem_size), S_OK);
  }

  printf("confirming updated values at index 10 to 19...\n");
  for (size_t i = 10; i < 20; i++) {
    int val = 0;
    ASSERT_EQ(array_get(&arr, i, &val, int_elem_size), S_OK);
    ASSERT_EQ(val, (int)(i * 2 + 1));
  }

  printf("cleaning up...\n");
  array_deinit(&arr);

  printf("array simulation test passed.\n");
  PASS();
}

typedef struct {
  int x, y, z;
} Vec3;

TEST test_array_struct_simulation(void) {
  const Vec3 a = {1, 2, 3};
  const Vec3 b = {4, 5, 6};
  const Vec3 c = {7, 8, 9};

  array_t arr;
  const size_t length = 3;
  array_init(&arr, length, sizeof(Vec3));

  ASSERT_EQ(array_set(&arr, 0, &a, sizeof(Vec3)), S_OK);
  ASSERT_EQ(array_set(&arr, 1, &b, sizeof(Vec3)), S_OK);
  ASSERT_EQ(array_set(&arr, 2, &c, sizeof(Vec3)), S_OK);

  Vec3 a2, b2, c2;
  ASSERT_EQ(array_get(&arr, 0, &a2, sizeof(Vec3)), S_OK);
  ASSERT_EQ(array_get(&arr, 1, &b2, sizeof(Vec3)), S_OK);
  ASSERT_EQ(array_get(&arr, 2, &c2, sizeof(Vec3)), S_OK);

  ASSERT_EQ(a.x, a2.x);
  ASSERT_EQ(a.y, a2.y);
  ASSERT_EQ(a.z, a2.z);

  ASSERT_EQ(b.x, b2.x);
  ASSERT_EQ(b.y, b2.y);
  ASSERT_EQ(b.z, b2.z);

  ASSERT_EQ(c.x, c2.x);
  ASSERT_EQ(c.y, c2.y);
  ASSERT_EQ(c.z, c2.z);

  array_deinit(&arr);
  PASS();
}

SUITE(array_tests) {
  RUN_TEST(test_array_init_and_free);
  RUN_TEST(test_array_init_invalid_inputs);
  RUN_TEST(test_array_new_and_free);
  RUN_TEST(test_array_new_zero_length);
  RUN_TEST(test_array_free_null_pointer);
  RUN_TEST(test_array_free_double_free);
  RUN_TEST(test_array_new_large_allocation);
  RUN_TEST(test_array_set_and_get);
  RUN_TEST(test_array_bounds_check);
  RUN_TEST(test_array_len);
  RUN_TEST(test_array_zero_values);
}

SUITE(array_simulation) {
  RUN_TEST(test_array_int_simulation);
  RUN_TEST(test_array_struct_simulation);
}

GREATEST_MAIN_DEFS();

int main(const int argc, char **argv) {
  GREATEST_MAIN_BEGIN();
  RUN_SUITE(array_tests);
  RUN_SUITE(array_simulation);
  GREATEST_MAIN_END();
}
