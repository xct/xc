package vulns

import (
	"strings"
)

// https://portal.msrc.microsoft.com/en-US/security-guidance/advisory/CVE-2019-1315
// This one is a commented example so you can do create your own

// Check if vulnerable
func (cve *CVE20191315) Check(build string, version string, kbs []string) bool {
	major := strings.Split(build, ".")[2]
	// supersedence, lowest patch + all patches that replaced it over the time
	supersedence := []string{}
	switch major {
	case "10240":
		supersedence = append(supersedence, "KB4520011", "KB4540693", "KB4550930", "KB4525232", "KB4530681", "KB4537776", "KB4534306")
	case "10586":
		supersedence = append(supersedence, "KB4520011", "KB4534306", "KB4550930", "KB4537776", "KB4530681", "KB4525232", "KB4540693")
	case "14393":
		supersedence = append(supersedence, "KB4519998", "KB4519979", "KB4541329", "KB4540670", "KB4550929", "KB4525236", "KB4534271", "KB4537806", "KB4550947", "KB4530689", "KB4537764", "KB4534307")
	case "15063":
		supersedence = append(supersedence, "KB4520010")
	case "16299":
		supersedence = append(supersedence, "KB4520004", "KB4554342", "KB4540681", "KB4550927", "KB4520006", "KB4541330", "KB4534318", "KB4537816", "KB4537789", "KB4525241", "KB4534276", "KB4530714")
	case "17134":
		supersedence = append(supersedence, "KB4520008", "KB4530717", "KB4534308", "KB4541333", "KB4550944", "KB4534293", "KB4537795", "KB4554349", "KB4550922", "KB4540689", "KB4537762", "KB4519978", "KB4525237")
	case "17763":
		supersedence = append(supersedence, "KB4519338", "KB4523205", "KB4538461", "KB4541331", "KB4549949", "KB4554354", "KB4534321", "KB4537818", "KB4534273", "KB4532691", "KB4530715", "KB4520062", "KB4550969")
	case "18362":
		supersedence = append(supersedence, "KB4517389", "KB4541335", "KB4524570", "KB4528760", "KB4532693", "KB4554364", "KB4532695", "KB4551762", "KB4549951", "KB4535996", "KB4550945", "KB4530684", "KB4522355", "KB4540673")
	default:
		// higher versions should not be vulnerable and not need a specific patch installed
		// lower verions might be, but at some point it just gets too much work adding everything
		return false
	}
	// we are install in the function, so we have a potential vulnerable version
	// and need to have a kb installed, which means we need to have at least one intersecting
	// value between installed kbs and kbs that can fix this vuln
	for _, i := range kbs {
		for _, j := range supersedence {
			if i == j {
				// intersect exists, not vulnerable
				return false
			}
		}
	}
	// No suitable patch installed, therefore vulnerable
	return true
}

// Name returns the CVE Number
func (cve *CVE20191315) Name() string {
	return "CVE-2019-1315"
}

// Description returns the Description
func (cve *CVE20191315) Description() string {
	return "Windows Error Reporting: Arbitrary File Write as SYSTEM"
}

// CVE20191315 is the cve type
type CVE20191315 struct {
}
