# 요약 Summary
Go과 gRPC를 GCP에서 간단하게 서비스 하기 위한 프레임워크 구성

# 특징 Feature
- Proto 정의를 통한 코드 생성 자동화 (Go, Swift)
- GCP Cloud Run 배포 자동화
- Dockerfile의 빌드 대상 경로만 수정하여 복수의 서비스를 빌드

# 기술 스택 Stack
- gRPC
- Github action
- GCP Cloud Build
- GCP Cloud Storage
- GCP Cloud Run
- GCP Datastore

# 자동화 Automation
- gRPC 코드 생성
- Container build

### Protoc
```
.github/workflows/automation_proto_go.yml
```
proto 파일에 변경점 발생시 발동
```
on:
  push:
    paths:
      - '**.proto'
```
Code 생성 후 신규 커밋으로 추가
```
- name: Generate Go code from proto files
      run: |
        go mod download
        find . -name '*.proto' -exec protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative {} \;

    - name: Commit changes
      run: |
        git config --local user.email "actions@github.com"
        git config --local user.name "GitHub Actions"
        git add .
        git commit -m "Auto-generate Go code from proto files" || echo "No changes to commit"
```
### Container Build
`v*` 조건으로 태그가 발생시 빌드 후 github packages로 push
```
on:
  push:
    tags:
      - v*
```

### Service Deploy

```
gcloud builds submit .
gcloud run ...
```

## Run locally
```bash
go run sample-server/main.go
```

## CloudRun
<img width="313" alt="image" src="https://github.com/user-attachments/assets/5b41d9bf-1724-4142-a4f7-bface4f777b5" />

