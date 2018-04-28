# GOOS=windows GOARCH=386
mkdir -p mac-linux/app

cd cmd

go build -o ../mac-linux/app/app

cd ..

cp -r assets mac-linux/assets

zip -r mac-linux-quoridor.zip mac-linux 
rm -r mac-linux

