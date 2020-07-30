Set-Variable CGO_ENABLED=0
go build -o dist/peepingbot.exe
Copy-Item .\.peepingbot.yaml .\dist
