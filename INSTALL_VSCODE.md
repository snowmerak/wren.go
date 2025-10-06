# VS Code에 Wren LSP Extension 설치하기

## ✅ VSIX 파일로 설치 (추천 - 가장 쉬움!)

VSIX 파일이 이미 생성되었습니다: `vscode-extension/wren-lsp-vscode-0.1.0.vsix`

### 방법 A: VS Code UI에서 설치

1. **VS Code 열기**
2. **Extensions 패널** (`Ctrl+Shift+X`)
3. **... 메뉴** (우측 상단) 클릭
4. **"Install from VSIX..."** 선택
5. `C:\Users\snowm\Projects\snowmerak\wren.go\vscode-extension\wren-lsp-vscode-0.1.0.vsix` 선택
6. **완료!** VS Code 재시작

### 방법 B: 명령어로 설치

```powershell
code --install-extension "C:\Users\snowm\Projects\snowmerak\wren.go\vscode-extension\wren-lsp-vscode-0.1.0.vsix"
```

## 설치 후 설정

1. **VS Code 열기**
2. **Settings** (`Ctrl+,`)
3. **"wren lsp"** 검색
4. **"Server Path"**에 입력:
   ```
   C:/Users/snowm/Projects/snowmerak/wren.go/bin/wren-lsp-std.exe
   ```

또는 `settings.json`에 직접 추가:

```json
{
  "wren.lsp.serverPath": "C:/Users/snowm/Projects/snowmerak/wren.go/bin/wren-lsp-std.exe"
}
```

## 테스트

1. 새 파일 만들기: `test.wren`
2. 코드 작성:
   ```wren
   class Test {
     construct new() {
       _value = 0
     }
     
     add(x) {
       // Ctrl+Space로 자동완성 테스트
     }
   }
   ```

3. 확인:
   - ✅ 문법 강조
   - ✅ 자동완성 (`Ctrl+Space`)
   - ✅ 오류 진단 (문법 오류 시)

## 삭제하려면?

1. Extensions 패널 (`Ctrl+Shift+X`)
2. "Wren Language Support" 찾기
3. Uninstall 클릭

---

## 다른 설치 방법들 (선택사항)

### 방법 1: Symbolic Link (개발자용)

```powershell
# PowerShell 관리자 권한
New-Item -ItemType SymbolicLink -Path "$env:USERPROFILE\.vscode\extensions\wren-lsp-vscode-0.1.0" -Target "C:\Users\snowm\Projects\snowmerak\wren.go\vscode-extension"
```

장점: 수정 사항이 바로 반영됨

### 방법 2: 복사 (간단)

```powershell
Copy-Item -Recurse vscode-extension "$env:USERPROFILE\.vscode\extensions\wren-lsp-vscode-0.1.0"
```

단점: 수정할 때마다 다시 복사해야 함

---

## 문제 해결

### Extension이 활성화되지 않음
1. VS Code 재시작
2. Output 패널 (`Ctrl+Shift+U`) → "Wren Language Server" 확인

### LSP가 시작하지 않음
1. `wren.lsp.serverPath` 설정 확인
2. 절대 경로로 입력했는지 확인
3. `wren-lsp-std.exe`가 실행 가능한지 확인

### 자동완성이 안 됨
1. 파일 확장자가 `.wren`인지 확인
2. `Ctrl+Space` 눌러보기
3. Output 패널에서 오류 확인
