// Copyright (C) 2019-2020 Siodb GmbH. All rights reserved.
// Use of this source code is governed by a license that can be found
// in the LICENSE file.

package siodb

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"
)

// TLS connection (default)
var uri string = "siodbs://root@localhost:50000?identity_file=/home/siodb/.ssh/id_rsa"

// Plain text connection
// var uri string = "siodb://root@localhost:50000?identity_file=/home/siodb/.ssh/id_rsa"

// Local Unix socket connection
//var uri string = "siodbu://root@/run/siodb/siodb.socket?identity_file=/home/siodb/.ssh/id_rsa"

type testVars struct {
	db           *sql.DB
	databaseName string
	tableName    string
	tableColsDef string
	ctx          context.Context
	testing      *testing.T
}

func TestDatabase(t *testing.T) {

	db, err := sql.Open("siodb", uri)
	if err != nil {
		t.Fatalf("Error occurred %s", err.Error())
	}
	defer db.Close()

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("Error occurred %s", err.Error())
	}

	var name string
	err = db.QueryRowContext(ctx, "select name from sys_databases where name = 'TEST'").Scan(&name)
	switch {
	case err == sql.ErrNoRows:
		if _, err := db.ExecContext(ctx, "CREATE DATABASE test"); err != nil {
			t.Fatalf("Error occurred %s", err.Error())
		}
	case err != nil:
		t.Fatalf("Error occurred %s", err.Error())
	case err == nil:
		break
	default:
		t.Fatalf("Error occurred %s", err.Error())
	}
}

func TestAllDataTypes(t *testing.T) {

	db, err := sql.Open("siodb", uri)
	if err != nil {
		t.Fatalf("Error occurred %s", err.Error())
	}
	defer db.Close()

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("Error occurred %s", err.Error())
	}

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
		if _, err := db.ExecContext(ctx, "delete from test.TABLEALLDATATYPES"); err != nil {
			t.Fatalf("Error occurred %s", err.Error())
		}
	default:
		t.Fatalf("Error occurred: %s", err.Error())
	}

	// Data type values
	tinyintmin := int8(-128)
	tinyintmax := int8(127)
	tinyuint := uint8(255)
	smallintmin := int16(-32768)
	smallintmax := int16(32767)
	smalluint := uint16(65535)
	intmin := int32(-2147483648)
	intmax := int32(2147483647)
	uint := uint32(4294967295)
	bigintmin := int64(-9223372036854775808)
	bigintmax := int64(9223372036854775807)
	biguint := uint64(18446744073709551615)
	floatmin := float32(222.222)
	floatmax := float32(222.222)
	doublemin := float64(222.222)
	doublemax := float64(222.222)
	text := "汉字"
	ts := time.Now()

	// Returned values
	var retunedTrid uint64
	var retunedTinyintmin int8
	var retunedTinyintmax int8
	var retunedTinyuint uint8
	var retunedSmallintmin int16
	var retunedSmallintmax int16
	var retunedSmalluint uint16
	var retunedIntmin int32
	var retunedIntmax int32
	var retunedUint uint32
	var retunedBigintmin int64
	var retunedBigintmax int64
	var retunedBiguint uint64
	var retunedFloatmin float32
	var retunedFloatmax float32
	var retunedDoublemin float64
	var retunedDoublemax float64
	var retunedStringValue string
	var retunedTimestamp time.Time

	if _, err := db.ExecContext(ctx,
		fmt.Sprintf(
			`INSERT INTO test.tablealldatatypes
	                            VALUES  ( %v,
									%v,
									%v,
									%v,
									%v,
									%v,
									%v,
									%v,
									%v,
									%v,
									%v,
									%v,
									%v,
									%v,
									%v,
									%v,
									'%v',
									'%v' )`,
			tinyintmin,
			tinyintmax,
			tinyuint,
			smallintmin,
			smallintmax,
			smalluint,
			intmin,
			intmax,
			uint,
			bigintmin,
			bigintmax,
			biguint,
			floatmin,
			floatmax,
			doublemin,
			doublemax,
			text,
			ts.Format("2006-01-02 15:04:05.999999999"))); err != nil {
		log.Fatal("An Error occurred: ", err)
	}

	rows, err := db.QueryContext(ctx, "SELECT * FROM test.tablealldatatypes")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&retunedTrid,
			&retunedTinyintmin,
			&retunedTinyintmax,
			&retunedTinyuint,
			&retunedSmallintmin,
			&retunedSmallintmax,
			&retunedSmalluint,
			&retunedIntmin,
			&retunedIntmax,
			&retunedUint,
			&retunedBigintmin,
			&retunedBigintmax,
			&retunedBiguint,
			&retunedFloatmin,
			&retunedFloatmax,
			&retunedDoublemin,
			&retunedDoublemax,
			&retunedStringValue,
			&retunedTimestamp,
		); err != nil {
			log.Fatal(err)
		}
		fmt.Println(fmt.Sprintf("Row %v: %v | %v | %v | %v | %v | %v | %v | %v | %v | %v | %v | %v | %v | %v | %v | %v | %v | %v",
			retunedTrid,
			retunedTinyintmin,
			retunedTinyintmax,
			retunedTinyuint,
			retunedSmallintmin,
			retunedSmallintmax,
			retunedSmalluint,
			retunedIntmin,
			retunedIntmax,
			retunedUint,
			retunedBigintmin,
			retunedBigintmax,
			retunedBiguint,
			retunedFloatmin,
			retunedFloatmax,
			retunedDoublemin,
			retunedDoublemax,
			retunedStringValue,
			retunedTimestamp))
	}

	if retunedTinyintmin != tinyintmin {
		t.Fatalf("ERROR: retunedTinyintmin != tinyintmin => %v != %v", retunedTinyintmin, tinyintmin)
	}
	if retunedTinyintmax != tinyintmax {
		t.Fatalf("ERROR: retunedTinyintmax != tinyintmax => %v != %v", retunedTinyintmax, tinyintmax)
	}
	if retunedTinyuint != tinyuint {
		t.Fatalf("ERROR: retunedTinyuint != tinyuint => %v != %v", retunedTinyuint, tinyuint)
	}
	if retunedSmallintmin != smallintmin {
		t.Fatalf("ERROR: retunedSmallintmin != smallintmin => %v != %v", retunedSmallintmin, smallintmin)
	}
	if retunedSmallintmax != smallintmax {
		t.Fatalf("ERROR: retunedSmallintmax != smallintmax => %v != %v", retunedSmallintmax, smallintmax)
	}
	if retunedSmalluint != smalluint {
		t.Fatalf("ERROR: retunedSmalluint != smalluint => %v != %v", retunedSmalluint, smalluint)
	}
	if retunedIntmin != intmin {
		t.Fatalf("ERROR: retunedIntmin != intmin => %v != %v", retunedIntmin, intmin)
	}
	if retunedIntmax != intmax {
		t.Fatalf("ERROR: retunedIntmax != intmax => %v != %v", retunedIntmax, intmax)
	}
	if retunedUint != uint {
		t.Fatalf("ERROR: retunedUint != uint => %v != %v", retunedUint, uint)
	}
	if retunedBigintmin != bigintmin {
		t.Fatalf("ERROR: retunedBigintmin != bigintmin => %v != %v", retunedBigintmin, bigintmin)
	}
	if retunedBigintmax != bigintmax {
		t.Fatalf("ERROR: retunedBigintmax != bigintmax => %v != %v", retunedBigintmax, bigintmax)
	}
	if retunedBiguint != biguint {
		t.Fatalf("ERROR: retunedBiguint != biguint => %v != %v", retunedBiguint, biguint)
	}
	if retunedFloatmin != floatmin {
		t.Fatalf("ERROR: retunedFloatmin != floatmin => %v != %v", retunedFloatmin, floatmin)
	}
	if retunedFloatmax != floatmax {
		t.Fatalf("ERROR: retunedFloatmax != floatmax => %v != %v", retunedFloatmax, floatmax)
	}
	if retunedDoublemin != doublemin {
		t.Fatalf("ERROR: retunedDoublemin != doublemin => %v != %v", retunedDoublemin, doublemin)
	}
	if retunedDoublemax != doublemax {
		t.Fatalf("ERROR: retunedDoublemax != doublemax => %v != %v", retunedDoublemax, doublemax)
	}
	if retunedStringValue != text {
		t.Fatalf("ERROR: retunedStringValue != text => %v != %v", retunedStringValue, text)
	}
	if retunedTimestamp.Format("2006-01-02 15:04:05.999999999") != ts.Format("2006-01-02 15:04:05.999999999") {
		t.Fatalf("ERROR: retunedTimestamp != ts => %v != %v", retunedTimestamp.Format("2006-01-02 15:04:05.999999999"), ts.Format("2006-01-02 15:04:05.999999999"))
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
