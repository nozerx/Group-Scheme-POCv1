package group

import (
	"fmt"
)

func (gk *GroupKeyShare) GenerateGroupKey() string {
	return fmt.Sprintf("%s/%s/%s", gk.GroupName, gk.Host.Pretty(), gk.Key)
}
