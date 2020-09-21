package call

import (
	"fmt"
)

func ExampleProcessFile() {
	//path := `d:\tmp\bill\20200902\b0000440816.bil`
	//path := `d:\tmp\bill\20200903\b0000440817.bil`
	//path := `d:\tmp\bill\20200903\b0000440885.bil`
	path := `d:\tmp\bill\20200903\b0000440995.bil`
	failCdr, normCdr := ProcessFile(path)
	for i, c := range normCdr {
		fmt.Println(i, "====", c.Conversation_time)
		fmt.Printf("                %d=%08b\r\n", c.Partial_record_indicator, c.Partial_record_indicator)
		// fmt.Println("2=", c.Partial_record_indicator&4)
		// fmt.Println("#=", c.Partial_record_indicator&4)
	}

	fmt.Println(len(failCdr), len(normCdr))
	//Output: 2
}
