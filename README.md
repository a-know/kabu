# kabu
## About
`kabu` は、【前株／後株判定ツール】です。

- 前株
    - 株式会社ほげほげ
- 後株
    - ほげほげ株式会社

## Usage

`kabu` は、[Google Custom Search API](https://developers.google.com/custom-search/json-api/v1/overview) を使用するため、予め [Goole Cloud Console](https://cloud.google.com/console/project) から API キーを取得し、環境変数にセットしておく必要があります。

```console
$ export GOOGLE_APIKEY=ABCDE12345
```

`kabu` は [Custom search engine](https://www.google.com/cse/) も使用するため、予め Custom Search Engine を作成し ID を取得、環境変数にセットしておく必要があります。

```console
$ export GOOGLE_CSE_ID=ABCDE12345
```

インストールは以下の方法で可能です。

```console
$ go get github.com/a-know/kabu
```

## Design

「会社名 株式会社」で Google 検索をおこなった結果に対し、「株式会社（会社名）」「（会社名）株式会社」それぞれのパターンでマッチした数を数え、カウントの多かった方を正として判定結果を表示します。

## Example
```console
$ kabu はてな
前株マッチ数:9
後株マッチ数:0
前株です！
株式会社はてな
```