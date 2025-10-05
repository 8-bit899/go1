package events

import "testing"

func TestIsValidDate(t *testing.T) {

	data := "2025-11-16 20:00"
	_, err := IsValidDate(data)
	if err != nil {
		t.Errorf("ожидалась дата 16/11/2025 и время 20:00 получено " + err.Error())
	}

	data = "202511-16 20:00"
	_, err = IsValidDate(data)
	if err == nil {

		t.Errorf("ожидалась ошибка данных даты\n получено " + err.Error())
	}
}
func TestIsValidTitle(t *testing.T) {

	str := "qw"
	err := IsValidTitle(str)
	if err != false {
		t.Errorf("ожидалoсь false")
	}

	str = "sfqw"
	err = IsValidTitle(str)
	if err != true {
		t.Errorf("ожидалoсь true")
	}
	str = "qw!@"
	err = IsValidTitle(str)
	if err != false {
		t.Errorf("ожидалoсь false")
	}

}
