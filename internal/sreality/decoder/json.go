package decoder

import (
	"encoding/json"
	"fmt"

	"github.com/bliuchak/flatsearch/internal/platform/sreality"
)

type JSON struct{}

func (JSON) Decode(data []byte, resp *sreality.ClientResponse) error {
	var input sreality.ClientResponse

	if err := json.Unmarshal(data, &input); err != nil {
		return fmt.Errorf("estate unmarshal: %w", err)
	}

	*resp = input

	return nil
}
