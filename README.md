# nicopush

ニコニコ通知サーバー

Web Pushをサブスクライブして受信した通知を標準出力します。

## 実行手順

### 0. ソースコードを落とす

```sh
git clone git@github.com:genkaieng/nicopush.git
```

### 1. .env ファイルを作成

```sh
cp .env.example .env
```

#### 1.1 ニコニコのセッションキーを環境変数に設定

ニコニコのページ https://www.nicovideo.jp/ をブラウザで開きデベロッパーツールを開きます。(F12キー押下)

デベロッパーツールの**アプリケーションタブ**の左ペインの `ストレージ > Cookie > https://www.nicovideo.jp` から`user_session`セッションキーを見つけます。
(`user_session_xxx`という形式の値があります。)

**※セッションキーは他人に教えないでください**

これを.envファイルの `SESSION=` の後ろに貼り付け。

https://github.com/genkaieng/nicopush/blob/ef4ceb4ac7121c7472b1a5dbf613887546c08690/.env.example#L1-L2

#### 1.2 キーを生成&環境変数に設定

暗号化まわりのキーを生成します。

```sh
go run cmd/genkeys/main.go
```

出力されたキーを.envファイルの以下の部分に貼り付け。

https://github.com/genkaieng/nicopush/blob/ef4ceb4ac7121c7472b1a5dbf613887546c08690/.env.example#L4-L8

### 2. サーバー起動

```sh
go run cmd/subscribe/main.go
```

※起動すると `UAID` が出力されるので、その値を.envに貼り付け。(WebPushサーバーから吐き出されるユーザー識別子。セッションを保持してくれます。)

https://github.com/genkaieng/nicopush/blob/ef4ceb4ac7121c7472b1a5dbf613887546c08690/.env.example#L10-L11
