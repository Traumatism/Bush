export GB_OUTPUT="out/"

rm -rf $GB_OUTPUT
mkdir -p $GB_OUTPUT

go build -o $GB_OUTPUT/minecraft utils/check/minecraft.go
go build -o $GB_OUTPUT/generate utils/generate/generate.go
go build -o $GB_OUTPUT/randomize utils/randomize/randomize.go
go build -o $GB_OUTPUT/rcon-brute utils/rcon/brute.go

curl -O https://raw.githubusercontent.com/ignis-sec/Pwdb-Public/master/wordlists/ignis-10K.txt
curl -O https://raw.githubusercontent.com/ignis-sec/Pwdb-Public/master/wordlists/ignis-100K.txt
