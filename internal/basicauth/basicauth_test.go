package basicauth

import (
    "testing"
)

func TestIsKnownCredential(t *testing.T){
	if IsKnownCredential("test", "secret") != true {
		t.Error("Expected true, not false")
	}	
}

