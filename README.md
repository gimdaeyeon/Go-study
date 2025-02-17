# Go-study

## 1. 기본 코드구성

```go
package main

func main() {
}
```

- java처럼 기본적으로 main함수 안에서 실행

## 2. 코드 실행

```bash
go run <파일명>
go run test.go
```

## 3. 코드 빌드

```bash
go build <파일명>
go build test.go
```

- 코드를 빌드하게 되면 .exe파일이 생성된다.
- 해당 exe를 실행하게 되면 기존에 코드로 작성한 내용이 실행된다.

## 4. 기본 출력

```go
package main

func main() {
	fmt.Println("Hello, World!")
}
```

- fmt 패키지를 사용하여 출력
- 출력 후 줄바꿈이 된다.

## 5. 변수 선언

변수는 var 키워드 뒤에 변수명을 적고, 그 뒤에 변수타입을 적는다.   
예를 들어, 아래는 a라는 정수(int) 변수를 선언한 것이다.

```go
var a int = 10
```
go 에서는 할당되는 값을 보고 그 타입을 추론하는 기능이 자주 사용된다. 즉, 아래 코드에서 i는 정수형으로 1이 할당되고, s는 문자열로 Hi가 할당된다

```go
var i = 1
var s = "Hi"
```

변수를 선언하는 또 다른 방식으로 Short Assignment Statement(:=)를 사용할 수 있다.   
즉, `var i = 1` 대신 `i := 1` 처럼 var 를생략하고 사용할 수 있다.   
하지만 이러한 표현은 함수 내에서만 사용할 수 있으며, 함수 밖에서는 var를 사용해야한다. go에서 볏와 상수는 함수 밖에서도 사용할 수 있다.

## 6. 상수
상수는 const 키워드를 사용하여 선언한다. const 키우드 뒤에 상수명을 적고, 그 뒤에 상수타입, 그리고 상수 값을 할당한다.

```go
const i = 1
const s = "Hi"
```
go 에서는 할당되는 값을 보고 그 타입을 추론하는 기능이 자주 사용된다.


여러 개의 상수를 묶어서 지정할 수 있는데 아래와 같이 괄호 안에 상수들을 나열하여 사용할 수 있다.

```go
const (
	Visa = "Visa"
    Master = "MasterCard"
    Amex = "American Express"
)
```
한가지 유용한 팁으로 상수값을 0부터 순차적으로 부여하기 위해 `iota` 키워드를 사용할 수 있다.
이 경우 iota가 지정된 Apple에는 0이 할당되고, 나머지 상수들을 순서대로 1씩 증가된 값을 부여받는다.
```go
const (
	Apple = iota // 0
    Grape        // 1
    Orange       // 2
)
```

## 7. Go 데이터 타입

1. 부울린 타입
    - bool
2. 문자열 타입
    - string: string은 한 번 생성되면 수정될 수 없는 Immutable 타입
3. 정수형 타입
   - int, int8, int16, int32, int64
   - uint, uint8, uint16, uint32, uint64, uintptr
4. Float 및 복소수 타입
   - float32, float64, complex64, complex128
5. 기타 타입
   - byte: unit8과 동일하며 바이트 코드에 사용
   - rune: int32와 동일하며 유니코드 코드포인트에 사용한다.

## 8. if 문
Go의 if 조건문은 조건식을 괄호()로 둘러싸지 않아도 된다. 그리고 반드시 조건 블럭 시작 브레이스`{`를 if문과 같은 라인에 두어야 한다.

그리고 if문의 조건식은 반드시 Boolean 식으로 표현되어야 한다. 이점은 c/c++같은 언어들이 조건식에 1, 0 과 같은 숫자를 쓸 수 있는것과 대조적이다.
```go
if k == 1 {  //같은 라인
    println("One")
}
```

## 9. for 문
Go에서 반복문은 for 루프를 이용한다. go는 반복문이 for 하나 밖에 없다.
```go
func main() {
    sum := 0
    for i := 1; i <= 100; i++ {
        sum += i
    }
    println(sum)
}
```
Go에서 for루프는 초기값과 증감식을 생략하고 조건식만을 사용할 수 있는데, 이는 마치 for 루프가 다른 언어의 while 루프와 같이 쓰이도록 한다.
```go
func main() {
    n := 1
    for n < 100 {
        n *= 2      
        //if n > 90 {
        //   break 
        //}     
    }
    println(n)
}
```
for 루프로 무한루프를 만들려면 초기값, 조건식, 증감 모두를 생략하면 된다.

### for range
for range문은 컬렉션으로 부터 한 요소씩 가져와 차례로 for 블럭의 문장들을 실행한다. 이는 다른 언어의 foreach와 비슷하다.
```go
names := []string{"홍길동", "이순신", "강감찬"}
 
for index, name := range names {
    println(index, name)
}
```

## 10. 함수
go에서 함수는 func 키워드를 사용하여 정의한다. func 뒤에 함수명을 적고 괄호 () 안에 그 함함에 전달하는 파라미터들을 적게 된다. 함수 파라미터는 0개 이상 사용할 수 있는데, 각 파라미터는 파라미터 명 뒤에 int, stringt 등의 파라미터 타입을 적어서 정의한다.   
함수의 리턴 타입은 파라미터 괄호 () 뒤에 적게 되는데, 이는 c와 같은 다른 언어에서 리턴 타입을 함수명 앞에 쓰는 것과 대조적이다.(typescript와 비슷)
```go
func say(msg string) {
    println(msg)
}
```