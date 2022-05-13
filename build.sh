export GB_OUTPUT="out/"

rm -rf $GB_OUTPUT
mkdir -p $GB_OUTPUT

go build -o $GB_OUTPUT/minecraft check/minecraft.go
go build -o $GB_OUTPUT/minecraft-parse check/parse/parse.go
go build -o $GB_OUTPUT/generate generate/generate.go
go build -o $GB_OUTPUT/randomize randomize/randomize.go
go build -o $GB_OUTPUT/rcon-brute rcon/brute.go

curl -O https://raw.githubusercontent.com/ignis-sec/Pwdb-Public/master/wordlists/ignis-10K.txt
curl -O https://raw.githubusercontent.com/ignis-sec/Pwdb-Public/master/wordlists/ignis-100K.txt
