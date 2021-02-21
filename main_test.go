package main

import "testing"

func TestCutFront(t *testing.T) {
	initial := "152:menuentry 'Ubuntu, with Linux 5.11.0' --class ubuntu --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-5.11.0-advanced-b70cb823-9505-4ab6-bc0a-ca359515bf03' {"

	expected := "'gnulinux-5.11.0-advanced-b70cb823-9505-4ab6-bc0a-ca359515bf03' {"
	actual, err := cutFront(initial)

	if err != nil {
		t.Errorf("expected error to be nil, got %s", err)
	}

	if actual != expected {
		t.Errorf("expected = %s, actual = %s", expected, actual)
	}

	initial = ""
	_, err = cutFront(initial)

	if err == nil {
		t.Errorf("expected error to be nil, got %s", err)
	}
}

func TestCutRear(t *testing.T) {
	initial := "152:menuentry 'Ubuntu, with Linux 5.11.0' --class ubuntu --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-5.11.0-advanced-b70cb823-9505-4ab6-bc0a-ca359515bf03' {"

	expected := "152:menuentry 'Ubuntu, with Linux 5.11.0' --class ubuntu --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-5.11.0-advanced-b70cb823-9505-4ab6-bc0a-ca359515bf03'"

	actual, err := cutRear(initial)

	if err != nil {
		t.Errorf("got error %s, expected nil", err)
	}

	if actual != expected {
		t.Errorf("expected %s, actual = %s", expected, actual)
	}

	initial = ""
	_, err = cutRear(initial)

	if err == nil {
		t.Errorf("expected error to be nil, got %s", err)
	}
}

func TestProcess(t *testing.T) {
	initial := "152:menuentry 'Ubuntu, with Linux 5.11.0' --class ubuntu --class gnu-linux --class gnu --class os $menuentry_id_option 'gnulinux-5.11.0-advanced-b70cb823-9505-4ab6-bc0a-ca359515bf03' {"

	expected := "'gnulinux-5.11.0-advanced-b70cb823-9505-4ab6-bc0a-ca359515bf03'"
	actual, err := process(initial)

	if err != nil {
		t.Errorf("expected nil error, got %s", err)
	}

	if actual != expected {
		t.Errorf("expected = %s, actual = %s", expected, actual)
	}
}
