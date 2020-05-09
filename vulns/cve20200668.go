package vulns

import (
	"strings"
)

// Check if vulnerable
func (cve *CVE20200668) Check(build string, version string, kbs []string) bool {
	major := strings.Split(build, ".")[2]
	supersedence := []string{}
	switch major {
	case "14393":
		supersedence = append(supersedence, "KB4537764", "KB4550947", "KB4537806", "KB4541329", "KB4540670", "KB4550929")
	case "15063":
		supersedence = append(supersedence, "KB4537789", "KB4550927", "KB4537816", "KB4554342", "KB4540681", "KB4541330")
	case "16299":
		supersedence = append(supersedence, "KB4537789", "KB4550927", "KB4554342", "KB4537816", "KB4541330", "KB4540681")
	case "17134":
		supersedence = append(supersedence, "KB4537762", "KB4541333", "KB4540689", "KB4554349", "KB4550922", "KB4537795", "KB4550944")
	case "17763":
		supersedence = append(supersedence, "KB4532691", "KB4541331", "KB4538461", "KB4554354", "KB4549949", "KB4537818", "KB4550969")
	case "18362":
		supersedence = append(supersedence, "KB4532693", "KB4551762", "KB4540673", "KB4541335", "KB4549951", "KB4550945", "KB4554364", "KB4535996")
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
func (cve *CVE20200668) Name() string {
	return "CVE-2020-0668"
}

// Description returns the Description
func (cve *CVE20200668) Description() string {
	return "Windows Service Tracing: Arbitrary File Write as SYSTEM"
}

// CVE20200668 is the cve type
type CVE20200668 struct {
	name string
}
