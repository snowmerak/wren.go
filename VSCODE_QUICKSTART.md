# 🚀 가장 빠른 방법: F5로 바로 테스트!

## 1단계: Extension 빌드 (1분)

```bash
cd vscode-extension
npm install
npm run compile
```

## 2단계: VS Code에서 F5 (10초)

1. **VS Code에서 `vscode-extension` 폴더 열기**
2. **F5 누르기** → 새 VS Code 창이 열림 (Extension Development Host)
3. 새 창에서 **test_lsp.wren** 파일 열기
4. **완료!** LSP가 자동으로 작동합니다

## 3단계: LSP 서버 경로 설정 (선택사항)

새로 열린 VS Code 창에서:

1. `Ctrl+Shift+P` → "Preferences: Open Settings (JSON)"
2. 추가:
```json
{
  "wren.lsp.serverPath": "C:/Users/snowm/Projects/snowmerak/wren.go/bin/wren-lsp-std.exe"
}
```

> **참고**: 절대 경로를 입력하세요!

## 테스트해보기

새 창에서 `.wren` 파일을 만들고:

```wren
class Test {
  construct new() {
    _value = 0
  }
  
  // 'v' 타이핑 → 자동완성에 'var', '_value' 등 나타남
  // 'c' 타이핑 → 'class', 'construct' 등 나타남
}
```

## 확인 사항

- ✅ **문법 강조**: 키워드가 색상으로 표시됨
- ✅ **자동완성**: `Ctrl+Space`로 completion
- ✅ **오류 진단**: 문법 오류 시 빨간 밑줄
- ✅ **Hover**: 심볼 위에 마우스 → 정보 표시 (foreign methods만)

## 문제가 있다면?

1. **Output 패널 열기**: `Ctrl+Shift+U`
2. **드롭다운에서 "Wren Language Server" 선택**
3. 오류 메시지 확인

---

# 영구 설치하려면? (선택)

Extension을 계속 사용하려면:

```bash
# Windows
mklink /D "%USERPROFILE%\.vscode\extensions\wren-lsp-vscode-0.1.0" "C:\Users\snowm\Projects\snowmerak\wren.go\vscode-extension"

# VS Code 재시작
```

또는 vscode-extension 폴더를 `~/.vscode/extensions/` 에 복사하세요.
