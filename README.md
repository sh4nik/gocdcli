# GoPipe

GoPipe is a simple CLI tool for viewing GoCD pipeline information

```
Usage:
  gopipe [command]

Available Commands:
  compare     Compare two pipelines
  diff        Diff two pipelines
  help        Help about any command
  list        List all pipelines

Flags:
      --config string     config file (default is $HOME/.gopipe.json)
  -h, --help              help for gopipe
      --password string   Password
      --url string        URL of GoCD application
      --username string   Username
```

## List Pipelines

List all pipelines using the following command

```sh
gopipe list
```

## Compare Pipelines

Compare two pipelines using the following command

```sh
gopipe compare <pipeline-one> <pipeline-two>
```

## Diff Pipelines

View differences between two pipelines using the following command

```sh
gopipe diff <pipeline-one> <pipeline-two>
```

## Configurations

Configurations can be passed in as flags:

- `--url http://.../go/api/`
- `--username sam`
- `--password abc123`

Or read from `$HOME/.gopipe.json` in the following format:

```json
{
  "url": "http://.../go/api/",
  "username": "sam",
  "password": "abc123"
}
```
