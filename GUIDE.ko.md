# Wren.go 런타임 기능 가이드

wren.go 런타임 시스템과 고급 기능 사용에 대한 종합 가이드입니다.

## 목차

1. [시작하기](#시작하기)
2. [기본 VM 사용법](#기본-vm-사용법)
3. [외래 함수(Foreign Functions)](#외래-함수foreign-functions)
4. [코드 생성기 (Wrengen)](#코드-생성기-wrengen)
5. [비동기 태스크](#비동기-태스크)
6. [다중 VM 인스턴스](#다중-vm-인스턴스)
7. [슬롯을 통한 데이터 교환](#슬롯을-통한-데이터-교환)
8. [모범 사례](#모범-사례)

---

## 시작하기

### 설치

```bash
go get github.com/snowmerak/wren.go
cd $(go list -f '{{.Dir}}' github.com/snowmerak/wren.go)
# Wren 라이브러리 빌드
./build_wren.sh  # Linux/Mac
# 또는
build_wren.bat   # Windows
```

### 기본 예제

```go
package main

import (
    "fmt"
    wrengo "github.com/snowmerak/wren.go"
)

func main() {
    vm := wrengo.NewVM()
    defer vm.Free()
    
    code := `System.print("Hello from Wren!")`
    result, err := vm.Interpret("main", code)
    
    if err != nil {
        fmt.Println("에러:", err)
        return
    }
    
    if result == wrengo.ResultSuccess {
        fmt.Println("성공!")
    }
}
```

---

## 기본 VM 사용법

### VM 생성하기

**기본 설정:**

```go
vm := wrengo.NewVM()
defer vm.Free()  // 항상 리소스를 해제하세요!
```

**커스텀 설정:**

```go
config := wrengo.DefaultConfiguration()
config.InitialHeapSize = 10 * 1024 * 1024  // 10MB
config.MinHeapSize = 1 * 1024 * 1024       // 1MB
config.HeapGrowthPercent = 50              // 50% 증가율

vm := wrengo.NewVMWithConfig(config)
defer vm.Free()
```

**외래 메서드 포함:**

```go
// 모든 외래 메서드를 자동으로 등록합니다
vm := wrengo.NewVMWithForeign()
defer vm.Free()
```

### Wren 코드 실행하기

```go
result, err := vm.Interpret("main", `
    var x = 42
    System.print("답은 %(x)입니다")
`)

if err != nil {
    fmt.Println("에러:", err)
}
```

### 메모리 관리

```go
// 수동 가비지 컬렉션
vm.CollectGarbage()
```

---

## 외래 함수(Foreign Functions)

외래 함수를 사용하면 Wren 스크립트에서 Go 코드를 호출할 수 있습니다.

### 수동 등록

```go
// 외래 메서드 등록
wrengo.RegisterForeignMethod(
    "main",           // 모듈 이름
    "Calculator",     // 클래스 이름
    true,             // 정적 메서드 여부
    "add(_,_)",       // 시그니처
    func(vm *wrengo.WrenVM) {
        // 슬롯에서 파라미터 가져오기
        a := vm.GetSlotDouble(1)
        b := vm.GetSlotDouble(2)
        
        // 결과 반환
        vm.SetSlotDouble(0, a + b)
    },
)
```

### Wren에서 사용하기

```wren
foreign class Calculator {
    foreign static add(a, b)
}

System.print(Calculator.add(10, 20))  // 30
```

---

## 코드 생성기 (Wrengen)

어노테이션이 달린 Go 코드에서 외래 함수 바인딩을 자동으로 생성합니다.

### 기본 사용법

**1. Go 코드에 어노테이션 추가:**

```go
//go:generate go run github.com/snowmerak/wren.go/wrengen -dir .

//wren:bind module=main
type Math struct{}

//wren:bind static
func (m *Math) Multiply(a, b float64) float64 {
    return a * b
}

//wren:bind name=add(_,_) static
func (m *Math) Add(a, b float64) float64 {
    return a + b
}
```

**2. 바인딩 생성:**

```bash
go generate
```

이렇게 하면 자동 바인딩이 포함된 `<패키지명>_wren.go` 파일이 생성됩니다.

**3. Wren에서 사용:**

```wren
foreign class Math {
    foreign static multiply(a, b)
    foreign static add(a, b)
}

System.print(Math.multiply(7, 6))  // 42
System.print(Math.add(10, 20))     // 30
```

### 어노테이션 옵션

```go
//wren:bind                        // 기본 설정 사용
//wren:bind module=mymodule        // 모듈 이름 지정
//wren:bind class=MyClass          // 클래스 이름 지정
//wren:bind name=customName(_)     // 커스텀 메서드 이름과 시그니처
//wren:bind static                 // 정적 메서드로 표시
```

### 지원되는 타입

| Go 타입 | Wren 타입 | Get 메서드 | Set 메서드 |
|---------|-----------|------------|------------|
| `float64` | Num | `GetSlotDouble()` | `SetSlotDouble()` |
| `float32` | Num | `GetSlotDouble()` | `SetSlotDouble()` |
| `int`, `int64` | Num | `GetSlotDouble()` | `SetSlotDouble()` |
| `string` | String | `GetSlotString()` | `SetSlotString()` |
| `bool` | Bool | `GetSlotBool()` | `SetSlotBool()` |
| `error` | String (abort) | - | `AbortFiber()` |

### 완전한 예제

```go
package example

//go:generate go run github.com/snowmerak/wren.go/wrengen -dir .

//wren:bind module=geometry
type Circle struct{}

//wren:bind name=area(_) static
func (c *Circle) Area(radius float64) float64 {
    return 3.14159 * radius * radius
}

//wren:bind name=circumference(_) static
func (c *Circle) Circumference(radius float64) float64 {
    return 2 * 3.14159 * radius
}
```

```wren
// geometry.wren
foreign class Circle {
    foreign static area(radius)
    foreign static circumference(radius)
}

var r = 5
System.print("면적: %(Circle.area(r))")
System.print("둘레: %(Circle.circumference(r))")
```

---

## 비동기 태스크

비동기 시스템은 오래 걸리는 작업의 논블로킹 실행을 가능하게 합니다.

### 개요

- **Future 기반**: 태스크를 제출하고 `Future`를 받아 결과를 추적
- **워커 풀**: 백그라운드 워커가 태스크를 동시에 실행
- **Context 지원**: 타임아웃과 취소 기능
- **스레드 안전**: 동시 사용에 안전

### 기본 비동기 태스크

**Go 쪽:**

```go
// 비동기 외래 메서드 등록
wrengo.RegisterAsyncMethod("slowComputation", func(args ...interface{}) (interface{}, error) {
    // 오래 걸리는 계산 시뮬레이션
    time.Sleep(2 * time.Second)
    return 42, nil
})
```

**Wren 쪽:**

```wren
foreign class SlowTask {
    foreign static start()
}

// 태스크 제출 및 future ID 받기
var futureId = SlowTask.start()

// 결과 대기 (완료될 때까지 블록)
var result = Async.await(futureId)
System.print("결과: %(result)")

// 정리
Async.cleanup(futureId)
```

### 논블로킹 확인

```wren
var futureId = SlowTask.start()

// 블로킹 없이 확인
while (!Async.isReady(futureId)) {
    System.print("아직 작업 중...")
    Fiber.yield()
}

var result = Async.get(futureId)
System.print("완료: %(result)")
Async.cleanup(futureId)
```

### 태스크 취소

```wren
var futureId = LongTask.start()

// 태스크 취소
Async.cancel(futureId)

// 상태 확인
var state = Async.getState(futureId)
System.print("상태: %(state)")  // "cancelled"

Async.cleanup(futureId)
```

### 다중 동시 태스크

```wren
var tasks = []
for (i in 1..5) {
    tasks.add(Task.start(i))
}

// 모든 태스크 대기
for (futureId in tasks) {
    var result = Async.await(futureId)
    System.print("태스크 결과: %(result)")
    Async.cleanup(futureId)
}
```

### 비동기 API 레퍼런스

| 메서드 | 설명 | 반환값 |
|--------|------|--------|
| `Async.await(id)` | Future 완료 대기 | 결과값 또는 에러 |
| `Async.isReady(id)` | Future 준비 여부 확인 | `true` 또는 `false` |
| `Async.get(id)` | 결과 가져오기 (준비된 상태여야 함) | 결과값 또는 에러 |
| `Async.cancel(id)` | 태스크 취소 | 취소 성공 시 `true` |
| `Async.getState(id)` | 태스크 상태 가져오기 | `"pending"`, `"completed"`, `"failed"`, `"cancelled"` |
| `Async.cleanup(id)` | 매니저에서 Future 제거 | - |

### Go 비동기 매니저 API

```go
// 전역 비동기 매니저 가져오기
mgr := wrengo.GetAsyncManager()

// 태스크 제출
future := mgr.Submit(func(ctx context.Context) (interface{}, error) {
    // 취소 확인
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }
    
    // 작업 수행
    result := expensiveComputation()
    return result, nil
})

// 결과 대기
result, err := future.Wait()

// 또는 타임아웃과 함께 대기
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result, err := future.WaitContext(ctx)
```

### 완전한 비동기 예제

```go
package main

import (
    "context"
    "time"
    wrengo "github.com/snowmerak/wren.go"
)

func init() {
    // 비동기 외래 메서드 등록
    wrengo.RegisterForeignMethod("main", "SlowTask", true, "start()", 
        func(vm *wrengo.WrenVM) {
            // 비동기 태스크 제출
            future := wrengo.GetAsyncManager().Submit(
                func(ctx context.Context) (interface{}, error) {
                    // 작업 시뮬레이션
                    for i := 0; i < 10; i++ {
                        select {
                        case <-ctx.Done():
                            return nil, ctx.Err()
                        default:
                            time.Sleep(200 * time.Millisecond)
                        }
                    }
                    return 42.0, nil
                },
            )
            
            // Future 저장 및 ID 반환
            id := wrengo.GetAsyncManager().StoreFuture(future)
            vm.SetSlotDouble(0, float64(id))
        },
    )
}

func main() {
    vm := wrengo.NewVMWithForeign()
    defer vm.Free()
    
    wrenScript := `
        foreign class SlowTask {
            foreign static start()
        }
        
        System.print("태스크 시작...")
        var id = SlowTask.start()
        
        System.print("다른 작업 수행 중...")
        
        while (!Async.isReady(id)) {
            System.print("아직 계산 중...")
            Fiber.yield()
        }
        
        var result = Async.get(id)
        System.print("결과: %(result)")
        Async.cleanup(id)
    `
    
    vm.Interpret("main", wrenScript)
}
```

---

## 다중 VM 인스턴스

독립적으로 실행되는 여러 VM 인스턴스를 생성할 수 있습니다.

### 주요 개념

- **레지스트리는 전역**: 외래 메서드는 모든 VM에 대해 한 번만 등록
- **실행은 VM별**: 각 VM은 자신만의 상태, 변수, 메모리를 가짐
- **스레드 안전**: 여러 VM을 동시에 사용해도 안전
- **래퍼 제한**: 각 VM은 최대 300개의 외래 메서드 사용 가능

### 예제

```go
// 두 개의 독립적인 VM 생성
vm1 := wrengo.NewVMWithForeign()
vm2 := wrengo.NewVMWithForeign()
defer vm1.Free()
defer vm2.Free()

// VM1은 자체 상태를 가짐
vm1.Interpret("main", `
    var x = 10
    System.print("VM1: x = %(x)")
`)

// VM2는 완전히 독립적인 상태를 가짐
vm2.Interpret("main", `
    var x = 20
    System.print("VM2: x = %(x)")
`)

// 둘 다 같은 외래 메서드를 호출 가능
vm1.Interpret("main", `System.print(Math.multiply(3, 4))`)
vm2.Interpret("main", `System.print(Math.multiply(5, 6))`)
```

### 사용 사례

- **샌드박싱**: 스크립트 실행 격리
- **멀티 테넌시**: 다른 사용자를 위한 스크립트 실행
- **테스트**: 병렬로 다른 시나리오 테스트
- **모듈 격리**: 모듈 상태를 분리해서 유지

---

## 슬롯을 통한 데이터 교환

슬롯은 Go와 Wren 사이에 데이터를 전달하는 데 사용됩니다.

### 슬롯 이해하기

- **슬롯 0**: 반환값 (Go → Wren)
- **슬롯 1+**: 파라미터 (Wren → Go)
- 사용 전에 슬롯이 존재하는지 확인해야 함

### Wren에서 값 가져오기

```go
func myForeignMethod(vm *wrengo.WrenVM) {
    // 슬롯 존재 확인
    vm.EnsureSlots(3)
    
    // 파라미터 가져오기
    numValue := vm.GetSlotDouble(1)
    strValue := vm.GetSlotString(2)
    boolValue := vm.GetSlotBool(3)
    
    // 처리...
}
```

### Wren에 값 반환하기

```go
func myForeignMethod(vm *wrengo.WrenVM) {
    result := 42.0
    
    // 슬롯 0을 통해 반환
    vm.SetSlotDouble(0, result)
}
```

### 리스트

```go
// Wren에서 리스트 가져오기
vm.EnsureSlots(2)
count := vm.GetListCount(1)

for i := 0; i < count; i++ {
    vm.GetListElement(1, i, 0)
    value := vm.GetSlotDouble(0)
    // 값 처리...
}

// Go에서 리스트 생성
vm.SetSlotNewList(0)
vm.SetSlotDouble(1, 10.0)
vm.InsertInList(0, -1, 1)  // 10.0 추가
vm.SetSlotDouble(1, 20.0)
vm.InsertInList(0, -1, 1)  // 20.0 추가
```

### 맵

```go
// 맵 생성
vm.SetSlotNewMap(0)

// 키-값 쌍 설정
vm.SetSlotString(1, "name")
vm.SetSlotString(2, "Alice")
vm.SetMapValue(0, 1, 2)

vm.SetSlotString(1, "age")
vm.SetSlotDouble(2, 30)
vm.SetMapValue(0, 1, 2)
```

### 에러 처리

```go
func myForeignMethod(vm *wrengo.WrenVM) {
    if errorCondition {
        vm.SetSlotString(0, "에러 메시지")
        vm.AbortFiber(0)
        return
    }
    
    // 정상 실행...
}
```

---

## 모범 사례

### 1. 항상 VM 해제하기

```go
vm := wrengo.NewVM()
defer vm.Free()  // ✅ 항상 defer 사용
```

### 2. 코드 생성기 사용하기

```go
// ✅ 좋음: wrengen 어노테이션 사용
//wren:bind static
func (m *Math) Add(a, b float64) float64 {
    return a + b
}

// ❌ 나쁨: 수동 등록 (에러 발생 가능)
wrengo.RegisterForeignMethod("main", "Math", true, "add(_,_)", ...)
```

### 3. 비동기 태스크 상태 확인하기

```wren
// ✅ 좋음: 가져오기 전에 항상 준비 여부 확인
if (Async.isReady(id)) {
    var result = Async.get(id)
}

// ❌ 나쁨: 확인 없이 가져오기 (블록되거나 에러 발생 가능)
var result = Async.get(id)
```

### 4. Future 정리하기

```wren
var id = Task.start()
var result = Async.await(id)
Async.cleanup(id)  // ✅ 완료 후 항상 정리
```

### 5. 에러 처리하기

```go
result, err := vm.Interpret("main", code)
if err != nil {
    log.Printf("에러: %v", err)  // ✅ 에러 처리
    return
}
```

### 6. 취소를 위한 Context 사용하기

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

future := mgr.Submit(func(ctx context.Context) (interface{}, error) {
    select {
    case <-ctx.Done():
        return nil, ctx.Err()  // ✅ 취소 존중
    default:
        // 작업 수행...
    }
})
```

### 7. VM 개수 제한하기

```go
// ✅ 좋음: VM 재사용
var vmPool = make(chan *wrengo.WrenVM, 10)

// ❌ 나쁨: 무제한 VM 생성
for i := 0; i < 10000; i++ {
    vm := wrengo.NewVM()  // 메모리 집약적!
}
```

### 8. 모듈 구성

```
project/
├── main.go              # 메인 애플리케이션
├── math.go              # Math 외래 함수
├── geometry.go          # Geometry 외래 함수
├── main_wren.go         # 생성된 바인딩 (wrengen에서)
└── scripts/
    ├── main.wren        # 메인 스크립트
    ├── math.wren        # Math 모듈
    └── geometry.wren    # Geometry 모듈
```

---

## 예제

완전한 작동 예제는 `example/` 디렉토리를 참조하세요:

- `example/main.go` - 종합적인 기능 데모
- `example/async_visual/` - 비주얼 비동기 실행 데모
- `example/README.md` - 상세한 예제 문서

---

## 더 읽어보기

- [Wren 언어 가이드](https://wren.io/)
- [Wrengen 문서](./wrengen/README.md)
- [예제 문서](./example/README.md)
- [Wren 임베딩 가이드](https://wren.io/embedding/)

---

## 지원

- **이슈**: https://github.com/snowmerak/wren.go/issues
- **토론**: https://github.com/snowmerak/wren.go/discussions

---

**즐거운 Wren 코딩 되세요! 🚀**
