# Async Future를 반환하는 Go 메서드 작성 가이드

## 개요

Go에서 비동기 작업을 수행하고 Future ID를 Wren으로 반환하는 메서드를 작성하는 방법을 설명합니다.

## 기본 패턴

### 1. 단순 비동기 작업

```go
// Sleep creates a future that completes after specified seconds
//wren:bind name=sleep(_) static
func (a *Async) Sleep(vm *WrenVM) error {
    seconds := vm.GetSlotDouble(1)  // Wren에서 전달받은 파라미터
    
    // 비동기 작업 생성
    future := GetAsyncManager().Submit(func(ctx context.Context) (interface{}, error) {
        time.Sleep(time.Duration(seconds) * time.Second)
        return fmt.Sprintf("Slept for %.1f seconds", seconds), nil
    })
    
    // Future ID를 Wren으로 반환
    vm.SetSlotDouble(0, float64(future.ID()))
    return nil
}
```

**Wren에서 사용:**
```wren
import "main" for Async

var futureId = Async.sleep(2.5)
var result = Async.await(futureId)
System.print(result)  // "Slept for 2.5 seconds"
Async.cleanup(futureId)
```

### 2. 복잡한 계산 작업

```go
// Calculate performs heavy computation asynchronously
//wren:bind name=calculate(_,_) static
func (a *Async) Calculate(vm *WrenVM) error {
    start := int(vm.GetSlotDouble(1))
    end := int(vm.GetSlotDouble(2))
    
    future := GetAsyncManager().Submit(func(ctx context.Context) (interface{}, error) {
        // 무거운 계산 시뮬레이션
        sum := 0
        for i := start; i <= end; i++ {
            // context 취소 확인
            select {
            case <-ctx.Done():
                return nil, ctx.Err()
            default:
            }
            
            sum += i * i
            time.Sleep(1 * time.Millisecond) // 계산 시간 시뮬레이션
        }
        return sum, nil
    })
    
    vm.SetSlotDouble(0, float64(future.ID()))
    return nil
}
```

**Wren에서 사용:**
```wren
var futureId = Async.calculate(1, 1000)
System.print("Calculation started...")

// 논블로킹으로 체크
while (!Async.isReady(futureId)) {
    System.print("Still calculating...")
    // 다른 작업 수행 가능
}

var result = Async.get(futureId)
System.print("Sum of squares: %(result)")
Async.cleanup(futureId)
```

### 3. HTTP 요청 예제

```go
import (
    "context"
    "io"
    "net/http"
    "time"
)

// HttpGet performs an HTTP GET request asynchronously
//wren:bind name=httpGet(_) static
func (a *Async) HttpGet(vm *WrenVM) error {
    url := vm.GetSlotString(1)  // URL 파라미터
    
    future := GetAsyncManager().Submit(func(ctx context.Context) (interface{}, error) {
        // 타임아웃 설정
        client := &http.Client{
            Timeout: 10 * time.Second,
        }
        
        req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
        if err != nil {
            return nil, err
        }
        
        resp, err := client.Do(req)
        if err != nil {
            return nil, err
        }
        defer resp.Body.Close()
        
        body, err := io.ReadAll(resp.Body)
        if err != nil {
            return nil, err
        }
        
        return map[string]interface{}{
            "status": resp.StatusCode,
            "body":   string(body),
            "length": len(body),
        }, nil
    })
    
    vm.SetSlotDouble(0, float64(future.ID()))
    return nil
}
```

**Wren에서 사용:**
```wren
var futureId = Async.httpGet("https://api.github.com")
System.print("Request sent...")

var result = Async.await(futureId)
// result는 Map 형태로 반환됨
System.print("Status: %(result["status"])")
System.print("Body length: %(result["length"])")
Async.cleanup(futureId)
```

### 4. 파일 처리 예제

```go
import (
    "context"
    "os"
    "path/filepath"
)

// ReadFile reads a file asynchronously
//wren:bind name=readFile(_) static
func (a *Async) ReadFile(vm *WrenVM) error {
    filename := vm.GetSlotString(1)
    
    future := GetAsyncManager().Submit(func(ctx context.Context) (interface{}, error) {
        content, err := os.ReadFile(filename)
        if err != nil {
            return nil, err
        }
        
        return map[string]interface{}{
            "content": string(content),
            "size":    len(content),
            "name":    filepath.Base(filename),
        }, nil
    })
    
    vm.SetSlotDouble(0, float64(future.ID()))
    return nil
}

// WriteFile writes content to a file asynchronously
//wren:bind name=writeFile(_,_) static
func (a *Async) WriteFile(vm *WrenVM) error {
    filename := vm.GetSlotString(1)
    content := vm.GetSlotString(2)
    
    future := GetAsyncManager().Submit(func(ctx context.Context) (interface{}, error) {
        err := os.WriteFile(filename, []byte(content), 0644)
        if err != nil {
            return nil, err
        }
        
        return map[string]interface{}{
            "filename": filename,
            "written":  len(content),
        }, nil
    })
    
    vm.SetSlotDouble(0, float64(future.ID()))
    return nil
}
```

**Wren에서 사용:**
```wren
// 파일 쓰기
var writeId = Async.writeFile("test.txt", "Hello World!")
Async.await(writeId)
Async.cleanup(writeId)

// 파일 읽기
var readId = Async.readFile("test.txt")
var fileData = Async.await(readId)
System.print("File: %(fileData["name"])")
System.print("Size: %(fileData["size"]) bytes")
System.print("Content: %(fileData["content"])")
Async.cleanup(readId)
```

## 중요한 패턴들

### 1. Context 취소 확인

긴 작업에서는 반드시 context 취소를 확인하세요:

```go
future := GetAsyncManager().Submit(func(ctx context.Context) (interface{}, error) {
    for i := 0; i < 1000000; i++ {
        // 주기적으로 취소 확인
        select {
        case <-ctx.Done():
            return nil, ctx.Err() // 취소됨
        default:
        }
        
        // 실제 작업
        doSomeWork(i)
    }
    return "completed", nil
})
```

### 2. 에러 처리

Future에서 에러가 발생하면 Wren의 `await()`나 `get()`에서 예외가 발생합니다:

```go
future := GetAsyncManager().Submit(func(ctx context.Context) (interface{}, error) {
    if someCondition {
        return nil, errors.New("something went wrong")
    }
    return "success", nil
})
```

**Wren에서 에러 확인:**
```wren
var state = Async.getState(futureId)
if (state == 2) {  // Failed
    System.print("Future failed!")
    Async.cleanup(futureId)
} else if (state == 1) {  // Completed
    var result = Async.get(futureId)
    System.print("Success: %(result)")
    Async.cleanup(futureId)
}
```

### 3. 복잡한 데이터 반환

Go의 map, slice 등을 반환하면 Wren에서 자동 변환됩니다:

```go
return map[string]interface{}{
    "users": []interface{}{
        map[string]interface{}{"name": "Alice", "age": 30},
        map[string]interface{}{"name": "Bob", "age": 25},
    },
    "total": 2,
    "timestamp": time.Now().Unix(),
}
```

**Wren에서:**
```wren
var result = Async.await(futureId)
System.print("Total users: %(result["total"])")
for (user in result["users"]) {
    System.print("User: %(user["name"]), Age: %(user["age"])")
}
```

## 새로운 클래스 만들기

Async가 아닌 별도 클래스도 가능합니다:

```go
// HttpClient provides HTTP operations
//wren:bind module=main
type HttpClient struct{}

// Get performs GET request
//wren:bind name=get(_) static
func (h *HttpClient) Get(vm *WrenVM) error {
    url := vm.GetSlotString(1)
    
    future := GetAsyncManager().Submit(func(ctx context.Context) (interface{}, error) {
        // HTTP GET 로직
        return httpGetResult, nil
    })
    
    vm.SetSlotDouble(0, float64(future.ID()))
    return nil
}
```

**Wren에서:**
```wren
import "main" for HttpClient, Async

var futureId = HttpClient.get("https://example.com")
var result = Async.await(futureId)
Async.cleanup(futureId)
```

## 빌드 및 등록

1. **wrengen 실행**: `go generate` 또는 `go run ./wrengen -dir .`
2. **빌드**: 새로운 메서드가 자동으로 등록됨
3. **Wren에서 사용**: `import "main" for ClassName`

## 베스트 프랙티스

1. **항상 cleanup 호출**: 메모리 누수 방지
2. **Context 확인**: 긴 작업에서 취소 가능하게
3. **적절한 타임아웃**: HTTP 요청 등에 타임아웃 설정
4. **에러 처리**: 의미있는 에러 메시지 반환
5. **문서화**: wren:bind 주석에 설명 추가

```go
// DownloadFile downloads a file from URL to local path
// Returns a map with download info: {filename, size, duration}
//wren:bind name=downloadFile(_,_) static
func (a *Async) DownloadFile(vm *WrenVM) error {
    // 구현...
}
```

이제 원하는 비동기 기능을 만들어보세요! 🚀