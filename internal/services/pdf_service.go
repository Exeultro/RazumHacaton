package services

import (
	"bytes"
	"fmt"

	"github.com/google/uuid"
)

type PDFService struct {
	cadreService *CadreService
}

func NewPDFService(cadreService *CadreService) *PDFService {
	return &PDFService{
		cadreService: cadreService,
	}
}

// GenerateCandidatePDF генерирует PDF с информацией о кандидате
func (s *PDFService) GenerateCandidatePDF(userID uuid.UUID) ([]byte, error) {
	candidate, err := s.cadreService.GetCandidateByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get candidate: %w", err)
	}
	if candidate == nil {
		return nil, fmt.Errorf("candidate not found")
	}

	// Формируем PDF (простейший вариант)
	var buf bytes.Buffer

	buf.WriteString("%PDF-1.4\n")
	buf.WriteString("1 0 obj\n")
	buf.WriteString("<< /Type /Catalog /Pages 2 0 R >>\n")
	buf.WriteString("endobj\n")
	buf.WriteString("2 0 obj\n")
	buf.WriteString("<< /Type /Pages /Kids [3 0 R] /Count 1 >>\n")
	buf.WriteString("endobj\n")
	buf.WriteString("3 0 obj\n")
	buf.WriteString("<< /Type /Page /Parent 2 0 R /MediaBox [0 0 595 842] /Contents 4 0 R >>\n")
	buf.WriteString("endobj\n")
	buf.WriteString("4 0 obj\n")
	buf.WriteString("<< /Length 100 >>\n")
	buf.WriteString("stream\n")
	buf.WriteString("BT /F1 24 Tf 100 700 Td (Отчет о кандидате) Tj ET\n")
	buf.WriteString(fmt.Sprintf("BT /F1 12 Tf 100 650 Td (ФИО: %s) Tj ET\n", candidate.FullName))
	buf.WriteString(fmt.Sprintf("BT /F1 12 Tf 100 630 Td (Email: %s) Tj ET\n", candidate.Email))
	buf.WriteString(fmt.Sprintf("BT /F1 12 Tf 100 610 Td (Баллов: %d) Tj ET\n", candidate.TotalPoints))
	buf.WriteString(fmt.Sprintf("BT /F1 12 Tf 100 590 Td (Мероприятий: %d) Tj ET\n", candidate.EventsCount))
	buf.WriteString(fmt.Sprintf("BT /F1 12 Tf 100 570 Td (Средний балл: %.2f) Tj ET\n", candidate.AvgPoints))
	if candidate.City != nil {
		buf.WriteString(fmt.Sprintf("BT /F1 12 Tf 100 550 Td (Город: %s) Tj ET\n", *candidate.City))
	}
	if candidate.Direction != nil {
		buf.WriteString(fmt.Sprintf("BT /F1 12 Tf 100 530 Td (Направление: %s) Tj ET\n", *candidate.Direction))
	}
	buf.WriteString("endstream\n")
	buf.WriteString("endobj\n")
	buf.WriteString("xref\n")
	buf.WriteString("0 5\n")
	buf.WriteString("0000000000 65535 f\n")
	buf.WriteString("0000000009 00000 n\n")
	buf.WriteString("0000000065 00000 n\n")
	buf.WriteString("0000000119 00000 n\n")
	buf.WriteString("0000000200 00000 n\n")
	buf.WriteString("trailer\n")
	buf.WriteString("<< /Size 5 /Root 1 0 R >>\n")
	buf.WriteString("startxref\n")
	buf.WriteString("400\n")
	buf.WriteString("%%EOF\n")

	return buf.Bytes(), nil
}
