package v1

import "time"

// ServerVersion provides server version info
type ServerVersion struct {
	// Build is the server build version.
	Build string `protobuf:"bytes,1,opt,name=build,proto3" json:"build,omitempty"`
	// Runtime is the runtime version.
	Runtime string `protobuf:"bytes,2,opt,name=runtime,proto3" json:"runtime,omitempty"`
}

// ServerStatus provides response about current server
type ServerStatus struct {
	// Version is the app build version
	Version string `protobuf:"bytes,1,name=version,proto3" json:"version"`
	// HostName is operating system's host name
	HostName string `protobuf:"bytes,2,name=hostname,proto3" json:"hostname"`
	// Port is the listening port
	Port string `protobuf:"bytes,3,name=port,proto3" json:"port"`
	// StartedAt is the time when the server started
	StartedAt time.Time `protobuf:"bytes,4,name=started_at,proto3" json:"started_at"`
	// Uptime is the total time elapsed since the server started
	Uptime time.Duration `protobuf:"bytes,5,name=uptime,proto3" json:"uptime"`
}

// ServerStatusResponse is response for /v1/status
type ServerStatusResponse struct {
	// Status of the server.
	Status *ServerStatus `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
	// Version of the server.
	Version *ServerVersion `protobuf:"bytes,2,opt,name=version" json:"version,omitempty"`
}
