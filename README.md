# Azd Java Project Analyzer

This project outputs a binary executable, which can analyze Java project and output azure.yaml. 
The generated azure.yaml can be used to deploy project by [azd](https://github.com/Azure/azure-dev).

## Build binary executable

```shell
go build
```

## Executable binary executable

### 1. Example 1: Change directory then run ajpa

```shell
cd ${WORKING_DIRECTORY}
./ajpa
```

### 2. Example 2: Run ajpa with `-cwd` parameter

```shell
./ajpa -cwd ${WORKING_DIRECTORY}
```

