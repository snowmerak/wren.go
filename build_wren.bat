@echo off
REM build_wren.bat - Build Wren static library on Windows

setlocal

set BUILD_DIR=%~dp0build
set WREN_SRC=%~dp0deps\wren\src

REM Create build directory
if not exist "%BUILD_DIR%" mkdir "%BUILD_DIR%"

echo Compiling Wren sources...

gcc -c -I "%WREN_SRC%\include" -I "%WREN_SRC%\vm" -I "%WREN_SRC%\optional" -std=c99 -O2 -o "%BUILD_DIR%\wren_vm.o" "%WREN_SRC%\vm\wren_vm.c"
gcc -c -I "%WREN_SRC%\include" -I "%WREN_SRC%\vm" -I "%WREN_SRC%\optional" -std=c99 -O2 -o "%BUILD_DIR%\wren_compiler.o" "%WREN_SRC%\vm\wren_compiler.c"
gcc -c -I "%WREN_SRC%\include" -I "%WREN_SRC%\vm" -I "%WREN_SRC%\optional" -std=c99 -O2 -o "%BUILD_DIR%\wren_core.o" "%WREN_SRC%\vm\wren_core.c"
gcc -c -I "%WREN_SRC%\include" -I "%WREN_SRC%\vm" -I "%WREN_SRC%\optional" -std=c99 -O2 -o "%BUILD_DIR%\wren_debug.o" "%WREN_SRC%\vm\wren_debug.c"
gcc -c -I "%WREN_SRC%\include" -I "%WREN_SRC%\vm" -I "%WREN_SRC%\optional" -std=c99 -O2 -o "%BUILD_DIR%\wren_primitive.o" "%WREN_SRC%\vm\wren_primitive.c"
gcc -c -I "%WREN_SRC%\include" -I "%WREN_SRC%\vm" -I "%WREN_SRC%\optional" -std=c99 -O2 -o "%BUILD_DIR%\wren_utils.o" "%WREN_SRC%\vm\wren_utils.c"
gcc -c -I "%WREN_SRC%\include" -I "%WREN_SRC%\vm" -I "%WREN_SRC%\optional" -std=c99 -O2 -o "%BUILD_DIR%\wren_value.o" "%WREN_SRC%\vm\wren_value.c"
gcc -c -I "%WREN_SRC%\include" -I "%WREN_SRC%\vm" -I "%WREN_SRC%\optional" -std=c99 -O2 -o "%BUILD_DIR%\wren_opt_meta.o" "%WREN_SRC%\optional\wren_opt_meta.c"
gcc -c -I "%WREN_SRC%\include" -I "%WREN_SRC%\vm" -I "%WREN_SRC%\optional" -std=c99 -O2 -o "%BUILD_DIR%\wren_opt_random.o" "%WREN_SRC%\optional\wren_opt_random.c"

echo Creating static library...
ar rcs "%BUILD_DIR%\libwren.a" "%BUILD_DIR%\*.o"

echo Build complete: %BUILD_DIR%\libwren.a
