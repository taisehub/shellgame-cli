package shellgame

import (
	"github.com/taise-hub/shellgame-cli/common"
)

var myProfile *common.Profile

func SetMyProfile(profile *common.Profile) {
	myProfile = profile
} 

func GetMyProfile() *common.Profile {
	if myProfile == nil {
		return nil
	}
	return myProfile
}