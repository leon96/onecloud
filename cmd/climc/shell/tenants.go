package shell

import (
	"github.com/yunionio/onecloud/pkg/mcclient"
	"github.com/yunionio/onecloud/pkg/mcclient/modules"
)

func init() {
	type TenantListOptions struct {
	}
	R(&TenantListOptions{}, "tenant-list", "List tenants", func(s *mcclient.ClientSession, args *TenantListOptions) error {
		result, err := modules.Tenants.List(s, nil)
		if err != nil {
			return err
		}
		printList(result, modules.Tenants.GetColumns(s))
		return nil
	})
}
