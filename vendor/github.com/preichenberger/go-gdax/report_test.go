package gdax

import (
	"errors"
	"testing"
	"time"
)

func TestCreateReportAndStatus(t *testing.T) {
	// # DISABLED in sandbox
	return
	client := NewTestClient()
	newReport := Report{
		Type:      "fill",
		StartDate: time.Now().Add(-24 * 4 * time.Hour),
		EndDate:   time.Now().Add(-24 * 2 * time.Hour),
	}

	report, err := client.CreateReport(&newReport)
	if err != nil {
		t.Error(err)
	}

	currentReport, err := client.GetReportStatus(report.Id)
	if err != nil {
		t.Error(err)
	}
	if StructHasZeroValues(currentReport) {
		t.Error(errors.New("Zero value"))
	}
}
