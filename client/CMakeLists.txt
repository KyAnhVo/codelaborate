cmake_minimum_required(VERSION 3.16)
project(artscii_desktop)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# find packages (if needed just add here)

find_package(Qt6 REQUIRED COMPONENTS Widgets Core Network)

# Qt-specific setup

qt_standard_project_setup()

#include dirs (if more include dirs add here, no comma)

include_directories(
    ${CMAKE_SOURCE_DIR}/include
)

#source dirs

file(GLOB RECURSE SRC_FILES
    "${CMAKE_SOURCE_DIR}/src/*.cpp"
)
file(GLOB RECURSE HEADER_FILES
    "${CMAKE_SOURCE_DIR}/include/*.h"
)
add_executable(artscii_desktop ${SRC_FILES} ${HEADER_FILES})

# Link libraries from outside installs (most of these are dynamically linked libraries)

target_link_libraries(artscii_desktop PRIVATE
    Qt6::Widgets
    Qt6::Core
    Qt6::Network
)

