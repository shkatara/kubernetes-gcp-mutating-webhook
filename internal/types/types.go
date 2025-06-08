package types

type AdmissionReview struct {
	APIVersion string             `json:"apiVersion,omitempty"`
	Kind       string             `json:"kind,omitempty"`
	Request    *AdmissionRequest  `json:"request,omitempty"`
	Response   *AdmissionResponse `json:"response,omitempty"`
}

type AdmissionRequest struct {
	UID       string `json:"uid"`
	Kind      Kind   `json:"kind"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Object    any    `json:"object"`
	// Plus other fields you might need, but these are the basics
}

type AdmissionResponse struct {
	UID       string  `json:"uid"`
	Allowed   bool    `json:"allowed"`
	Status    *Status `json:"status,omitempty"`
	PatchType *string `json:"patchType,omitempty"`
	Patch     []byte  `json:"patch,omitempty"`
}

type Status struct {
	Code    int32  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type Kind struct {
	Group   string `json:"group"`
	Version string `json:"version"`
	Kind    string `json:"kind"`
}
