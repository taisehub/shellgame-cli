# シェルゲ-(Work in Progress)
シェルゲーは対戦型コマンドラインクイズゲームです。  
プレイヤーごとに独立したシェルゲーコンテナを用意しており、WebSocketを用いてシェルゲーコンテナに接続します。

## 🚧DEMO🚧
(未完成なので、できてるところまで雰囲気DEMOです)

https://user-images.githubusercontent.com/59153204/199500110-e39afc88-5510-4039-96a8-1e212ebe18f7.mp4


## 🚧Usage🚧
シェルゲーサーバを立てる
```
$ git clone https://github.com/taise-hub/shellgame-cli
$ cd server
$ go run cmd/shellgame/main.go
```
 
シェルゲークライアントを実行する
```bash
$ git clone https://github.com/taise-hub/shellgame-cli
$ cd client
$ go run cmd/shellgame/main.go
```
