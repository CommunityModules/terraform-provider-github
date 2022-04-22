package main

import (
	"github.com/CommunityModules/terraform-provider-github/git-release"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: git_release.Provider})
}
