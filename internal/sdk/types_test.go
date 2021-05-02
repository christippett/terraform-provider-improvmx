package improvmx

import (
	"testing"
)

func TestRecordValues_ValueSlice(t *testing.T) {
	var values RecordValues
	json := []byte(`["val1", "val2"]`)
	if err := values.UnmarshalJSON(json); err != nil {
		t.Errorf("Error processing value slice: %v", err)
	}
	if len(values) != 2 {
		t.Errorf("RecordValues.UnmarshalJSON() returned invalid value: %s", values)
	}
}

func TestRecordValues_ReturnSlice(t *testing.T) {
	var values RecordValues
	json := []byte(`"val1"`)
	if err := values.UnmarshalJSON(json); err != nil {
		t.Errorf("Error processing value slice: %v", err)
	}
	if len(values) != 1 || values[0] != "val1" {
		t.Errorf("RecordValues.UnmarshalJSON() returned invalid value: %s", values)
	}
}
