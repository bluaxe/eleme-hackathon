rm -rf go15
mkdir go15
tar -C ./go15 -zxf go.tar.gz
P=`pwd`/go15/go
# echo $P
export GOROOT=$P
# echo $GOROOT
export PATH=$GOROOT/bin:$PATH
echo $PATH
go version
go build main.go