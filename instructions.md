 set GOOS=linux&& set GOARCH=amd64&& go build -o ..\build\bootstrap main.go  //compiling in windows
 zip -jrm bootstrap.zip bootstrap
 upload to lambda



 
