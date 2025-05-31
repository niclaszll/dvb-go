package dvb

// Status represents the response status from DVB API calls.
// It contains information about whether the request was successful
// and provides error details when applicable.
type Status struct {
	// Code is the status code returned by the DVB API.
	// Typical values include success codes and various error codes
	// that indicate different types of failures or issues.
	Code string `json:"Code"`

	// Message provides human-readable information about the status.
	// This field contains error descriptions when the request fails,
	// or may be empty for successful requests.
	Message string `json:"Message,omitempty"`
}

// Diva contains DVB-specific identifiers
type Diva struct {
	// Number is the DIVA line number identifier.
	// This is a standardized identifier used across the transport network
	// for uniquely identifying transport lines.
	Number string `json:"Number"`

	// Network indicates which transport network this identifier belongs to.
	// For Dresden, this is typically "vvo" (Verkehrsverbund Oberelbe),
	// but may include other regional networks for connecting services.
	Network string `json:"Network"`
}

// Platform represents information about a physical platform or stop position
// where passengers board or alight from public transport vehicles.
type Platform struct {
	// Name is the display name or identifier of the platform.
	// This could be a platform number (e.g., "1", "2A") for train stations,
	// or a position identifier for bus/tram stops.
	Name string `json:"Name"`

	// Type indicates the kind of platform or stop.
	// Examples include "Platform" for railway platforms,
	// "Stop" for bus/tram stops, or other location-specific types.
	Type string `json:"Type"`
}
