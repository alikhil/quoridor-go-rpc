cd cmd

GOOS=windows GOARCH=386 go build -o app.exe

cd ..

mkdir -p windows/app
cp cmd/app.exe windows/app/app.exe
cp -r assets windows/assets

zip -r windows-quoridor.zip windows 
rm -r windows

