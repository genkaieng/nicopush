# niconico-notification

ニコニコ通知サーバー

Web Pushをサブスクライブして受信した通知を標準出力します。

## 実行手順

### .env ファイルを作成

```sh
cp .env.example .env
```

#### ニコニコのセッションキーを環境変数に設定

ニコニコのページ https://www.nicovideo.jp/ をブラウザで開きデベロッパーツールを開きます。(F12キー押下)

デベロッパーツールのアプリケーションタブの左ペインの `Cookie > https://www.nicovideo.jp` から`user_session`セッションキーを見つけます。
(`user_session_xxx`という形式の値があります。)

**※セッションキーは他人に教えないでください**

これを.envファイルの `SESSION=` の後ろに貼り付け。

https://github.com/genkaieng/niconico-notification/blob/ef4ceb4ac7121c7472b1a5dbf613887546c08690/.env.example#L1-L2

#### キーを環境変数に設定

暗号化まわりのキーを生成します。

``sh
go run cmd/genkeys/main.go
``

出力されたキーを.envファイルに貼り付け。

https://github.com/genkaieng/niconico-notification/blob/ef4ceb4ac7121c7472b1a5dbf613887546c08690/.env.example#L4-L8

### サーバー起動

```sh
go run cmd/subscribe/main.go
```

※起動すると `UAID` が出力されるので、その値を.envに貼り付け。(次回起動時にセッションが保持されます。)

https://github.com/genkaieng/niconico-notification/blob/ef4ceb4ac7121c7472b1a5dbf613887546c08690/.env.example#L10-L11
