package broker

import "testing"

// FuzzParseServiceEnvelope verifica que bytes arbitrários não causam panic
// durante o parsing protobuf e a validação estrutural básica do envelope.
func FuzzParseServiceEnvelope(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{0xff})
	f.Add([]byte{0x0a, 0x01, 'x'})
	f.Add([]byte{0x12, 0x00})

	f.Fuzz(func(t *testing.T, payload []byte) {
		processor := newPacketProcessor(nil)
		envelope, err := processor.parseServiceEnvelope(payload)
		if err != nil {
			return
		}
		if envelope == nil {
			t.Fatal("parseServiceEnvelope() retornou envelope nil sem erro")
		}
		_ = processor.validateServiceEnvelope(envelope)
	})
}
