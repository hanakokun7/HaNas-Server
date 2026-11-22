# HaNas

[日本語](README.ja.md) | [English](README.md)

HaNas는 Go와 SQLite를 사용하는 셀프 호스팅 개인 파일 저장소입니다. 웹 브라우저에서 계정 관리, 파일 및 폴더 작업, 휴지통, 공유 링크와 미디어 미리보기를 사용할 수 있습니다.

## 주요 기능

- 사용자 회원가입 및 로그인
- JWT 토큰 기반 인증 쿠키
- 폴더 및 파일 생성
- 파일 업로드와 덮어쓰기
- 복사, 이동, 이름 변경 및 삭제
- 휴지통 및 30일 후 영구 삭제
- 파일 공유 링크 생성 및 삭제
- 이미지와 동영상 썸네일 생성
- 이미지, 동영상, 오디오와 문서 미리보기
- 계정 및 계정의 모든 파일 삭제

## 기술 스택

- Go `1.25.1`
- `net/http`
- GORM과 glebarez SQLite
- JWT
- bcrypt
- React 기반 정적 프런트엔드
- 선택 사항: `ffmpeg`

프런트엔드 정적 파일은 Go 바이너리에 포함되므로 별도의 프런트엔드 빌드 과정이 필요 없습니다.

## 요구 사항

- Go 1.25 이상
- 포트 `80`을 사용할 수 있는 권한
- 이미지 썸네일 처리를 제외한 기본 기능에는 별도 시스템 의존성이 없습니다
- 동영상 썸네일을 사용하려면 `ffmpeg`가 `PATH`에 있어야 합니다

macOS 또는 Linux에서 포트 80을 사용하려면 일반적으로 관리자 권한이 필요할 수 있습니다.

```bash
brew install go ffmpeg
```

위 명령은 macOS의 선택 예시이며, 운영체제의 패키지 관리자로 Go와 ffmpeg를 설치할 수 있습니다.

## 빠른 시작

프로젝트 루트에서 빌드합니다.

```bash
cd /Users/user/Desktop/server
go mod download
go build -o hanas .
```

실행 파일과 같은 폴더에 `config.js`를 두고 서버를 시작합니다.

```bash
./hanas
```

기본 설정은 다음과 같습니다.

```javascript
window.APP_CONFIG = {
  BACKEND_URL: 'http://127.0.0.1:80'
}
```

브라우저에서 `http://127.0.0.1` 또는 서버의 호스트 주소로 접속합니다.

포트 80 사용 권한이 없다면 임시로 권한을 부여해 실행할 수 있습니다.

```bash
sudo ./hanas
```

운영 환경에서는 TLS 리버스 프록시 뒤에 두고 `config.js`의 `BACKEND_URL`을 공개 URL로 변경하는 것을 권장합니다.

## 실행 설정

`config.js`는 실행 파일이 위치한 디렉터리에서 읽습니다. 바이너리를 다른 경로로 이동하거나 배포할 때는 프로젝트 루트의 설정 파일을 실행 파일 옆에 함께 복사해야 합니다.

예:

```text
/opt/hanas/
├── hanas
├── config.js
├── database.db
├── data/
└── thumbnails/
```

서버는 고정으로 `:80`에서 수신하며 환경 변수로 포트를 변경할 수 없습니다.

## 데이터 저장 위치

애플리케이션은 실행 위치가 아닌 현재 작업 디렉터리를 기준으로 다음 경로를 사용합니다.

| 경로 | 용도 |
| --- | --- |
| `database.db` | 사용자, 노드, 공유 및 휴지통 메타데이터 |
| `data/` | 업로드된 원본 파일 |
| `thumbnails/` | 생성된 썸네일 |

데이터를 배포할 때는 이 세 가지 위치를 모두 백업해야 합니다.

```bash
mkdir -p backups
cp database.db backups/database.db
cp -R data backups/data
cp -R thumbnails backups/thumbnails
```

서버가 실행 중인 상태에서는 SQLite 데이터베이스를 안전하게 백업해야 합니다. 운영 환경에서는 정지 후 백업하거나 데이터베이스 전용 백업 절차를 사용하세요.

## 프로젝트 구조

```text
.
├── app.go            # 서버 시작 및 HTTP 라우트
├── logic.go          # 인증, 파일 및 공유 API 처리
├── exec.go           # 파일 시스템 작업 및 데이터 조작
├── object.go         # 데이터 모델
├── go.mod            # Go 의존성
├── index.html        # 웹 애플리케이션 진입점
├── config.js         # 프런트엔드 백엔드 주소 설정
├── assets/
│   ├── index.js      # 빌드된 프런트엔드
│   ├── index.css     # 프런트엔드 스타일
│   ├── icon.ico
│   └── apple-touch-icon.png
├── data/             # 업로드 파일 저장소
├── thumbnails/       # 썸네일 저장소
└── database.db       # SQLite 데이터베이스
```

## 동작 방식

1. 서버 시작 시 SQLite 데이터베이스를 열고 필요한 테이블을 자동으로 생성합니다.
2. 최초 시작 시 JWT 서명을 위한 비밀 키를 생성해 데이터베이스에 저장합니다.
3. 회원가입 시 사용자별 루트 폴더가 생성됩니다.
4. 로그인 성공 시 `HttpOnly`, `SameSite=Strict` 쿠키에 JWT를 저장합니다.
5. 업로드된 파일은 `data/` 아래의 숫자형 파일 ID로 저장되고, 노드 메타데이터는 SQLite에 저장됩니다.
6. 휴지통으로 이동한 항목은 30일 뒤 백그라운드 작업을 통해 영구 삭제됩니다.

JWT 서명 키는 데이터베이스에 저장되므로 재배포 시 기존 토큰과 세션 호환성을 유지하려면 `database.db`를 함께 보존해야 합니다.

## 주요 API

대부분의 사용자 API는 `token` 쿠키가 필요합니다.

| 메서드 | 경로 | 설명 |
| --- | --- | --- |
| `POST` | `/register` | 계정 생성 |
| `POST` | `/login` | 로그인 |
| `POST` | `/logout` | 로그아웃 |
| `GET` | `/me` | 현재 사용자 조회 |
| `GET` | `/node/` | 루트 폴더 조회 |
| `GET` | `/node/{id}` | 폴더 조회 |
| `GET` | `/file/{nodeId}` | 파일 다운로드 |
| `GET` | `/thumbnail/{nodeId}` | 썸네일 조회 |
| `POST` | `/upload` | 파일 또는 폴더 업로드 |
| `POST` | `/copy` | 항목 복사 |
| `POST` | `/move` | 항목 이동 |
| `POST` | `/rename` | 항목 이름 변경 |
| `POST` | `/delete` | 항목 영구 삭제 |
| `POST` | `/trash/in` | 휴지통으로 이동 |
| `GET` | `/trash` | 휴지통 목록 |
| `GET` | `/trash/info/{id}` | 휴지통 항목 상세 |
| `POST` | `/trash/restore/{id}` | 휴지통 항목 복원 |
| `POST` | `/share/create` | 공유 링크 생성 |
| `POST` | `/share/delete` | 공유 링크 삭제 |
| `GET` | `/s/{token}` | 공유 파일 접근 |
| `POST` | `/delete-account` | 계정 및 데이터 삭제 |

## 업로드 요청 예시

JSON 요청으로 작은 파일을 업로드할 수 있습니다.

```bash
curl -b cookie.txt \
  -H 'Content-Type: application/json' \
  -d '{
    "filename": "hello.txt",
    "is_dir": false,
    "data_base64": "aGVsbG8="
  }' \
  http://127.0.0.1/upload
```

multipart 업로드 예시:

```bash
curl -b cookie.txt \
  -F 'filename=photo.jpg' \
  -F 'oya_id=1' \
  -F 'file=@photo.jpg' \
  http://127.0.0.1/upload
```

로그인 쿠키는 실제 로그인 응답에서 저장한 쿠키를 사용해야 합니다.

## 개발 및 검증

의존성을 내려받고 빌드합니다.

```bash
go mod download
go build ./...
```

현재 프로젝트에는 기본 Go 테스트가 포함되어 있지 않습니다.
