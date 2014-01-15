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

type OkCoinTrade struct {
        Amount        float64s        `json:"amount"`
        Timestamp     Timestamp       `json:"date"`
        Price         float64s        `json:"price"`
        TransactionId int64           `json:"tid"`
        Type          string          `json:"type"`
}

func (self OKCOIN) FromJson(b []byte) (resp []OkCoinTrade, err error)  {
    err = json.Unmarshal(b, &resp)
    return 
}

type OKCOIN struct {
    EXCHANGE
    pair        string
}


func NewOKCOIN(db *sql.DB, pair string) *OKCOIN {
    m := new(OKCOIN)
    m.db = db
    m.transactionId=int64(0)
    m.key = "okcoin"
    m.pair = pair
    
    {
        stmt, err := db.Prepare("SELECT MAX(btceid) AS btceid FROM trades WHERE (exchange=? AND pair=?)")
        if err != nil {
            log.Fatal(err)
        }

        var transactionId sql.NullInt64
        if err := stmt.QueryRow(m.key, m.pair).Scan(&transactionId); err != nil {
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


func (self *OKCOIN) Key() string {
    return (self.key)
}

func (self *OKCOIN) Collect() {
    response, err := http.Get(fmt.Sprintf("http://www.okcoin.com/api/trades.do?symbol=%s", self.pair))
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
        log.Fatalf("%s", err)
        return
    }
    
    stmt, err := tx.Prepare("INSERT INTO trades (tradeid, btceid, date, amount, price, type, pair, exchange) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
    if err != nil {
        log.Fatalf("%s", err)
    }
    
    transactionId:=self.transactionId;
    
    count := 0
    
    for _, trade := range trades {
        if (trade.TransactionId <= self.transactionId) {
            continue
        }
        
        transactionId = MaxInt(transactionId, trade.TransactionId)
        
        u4, err := uuid.NewV4()
        if err != nil {
            log.Fatalf("%s", err)
        }
        
        // unique constraint guards uniqueness
        
        _, err = stmt.Exec(u4.String(), trade.TransactionId, trade.Timestamp.UTC(), trade.Amount.float64, trade.Price.float64, trade.Type, self.pair, self.key)
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