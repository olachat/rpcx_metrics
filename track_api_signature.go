package prom

import (
	"strconv"
)

type APISignatureType string

var (
	APISignatureTypeUnknown       APISignatureType = "unknown"
	APISignatureTypeError         APISignatureType = "error"
	APISignatureTypeValid         APISignatureType = "valid"
	APISignatureTypeInvalid       APISignatureType = "invalid"
	APISignatureTypeNoNonce       APISignatureType = "no_nonce"
	APISignatureTypeRepeatedNonce APISignatureType = "repeated_nonce"
)

func (o *Registerer) TrackAPISignature(version uint64, validType APISignatureType, path string) {
	verStr := strconv.FormatUint(version, 10)
	isValidStr := string(validType)

	// project, version, validity, path
	o.apiSignatureCounter.WithLabelValues(verStr, isValidStr, path).Inc()
}
