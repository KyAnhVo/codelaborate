cmake_minimum_required(VERSION 3.16)
project(codelaborate)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# find packages
find_package(Qt6 REQUIRED COMPONENTS Widgets Network)

# Qt-specific setup

qt_standard_project_setup()

#include dirs
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

add_executable(codelaborate ${SRC_FILES} ${HEADER_FILES})

# Link libraries from outside installs (most of these are dynamically linked libraries)
target_link_libraries(codelaborate PRIVATE
    Qt6::Widgets
    Qt6::Network
)

