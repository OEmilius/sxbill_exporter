/*

 */

package call

import (
	"encoding/binary"
	"fmt"

	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type Call struct {
	//net_type                 byte
	Bill_type                uint16
	Partial_record_indicator uint16
	Ans_time                 time.Time
	End_time                 time.Time
	Conversation_time        uint32 // unit is 10 ms
	Caller_number            string
	Called_number            string
	Trunk_group_in           uint16
	Trunk_group_out          uint16
	Termination_code         uint16
	//Local_csn                uint16
}

func FindNumber(s string) string {
	index_f := strings.Index(s, "f")
	if index_f > 0 {
		return s[:index_f]
	}
	return "fffffffffff"
}

func Decode(b []byte) (Call, error) {
	o := -5 //offset
	//fmt.Printf("Decode= %x \r\n", b)
	c := Call{}
	c.Bill_type = uint16(b[o+7])
	//fmt.Printf("Bill_type %x\r\n", c.Bill_type)
	c.Partial_record_indicator = uint16(b[o+9])
	var err error
	c.Ans_time, err = DecodeTime(b[o+11 : o+11+6])
	if err != nil {
		return c, err
	}
	c.End_time, err = DecodeTime(b[o+17 : o+17+6])
	if err != nil {
		return c, err
	}
	c.Conversation_time = binary.LittleEndian.Uint32(b[o+23 : o+23+4])
	c.Caller_number = FindNumber(fmt.Sprintf("%x", b[o+30:o+30+16]))
	c.Called_number = FindNumber(fmt.Sprintf("%x", b[o+49:o+49+16]))
	c.Trunk_group_in = binary.LittleEndian.Uint16(b[o+77 : o+77+2])
	c.Trunk_group_out = binary.LittleEndian.Uint16(b[o+79 : o+79+2])
	c.Termination_code = uint16(b[o+87])
	//fmt.Printf("%+v\r\n", c)
	return c, nil
}

func ProcessFile(path string) (failCdr, normCdr []Call) {
	log.Println("start process file ", path)
	buf, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = buf.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	d, err := ioutil.ReadAll(buf)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	for i := 0; i < len(d)-200; i++ {
		if d[i] == byte(0) {
			if d[i+1] == byte(11) && d[i+2] == byte(85) {
				// log.Println("failed call ticket")
				// log.Printf("%x\r\n", d[i:i+235])
				c, err := Decode(d[i : i+235])
				if err != nil {
					fmt.Println(err)
					continue
				}
				failCdr = append(failCdr, c)
				//log.Printf("%+v\r\n--", c)
				i = i + 234
				continue
			}
			if d[i+1] == byte(11) && d[i+2] == byte(1) {
				//log.Println("detailed ticket")
				//log.Printf("%x\r\n", d[i:i+235])
				c, err := Decode(d[i : i+235])
				if err != nil {
					fmt.Println(err)
					continue
				}
				normCdr = append(normCdr, c)
				//log.Printf("%+v\r\n--", c)
				i = i + 234
			}
		}
	}
	log.Println("end call ProcessFile", path)
	return failCdr, normCdr
}

func DecodeTime(b []byte) (time.Time, error) {
	YYMMDDHHMMSS := ""
	for i := 0; i < 6; i++ {
		YYMMDDHHMMSS = YYMMDDHHMMSS + fmt.Sprintf("%02d", int(b[i]))
	}
	const layout = "060102150405"
	t, err := time.Parse(layout, YYMMDDHHMMSS)
	if err != nil {
		fmt.Println(err)
		return time.Time{}, err
	}
	return t, nil
}
