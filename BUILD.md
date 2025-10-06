# Build System for wren.go

Python-based build system for creating statically-linked binaries.

## Status

✅ **Windows**: Dynamic linking with MinGW DLLs (automatically copied)  
✅ **Linux**: Full static linking support  
✅ **macOS**: Dynamic linking (system limitation)

## Quick Start

### Build Everything (Recommended)

```bash
python build.py all
```

This will:
1. Build the Wren C static library
2. Build the wren-std CLI with static linking
3. Create binaries in `bin/` directory

### Build Commands

```bash
# Build everything (Wren library + CLI)
python build.py all

# Build only Wren C library
python build.py wren

# Build only Go CLI
python build.py cli

# Build only LSP server
python build.py lsp

# Run tests
python build.py test

# Clean build artifacts
python build.py clean
```

### Advanced Options

```bash
# Build without static linking
python build.py all --no-static

# Build a specific target
python build.py cli --target cmd/wren-std --output my-wren

# Custom output name
python build.py cli --output my-cli
```

## Platform Support

- ✅ **Linux**: Full static linking support
- ✅ **Windows**: Dynamic linking + MinGW DLLs (auto-copied to bin/)
- ✅ **macOS**: Dynamic linking (system limitation)

## Requirements

### All Platforms
- Python 3.6+
- Go 1.25.1+
- CGO enabled

### Linux
- GCC
- Make

### Windows
- MinGW-w64 (GCC for Windows)
- Git Bash (optional, for shell scripts)

### macOS
- Xcode Command Line Tools
- Make

## Output

Binaries are created in the `bin/` directory:

**Linux/macOS:**
```
bin/
├── wren-std       # Standard CLI (single file)
└── wren-lsp-std   # LSP server (single file)
```

**Windows:**
```
bin/
├── wren-std.exe              # Standard CLI
├── wren-lsp-std.exe          # LSP server
├── libgcc_s_seh_64-1.dll     # GCC runtime (shared)
├── libwinpthread_64-1.dll    # pthread support (shared)
└── libstdc++_64-6.dll        # C++ stdlib (shared)
```

**Note**: On Windows, distribute the entire `bin/` directory with all DLLs.

## Linking Strategy

### Linux (Static)
- Fully static binary
- No external dependencies
- Single file distribution

### Windows (Dynamic + DLLs)
- Dynamic linking with MinGW runtime
- DLLs automatically copied to `bin/`
- Distribute entire `bin/` folder

### macOS (Dynamic)
- Dynamic linking (system limitation)
- Uses system libraries
- Single binary file

### Size Comparison

| Platform | Binary | DLLs | Total |
|----------|--------|------|-------|
| Linux (static) | 8-12 MB | - | 8-12 MB |
| Windows (dynamic) | ~2 MB | ~2 MB | ~4 MB |
| macOS (dynamic) | 3-5 MB | - | 3-5 MB ||## Linking Strategy

### Linux (Static)
- Fully static binary
- No external dependencies
- Single file distribution

### Windows (Dynamic + DLLs)
- Dynamic linking with MinGW runtime
- DLLs automatically copied to `bin/`
- Distribute entire `bin/` folder

### macOS (Dynamic)
- Dynamic linking (system limitation)
- Uses system libraries
- Single binary file

### Size Comparison

| Platform | Binary | DLLs | Total |
|----------|--------|------|-------|
| Linux (static) | 8-12 MB | - | 8-12 MB |
| Windows (dynamic) | ~2 MB | ~2 MB | ~4 MB |
| macOS (dynamic) | 3-5 MB | - | 3-5 MB |

## Troubleshooting

### "CGO_ENABLED not set"

```bash
export CGO_ENABLED=1  # Linux/macOS
set CGO_ENABLED=1     # Windows
```

### "gcc not found"

Install GCC:
- **Linux**: `sudo apt install build-essential`
- **Windows**: Install MinGW-w64
- **macOS**: `xcode-select --install`

### Static linking fails on macOS

This is expected. macOS doesn't support full static linking. The build script will create a smaller dynamic binary instead.

### Build script permission denied (Linux/macOS)

```bash
chmod +x build.py
chmod +x build_wren.sh
```

## Distribution

### Windows Distribution

```bash
# Build with DLLs
python build.py all

# Create distributable package
mkdir wren-std-windows
xcopy /E bin wren-std-windows\
# Now you have:
# wren-std-windows/
# ├── wren-std.exe
# ├── libgcc_s_seh_64-1.dll
# ├── libwinpthread_64-1.dll
# └── libstdc++_64-6.dll

# Zip it
powershell Compress-Archive -Path wren-std-windows -DestinationPath wren-std-windows.zip
```

### Linux Distribution

```bash
# Build static binary
python build.py all

# Single file distribution!
cp bin/wren-std wren-std-linux
# or
tar czf wren-std-linux.tar.gz -C bin wren-std
```

### Skipping DLL Copy

```bash
# Build without copying DLLs (if you want to handle them manually)
python build.py all --no-copy-dlls
```

## Examples

### Build for Distribution

```bash
# Clean previous builds
python build.py clean

# Build with static linking
python build.py all

# Test the binary
./bin/wren-std version
./bin/wren-std run script.wren
```

### Development Build

```bash
# Fast build without static linking
python build.py all --no-static

# Run tests
python build.py test
```

### Cross-Platform Build

```bash
# On Linux, build for Linux
python build.py all

# On Windows, build for Windows
python build.py all

# On macOS, build for macOS
python build.py all
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Build

on: [push, pull_request]

jobs:
  build:
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
    
    runs-on: ${{ matrix.os }}
    
    steps:
    - uses: actions/checkout@v3
      with:
        submodules: recursive
    
    - uses: actions/setup-go@v4
      with:
        go-version: '1.25'
    
    - name: Build
      run: python build.py all
    
    - name: Test
      run: python build.py test
    
    - name: Upload artifacts
      uses: actions/upload-artifact@v3
      with:
        name: wren-std-${{ matrix.os }}
        path: bin/
```

## License

MIT License - see main repository for details.
