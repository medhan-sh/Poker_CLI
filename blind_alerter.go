package poker

import (
	"fmt"
	"os"
	"time"
)

type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, ammount int)
}
type BlindAlerterfunc func(duration time.Duration, ammount int)

func (a BlindAlerterfunc) ScheduleAlertAt(duration time.Duration, ammount int) {
	a(duration, ammount)
}
func StdOutAlerter(duration time.Duration, ammount int) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(os.Stdout, "Blind is now %v\n", ammount)
	})
}
