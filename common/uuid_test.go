package common

import (
	"testing"
)

func TestUUID4(t *testing.T) {

    t.Run("UUID generation test", func(t *testing.T) {
	    uuid := UUID4()
        if len(uuid) != 36 {
            t.Errorf("Incorrect UUID4: %s", uuid)
        }
    })
}
