/*
0×00   0       No errors occurred, and no copying was done.
               The source and destination directory trees are completely synchronized.

0×01   1       One or more files were copied successfully (that is, new files have arrived).

0×02   2       Some Extra files or directories were detected. No files were copied
               Examine the output log for details.

0×04   4       Some Mismatched files or directories were detected.
               Examine the output log. Housekeeping might be required.

0×08   8       Some files or directories could not be copied
               (copy errors occurred and the retry limit was exceeded).
               Check these errors further.

0×10  16       Serious error. Robocopy did not copy any files.
               Either a usage error or an error due to insufficient access privileges
               on the source or destination directories.
*/

package rcopy

import (
	"fmt"
	"os/exec"
)

/*
robocopy \\1.1.1.1\d$\frontsave\X3KF\ d:\tmp\bill /MIR /FFT /maxage:1 /min:524287
*/
func Execute(pathBatFile string) string {
	//out, err := exec.Command("cmd.exe", "/C", "d:\\gocode\\src\\sxbill\\cmd\\rcopy\\rcopy.bat").CombinedOutput()
	out, err := exec.Command("cmd.exe", "/C", pathBatFile).CombinedOutput()
	return fmt.Sprintf("%s, err=%s\r\n", out, err)
}

// func ExecutePeriodically(periodMin int, pathBatFile string) {
// 	ticker := time.NewTicker(time.Minute * time.Duration(periodMin))
// 	for {
// 		<-ticker.C
// 		fmt.Println(Execute(pathBatFile))
// 	}
// }
