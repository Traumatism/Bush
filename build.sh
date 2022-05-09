export GB_OUTPUT="out/"

rm -rf $GB_OUTPUT
mkdir -p $GB_OUTPUT

go build -o $GB_OUTPUT/check utils/check/check.go
go build -o $GB_OUTPUT/generate utils/generate/generate.go
go build -o $GB_OUTPUT/randomize utils/randomize/randomize.go

chmod +x $GB_OUTPUT/*
