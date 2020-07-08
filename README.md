[![Go Report Card](https://goreportcard.com/badge/github.com/siodb/siodb-go-driver)](https://goreportcard.com/report/github.com/siodb/siodb-go-driver)

# Go driver for Siodb

A simple driver for Siodb written in pure Go.

## Features

- Support of URI
- Connections to Siodb (TLS, TCP, Unix socket)
- Authentication to Siodb
- Query execution
- DML execution

## Quick start

### Docker

```bash
# Pull the latest container of Siodb
docker run -p 127.0.0.1:50000:50000/tcp --name siodb siodb/siodb
# Get a copy of the private key of the Siodb root user
docker exec -it siodb cat /home/siodb/.ssh/id_rsa > ~/root_id_rsa
```

### Cloud

[![Deploy to Hidora](https://raw.githubusercontent.com/siodb/siodb-jelastic/master/images/deploy-to-hidora.png)](https://siodb.hidora.com)

*Free Trial. Requires only an email address.*

### Driver installation

Get the driver into you Go project:

```bash
go get -u https://github.com/siodb/siodb-go-driver
```

You're ready to Go!

## Example

### Import

```go
package main

import (
    "context"
    "database/sql"
    _ "bitbucket.org/siodb-squad/siodb-go-driver"
)
```

### Siodb Connection

```go
    db, err := sql.Open("siodb", "siodbs://root@localhost:50000?identity_file=/home/nico/root_id_rsa")
    if err != nil {
        t.Fatalf("Error occurred %s", err.Error())
    }
    defer db.Close()

    ctx, stop := context.WithCancel(context.Background())
    defer stop()

    if err := db.PingContext(ctx); err != nil {
        t.Fatalf("Error occurred %s", err.Error())
    }
```

### DDL

```go
    var name string
    err = db.QueryRowContext(ctx, "select name from test.sys_tables where name = 'TABLEALLDATATYPES'").Scan(&name)
    switch {
    case err == sql.ErrNoRows:
        if _, err := db.ExecContext(ctx,
            `CREATE TABLE test.tablealldatatypes
                                    (
                                        ctinyintmin  TINYINT,
                                        ctinyintmax  TINYINT,
                                        ctinyuint    TINYUINT,
                                        csmallintmin SMALLINT,
                                        csmallintmax SMALLINT,
                                        csmalluint   SMALLUINT,
                                        cintmin      INT,
                                        cintmax      INT,
                                        cuint        UINT,
                                        cbigintmin   BIGINT,
                                        cbigintmax   BIGINT,
                                        cbiguint     BIGUINT,
                                        cfloatmin    FLOAT,
                                        cfloatmax    FLOAT,
                                        cdoublemin   DOUBLE,
                                        cdoublemax   DOUBLE,
                                        ctext        TEXT,
                                        cts          TIMESTAMP
                                    )
                                    `); err != nil {
            t.Fatalf("Error occurred: %s", err.Error())
        }
    case err != nil:
        t.Fatalf("Error occurred: %s", err.Error())
    case err == nil:
        break
    default:
        t.Fatalf("Error occurred: %s", err.Error())
    }
```

### DML

```go
    if _, err := db.ExecContext(ctx,
        `INSERT INTO test.tablealldatatypes
                                VALUES  ( -128,
                                          127,
                                          255,
                                          -32768,
                                          32767,
                                          65535,
                                          -2147483648,
                                          2147483647,
                                          4294967295,
                                          -9223372036854775808,
                                          9223372036854775807,
                                          18446744073709551615,
                                          222.222,
                                          222.222,
                                          222.222,
                                          222.222,
                                          '汉字',
                                          CURRENT_TIMESTAMP )`,
    ); err != nil {
        log.Fatal("An Error occurred: ", err)
    }
```

### Query

```go
    rows, err := db.QueryContext(ctx, "SELECT trid, cbigintmin, ctext, cts FROM test.tablealldatatypes")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var Trid uint64
    var AnyValue interface{}
    var StringValue string
    var TiemstampValue time.Time

    for rows.Next() {
        if err := rows.Scan(&Trid, &AnyValue, &StringValue, &TiemstampValue); err != nil {
            log.Fatal(err)
        }
        fmt.Println(fmt.Sprintf("Row: %v | %v | %v | %v ", Trid, AnyValue, StringValue, TiemstampValue))
    }
    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }
```

## URI

To identify a Siodb resource, the driver use the
[URI format](https://en.wikipedia.org/wiki/Uniform_Resource_Identifier).

For TLS connection (default):

```golang
siodbs://root@localhost:50000?identity_file=/home/siodb/.ssh/id_rsa
```

For TCP plain text connection:

```golang
siodb://root@localhost:50000?identity_file=/home/siodb/.ssh/id_rsa
```

For Unix socket connection:

```golang
siodbu://root@/run/siodb/siodb.socket?identity_file=/home/siodb/.ssh/id_rsa
```

The above examples will connect you to the localhost with port number `50000`.
The driver will do the authentication with the Siodb user root and the identity file `/home/siodb/.ssh/id_rsa`.

### Options

- identity_file: the path to the RSA private key.
- trace: to trace everything within the driver to sdtout.

## Support Siodb

Do you like this project? Tell it by clicking the star 🟊 on the top right of this page ☝☝

## Documentation

We write the Siodb documentation in Markdow and it is available in the folder `docs/users/docs`.
If you prefer a more user friendly format, the same documentation is
available online [here]( https://docs.siodb.io).

## Contribution

Please refer to the [Contributing file](CONTRIBUTING.md).

## Support

- Report your issue with Siodb 👉 [here](https://github.com/siodb/siodb/issues/new).
- Report your issue with the driver 👉 [here](https://github.com/siodb/siodb-go-driver/issues/new).
- Ask a question 👉 [here](https://stackoverflow.com/questions/tagged/siodb).
- Siodb Slack space 👉 [here](https://join.slack.com/t/siodb-squad/shared_invite/zt-e766wbf9-IfH9WiGlUpmRYlwCI_28ng).

## Follow Siodb

- [Twitter](https://twitter.com/Sio_db)
- [Linkedin](https://www.linkedin.com/company/siodb)

## License

Licensed under [Apache License version 2.0](https://www.apache.org/licenses/LICENSE-2.0).
