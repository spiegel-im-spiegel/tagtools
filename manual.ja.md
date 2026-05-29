# tagtools

Hugo タグ運用向けの小さな CLI ツールです。

このリポジトリは `spiegel-im-spiegel/github-pages-env` から独立化した `tagtools` の本体です。

## 機能

- front matter の tags から `.github/workflows/tagslist.csv` を生成します。
- 直近1年の記事から `data/toptags.json` を生成します。
- 生成結果を golden ファイルと比較する verify を提供します。
- verify は `--debug` 指定で JSON 形式の結果を出力できます。

## 前提

- Go 1.26 以上

## ビルド

```bash
go build -o ./bin/tagtools .
```

または直接インストール:

```bash
go install github.com/spiegel-im-spiegel/tagtools@latest
```

## コマンド

### tagslist

タグ件数を集計し，既存 CSV の `means` を引き継いで出力します。

```bash
./bin/tagtools tagslist \
  --content-dir ./content \
  --out ./.github/workflows/tagslist.csv \
  --inherit-means-from ./.github/workflows/tagslist.csv
```

### toptags

直近1年の記事から上位タグを集計し，タグ名配列を出力します。

```bash
./bin/tagtools toptags \
  --content-dir ./content \
  --out ./data/toptags.json \
  --top-n 15 \
  --today 2026-05-18
```

### all

既定値で `toptags` -> `tagslist` の順に実行します。

```bash
./bin/tagtools all
```

### verify

生成結果を期待ファイルと比較します。

```bash
./bin/tagtools verify \
  --content-dir testdata/content \
  --expected-tagslist testdata/golden/tagslist.csv \
  --expected-toptags testdata/golden/toptags.json \
  --inherit-means-from testdata/inherit/tagslist.csv \
  --today 2026-05-18
```

## 環境変数

`TAGTOOLS_` プレフィックス付きの環境変数で各オプションを上書きできます。
優先順位は `flag > env > default` です。

例:

- `TAGTOOLS_TOP_N=20`
- `TAGTOOLS_TODAY=2026-05-18`
- `TAGTOOLS_CONTENT_DIR=./content`
- `TAGTOOLS_INHERIT_MEANS_FROM=./.github/workflows/tagslist.csv`
- `TAGTOOLS_DEBUG=true`（`verify` 用）

失敗時は次の情報を表示します。

- 失敗ステージ（`tagslist` または `toptags`）
- expected/actual のファイルパス
- 最初に差分が出た行番号
- expected/actual 行の抜粋
- `reproduce:` で始まる再現コマンド

## --debug の JSON 出力

`verify` に `--debug` を付けると JSON 形式で出力します。

```bash
./bin/tagtools verify --debug \
  --content-dir testdata/content \
  --expected-tagslist testdata/golden/tagslist.csv \
  --expected-toptags testdata/golden/toptags.json \
  --inherit-means-from testdata/inherit/tagslist.csv \
  --today 2026-05-18
```

成功時の例:

```json
{"status":"ok","reproduce_command":"./bin/tagtools verify ... --debug"}
```

失敗時の例:

```json
{
  "status": "error",
  "stage": "toptags",
  "message": "toptags verification failed (first diff line: 1)",
  "expected_path": "testdata/golden/toptags.json",
  "actual_path": "/tmp/tagtools-verify-xxxx/toptags.json",
  "line": 1,
  "expected_line": "...",
  "actual_line": "...",
  "reproduce_command": "./bin/tagtools verify ... --debug"
}
```

## テスト

```bash
task test
```
