#!/usr/bin/env python3
"""
build.py - Build wren.go projects with static linking

This script:
1. Builds the Wren C static library
2. Builds Go CLI tools with static linking
3. Creates distributable binaries
"""

import os
import sys
import platform
import subprocess
import argparse
from pathlib import Path

# Color output
class Colors:
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKCYAN = '\033[96m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'

def print_step(msg):
    print(f"{Colors.OKBLUE}{Colors.BOLD}==>{Colors.ENDC} {msg}")

def print_success(msg):
    print(f"{Colors.OKGREEN}[OK]{Colors.ENDC} {msg}")

def print_error(msg):
    print(f"{Colors.FAIL}[ERROR]{Colors.ENDC} {msg}")

def print_warning(msg):
    print(f"{Colors.WARNING}[WARN]{Colors.ENDC} {msg}")

def find_gcc():
    """Find GCC in common locations"""
    common_paths = [
        "C:/TDM-GCC-64/bin",
        "C:/MinGW/bin",
        "C:/mingw64/bin",
        os.path.expanduser("~/scoop/apps/tdm-gcc/current/bin"),
        os.path.expanduser("~/scoop/apps/gcc/current/bin"),
    ]
    
    for path in common_paths:
        gcc_path = os.path.join(path, "gcc.exe")
        if os.path.exists(gcc_path):
            return path
    
    return None

def run_command(cmd, cwd=None, env=None):
    """Run a command and return success status"""
    try:
        # Add GCC to PATH on Windows if not already there
        if platform.system().lower() == 'windows':
            if env is None:
                env = os.environ.copy()
            gcc_bin = find_gcc()
            if gcc_bin:
                env['PATH'] = f"{gcc_bin};{env.get('PATH', os.environ.get('PATH', ''))}"
        
        result = subprocess.run(
            cmd,
            cwd=cwd,
            env=env,
            check=True,
            capture_output=True,
            text=True
        )
        return True, result.stdout
    except subprocess.CalledProcessError as e:
        return False, e.stderr

def get_platform_info():
    """Get platform-specific information"""
    system = platform.system().lower()
    machine = platform.machine().lower()
    
    if system == 'windows':
        return 'windows', 'amd64', '.exe'
    elif system == 'darwin':
        return 'darwin', 'arm64' if 'arm' in machine else 'amd64', ''
    elif system == 'linux':
        return 'linux', 'amd64', ''
    else:
        return system, machine, ''

def build_wren_library():
    """Build the Wren C static library"""
    print_step("Building Wren C library...")
    
    system, _, _ = get_platform_info()
    
    if system == 'windows':
        script = 'build_wren.bat'
    else:
        script = './build_wren.sh'
    
    if not os.path.exists(script):
        print_error(f"Build script not found: {script}")
        return False
    
    # Make script executable on Unix
    if system != 'windows':
        os.chmod(script, 0o755)
    
    success, output = run_command([script] if system == 'windows' else [script])
    
    if success:
        print_success("Wren library built successfully")
        return True
    else:
        print_error(f"Failed to build Wren library:\n{output}")
        return False

def generate_builtin_code():
    """Generate builtin FFI bindings and LSP symbols"""
    print_step("Generating builtin code...")
    
    # Get current working directory for absolute paths
    root_dir = os.getcwd()
    
    # Step 1: Generate builtin_wren.go using wrengen
    print_step("Running wrengen to generate FFI bindings...")
    cmd = ['go', 'run', './cmd/wrengen', '-dir', 'builtin']
    success, output = run_command(cmd)
    
    if not success:
        print_error(f"Failed to generate FFI bindings:\n{output}")
        return False
    
    print_success("FFI bindings generated")
    
    # Step 2: Build wrenlsp-gen tool if it doesn't exist
    system, _, ext = get_platform_info()
    wrenlsp_gen_path = os.path.join(root_dir, f"wrenlsp-gen{ext}")
    
    if not os.path.exists(wrenlsp_gen_path):
        print_step("Building wrenlsp-gen tool...")
        cmd = ['go', 'build', '-o', wrenlsp_gen_path, './cmd/wrenlsp-gen/main.go']
        success, output = run_command(cmd)
        
        if not success:
            print_error(f"Failed to build wrenlsp-gen:\n{output}")
            return False
        
        print_success("wrenlsp-gen tool built")
    
    # Step 3: Generate LSP builtin symbols
    print_step("Generating LSP builtin symbols...")
    builtin_wren_path = os.path.join(root_dir, 'cmd', 'gwen', 'builtin_wren.go')
    cmd = [wrenlsp_gen_path, builtin_wren_path]
    success, output = run_command(cmd)
    
    if not success:
        print_error(f"Failed to generate LSP symbols:\n{output}")
        return False
    
    print_success("LSP builtin symbols generated")
    print_success("All builtin code generated successfully")
    return True

def copy_mingw_dlls(output_dir):
    """Copy required MinGW DLLs to output directory on Windows"""
    system = platform.system().lower()
    if system != 'windows':
        return True
    
    print_step("Copying MinGW DLLs...")
    
    gcc_bin = find_gcc()
    if not gcc_bin:
        print_warning("Could not find MinGW installation")
        return False
    
    # Required DLLs for 64-bit
    required_dlls = [
        'libgcc_s_seh_64-1.dll',
        'libwinpthread_64-1.dll',
        'libstdc++_64-6.dll',
    ]
    
    # Try alternative names if _64 versions not found
    alternative_dlls = [
        'libgcc_s_seh-1.dll',
        'libwinpthread-1.dll', 
        'libstdc++-6.dll',
    ]
    
    copied = 0
    for i, dll_name in enumerate(required_dlls):
        src = os.path.join(gcc_bin, dll_name)
        
        # Try alternative name if primary not found
        if not os.path.exists(src):
            alt_name = alternative_dlls[i]
            src = os.path.join(gcc_bin, alt_name)
            if not os.path.exists(src):
                print_warning(f"DLL not found: {dll_name} or {alt_name}")
                continue
        
        dst = os.path.join(output_dir, os.path.basename(src))
        
        try:
            import shutil
            shutil.copy2(src, dst)
            print_success(f"Copied {os.path.basename(src)}")
            copied += 1
        except Exception as e:
            print_warning(f"Failed to copy {dll_name}: {e}")
    
    if copied > 0:
        print_success(f"Copied {copied} DLL(s)")
        return True
    else:
        print_error("No DLLs were copied")
        return False

def build_go_binary(target, output_name=None, static=True, copy_dlls=True):
    """Build a Go binary with optional static linking"""
    print_step(f"Building {target}...")
    
    system, arch, ext = get_platform_info()
    
    if output_name is None:
        output_name = os.path.basename(target)
    
    output_path = f"bin/{output_name}{ext}"
    
    # Create bin directory
    os.makedirs("bin", exist_ok=True)
    
    # Build command
    cmd = ['go', 'build', '-o', output_path]
    
    # Add build flags
    if static:
        print_warning("Building with static linking...")
        
        # Platform-specific flags
        if system == 'linux':
            cmd.extend([
                '-ldflags',
                '-extldflags "-static"',
                '-tags', 'netgo'
            ])
        elif system == 'windows':
            # Windows: build with all static libraries
            cmd.extend([
                '-ldflags',
                '-s -w'  # Just strip debug info, use dynamic linking
            ])
        elif system == 'darwin':
            # macOS doesn't support full static linking
            print_warning("Note: Full static linking not available on macOS")
            cmd.extend([
                '-ldflags',
                '-s -w'  # Strip debug info to reduce size
            ])
    
    cmd.append(f'./{target}')
    
    # Set CGO environment
    env = os.environ.copy()
    env['CGO_ENABLED'] = '1'
    
    # Set CGO paths for wren library
    root_dir = os.getcwd()
    build_dir = os.path.join(root_dir, 'build')
    
    if system == 'windows':
        env['CGO_LDFLAGS'] = f'-L{build_dir} -lwren -lm'
    else:
        env['CGO_LDFLAGS'] = f'-L{build_dir} -lwren -lm'
    
    success, output = run_command(cmd, env=env)
    
    if success:
        size = os.path.getsize(output_path)
        size_mb = size / (1024 * 1024)
        print_success(f"Built {output_path} ({size_mb:.2f} MB)")
        
        # Copy DLLs on Windows if requested
        if copy_dlls and system == 'windows':
            output_dir = os.path.dirname(output_path)
            copy_mingw_dlls(output_dir)
        
        return True
    else:
        print_error(f"Failed to build {target}:\n{output}")
        return False

def clean():
    """Clean build artifacts"""
    print_step("Cleaning build artifacts...")
    
    dirs_to_clean = ['build', 'bin']
    
    for dir_name in dirs_to_clean:
        if os.path.exists(dir_name):
            import shutil
            shutil.rmtree(dir_name)
            print_success(f"Removed {dir_name}/")
    
    print_success("Clean complete")

def run_tests(package=None):
    """Run Go tests"""
    print_step("Running tests...")
    
    cmd = ['go', 'test', '-v']
    
    if package:
        cmd.append(f'./{package}')
    else:
        cmd.append('./...')
    
    success, output = run_command(cmd)
    
    if success:
        print_success("All tests passed")
        print(output)
        return True
    else:
        print_error("Tests failed")
        print(output)
        return False

def main():
    parser = argparse.ArgumentParser(description='Build wren.go projects')
    parser.add_argument('command', nargs='?', default='all',
                        choices=['all', 'wren', 'cli', 'lsp', 'generate', 'test', 'clean'],
                        help='Build command (all: full build, wren: C library only, cli: Go CLI only, lsp: Go LSP only, generate: code generation only, test: run tests, clean: clean artifacts)')
    parser.add_argument('--no-static', action='store_true',
                       help='Disable static linking')
    parser.add_argument('--target', default='cmd/gwen',
                       help='Target to build (default: cmd/gwen)')
    parser.add_argument('--output', help='Output binary name')
    parser.add_argument('--no-copy-dlls', action='store_true',
                       help='Do not copy MinGW DLLs (Windows only)')
    
    args = parser.parse_args()
    
    print(f"{Colors.HEADER}{Colors.BOLD}")
    print("=" * 50)
    print("   Wren.go Build System")
    print("=" * 50)
    print(f"{Colors.ENDC}")
    
    # Show platform info
    system, arch, ext = get_platform_info()
    print(f"Platform: {system}/{arch}")
    print()
    
    if args.command == 'clean':
        clean()
        return 0
    
    if args.command == 'test':
        success = run_tests()
        return 0 if success else 1
    
    if args.command == 'generate':
        success = generate_builtin_code()
        return 0 if success else 1
    
    # Build Wren library
    if args.command in ['all', 'wren']:
        if not build_wren_library():
            return 1
    
    # Generate builtin code
    if args.command in ['all', 'wren', 'cli', 'lsp']:
        if not generate_builtin_code():
            return 1
    
    # Build Go binaries
    if args.command in ['all', 'cli']:
        # Build CLI
        if not build_go_binary(args.target, args.output, static=not args.no_static, copy_dlls=not args.no_copy_dlls):
            return 1
    
    if args.command in ['all', 'lsp']:
        # Build LSP server
        lsp_target = 'cmd/gwen-lsp'
        lsp_output = f'gwen-lsp'
        if not build_go_binary(lsp_target, lsp_output, static=not args.no_static, copy_dlls=False):
            return 1
    
    print()
    print_success("Build complete!")
    print()
    print("Binaries are in the bin/ directory:")
    
    if os.path.exists('bin'):
        for item in os.listdir('bin'):
            path = os.path.join('bin', item)
            if os.path.isfile(path):
                size = os.path.getsize(path)
                size_mb = size / (1024 * 1024)
                print(f"  - {item} ({size_mb:.2f} MB)")
    
    return 0

if __name__ == '__main__':
    try:
        sys.exit(main())
    except KeyboardInterrupt:
        print()
        print_warning("Build interrupted")
        sys.exit(1)
    except Exception as e:
        print_error(f"Unexpected error: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)
