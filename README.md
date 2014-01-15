# install
gvm install go1.1.1
gvm use go1.1.1 --default

export GOPATH=$HOME/golibs/

go get github.com/go-sql-driver/mysql
go get launchpad.net/goyaml
go run harvest/harvest.go harvest/btce.go harvest/mtgox.go  harvest/okcoin.go harvest/btcchina.go  harvest/utils.go harvest/bitfinex.go

