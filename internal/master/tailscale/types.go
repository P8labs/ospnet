package tailscale

type CreateKeyRequest struct {
	Capabilities Capabilities `json:"capabilities"`
	Expiry       int          `json:"expirySeconds"`
}

type Capabilities struct {
	Devices Devices `json:"devices"`
}

type Devices struct {
	Create DeviceCreate `json:"create"`
}

type DeviceCreate struct {
	Reusable      bool     `json:"reusable"`
	Ephemeral     bool     `json:"ephemeral"`
	Preauthorized bool     `json:"preauthorized"`
	Tags          []string `json:"tags"`
}

type CreateKeyResponse struct {
	Key string `json:"key"`
}

type CreateKeyOptions struct {
	Tags      []string
	Ephemeral bool
	Reusable  bool
	Expiry    int
}
