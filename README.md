# langchaingo で RAGアプリを作ってみる

## Run

```
# ベクトルDBを起動
$ docker compose -f docker/compose.yml up -d
# アプリケーションを起動
# GOOGLEAI_API_KEYは、Google AI Studioから取得すること
$ GOOGLEAI_API_KEY=xxxxxxxx ENV=production PORT=8082 go run cmd/main.go
```

## Test Request

```
$ curl -X POST -d '{"question": "Go1.21で追加されたbuilt-insはなんですか"}' localhost:8082/v1/question
$ curl -X POST -d '{"page_cintent": "追加した情報をテキスト形式で追加する"}' localhost:8082/v1/add/document
```
