# tenco

## Synopsis

cron式はUTCを強制されるのでJSTでの動作を期待する場合、時差を考慮する必要があります

tencoは `minutes` `hours` `day_of_weeks` の３つの入力に絞り時差を考慮しつつcron式に直してtf.jsonを生成します

### e.g. JSTで月曜日の0時

```yaml
    schedule:
      minutes:      0
      hours:        0
      day_of_weeks: MON
```

cron式ではUTCで日曜15時
```
cron(0 15 ? * 1 *)
```

### e.g. JSTで月曜の8-10時
```yaml
# 範囲指定で曜日を跨ぐような場合は分割します
    schedule:
      minutes:      0
      hours:        8-10
      day_of_weeks: MON
```

cron式ではUTCで日曜23時,月曜0-1時
```
cron(0 23 ? * 1 *)
cron(0 0-1 ? * 2 *)
```

## Usage

```
Usage of tenco:
  tenco config1.yaml [config2.yaml ...]
  -o string
        Write to FILE (default "-")
  -offset int
        offset (default -9)
```

## install

### Homebrew (macOS and Linux)

```
$ brew install mix3/tap/tenco
```

### Binary packages

[Releases](https://github.com/mix3/tenco/releases)
