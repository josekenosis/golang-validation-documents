package utils

import "testing"

func TestValidateCPF(t *testing.T) {
	tests := []struct {
		name    string
		cpf     string
		want    bool
	}{
		{"valid CPF", "11144477735", true},
		{"valid CPF with formatting", "111.444.777-35", true},
		{"invalid CPF", "11144477736", false},
		{"invalid CPF same digits", "11111111111", false},
		{"invalid CPF short", "123456789", false},
		{"invalid CPF long", "123456789012", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateCPF(tt.cpf); got != tt.want {
				t.Errorf("ValidateCPF() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateCNPJ(t *testing.T) {
	tests := []struct {
		name    string
		cnpj    string
		want    bool
	}{
		{"valid CNPJ", "11222333000181", true},
		{"valid CNPJ with formatting", "11.222.333/0001-81", true},
		{"invalid CNPJ", "11222333000182", false},
		{"invalid CNPJ same digits", "11111111111111", false},
		{"invalid CNPJ short", "1234567890123", false},
		{"invalid CNPJ long", "123456789012345", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateCNPJ(tt.cnpj); got != tt.want {
				t.Errorf("ValidateCNPJ() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateDocument(t *testing.T) {
	tests := []struct {
		name    string
		doc     string
		want    bool
	}{
		{"valid CPF", "11144477735", true},
		{"valid CNPJ", "11222333000181", true},
		{"invalid document", "123456789", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateDocument(tt.doc); got != tt.want {
				t.Errorf("ValidateDocument() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDocumentType(t *testing.T) {
	tests := []struct {
		name    string
		doc     string
		want    string
	}{
		{"CPF", "11144477735", "CPF"},
		{"CNPJ", "11222333000181", "CNPJ"},
		{"invalid", "123", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDocumentType(tt.doc); got != tt.want {
				t.Errorf("GetDocumentType() = %v, want %v", got, tt.want)
			}
		})
	}
}

