cmake_minimum_required(VERSION 3.15)
project(ds C)

set(CMAKE_C_STANDARD 11)
set(CMAKE_C_STANDARD_REQUIRED ON)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

add_library(array src/array.c include/josestg/ds/array.h)
target_include_directories(array PUBLIC include)

if(CMAKE_C_COMPILER_ID MATCHES "Clang|GNU")
  target_compile_options(array PRIVATE -Wall -Wextra -Wpedantic)
endif()

option(BUILD_TESTS "Build tests" ON)
if(BUILD_TESTS)
  include(FetchContent)
  FetchContent_Declare(
    greatest
    GIT_REPOSITORY https://github.com/silentbicycle/greatest.git
    GIT_TAG v1.5.0
  )
  FetchContent_MakeAvailable(greatest)

  enable_testing()
  if(EXISTS ${CMAKE_CURRENT_SOURCE_DIR}/tests/test_array.c)
    add_executable(test_runner tests/test_array.c)
    target_link_libraries(test_runner PRIVATE array)
    target_include_directories(test_runner PRIVATE
      tests
      ${greatest_SOURCE_DIR}
    )
    if(TARGET greatest)
      target_link_libraries(test_runner PRIVATE greatest)
    endif()
    add_test(NAME array_tests COMMAND test_runner)
  endif()
endif()

find_program(CLANG_FORMAT clang-format)
if(CLANG_FORMAT)
  file(GLOB_RECURSE FORMAT_FILES
    ${CMAKE_CURRENT_SOURCE_DIR}/src/*.c
    ${CMAKE_CURRENT_SOURCE_DIR}/include/**/*.h
    ${CMAKE_CURRENT_SOURCE_DIR}/tests/*.c
    ${CMAKE_CURRENT_SOURCE_DIR}/tests/*.h
  )

  if(FORMAT_FILES)
    add_custom_target(format
      COMMAND ${CLANG_FORMAT} -i ${FORMAT_FILES}
      WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
      COMMENT "formatting sources")

    add_custom_target(lint
      COMMAND ${CLANG_FORMAT} --dry-run --Werror ${FORMAT_FILES}
      WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
      COMMENT "checking code format")
  endif()
endif()

find_program(VALGRIND valgrind)
if(VALGRIND AND TARGET test_runner)
  add_custom_target(valgrind
    COMMAND ${VALGRIND} --leak-check=full --track-origins=yes
            --show-leak-kinds=all --error-exitcode=1
            $<TARGET_FILE:test_runner>
    DEPENDS test_runner
    WORKING_DIRECTORY ${CMAKE_CURRENT_BINARY_DIR}
    COMMENT "exec tests with valgrind")
endif()
