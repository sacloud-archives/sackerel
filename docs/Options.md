# Options

![logo.png](images/logo.png "logo.png")

## 概要

`sackerel`ではオプションは以下のように指定します。

```bash
# 通常
$ sackerel [オプション名] [設定値]

# イコールを使う方法
$ sackerel [オプション名]=[設定値]

# 環境変数を使う方法
$ export [オプションの環境変数名]=設定値
$ saclerel 
```

`--help`と`--version`を除く全てのオプションは環境変数経由でも指定できるようになっています。

また、時間間隔(Duration)の指定では、`1h`や`10m`のような、[ParseDuration](https://golang.org/pkg/time/#ParseDuration)の仕様にそった文字列を指定できます。


## オプション一覧

### 必須項目

以下の項目は必須入力項目です。

|名称            | エイリアス                           | 環境変数                         | デフォルト値 | 説明                          |
|---------------------------------------------|-------------------------------------|---------------------------------|:----------:|-------------------------------|
| `--token`     | `--sakuracloud-access-token`        | SAKURACLOUD_ACCESS_TOKEN        | -          | さくらのクラウド APIトークン     |
| `--secret`    | `--sakuracloud-access-token-secret` | SAKURACLOUD_ACCESS_TOKEN_SECRET | -          | さくらのクラウド APIシークレット  |
| `--apikey`    | `--mackerel-apikey`                 | MACKEREL_APIKEY                 | -          | Mackerel APIキー              |


### 動作オプション


|名称                    | エイリアス                   | 環境変数                         | デフォルト値 | 説明                          |
|------------------------|----------------------------|---------------------------------|:----------:|-------------------------------|
| `--zones`              | `--sakuracloud-zones`      | SAKURACLOUD_ZONES               | `is1b,tk1a`| 情報取得対象のゾーン名(複数指定可) |
| `--interval`           | `--timer-job-interval`     | SACKEREL_TIMER_JOB_INTERVAL     | `2m`(2分)  | 情報収集ジョブの実行間隔     |
| `--period`             | `--metrics-history-period` | SACKEREL_METRICS-HISTORY-PERIOD | `15m`(15分)| メトリクスの収集範囲<br >/(実行時からどれだけ遡って取得するか) |
| `--port`               | `--healthcheck-port`       | SACKEREL_HEALTHCHECK_PORT       | `39700`    | ヘルスチェック用Webサーバーのポート番号 |
| `--disable-healthcheck`| -                          | SACKEREL_DISABLE_HEALTHCHECK    | `false`    | ヘルスチェックの無効化 |
| `--skip-init`          | -                          | SACKEREL_SKIP_INIT              | `false`    | 起動時の初期化処理(グラフ定義など)の無効化 |


### パフォーマンスチューニング用オプション

これらはパフォーマンスチューニングのためのオプションです。
設定値によっては動作が不安定になる可能性があります。
特段の事情がない場合、変更せずにご利用ください。

|名称                            | エイリアス | 環境変数                                  | デフォルト値 | 説明                                 |
|--------------------------------|----------|------------------------------------------|:----------:|--------------------------------------|
| `--api-call-interval`          | -        | SACKEREL_API_CALL_INTERVAL               | `500ms`    | スロットリング対象APIの待機時間          |
| `--job-queue-size`             | -        | SACKEREL_JOB_QUEUE_SIZE                  | `50`       | ジョブキューのバッファサイズ             |
| `--throttled-api-request-size` | -        | SACKEREL_THROTTLED_API_REQUST_QUEUE_SIZE | `0`        | スロットリング対象APIのバッファサイズ     |
| `--sakura-api-queue-size`      | -        | SACKEREL_SAKURA_API_REQEST_QUEUE_SIZE    | `5`        | さくらのクラウドAPI用キューのバッファサイズ|
| `--mackerel-api-queue-size`    | -        | SACKEREL_MACKEREL_API_REQEST_QUEUE_SIZE  | `5`        | MackerelAPI用キューのバッファサイズ      |
                 
### デバッグ用オプション

これらはデバッグ用にログ出力を制御します。
設定値によっては大量のログが出力されます。
特段の事情がない場合、変更せずにご利用ください。

|名称                            | エイリアス | 環境変数                         | デフォルト値 | 説明                                |
|--------------------------------|----------|---------------------------------|:----------:|------------------------------------|
| `--sakuracloud-trace-mode`     | -        | SAKURACLOUD_TRACE_MODE          | `false`    | さくらのクラウドAPI詳細ログの出力フラグ |
| `--mackerel-trace-mode`        | -        | MACKEREL_TRACE_MODE             | `false`    | MackerelAPI詳細ログの出力フラグ       |
| `--trace-log`                  | -        | SACKEREL_TRACE_LOG              | `false`    | sackerel内部動作の詳細ログ出力フラグ   |
| `--info-log`                   | -        | SACKEREL_INFO_LOG               | `true`     | INFOレベルログ出力フラグ              |
| `--warn-log`                   | -        | SACKEREL_WARN_LOG               | `true`     | WARNレベルログ出力フラグ              |
| `--error-log`                  | -        | SACKEREL_ERROR_LOG              | `true`     | ERRORレベルログ出力フラグ             |

詳細ログの中にはAPIキーなどのセキュアな情報が含まれる場合があります。
出力されるログの取り扱いについては十分にご注意ください。
