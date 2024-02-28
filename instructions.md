 set GOOS=linux&& set GOARCH=amd64&& go build -o ..\build\bootstrap main.go  //compiling in windows
 zip -jrm bootstrap.zip bootstrap
 upload to lambda



 My endpoint is https://u8pgaiegjj.execute-api.eu-north-1.amazonaws.com/staging/