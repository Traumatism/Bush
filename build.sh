export GB_OUTPUT="out/"

rm -rf $GB_OUTPUT
mkdir -p $GB_OUTPUT

go build -o $GB_OUTPUT/minecraft minecraft.go
go build -o $GB_OUTPUT/minecraft-parse minecraft-parse.go
go build -o $GB_OUTPUT/generate generate.go
go build -o $GB_OUTPUT/randomize randomize.go
go build -o $GB_OUTPUT/rcon-brute brute.go

curl -o out/10k.txt https://raw.githubusercontent.com/ignis-sec/Pwdb-Public/master/wordlists/ignis-10K.txt
curl -o out/100k.txt https://raw.githubusercontent.com/ignis-sec/Pwdb-Public/master/wordlists/ignis-100K.txt
