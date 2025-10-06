# Wren LSP VS Code Extension

간단한 VS Code extension으로 wren-lsp-std를 사용합니다.

## 빠른 시작

### 1. Extension 빌드

```bash
cd vscode-extension
npm install
npm run compile
```

### 2. Extension 설치 (로컬)

**방법 1: F5로 테스트 (가장 빠름)**
1. VS Code에서 `vscode-extension` 폴더 열기
2. F5 누르기 → 새 VS Code 창 열림
3. `.wren` 파일 열기 → LSP 자동 실행

**방법 2: 수동 설치**
```bash
# vscode-extension 폴더를 VS Code extensions 폴더에 복사
# Windows:
cp -r vscode-extension ~/.vscode/extensions/wren-lsp-vscode-0.1.0

# 또는 symbolic link
mklink /D "%USERPROFILE%\.vscode\extensions\wren-lsp-vscode-0.1.0" "C:\path\to\wren.go\vscode-extension"
```

### 3. LSP 서버 경로 설정

VS Code settings.json에 추가:

```json
{
  "wren.lsp.serverPath": "C:/Users/snowm/Projects/snowmerak/wren.go/bin/wren-lsp-std.exe"
}
```

또는 `wren-lsp-std.exe`를 PATH에 추가하면 자동으로 찾습니다.

### 4. 사용

1. `.wren` 파일 열기
2. LSP 기능 자동 활성화:
   - ✅ 문법 강조 (syntax highlighting)
   - ✅ 자동완성 (completion)
   - ✅ 오류 진단 (diagnostics)
   - ✅ Hover 정보

## 설정

```json
{
  // LSP 서버 경로 (기본값: PATH에서 wren-lsp-std 찾기)
  "wren.lsp.serverPath": "C:/path/to/wren-lsp-std.exe"
}
```

## 문제 해결

### LSP가 시작하지 않음
1. Output 패널에서 "Wren Language Server" 확인
2. `wren.lsp.serverPath` 설정 확인
3. `wren-lsp-std.exe` 실행 권한 확인

### Completion이 작동하지 않음
1. 파일 확장자가 `.wren`인지 확인
2. VS Code 재시작

## 개발

```bash
# TypeScript 컴파일 (watch mode)
npm run watch

# F5로 디버깅
```
