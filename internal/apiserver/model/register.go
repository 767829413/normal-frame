package model

func GetModels() (models []interface{}) {
	models = []interface{}{
		&User{},
	}
	return
}
