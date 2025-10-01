cmake_minimum_required(VERSION 3.16)
project(codelaborate)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# find packages (if needed just add here)

find_package(Qt6 REQUIRED COMPONENTS Widgets)

# Qt-specific setup

qt_standard_project_setup()

#include dirs (if more include dirs add here, no comma)

include_directories(
    ${CMAKE_SOURCE_DIR}/include
)

#source dirs

file(GLOB_RECURSE SRC_FILES CONFIGURE_DEPENDS
    "${CMAKE_SOURCE_DIR}/src/*.cpp"
)
file(GLOB_RECURSE HEADER_FILES CONFIGURE_DEPENDS
    "${CMAKE_SOURCE_DIR}/include/*.h"
)
message(STATUS "SRC_FILES=${SRC_FILES}")
message(STATUS "HDR_FILES=${HEADER_FILES}")
add_executable(codelaborate ${SRC_FILES} ${HEADER_FILES})

# Link libraries from outside installs (most of these are dynamically linked libraries)

target_link_libraries(codelaborate PRIVATE
    Qt6::Widgets
)

