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
    "log"
    "sync"
    "os"
    "io/ioutil"
    "launchpad.net/goyaml"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "time"
)

type EXCHANGE struct {
        transactionId int64
        db            *sql.DB
        key           string
}

type Exchange interface {
    Key() string
    Collect()
}

type CONFIG struct {
    DB string
}

func main() {
    var wg sync.WaitGroup
  
    data, err := ioutil.ReadFile("harvest.yaml")
    if err != nil {
        log.Panic(err)  
    }
    
    config := CONFIG{}
    
    err = goyaml.Unmarshal([]byte(data), &config)
    if err != nil {
        log.Panic(err) 
    }

    if (os.Getenv("DB")!="") {
        config.DB = os.Getenv("DB")
    }
    
    log.Printf("Connecting to database %s.", config.DB)
    
    db, err := sql.Open("mysql", config.DB)
    if err != nil {
        log.Panic(err)  
    }

    defer db.Close()

    exchanges  := [...]Exchange{
                            NewBTCE(db),
                            NewBITSTAMP(db),
                            NewBTCCHINA(db),
                            NewMTGOX(db),
                            NewOKCOIN(db, "ltc_cny"),
                            NewOKCOIN(db, "btc_cny"),
                        }
    
    for i, exchange:=range exchanges {
        wg.Add(1)
        
        go func(i int, exchange Exchange) {
            log.Printf("starting exchange %d %t....", i, exchange.Key)        
            defer wg.Done()
            
            // endless loop
            for {
                exchange.Collect()
                time.Sleep(time.Millisecond * 5000)
            }
            log.Printf("finished exchange %d %t....", i, exchange.Key)                    
        }(i, exchange)
    }
    
    wg.Wait()
}
