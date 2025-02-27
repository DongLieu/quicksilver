package types

import (
	"bytes"
	"encoding/json"
	encoding_json "encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

var validLiquidData string = `{
	"chainid": "somechain-1",
	"registeredzonechainid": "someotherchain-1",
	"ibcdenom": "ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3",
	"qassetdenom": "uqstake"
}`

func TestAddProtocolDataProposal_ValidateBasic(t *testing.T) {
	type fields struct {
		Title       string
		Description string
		Protocol    string
		Type        string
		Key         string
		Data        json.RawMessage
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"blank",
			fields{},
			true,
		},
		{
			"invalid_protocol",
			fields{
				Title:       "Add Test Protocol",
				Description: "A new protocol for testing protocols",
				Protocol:    "",
				Type:        "",
				Key:         "",
				Data:        nil,
			},
			true,
		},
		{
			"invalid_type",
			fields{
				Title:       "Add Test Protocol",
				Description: "A new protocol for testing protocols",
				Protocol:    "TestProtocol",
				Type:        "",
				Key:         "",
				Data:        nil,
			},
			true,
		},
		{
			"invalid_key",
			fields{
				Title:       "Add Test Protocol",
				Description: "A new protocol for testing protocols",
				Protocol:    "TestProtocol",
				Type:        "TestType",
				Key:         "",
				Data:        nil,
			},
			true,
		},
		{
			"invalid_data",
			fields{
				Title:       "Add Test Protocol",
				Description: "A new protocol for testing protocols",
				Protocol:    "TestProtocol",
				Type:        "TestType",
				Key:         "TestKey",
				Data:        nil,
			},
			true,
		},
		{
			"valid_liquid_data",
			fields{
				Title:       "Valid Protocol Data",
				Description: "A valid protocol that is valid",
				Protocol:    "ValidProtocol",
				Type:        "liquidtoken",
				Key:         "liquid",
				Data:        []byte(validLiquidData),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := AddProtocolDataProposal{
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Type:        tt.fields.Type,
				Data:        tt.fields.Data,
				Key:         tt.fields.Key,
			}
			err := m.ValidateBasic()
			if tt.wantErr {
				t.Logf("Error:\n%v\n", err)
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestAddProtocolDataProposal_String(t *testing.T) {
	type fields struct {
		Title       string
		Description string
		Protocol    string
		Type        string
		Key         string
		Data        encoding_json.RawMessage
	}

	tt := fields{
		Title:       "Valid Protocol Data",
		Description: "A valid protocol that is valid",
		Protocol:    "ValidProtocol",
		Type:        "liquidtoken",
		Key:         "liquid",
		Data:        []byte(validLiquidData),
	}

	want := `Add Protocol Data Proposal:
Title:			Valid Protocol Data
Description:	A valid protocol that is valid
Type:			liquidtoken
Data:			{
	"chainid": "somechain-1",
	"registeredzonechainid": "someotherchain-1",
	"ibcdenom": "ibc/3020922B7576FC75BBE057A0290A9AEEFF489BB1113E6E365CE472D4BFB7FFA3",
	"qassetdenom": "uqstake"
}
Key:			liquid
`

	t.Run("stringer", func(t *testing.T) {
		m := AddProtocolDataProposal{
			Title:       tt.Title,
			Description: tt.Description,
			Type:        tt.Type,
			Data:        tt.Data,
			Key:         tt.Key,
		}
		got := m.String()
		require.Equal(t, want, got)
	})
}

var sink interface{}

func BenchmarkUpdateZoneProposalString(b *testing.B) {
	adp := &AddProtocolDataProposal{
		Title:       "Testing right here",
		Description: "Testing description",
		Key:         "This is my key",
		Data: bytes.Join(
			[][]byte{
				[]byte(`{"box":`),
				bytes.Repeat([]byte("{"), 1<<10),
				bytes.Repeat([]byte("}"), 1<<10),
				[]byte(`}`),
			},
			[]byte("")),
	}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		str := adp.String()
		b.SetBytes(int64(len(str)))
		sink = str
	}

	if sink == nil {
		b.Fatal("Benchmark did not run")
	}
	sink = (interface{})(nil)
}
