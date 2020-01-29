# Backend

## Running a Server
1. `cp config.example.yml config.yml`
2. 相应位置填入石墨文档 Cookie 以及文档 ID
3. 使用 `go run .` 即可运行项目

## Usage
- `GET /accommodations` 返回 JSON 形式的住宿信息列表
- `GET /platforms/psychological` 返回 JSON 形式的心理咨询机构列表
- `GET /platforms/medical` 返回 JSON 形式的线上医疗平台列表
