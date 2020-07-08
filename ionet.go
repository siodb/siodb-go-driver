// Copyright (C) 2019-2020 Siodb GmbH. All rights reserved.
// Use of this source code is governed by a license that can be found
// in the LICENSE file.

package siodb

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"io"
	"time"
	"unicode/utf8"

	"github.com/golang/protobuf/proto"
)

func (sc *siodbConn) debug(message string, args ...interface{}) {
	if sc.cfg.trace {
		fmt.Printf("## SIODB DRIVER | "+message+"\n", args...)
	}
}

func (sc *siodbConn) cleanupBuffer() (cpt int16, err error) {

	var rowLength uint64 = 0

	sc.debug("cleanupBuffer | starts.")

	for {
		// Get Current Row Size
		if _, rowLength, err = sc.readVarint(); err != nil {
			return 0, &siodbDriverError{"Unable to read the row size."}
		}
		if rowLength == 0 {
			sc.debug("cleanupBuffer | Dropped %d rows so far.", cpt)
			sc.completed = true
			return cpt, err
		}
		sc.debug("cleanupBuffer | Row size detected: %d.", rowLength)
		buff := make([]byte, rowLength)
		_, err = io.ReadFull(sc.netConn, buff)

		cpt++
	}

}
func (sc *siodbConn) writeServerCommand(sqlText string) error {

	sc.RequestID = 1

	command := &Command{
		RequestID: sc.RequestID,
		Text:      sqlText,
	}

	if len(command.Text) < 300 {
		sc.debug("writeServerCommand | Message to send: %v", command)
	} else {
		sc.debug("writeServerCommand | Message too big to dump.")
	}

	var buf [binary.MaxVarintLen32]byte
	encodedLength := binary.PutUvarint(buf[:], uint64(1))
	sc.netConn.Write(buf[:encodedLength])
	writeMessage(sc.netConn, command)

	return nil
}

func (sc *siodbConn) readServer() (serverResponse ServerResponse, err error) {

	// Get Message
	if _, err = sc.ReadMessage(2, &serverResponse); err != nil {
		return serverResponse, err
	}

	sc.debug("Raw Proto Message: %v", serverResponse)

	// Check request ID
	sc.debug("readServer | Request Id: %d.", serverResponse.RequestID)
	if serverResponse.RequestID != sc.RequestID {
		return serverResponse, &siodbDriverError{"Wrong request ID in the server response."}
	}

	// Check dataset presence
	var columnCount int = len(serverResponse.ColumnDescription)
	if columnCount == 0 {

		sc.debug("readServer | No dataset in response (columnCount=%d).", columnCount)
		sc.completed = true

	} else {

		sc.debug("readServer | Number of Columns: %d.", columnCount)
		sc.completed = false

		// Check if one column can be null meaning that the stream contains the nullbitmask.
		for _, column := range serverResponse.ColumnDescription {
			sc.debug("readServer | Column %s of type %s (can be bull? => %t).", column.Name, column.Type, column.IsNull)
			if column.IsNull == true {
				sc.nullAllowed = true
			}
		}
		sc.debug("readServer | Null columns possible?: %t.", sc.nullAllowed)

		// Derive null Bitmask size if one column can be null.
		if sc.nullAllowed == true {
			if columnCount%8 == 0 {
				sc.nullBitmaskByteSize = columnCount / 8
			} else {
				sc.nullBitmaskByteSize = columnCount/8 + 1
			}
			sc.debug("readRow | Null Bitmask size in bytes: %d.", sc.nullBitmaskByteSize)
		}

	}

	return serverResponse, nil
}

func (sc *siodbConn) readRow(dest []driver.Value, columnDesc []*ColumnDescription) error {

	var rowLength uint64 = 0
	var err error

	// Get Current Row Size
	if _, rowLength, err = sc.readVarint(); err != nil {
		return &siodbDriverError{"Unable to read the row size."}
	}

	sc.debug("readRow | --------------------------------------------------")
	sc.debug("readRow | Row data length: %d.", rowLength)

	if rowLength == 0 {
		sc.debug("readRow | Last row reached; break.")
		sc.completed = true
		return io.EOF
	}

	// Read null Bitmask to figure out null value which are not streamed.
	var Bitmask []byte
	if sc.nullAllowed == true {
		Bitmask = make([]byte, sc.nullBitmaskByteSize)
		if _, err = io.ReadFull(sc.netConn, Bitmask); err != nil {
			return &siodbDriverError{"Fail to read the bitmask byte(s)."}
		}
		sc.debug("readRow | Bitmask value : %08b.", Bitmask)
	}

	// Read Row data
	for idx, column := range columnDesc {

		// Get Bit value for current column
		var IsNull byte
		if sc.nullAllowed == true {
			mask := byte(1 << (idx % 8))
			sc.debug("readRow | Bit %d with mask value: %08b", (idx % 8), mask)
			IsNull = (Bitmask[idx/8] & mask) >> (idx % 8)
			sc.debug("readRow | Bitmask in pos: %d. Is this cell null? => %d", Bitmask[idx/8], IsNull)
		}

		if IsNull == byte(0) { // If not null
			if dest[idx], err = sc.readFieldData(column.Type); err != nil {
				return &siodbDriverError{"Fail to read field " + column.Name + " from current row | " + err.Error()}
			}
		} else { // if null
			dest[idx] = nil
			sc.debug("readRow | NULL Value.")
		}

	}

	return nil
}

func (sc *siodbConn) readFieldData(ColumnType ColumnDataType) (dest driver.Value, err error) {

	sc.debug("readFieldData | Type detected: %s.", ColumnType)
	var bytesRead int

	switch ColumnType {

	case ColumnDataType_COLUMN_DATA_TYPE_BOOL:

		var Value bool
		err := binary.Read(sc.netConn, binary.LittleEndian, &Value)
		sc.debug("readFieldData   |--> bytesRead: %d | err: %d | Value: %t.", bytesRead, err, Value)
		return Value, err

	case ColumnDataType_COLUMN_DATA_TYPE_INT8:

		var Value int8
		err := binary.Read(sc.netConn, binary.LittleEndian, &Value)
		sc.debug("readFieldData   |--> bytesRead: %d | err: %d | Value: %d.", bytesRead, err, Value)
		return Value, err

	case ColumnDataType_COLUMN_DATA_TYPE_UINT8:

		var Value uint8
		err := binary.Read(sc.netConn, binary.LittleEndian, &Value)
		sc.debug("readFieldData   |--> bytesRead: %d | err: %d | Value: %d.", bytesRead, err, Value)
		return Value, err

	case ColumnDataType_COLUMN_DATA_TYPE_INT16:

		var Value int16
		err := binary.Read(sc.netConn, binary.LittleEndian, &Value)
		sc.debug("readFieldData   |--> bytesRead: %d | err: %d | Value: %d.", bytesRead, err, Value)
		return Value, err

	case ColumnDataType_COLUMN_DATA_TYPE_UINT16:

		var Value uint16
		err := binary.Read(sc.netConn, binary.LittleEndian, &Value)
		sc.debug("readFieldData   |--> bytesRead: %d | err: %d | Value: %d.", bytesRead, err, Value)
		return Value, err

	case ColumnDataType_COLUMN_DATA_TYPE_INT32:

		bytesRead, Value, err := sc.readVarint()
		sc.debug("readFieldData   |--> bytesRead: %d | err: %d | Value: %d.", bytesRead, err, int32(Value))
		return int32(Value), err

	case ColumnDataType_COLUMN_DATA_TYPE_UINT32:

		bytesRead, Value, err := sc.readVarint()
		sc.debug("readFieldData   |--> bytesRead: %d | err: %d | Value: %d.", bytesRead, err, uint32(Value))
		return uint32(Value), err

	case ColumnDataType_COLUMN_DATA_TYPE_INT64:

		bytesRead, Value, err := sc.readVarint()
		sc.debug("readFieldData   |--> bytesRead: %d | err: %d | Value: %d.", bytesRead, err, int64(Value))
		return int64(Value), err

	case ColumnDataType_COLUMN_DATA_TYPE_UINT64:

		bytesRead, Value, err := sc.readVarint()
		sc.debug("readFieldData   |--> bytesRead: %d | err: %d | Value: %d.", bytesRead, err, uint64(Value))
		return uint64(Value), err

	case ColumnDataType_COLUMN_DATA_TYPE_FLOAT:

		var Value float32
		err := binary.Read(sc.netConn, binary.LittleEndian, &Value)
		sc.debug("readFieldData   |--> bytesRead: %d | err: %d | Value: %d.", bytesRead, err, float32(Value))
		return float32(Value), err

	case ColumnDataType_COLUMN_DATA_TYPE_DOUBLE:

		var Value float64
		err := binary.Read(sc.netConn, binary.LittleEndian, &Value)
		sc.debug("readFieldData   |--> bytesRead: %d | err: %d | Value: %d.", bytesRead, err, float64(Value))
		return float64(Value), err

	case ColumnDataType_COLUMN_DATA_TYPE_TEXT:

		bytesRead, TextLength, err := sc.readVarint()
		sc.debug("readFieldData   | bytesRead: %d | err: %d | TextLength: %d.", bytesRead, err, TextLength)
		var Value string
		buff := make([]byte, TextLength)
		bytesRead, err = io.ReadFull(sc.netConn, buff)
		for len(buff) > 0 {
			r, size := utf8.DecodeRune(buff)
			Value = Value + fmt.Sprintf("%c", r)

			buff = buff[size:]
		}
		sc.debug("readFieldData   |--> bytesRead: %d | err: %d | Value: %s.", bytesRead, err, Value)
		return Value, err

	case ColumnDataType_COLUMN_DATA_TYPE_NTEXT:

		sc.debug("Data type '%q' is not supported yet.", ColumnType)
		return nil, &siodbDriverError{"Data type '" + ColumnType.String() + "' not supported yet."}

	case ColumnDataType_COLUMN_DATA_TYPE_BINARY:

		bytesRead, blobLength, err := sc.readVarint()
		sc.debug("readFieldData   | bytesRead: %d | err: %d | blobLength: %d.", bytesRead, err, blobLength)
		Value := make([]byte, blobLength)
		bytesRead, err = io.ReadFull(sc.netConn, Value)
		sc.debug(" |--> Value [disabled for BLOB]")
		return Value, err

	case ColumnDataType_COLUMN_DATA_TYPE_DATE:
		// TODO: implement type
		sc.debug("Data type '%q' is not supported yet.", ColumnType)
		return nil, &siodbDriverError{"Data type '" + ColumnType.String() + "' not supported yet."}
	case ColumnDataType_COLUMN_DATA_TYPE_TIME:
		// TODO: implement type
		sc.debug("Data type '%q' is not supported yet.", ColumnType)
		return nil, &siodbDriverError{"Data type '" + ColumnType.String() + "' not supported yet."}
	case ColumnDataType_COLUMN_DATA_TYPE_TIME_WITH_TZ:
		// TODO: implement type
		sc.debug("Data type '%q' is not supported yet.", ColumnType)
		return nil, &siodbDriverError{"Data type '" + ColumnType.String() + "' not supported yet."}

	case ColumnDataType_COLUMN_DATA_TYPE_TIMESTAMP:

		var Value time.Time

		// Get date part, 4 first bytes
		buff := make([]byte, 4)
		bytesRead, err = io.ReadFull(sc.netConn, buff)
		hasTimePart := int(buff[0] & byte(1))
		sc.debug(" |--> hasTimePart     : %t", hasTimePart)
		dayOfWeek := int(buff[0] & byte(14) >> 1)
		sc.debug(" |--> dayOfWeek       : %d", dayOfWeek)
		dayOfMonth := int(((buff[0] & byte(240) >> 4) + (buff[1] & byte(1) << 4)) + 1)
		sc.debug(" |--> dayOfMonth      : %d", dayOfMonth)
		month := int((buff[1] & byte(30) >> 1) + 1)
		sc.debug(" |--> month           : %d", month)
		sliceYear := []byte{
			byte(0),
			byte((buff[3] & byte(224) >> 5)),
			byte((buff[2] & byte(224) >> 5) + (buff[3] & byte(31) << 3)),
			byte((buff[1] & byte(224) >> 5) + (buff[2] & byte(31) << 3)),
		}
		year := int(binary.BigEndian.Uint32(sliceYear[:]))
		sc.debug(" |--> year            : %d", int(year))

		// Get time part if any, 6 next bytes
		if hasTimePart == 1 {
			buff := make([]byte, 6)
			bytesRead, err = io.ReadFull(sc.netConn, buff)

			reserved1 := buff[0] & byte(1)
			sc.debug(" |--> reserved1       : %d", reserved1)

			sliceNanos := []byte{
				byte((buff[3] & byte(126) >> 1)),
				byte((buff[2] & byte(254) >> 1) + (buff[3] & byte(1) << 7)),
				byte((buff[1] & byte(254) >> 1) + (buff[2] & byte(1) << 7)),
				byte((buff[0] & byte(254) >> 1) + (buff[1] & byte(1) << 7)),
			}
			nanos := int(binary.BigEndian.Uint32(sliceNanos[:]))
			sc.debug(" |--> nanos           : %d", nanos)

			seconds := int((buff[3] & byte(128) >> 7) + (buff[4] & byte(31) << 1))
			sc.debug(" |--> seconds         : %d", seconds)

			minutes := int((buff[4] & byte(224) >> 5) + (buff[5] & byte(7) << 3))
			sc.debug(" |--> minutes         : %d", minutes)

			hours := int((buff[5] & byte(248) >> 3))
			sc.debug(" |--> hours           : %d", hours)

			Value = time.Date(year, time.Month(month), dayOfMonth, hours, minutes, seconds, nanos, time.Local)

		} else {

			Value = time.Date(year, time.Month(month), dayOfMonth, 0, 0, 0, 0, time.Local)

		}

		return Value, err

	case ColumnDataType_COLUMN_DATA_TYPE_TIMESTAMP_WITH_TZ:
		// TODO: implement type
		sc.debug("Data type '%q' is not supported yet.", ColumnType)
		return nil, &siodbDriverError{"Data type '" + ColumnType.String() + "' not supported yet."}
	case ColumnDataType_COLUMN_DATA_TYPE_DATE_INTERVAL:
		// TODO: implement type
		sc.debug("Data type '%q' is not supported yet.", ColumnType)
		return nil, &siodbDriverError{"Data type '" + ColumnType.String() + "' not supported yet."}
	case ColumnDataType_COLUMN_DATA_TYPE_TIME_INTERVAL:
		// TODO: implement type
		sc.debug("Data type '%q' is not supported yet.", ColumnType)
		return nil, &siodbDriverError{"Data type '" + ColumnType.String() + "' not supported yet."}
	case ColumnDataType_COLUMN_DATA_TYPE_STRUCT:
		// TODO: implement type
		sc.debug("Data type '%q' is not supported yet.", ColumnType)
		return nil, &siodbDriverError{"Data type '" + ColumnType.String() + "' not supported yet."}
	case ColumnDataType_COLUMN_DATA_TYPE_XML:
		// TODO: implement type
		sc.debug("Data type '%q' is not supported yet.", ColumnType)
		return nil, &siodbDriverError{"Data type '" + ColumnType.String() + "' not supported yet."}
	case ColumnDataType_COLUMN_DATA_TYPE_JSON:
		// TODO: implement type
		sc.debug("Data type '%q' is not supported yet.", ColumnType)
		return nil, &siodbDriverError{"Data type '" + ColumnType.String() + "' not supported yet."}
	case ColumnDataType_COLUMN_DATA_TYPE_UUID:
		// TODO: implement type
		sc.debug("Data type '%q' is not supported yet.", ColumnType)
		return nil, &siodbDriverError{"Data type '" + ColumnType.String() + "' not supported yet."}
	case ColumnDataType_COLUMN_DATA_TYPE_MAX:
		// TODO: implement type
		sc.debug("Data type '%q' is not supported yet.", ColumnType)
		return nil, &siodbDriverError{"Data type '" + ColumnType.String() + "' not supported yet."}
	case ColumnDataType_COLUMN_DATA_TYPE_UNKNOWN:
		// TODO: implement type
		sc.debug("Data type '%q' is not supported yet.", ColumnType)
		return nil, &siodbDriverError{"Data type '" + ColumnType.String() + "' not supported yet."}
	default:

		sc.debug("Data type '%q' unknown.", ColumnType)
		return nil, &siodbDriverError{"Unknown data type."}

	}
}

func (sc *siodbConn) readVarint() (bytesRead int, n uint64, err error) {

	// Function readVarint()
	// source: https://github.com/stashed/stash/blob/master/vendor/github.com/matttproud/golang_protobuf_extensions/pbutil/decode.go
	//
	// Copyright 2013 Matt T. Proud
	//
	// Licensed under the Apache License, Version 2.0 (the "License");
	// you may not use this file except in compliance with the License.
	// You may obtain a copy of the License at
	//
	//     http://www.apache.org/licenses/LICENSE-2.0
	//
	// Unless required by applicable law or agreed to in writing, software
	// distributed under the License is distributed on an "AS IS" BASIS,
	// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	// See the License for the specific language governing permissions and
	// limitations under the License.

	var prefixBuf [binary.MaxVarintLen64]byte
	var varIntBytes int
	for varIntBytes == 0 { // i.e. no varint has been decoded yet.
		if bytesRead >= len(prefixBuf) {
			return bytesRead, n, fmt.Errorf("invalid varint64 encountered")
		}
		// We have to read byte by byte here to avoid reading more bytes
		// than required. Each read byte is appended to what we have
		// read before.
		newBytesRead, err := sc.netConn.Read(prefixBuf[bytesRead : bytesRead+1])
		if newBytesRead == 0 {
			if io.EOF == err {
				return bytesRead, n, nil
			} else if err != nil {
				return bytesRead, n, err
			}
			// A Reader should not return (0, nil), but if it does,
			// it should be treated as no-op (according to the
			// Reader contract). So let's go on...
			continue
		}
		bytesRead += newBytesRead
		// Now present everything read so far to the varint decoder and
		// see if a varint can be decoded already.
		n, varIntBytes = proto.DecodeVarint(prefixBuf[:bytesRead])
	}

	return bytesRead, n, err
}

func writeMessage(w io.Writer, m proto.Message) (int, error) {

	// Function ReadMessage()
	// source: https://github.com/stashed/stash/blob/master/vendor/github.com/matttproud/golang_protobuf_extensions/pbutil/encode.go
	//
	// Copyright 2013 Matt T. Proud
	//
	// Licensed under the Apache License, Version 2.0 (the "License");
	// you may not use this file except in compliance with the License.
	// You may obtain a copy of the License at
	//
	//     http://www.apache.org/licenses/LICENSE-2.0
	//
	// Unless required by applicable law or agreed to in writing, software
	// distributed under the License is distributed on an "AS IS" BASIS,
	// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	// See the License for the specific language governing permissions and
	// limitations under the License.

	em, err := proto.Marshal(m)
	if nil != err {
		return 0, err
	}

	var buf [binary.MaxVarintLen32]byte
	encodedLength := binary.PutUvarint(buf[:], uint64(proto.Size(m)))

	vib, err := w.Write(buf[:encodedLength])
	if nil != err {
		return vib, err
	}

	pbmmb, err := w.Write(em)

	return vib + pbmmb, err
}

func (sc *siodbConn) ReadMessage(messageTypeID uint64, m proto.Message) (n int, err error) {

	// Function ReadMessage()
	// source: https://github.com/stashed/stash/blob/master/vendor/github.com/matttproud/golang_protobuf_extensions/pbutil/decode.go
	//
	// Copyright 2013 Matt T. Proud
	//
	// Licensed under the Apache License, Version 2.0 (the "License");
	// you may not use this file except in compliance with the License.
	// You may obtain a copy of the License at
	//
	//     http://www.apache.org/licenses/LICENSE-2.0
	//
	// Unless required by applicable law or agreed to in writing, software
	// distributed under the License is distributed on an "AS IS" BASIS,
	// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	// See the License for the specific language governing permissions and
	// limitations under the License.

	var prefixBuf [binary.MaxVarintLen64]byte
	var bytesRead, varIntBytes int
	var messageLength uint64
	var readMessageTypeID uint64

	// Read and check Message Type Id
	_, readMessageTypeID, err = sc.readVarint()
	if messageTypeID != readMessageTypeID {
		return 0, &siodbDriverError{"Wrong message type id."}
	}
	sc.debug("readServerMessage | Message Type Id: %d.", readMessageTypeID)

	// Read Message
	for varIntBytes == 0 { // i.e. no varint has been decoded yet.
		if bytesRead >= len(prefixBuf) {
			return 0, fmt.Errorf("invalid varint64 encountered")
		}
		// We have to read byte by byte here to avoid reading more bytes
		// than required. Each read byte is appended to what we have
		// read before.
		newBytesRead, err := sc.netConn.Read(prefixBuf[bytesRead : bytesRead+1])
		if newBytesRead == 0 {
			if io.EOF == err {
				return 0, nil
			} else if err != nil {
				return 0, err
			}
			// A Reader should not return (0, nil), but if it does,
			// it should be treated as no-op (according to the
			// Reader contract). So let's go on...
			continue
		}
		bytesRead += newBytesRead
		// Now present everything read so far to the varint decoder and
		// see if a varint can be decoded already.
		messageLength, varIntBytes = proto.DecodeVarint(prefixBuf[:bytesRead])
	}

	sc.debug("readServerMessage | %d.", messageLength)
	messageBuf := make([]byte, messageLength)
	newBytesRead, err := io.ReadFull(sc.netConn, messageBuf)
	bytesRead += newBytesRead
	if err != nil {
		return 0, err
	}

	return bytesRead, proto.Unmarshal(messageBuf, m)
}
