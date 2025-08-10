# nicopush

ニコニコ通知サーバー

Web Pushをサブスクライブして受信した通知を標準出力します。

## インストール手順

以下のいずれかの方法

### go install

```sh
go install github.com/genkaieng/nicopush@latest
```

### ソースコードを落とす

```sh
git clone git@github.com:genkaieng/nicopush.git
```

## 実行手順

### 1. .env ファイルを作成

```sh
cp .env.example .env
```

#### 1.1 ニコニコのセッションキーを環境変数に設定

ニコニコのページ https://www.nicovideo.jp/ をブラウザで開きデベロッパーツールを開きます。(F12キー押下)

デベロッパーツールの**アプリケーションタブ**の左ペインの `ストレージ > Cookie > https://www.nicovideo.jp` から`user_session`セッションキーを見つけます。
(`user_session_xxx`という形式の値があります。)

**※セッションキーは他人に教えないでください**

環境変数 `NICOLIVE_SESSION=` を設定する。

https://github.com/genkaieng/nicopush/blob/2e0ef5cf1fa5e13a1a2b8e46cd727e9409b5d3d0/.env.example#L3-L4

#### 1.2 キーを生成&環境変数に設定

暗号化まわりのキーを生成します。

```sh
nicopush genkeys
```

出力されたキーを環境変数に設定する

https://github.com/genkaieng/nicopush/blob/2e0ef5cf1fa5e13a1a2b8e46cd727e9409b5d3d0/.env.example#L6-L10

### 2. サーバー起動

```sh
nicopush subscribe
```

※起動すると `NICOPUSH_UAID` が出力されるので、その値を環境変数に設定する。(WebPushサーバーから吐き出されるユーザー識別子。セッションを保持してくれるらしい。)

https://github.com/genkaieng/nicopush/blob/2e0ef5cf1fa5e13a1a2b8e46cd727e9409b5d3d0/.env.example#L12-L13
