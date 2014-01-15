/*
Open Source Initiative OSI - The MIT License (MIT):Licensing

The MIT License (MIT)
Copyright (c) 2013 DutchCoders <http://github.com/dutchcoders/>

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import(
    "encoding/json"
    "net/http"
    "fmt"
    "io/ioutil"
    "log"
    _ "math"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/nu7hatch/gouuid"
    )

type BtcChinaTrade struct {
        Timestamp     Timestamps       `json:"date"`
        Price         float64         `json:"price"`
        Amount        float64         `json:"amount"`
        TransactionId int64s           `json:"tid"`
}


func (self BTCCHINA) FromJson(b []byte) (resp []BtcChinaTrade, err error)  {
    err = json.Unmarshal(b, &resp)
    return 
}

type BTCCHINA struct {
    EXCHANGE
    pair        string
}


func NewBTCCHINA(db *sql.DB) *BTCCHINA {
    m := new(BTCCHINA)
    m.db = db
    m.transactionId=int64(0)
    m.key = "btcchina"
    m.pair = "btc_usd"
    
    {
        stmt, err := m.db.Prepare("SELECT MAX(btceid) AS btceid FROM trades WHERE (exchange=?)")
        if err != nil {
            log.Fatal(err)
        }

        var transactionId sql.NullInt64
        if err := stmt.QueryRow(m.key).Scan(&transactionId); err != nil {
            log.Fatal(err)
        }
        
        if transactionId.Valid {
            m.transactionId=transactionId.Int64
        }

        // fmt.Printf("transactionId: %d", transactionId)
        stmt.Close()
    }
    return m
}


func (self *BTCCHINA) Key() string {
    return (self.key)
}

func (self *BTCCHINA) Collect() {
    response, err := http.Get(fmt.Sprintf("https://data.btcchina.com/data/trades"))
    if err != nil {
        log.Printf("%s", err)
        return
    }
    
    defer response.Body.Close()
    
    contents, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Printf("%s", err)
        return
    }
    
    trades, err:=self.FromJson(contents)
    if err != nil {
        log.Printf("%s %s", err, contents)
        return
    }
    
    tx, err := self.db.Begin()
    if err != nil {
        log.Panicf("%s", err)
    }
    
    stmt, err := tx.Prepare("INSERT INTO trades (tradeid, btceid, date, amount, price, pair, exchange) VALUES (?, ?, ?, ?, ?, ?, ?)")
    if err != nil {
         log.Panicf("%s", err)
    }
    
    transactionId:=self.transactionId;
    count:=0;
    
    for _, trade := range trades {
        if (trade.TransactionId.int64 <= self.transactionId) {
            continue
        }
        
        transactionId = MaxInt(transactionId, trade.TransactionId.int64)
        
        u4, err := uuid.NewV4()
        if err != nil {
            log.Panicf("%s", err)
        }
        
        _, err = stmt.Exec(u4.String(), trade.TransactionId.int64, trade.Timestamp.UTC(), trade.Amount, trade.Price, self.pair, self.key)
        if err != nil {
            log.Printf("%s", err)
            continue
        }
        
        count++;        
    }
    
    stmt.Close()
    
    tx.Commit()

    self.transactionId = MaxInt(self.transactionId, transactionId);
    
    log.Printf("%s: Imported %d trades", self.key, count)
}