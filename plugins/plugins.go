package plugins

import (
	"fmt"
	"net"
)

// Plugin ...
type Plugin interface {
	Auto() bool          // execute automatically on startup
	Execute(c net.Conn)  // execute
	Name() string        // name
	Description() string //description
}

var plugins = map[string]Plugin{}

// Init Initializes the plugin system
func Init(c net.Conn) {

	/* Modify Start */
	//flagGrabber := &FlagGrabber{}
	plugins = map[string]Plugin{
		//flagGrabber.Name(): flagGrabber,
	}

	/* Modify  End */

	// execute plugins that run on startup (auto)
	c.Write([]byte("\n[*] Auto-Plugins:\n"))
	for _, plugin := range plugins {
		if plugin.Auto() {
			c.Write([]byte(fmt.Sprintf(" - [%s]\n", plugin.Name())))
			plugin.Execute(c)
		}
	}
}

// Execute ...
func Execute(pluginName string, c net.Conn) {
	if _, ok := plugins[pluginName]; ok {
    	plugins[pluginName].Execute(c)
	} else {
		c.Write([]byte("[!] Plugin does not exist\n"))
	}	
}

// List ...
func List() []string {
	result := []string{}
	for _, plugin := range plugins {
		result = append(result, plugin.Name())
	}
	return result
}
