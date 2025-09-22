# r2
Cloudflare R2를 Go에서 쉽게 다루기 위한 유틸리티 패키지입니다.
AWS SDK for Go v2를 기반으로 파일 업로드, 다운로드, 삭제, 리스트 조회, 객체 정보 확인 기능을 제공합니다.


## 설치
```bash
go get github.com/pro200/go-r2
```

## 초기화
먼저 Init() 함수를 통해 R2 클라이언트를 초기화해야 합니다.
```go
import "github.com/pro200/go-r2"

err := r2.Init(r2.Config{
    AccountId:       "your-account-id",
    AccessKeyID:     "your-access-key-id",
    SecretAccessKey: "your-secret-access-key",
})
if err != nil {
    panic(err)
}
```
- AccountId: Cloudflare R2 계정 ID
- AccessKeyID: R2 액세스 키
- SecretAccessKey: R2 비밀 키

## 기능

### 객체 정보 확인
```go
info, err := r2.Info("my-bucket", "path/to/object.txt")
if err != nil {
    panic(err)
}
fmt.Println("Object size:", *info.ContentLength)
```

### 객체 목록 조회
```go
files, nextToken, err := r2.List("my-bucket", "prefix/", 100)
if err != nil {
    panic(err)
}

fmt.Println("Files:", files)
fmt.Println("NextToken:", nextToken)
```
- 최대 1,000개의 객체를 조회할 수 있습니다.
- nextToken을 사용하여 다음 페이지를 조회할 수 있습니다.

### 파일 업로드
```go
err := r2.Upload("my-bucket", "./local.txt", "remote/path.txt")
if err != nil {
    panic(err)
}
```
- 기본적으로 파일 확장자를 기반으로 Content-Type을 자동 지정합니다.
- 강제로 Content-Type을 지정하려면 마지막 인자로 전달하세요:
```go
r2.Upload("my-bucket", "./local.txt", "remote/path.txt", "text/plain")
```

### 파일 다운로드
```go
err := r2.Download("my-bucket", "remote/path.txt", "./downloaded.txt")
if err != nil {
    panic(err)
}
```

### 객체 삭제
```go
err := r2.Delete("my-bucket", "remote/path.txt")
if err != nil {
    panic(err)
}
```

### 의존성
- AWS SDK for Go v2
- Cloudflare R2
- pro200/go-utils (Content-Type 판별)