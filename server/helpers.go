package main

func (app *application) getInt32(value int32, defaultValue int32) int32 {
	if value <= 0 {
		return defaultValue
	}
	return value
}

func (app *application) getString(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
