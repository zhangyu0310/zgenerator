# zgenerator
Generate Golang code according to the database table structure.

## Quick start

```shell
./zgenerator -dsn="root:123456@tcp(127.0.0.1:3306)/TestDB?charset=utf8mb4" -table-name="TestTable" -output="./example/MyTest.go" -first-to-upper=true
```

Table create SQL:
```sql
CREATE Table TestTable (ID BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT 'UserID', Name VARCHAR(30) COMMENT 'User Name', Money BIGINT COMMENT 'Money');
```

You can get Golang file like this: [MyTest.go](https://github.com/zhangyu0310/zgenerator/blob/main/example/MyTest.go)

## Usage
```shell
Usage of zgenerator:
  -charset string
        Charset of db connection (default "utf8mb4")
  -code-after-struct string
        Customize code after struct
  -code-before-struct string
        Customize code before struct
  -code-in-struct string
        Customize code in struct
  -date-string
        Use string store date type (default true)
  -db string
        Database name
  -db-driver string
        Driver of database (default "mysql")
  -dsn string
        Data Source Name
  -first-to-upper
        Convert first letter to upper
  -generate-func
        Generate some easy function to use (Depend sqlx) (default true)
  -help
        Show usage
  -host string
        Host of database (default "127.0.0.1")
  -json-tag
        Add json tag
  -output string
        Path of output file
  -package-name string
        package name of generate code (default "handler")
  -password string
        Password of user
  -post uint
        Port of database (default 3306)
  -select-key string
        Key of select SQL (Default: PrimaryKey)
  -struct-name string
        Struct name of generate code (Default: same with table-name)
  -table-name string
        Table name in database
  -user string
        User of database (default "root")
  -version
        Version info
```