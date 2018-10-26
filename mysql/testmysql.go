package main

import (
    "strconv"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "fmt"
    "time"
    "log"
)
var db = &sql.DB{}

func init(){
    db,_ = sql.Open("mysql",  "root:mysqlmanager@tcp(127.0.0.1:3306)/test")
} 
func main() {
	insert()
    query()	
}

func insert() {
    
    //方式1 insert
    //strconv,int转string:strconv.Itoa(i)
    start := time.Now()
    for i := 1001;i<=1011;i++{
        //每次循环内部都会去连接池获取一个新的连接，效率低下
        db.Exec("INSERT INTO person(personid,personname) values(?,?)","a"+strconv.Itoa(i),"xnn"+strconv.Itoa(i))
    }
    end := time.Now()
    fmt.Println("方式1 insert total time:",end.Sub(start).Seconds())
    
    //方式2 insert
    start = time.Now()
    for i := 1101;i<=1110;i++{
        //Prepare函数每次循环内部都会去连接池获取一个新的连接，效率低下
        stm,_ := db.Prepare("INSERT INTO person(personid,personname) values(?,?)")
        stm.Exec("b"+strconv.Itoa(i),"xnn"+strconv.Itoa(i))
        stm.Close()
    }
    end = time.Now()
    fmt.Println("方式2 insert total time:",end.Sub(start).Seconds())
    
    //方式3 insert
    start = time.Now()
    stm,_ := db.Prepare("INSERT INTO person(personid,personname) values(?,?)")
    for i := 1201;i<=1211;i++{
        //Exec内部并没有去获取连接，为什么效率还是低呢？
        stm.Exec("c"+strconv.Itoa(i),"xnn"+strconv.Itoa(i))
    }
    stm.Close()
    end = time.Now()
    fmt.Println("方式3 insert total time:",end.Sub(start).Seconds())
    
    //方式4 insert
    start = time.Now()
    //Begin函数内部会去获取连接
    tx,_ := db.Begin()
    for i := 1301;i<=1311;i++{
        //每次循环用的都是tx内部的连接，没有新建连接，效率高
        tx.Exec("INSERT INTO person(personid,personname) values(?,?)","d"+strconv.Itoa(i),"xnn"+strconv.Itoa(i),i-1000)
    }
    //最后释放tx内部的连接
    tx.Commit()
    
    end = time.Now()
    fmt.Println("方式4 insert total time:",end.Sub(start).Seconds())
    
    //方式5 insert
    start = time.Now()
    for i := 1401;i<=1411;i++{
        //Begin函数每次循环内部都会去连接池获取一个新的连接，效率低下
        tx,_ := db.Begin()
        tx.Exec("INSERT INTO person(personid,personname) values(?,?)","e"+strconv.Itoa(i),"xnn"+strconv.Itoa(i),i-1000)
        //Commit执行后连接也释放了
        tx.Commit()
    }
    end = time.Now()
    fmt.Println("方式5 insert total time:",end.Sub(start).Seconds())
}

func query(){
    
    //方式1 query
    start := time.Now()
    rows,_ := db.Query("select personid,personname from person")
    defer rows.Close()
    for rows.Next(){
         var name string
         var id string
        if err := rows.Scan(&id,&name); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("name:%s ,id:%s \n", name, id)
    }
    end := time.Now()
    fmt.Println("方式1 query total time:",end.Sub(start).Seconds())
    
    //方式2 query
    start = time.Now()
    stm,_ := db.Prepare("select personid,personname from person")
    defer stm.Close()
    rows,_ = stm.Query()
    defer rows.Close()
    for rows.Next(){
         var name string
         var id string
        if err := rows.Scan(&id,&name); err != nil {
            log.Fatal(err)
        }
		fmt.Printf("name:%s ,id:%s \n", name, id)
    }
    end = time.Now()
    fmt.Println("方式2 query total time:",end.Sub(start).Seconds())
    
    
    //方式3 query
    start = time.Now()
    tx,_ := db.Begin()
    defer tx.Commit()
    rows,_ = tx.Query("select personid,personname from person")
    defer rows.Close()
    for rows.Next(){
         var name string
         var id string
        if err := rows.Scan(&id,&name); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("name:%s ,id:%s \n", name, id)
    }
    end = time.Now()
    fmt.Println("方式3 query total time:",end.Sub(start).Seconds())
}