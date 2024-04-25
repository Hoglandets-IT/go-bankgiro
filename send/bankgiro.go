package to_bg

// import (
// 	"fmt"
// 	"time"
// )

// type BankgiroFile struct {
// 	SealDate string
// 	Content  string
// 	Sealer   HmacSealer
// }

// func (bgf *BankgiroFile) SetDateToday() {
// 	bgf.SealDate = time.Now().Format("060102")
// }

// func (bgf *BankgiroFile) SetDate(date string) error {
// 	if len(date) != 6 {
// 		return fmt.Errorf("set date failed, got date of length %d, expected 6", len(date))
// 	}

// 	bgf.SealDate = date
// 	return nil
// }

// // func (bg)
