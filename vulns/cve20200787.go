package vulns

import (
	"strings"
)

// Check if vulnerable
func (cve *CVE20200787) Check(build string, version string, kbs []string) bool {
	major := strings.Split(build, ".")[2]
	supersedence := []string{}
	switch major {
	case "14393":
		supersedence = append(supersedence, "KB4540670", "KB4541329", "KB4550947", "KB4550929")
	case "15063":
		// missing on msrc
	case "16299":
		supersedence = append(supersedence, "KB4540681", "KB4550927", "KB4541330", "KB4554342")
	case "17134":
		supersedence = append(supersedence, "KB4540689", "KB4541333", "KB4550944", "KB4550922", "KB4554349")
	case "17763":
		supersedence = append(supersedence, "KB4538461", "KB4549949", "KB4550969", "KB4541331", "KB4554354")
	case "18362", "18363":
		supersedence = append(supersedence, "KB4540673", "KB4551762", "KB4550945", "KB4541335", "KB4549951", "KB4554364")
	default:
		return false
	}
	for _, i := range kbs {
		for _, j := range supersedence {
			if i == j {
				return false
			}
		}
	}
	return true
}

// Name returns the CVE Number
func (cve *CVE20200787) Name() string {
	return "CVE-2020-0787"
}

// Description returns the Description
func (cve *CVE20200787) Description() string {
	return "Windows BITS: Arbitrary File Move as SYSTEM"
}

// CVE20200787 is the cve type
type CVE20200787 struct {
	name string
}
