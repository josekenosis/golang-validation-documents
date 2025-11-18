package utils

import (
	"regexp"
	"strconv"
)

func ValidateCPF(cpf string) bool {
	cpf = cleanDocument(cpf)

	if len(cpf) != 11 {
		return false
	}

	if allSameDigits(cpf) {
		return false
	}

	firstDigit := calculateCPFDigit(cpf[:9], 10)
	secondDigit := calculateCPFDigit(cpf[:10], 11)

	return cpf[9] == firstDigit && cpf[10] == secondDigit
}

func ValidateCNPJ(cnpj string) bool {
	cnpj = cleanDocument(cnpj)

	if len(cnpj) != 14 {
		return false
	}

	if allSameDigits(cnpj) {
		return false
	}

	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	firstDigit := calculateCNPJDigit(cnpj[:12], weights1)

	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	secondDigit := calculateCNPJDigit(cnpj[:13], weights2)

	return cnpj[12] == firstDigit && cnpj[13] == secondDigit
}

func cleanDocument(doc string) string {
	re := regexp.MustCompile(`\D`)
	return re.ReplaceAllString(doc, "")
}

func CleanDocument(doc string) string {
	return cleanDocument(doc)
}

func allSameDigits(doc string) bool {
	if len(doc) == 0 {
		return true
	}
	first := doc[0]
	for i := 1; i < len(doc); i++ {
		if doc[i] != first {
			return false
		}
	}
	return true
}

func calculateCPFDigit(doc string, multiplier int) byte {
	sum := 0
	for i := 0; i < len(doc); i++ {
		digit, _ := strconv.Atoi(string(doc[i]))
		sum += digit * multiplier
		multiplier--
	}
	remainder := sum % 11
	if remainder < 2 {
		return '0'
	}
	return byte('0' + (11 - remainder))
}

func calculateCNPJDigit(doc string, weights []int) byte {
	sum := 0
	for i := 0; i < len(doc); i++ {
		digit, _ := strconv.Atoi(string(doc[i]))
		sum += digit * weights[i]
	}
	remainder := sum % 11
	if remainder < 2 {
		return '0'
	}
	return byte('0' + (11 - remainder))
}

func FormatDocument(doc string) string {
	cleaned := cleanDocument(doc)
	if len(cleaned) == 11 {
		return cleaned[:3] + "." + cleaned[3:6] + "." + cleaned[6:9] + "-" + cleaned[9:]
	}
	if len(cleaned) == 14 {
		return cleaned[:2] + "." + cleaned[2:5] + "." + cleaned[5:8] + "/" + cleaned[8:12] + "-" + cleaned[12:]
	}
	return cleaned
}

func GetDocumentType(doc string) string {
	cleaned := cleanDocument(doc)
	if len(cleaned) == 11 {
		return "CPF"
	}
	if len(cleaned) == 14 {
		return "CNPJ"
	}
	return ""
}

func ValidateDocument(doc string) bool {
	cleaned := cleanDocument(doc)
	if len(cleaned) == 11 {
		return ValidateCPF(doc)
	}
	if len(cleaned) == 14 {
		return ValidateCNPJ(doc)
	}
	return false
}

