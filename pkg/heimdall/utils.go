package heimdall

import (
	"fmt"
	"github.com/davidbanham/human_duration"
	"time"
)

func CalculateUptime(started string, finished string) string {
	startTimeStr, startTimeError := time.Parse(time.RFC3339, started)
	endTimeStr, endTimeError := time.Parse(time.RFC3339, finished)
	if startTimeError != nil || endTimeError != nil {
		Fatal(fmt.Sprintf("Errored parsing time: %s; %s", startTimeError, endTimeError))
	}
	return human_duration.String(endTimeStr.Sub(startTimeStr), human_duration.Second)
}
