package models

import (
	"strings"
)

// TrySave is wrapper function for pop.Connection#Save()
func TrySave(model interface{}) error {
	if err := DB.Eager().Save(model); err != nil {
		if strings.Contains(err.Error(), "Duplicate") { // mysql error 1062
			mlogger.Warnf("creation failed due to duplicated entry. ignore")
		} else {
			mlogger.Errorf("could not save new entry %v: %v", model, err)
			return err
		}
	}
	return nil
}
